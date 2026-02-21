# API Reference

Complete REST API documentation for QRAP (Quantum Risk Assessment Platform).

## Table of Contents

- [Base URL](#base-url)
- [Authentication](#authentication)
  - [JWT Bearer Tokens](#jwt-bearer-tokens)
  - [API Keys](#api-keys)
- [Rate Limiting](#rate-limiting)
- [Pagination](#pagination)
- [Error Responses](#error-responses)
- [Endpoints](#endpoints)
  - [Health](#health)
  - [Organizations](#organizations)
  - [Assessments](#assessments)
  - [Findings](#findings)
  - [ML Engine -- Risk Scoring](#ml-engine----risk-scoring)
  - [ML Engine -- HNDL Calculator](#ml-engine----hndl-calculator)
  - [ML Engine -- Migration Planner](#ml-engine----migration-planner)

---

## Base URL

```
http://localhost:8083
```

All API endpoints are prefixed with `/api/v1/` except for the health check endpoint.

## Authentication

QRAP supports two authentication methods. When neither `QUANTUN_JWT_SECRET` nor `QUANTUN_API_KEYS` is configured, all endpoints are accessible without authentication.

### JWT Bearer Tokens

Include a JWT token in the `Authorization` header:

```
Authorization: Bearer <jwt-token>
```

**Generate a token (for testing/development):**

```bash
# Using the shared Go library's CreateJWT function, or any JWT library
# that supports HS256 with the following claims:

# Required claims:
#   sub (string) -- Subject (user or service ID)
#   exp (int)    -- Expiration time (Unix timestamp)
#   iss (string) -- Issuer (must match QUANTUN_JWT_ISSUER)

# Optional claims:
#   role (string) -- Role for authorization
#   iat  (int)    -- Issued at time
#   nbf  (int)    -- Not before time
#   jti  (string) -- Unique token ID
```

**Example request with JWT:**

```bash
curl -X GET http://localhost:8083/api/v1/organizations \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Configuration:**

| Variable             | Description                          |
|----------------------|--------------------------------------|
| `QUANTUN_JWT_SECRET` | HMAC-SHA256 signing secret           |
| `QUANTUN_JWT_ISSUER` | Expected `iss` claim (default: `quantun`) |

### API Keys

Include an API key in the `Authorization` header:

```
Authorization: ApiKey <key>
```

**Example request with API key:**

```bash
curl -X GET http://localhost:8083/api/v1/organizations \
  -H "Authorization: ApiKey my-secret-api-key"
```

**Configuration:**

API keys are configured via the `QUANTUN_API_KEYS` environment variable as comma-separated entries in the format `key:subject:role`:

```bash
QUANTUN_API_KEYS="sk-prod-abc123:service-a:admin,sk-prod-def456:service-b:reader"
```

Each entry consists of:
- `key` -- The secret API key string
- `subject` -- The identity associated with this key
- `role` -- The role granted to this key

## Rate Limiting

All endpoints are rate-limited to **100 requests per minute per IP address**.

**Rate limit headers** are included in every response:

| Header                  | Description                                    |
|-------------------------|------------------------------------------------|
| `X-RateLimit-Limit`    | Maximum requests allowed per window             |
| `X-RateLimit-Remaining`| Requests remaining in the current window        |
| `X-RateLimit-Reset`    | Seconds until the rate limit window resets       |

**When the rate limit is exceeded:**

```http
HTTP/1.1 429 Too Many Requests
Content-Type: application/json
Retry-After: 45
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 60

{"error": "rate limit exceeded"}
```

## Pagination

List endpoints support pagination via query parameters:

| Parameter | Default | Range     | Description                   |
|-----------|---------|-----------|-------------------------------|
| `offset`  | `0`     | 0-1000000 | Number of items to skip       |
| `limit`   | `20`    | 1-100     | Maximum items to return       |

**Example:**

```bash
curl "http://localhost:8083/api/v1/organizations?offset=20&limit=10"
```

**Paginated response format:**

```json
{
  "organizations": [...],
  "total_count": 42,
  "offset": 20,
  "limit": 10
}
```

## Error Responses

All errors return a JSON object with a single `error` field:

```json
{
  "error": "description of the error"
}
```

**Standard HTTP status codes:**

| Code | Description             | Common Causes                                    |
|------|-------------------------|--------------------------------------------------|
| 400  | Bad Request             | Missing required fields, invalid UUID, malformed JSON |
| 401  | Unauthorized            | Missing/invalid token, expired JWT, invalid API key  |
| 403  | Forbidden               | Valid auth but insufficient role                 |
| 404  | Not Found               | Resource does not exist                          |
| 413  | Payload Too Large       | Request body exceeds 1 MB                        |
| 429  | Too Many Requests       | Rate limit exceeded                              |
| 500  | Internal Server Error   | Database error, unexpected failure               |

---

## Endpoints

### Health

#### `GET /health`

Health check endpoint. Always unauthenticated.

**Response:**

```bash
curl http://localhost:8083/health
```

```json
{
  "status": "ok",
  "service": "qrap-api"
}
```

| Status Code | Description     |
|-------------|-----------------|
| 200         | Service healthy |

---

### Organizations

#### `POST /api/v1/organizations`

Create a new organization.

**Request body:**

| Field         | Type   | Required | Description                       |
|---------------|--------|----------|-----------------------------------|
| `name`        | string | Yes      | Organization name (max 255 chars, unique) |
| `description` | string | No       | Organization description          |
| `created_by`  | string | No       | Creator identity (defaults to auth subject) |

**Example:**

```bash
curl -X POST http://localhost:8083/api/v1/organizations \
  -H "Content-Type: application/json" \
  -H "Authorization: ApiKey my-key" \
  -d '{
    "name": "Acme Corp",
    "description": "Financial services company"
  }'
```

**Response (201 Created):**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Acme Corp",
  "description": "Financial services company",
  "created_at": "2026-01-15T10:30:00Z",
  "updated_at": "2026-01-15T10:30:00Z"
}
```

**Errors:**

| Code | Condition                                |
|------|------------------------------------------|
| 400  | Missing `name`, name exceeds 255 chars, invalid JSON |
| 401  | Missing or invalid authentication        |
| 500  | Database error (e.g., duplicate name)    |

---

#### `GET /api/v1/organizations`

List all organizations with pagination.

**Query parameters:**

| Parameter | Default | Description           |
|-----------|---------|-----------------------|
| `offset`  | 0       | Pagination offset     |
| `limit`   | 20      | Pagination limit (max 100) |

**Example:**

```bash
curl "http://localhost:8083/api/v1/organizations?offset=0&limit=10" \
  -H "Authorization: ApiKey my-key"
```

**Response (200 OK):**

```json
{
  "organizations": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Acme Corp",
      "description": "Financial services company",
      "created_at": "2026-01-15T10:30:00Z",
      "updated_at": "2026-01-15T10:30:00Z"
    }
  ],
  "total_count": 1,
  "offset": 0,
  "limit": 10
}
```

---

#### `GET /api/v1/organizations/{id}`

Get a single organization by ID.

**Path parameters:**

| Parameter | Type | Description       |
|-----------|------|-------------------|
| `id`      | UUID | Organization UUID |

**Example:**

```bash
curl http://localhost:8083/api/v1/organizations/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: ApiKey my-key"
```

**Response (200 OK):**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Acme Corp",
  "description": "Financial services company",
  "created_at": "2026-01-15T10:30:00Z",
  "updated_at": "2026-01-15T10:30:00Z"
}
```

**Errors:**

| Code | Condition            |
|------|----------------------|
| 400  | Invalid UUID format  |
| 404  | Organization not found |

---

### Assessments

#### `POST /api/v1/assessments`

Create a new assessment in DRAFT status.

**Request body:**

| Field             | Type     | Required | Description                            |
|-------------------|----------|----------|----------------------------------------|
| `name`            | string   | Yes      | Assessment name (max 255 chars)        |
| `organization_id` | string  | Yes      | Organization UUID                      |
| `target_assets`   | string[] | No      | List of assets to scan                 |
| `created_by`      | string  | No       | Creator identity (defaults to auth subject) |

**Example:**

```bash
curl -X POST http://localhost:8083/api/v1/assessments \
  -H "Content-Type: application/json" \
  -H "Authorization: ApiKey my-key" \
  -d '{
    "name": "Q1 2026 Crypto Audit",
    "organization_id": "550e8400-e29b-41d4-a716-446655440000",
    "target_assets": ["api-gateway", "payment-service", "auth-service"]
  }'
```

**Response (201 Created):**

```json
{
  "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
  "name": "Q1 2026 Crypto Audit",
  "organization_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "DRAFT",
  "risk_score": 0,
  "target_assets": ["api-gateway", "payment-service", "auth-service"],
  "created_at": "2026-01-15T11:00:00Z",
  "updated_at": "2026-01-15T11:00:00Z"
}
```

**Errors:**

| Code | Condition                                       |
|------|-------------------------------------------------|
| 400  | Missing `name` or `organization_id`, invalid UUID, name too long |
| 401  | Missing or invalid authentication               |
| 500  | Database error                                  |

---

#### `GET /api/v1/assessments`

List assessments with optional filtering and pagination.

**Query parameters:**

| Parameter         | Default | Description                           |
|-------------------|---------|---------------------------------------|
| `offset`          | 0       | Pagination offset                     |
| `limit`           | 20      | Pagination limit (max 100)            |
| `organization_id` | --     | Filter by organization UUID           |
| `status`          | --     | Filter by status (DRAFT, IN_PROGRESS, COMPLETED, ARCHIVED) |

**Example:**

```bash
curl "http://localhost:8083/api/v1/assessments?organization_id=550e8400-e29b-41d4-a716-446655440000&status=COMPLETED" \
  -H "Authorization: ApiKey my-key"
```

**Response (200 OK):**

```json
{
  "assessments": [
    {
      "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
      "name": "Q1 2026 Crypto Audit",
      "organization_id": "550e8400-e29b-41d4-a716-446655440000",
      "status": "COMPLETED",
      "overall_risk": "CRITICAL",
      "risk_score": 75.0,
      "target_assets": ["api-gateway", "payment-service", "auth-service"],
      "created_at": "2026-01-15T11:00:00Z",
      "updated_at": "2026-01-15T11:05:00Z"
    }
  ],
  "total_count": 1,
  "offset": 0,
  "limit": 20
}
```

---

#### `GET /api/v1/assessments/{id}`

Get a single assessment with a finding summary.

**Path parameters:**

| Parameter | Type | Description     |
|-----------|------|-----------------|
| `id`      | UUID | Assessment UUID |

**Example:**

```bash
curl http://localhost:8083/api/v1/assessments/7c9e6679-7425-40de-944b-e07fc1f90ae7 \
  -H "Authorization: ApiKey my-key"
```

**Response (200 OK):**

```json
{
  "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
  "name": "Q1 2026 Crypto Audit",
  "organization_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "COMPLETED",
  "overall_risk": "CRITICAL",
  "risk_score": 75.0,
  "target_assets": ["api-gateway", "payment-service", "auth-service"],
  "summary": {
    "total_findings": 6,
    "critical_findings": 3,
    "high_findings": 3,
    "medium_findings": 0,
    "low_findings": 0,
    "pqc_readiness_percentage": 0.0,
    "assets_scanned": 3
  },
  "created_at": "2026-01-15T11:00:00Z",
  "updated_at": "2026-01-15T11:05:00Z"
}
```

**Errors:**

| Code | Condition            |
|------|----------------------|
| 400  | Invalid UUID format  |
| 404  | Assessment not found |

---

#### `POST /api/v1/assessments/{id}/run`

Execute an assessment. Analyzes all target assets, generates findings, calculates risk scores, and updates the assessment status to COMPLETED.

**Path parameters:**

| Parameter | Type | Description     |
|-----------|------|-----------------|
| `id`      | UUID | Assessment UUID |

**Preconditions:**
- Assessment must be in `DRAFT` or `COMPLETED` status

**Example:**

```bash
curl -X POST http://localhost:8083/api/v1/assessments/7c9e6679-7425-40de-944b-e07fc1f90ae7/run \
  -H "Authorization: ApiKey my-key"
```

**Response (200 OK):**

```json
{
  "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
  "name": "Q1 2026 Crypto Audit",
  "organization_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "COMPLETED",
  "overall_risk": "CRITICAL",
  "risk_score": 75.0,
  "target_assets": ["api-gateway", "payment-service", "auth-service"],
  "created_at": "2026-01-15T11:00:00Z",
  "updated_at": "2026-01-15T11:05:00Z"
}
```

**Errors:**

| Code | Condition                                      |
|------|------------------------------------------------|
| 400  | Invalid UUID format                            |
| 404  | Assessment not found                           |
| 500  | Assessment in invalid status, analysis failure |

---

### Findings

#### `GET /api/v1/findings`

List findings for a specific assessment. The `assessment_id` query parameter is required.

**Query parameters:**

| Parameter       | Default | Required | Description                       |
|-----------------|---------|----------|-----------------------------------|
| `assessment_id` | --     | Yes      | Assessment UUID to filter by      |
| `risk_level`    | --     | No       | Filter: CRITICAL, HIGH, MEDIUM, LOW, INFO |
| `category`      | --     | No       | Filter: WEAK_ALGORITHM, SHORT_KEY_LENGTH, DEPRECATED_PROTOCOL, MISSING_PQC, CERTIFICATE_EXPIRY, HARVEST_NOW_DECRYPT_LATER |
| `offset`        | 0       | No       | Pagination offset                 |
| `limit`         | 20      | No       | Pagination limit (max 100)        |

**Example:**

```bash
curl "http://localhost:8083/api/v1/findings?assessment_id=7c9e6679-7425-40de-944b-e07fc1f90ae7&risk_level=CRITICAL" \
  -H "Authorization: ApiKey my-key"
```

**Response (200 OK):**

```json
{
  "findings": [
    {
      "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "assessment_id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
      "category": "HARVEST_NOW_DECRYPT_LATER",
      "risk_level": "CRITICAL",
      "title": "HNDL risk on api-gateway",
      "description": "Asset api-gateway is vulnerable to harvest-now-decrypt-later attacks",
      "affected_asset": "api-gateway",
      "current_algorithm": "RSA-2048",
      "recommended_algorithm": "ML-KEM-768",
      "remediation": "Prioritise migration of long-lived secrets; data encrypted today can be captured and decrypted later by quantum computers",
      "discovered_at": "2026-01-15T11:05:00Z"
    }
  ],
  "total_count": 3,
  "offset": 0,
  "limit": 20
}
```

**Errors:**

| Code | Condition                                |
|------|------------------------------------------|
| 400  | Missing `assessment_id`, invalid UUID    |
| 401  | Missing or invalid authentication        |
| 500  | Database error                           |

---

#### `GET /api/v1/findings/{id}`

Get a single finding by ID.

**Path parameters:**

| Parameter | Type | Description  |
|-----------|------|--------------|
| `id`      | UUID | Finding UUID |

**Example:**

```bash
curl http://localhost:8083/api/v1/findings/a1b2c3d4-e5f6-7890-abcd-ef1234567890 \
  -H "Authorization: ApiKey my-key"
```

**Response (200 OK):**

```json
{
  "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "assessment_id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
  "category": "HARVEST_NOW_DECRYPT_LATER",
  "risk_level": "CRITICAL",
  "title": "HNDL risk on api-gateway",
  "description": "Asset api-gateway is vulnerable to harvest-now-decrypt-later attacks",
  "affected_asset": "api-gateway",
  "current_algorithm": "RSA-2048",
  "recommended_algorithm": "ML-KEM-768",
  "remediation": "Prioritise migration of long-lived secrets; data encrypted today can be captured and decrypted later by quantum computers",
  "discovered_at": "2026-01-15T11:05:00Z"
}
```

**Errors:**

| Code | Condition           |
|------|---------------------|
| 400  | Invalid UUID format |
| 404  | Finding not found   |

---

### ML Engine -- Risk Scoring

#### `POST /api/v1/score`

Calculate a composite risk score from a set of findings.

**Request body:**

| Field           | Type       | Required | Description                      |
|-----------------|------------|----------|----------------------------------|
| `findings`      | Finding[]  | Yes      | Array of finding objects          |
| `total_assets`  | integer    | Yes      | Total number of assets scanned    |

**Finding object:**

| Field                   | Type   | Required | Description                      |
|-------------------------|--------|----------|----------------------------------|
| `category`              | string | Yes      | Finding category enum value       |
| `risk_level`            | string | Yes      | Risk level enum value             |
| `affected_asset`        | string | Yes      | Name of affected asset            |
| `current_algorithm`     | string | No       | Currently used algorithm          |
| `recommended_algorithm` | string | No       | Recommended PQC replacement       |

**Example:**

```bash
curl -X POST http://localhost:8083/api/v1/score \
  -H "Content-Type: application/json" \
  -d '{
    "findings": [
      {
        "category": "MISSING_PQC",
        "risk_level": "HIGH",
        "affected_asset": "api-gateway",
        "current_algorithm": "RSA-2048",
        "recommended_algorithm": "ML-KEM-768"
      },
      {
        "category": "HARVEST_NOW_DECRYPT_LATER",
        "risk_level": "CRITICAL",
        "affected_asset": "api-gateway",
        "current_algorithm": "RSA-2048",
        "recommended_algorithm": "ML-KEM-768"
      }
    ],
    "total_assets": 5
  }'
```

**Response (200 OK):**

```json
{
  "risk_score": 66.67,
  "overall_risk": "HIGH",
  "pqc_readiness": 80.0,
  "finding_breakdown": {
    "HIGH": 1,
    "CRITICAL": 1
  }
}
```

**Score interpretation:**

| Score Range | Overall Risk | Description                         |
|-------------|-------------|-------------------------------------|
| 80-100      | CRITICAL    | Immediate action required           |
| 60-79       | HIGH        | Urgent migration recommended        |
| 30-59       | MEDIUM      | Plan migration within 12 months     |
| 0-29        | LOW         | Acceptable risk posture             |

---

### ML Engine -- HNDL Calculator

#### `POST /api/v1/hndl`

Calculate Harvest Now, Decrypt Later risk for a specific algorithm.

**Request body:**

| Field                   | Type    | Required | Default | Description                          |
|-------------------------|---------|----------|---------|--------------------------------------|
| `algorithm`             | string  | Yes      | --     | Algorithm to evaluate                |
| `data_shelf_life_years` | integer | No       | 10      | How many years the data must remain secret |

**Supported algorithms:**

| Algorithm    | Estimated Break Year | Type      |
|-------------|---------------------|-----------|
| RSA-2048    | 2030                | Classical |
| RSA-3072    | 2032                | Classical |
| RSA-4096    | 2035                | Classical |
| ECDSA-P256  | 2030                | Classical |
| ECDSA-P384  | 2032                | Classical |
| Ed25519     | 2030                | Classical |
| X25519      | 2030                | Classical |
| DH-2048     | 2030                | Classical |
| AES-128     | 2040                | Symmetric |
| AES-256     | 2060                | Symmetric |
| ML-KEM-*    | 2080                | PQC       |
| ML-DSA-*    | 2080                | PQC       |

**Example:**

```bash
curl -X POST http://localhost:8083/api/v1/hndl \
  -H "Content-Type: application/json" \
  -d '{
    "algorithm": "RSA-2048",
    "data_shelf_life_years": 15
  }'
```

**Response (200 OK):**

```json
{
  "algorithm": "RSA-2048",
  "estimated_break_year": 2030,
  "data_shelf_life_years": 15,
  "risk_window_years": 11,
  "is_at_risk": true,
  "urgency": "CRITICAL"
}
```

**Urgency levels:**

| Risk Window (years) | Urgency  | Action                          |
|---------------------|----------|---------------------------------|
| >= 10               | CRITICAL | Migrate immediately             |
| >= 5                | HIGH     | Migrate within 6 months         |
| > 0                 | MEDIUM   | Plan migration within 12 months |
| <= 0                | LOW      | Monitor and plan ahead          |

---

### ML Engine -- Migration Planner

#### `POST /api/v1/migration-plan`

Generate a prioritized PQC migration plan for a set of assets.

**Request body:**

| Field    | Type    | Required | Description                |
|----------|---------|----------|----------------------------|
| `assets` | Asset[] | Yes      | Array of asset objects      |

**Asset object:**

| Field       | Type   | Required | Default  | Description                       |
|-------------|--------|----------|----------|-----------------------------------|
| `asset`     | string | Yes      | --      | Asset name or identifier          |
| `algorithm` | string | Yes      | --      | Current algorithm in use          |
| `urgency`   | string | No       | MEDIUM   | Urgency level (CRITICAL/HIGH/MEDIUM/LOW) |

**Example:**

```bash
curl -X POST http://localhost:8083/api/v1/migration-plan \
  -H "Content-Type: application/json" \
  -d '{
    "assets": [
      {"asset": "payment-api", "algorithm": "RSA-2048", "urgency": "CRITICAL"},
      {"asset": "auth-service", "algorithm": "ECDSA-P256", "urgency": "HIGH"},
      {"asset": "internal-api", "algorithm": "RSA-4096", "urgency": "MEDIUM"},
      {"asset": "vpn-gateway", "algorithm": "X25519", "urgency": "HIGH"}
    ]
  }'
```

**Response (200 OK):**

```json
{
  "steps": [
    {
      "asset": "payment-api",
      "current_algorithm": "RSA-2048",
      "target_algorithm": "ML-KEM-768",
      "priority": "CRITICAL",
      "estimated_effort": "MEDIUM",
      "notes": "Migrate from RSA-2048 to ML-KEM-768"
    },
    {
      "asset": "auth-service",
      "current_algorithm": "ECDSA-P256",
      "target_algorithm": "ML-DSA-65",
      "priority": "HIGH",
      "estimated_effort": "MEDIUM",
      "notes": "Migrate from ECDSA-P256 to ML-DSA-65"
    },
    {
      "asset": "vpn-gateway",
      "current_algorithm": "X25519",
      "target_algorithm": "X25519-ML-KEM-768",
      "priority": "HIGH",
      "estimated_effort": "HIGH",
      "notes": "Migrate from X25519 to X25519-ML-KEM-768"
    },
    {
      "asset": "internal-api",
      "current_algorithm": "RSA-4096",
      "target_algorithm": "ML-KEM-1024",
      "priority": "MEDIUM",
      "estimated_effort": "MEDIUM",
      "notes": "Migrate from RSA-4096 to ML-KEM-1024"
    }
  ],
  "total_assets": 4,
  "critical_count": 1,
  "estimated_phases": 1
}
```

**Notes:**
- Steps are sorted by priority (CRITICAL first, LOW last)
- Assets already using PQC algorithms are excluded from the plan
- `estimated_phases` is calculated as 1 phase per 5 assets (minimum 1)
- Hybrid migrations (e.g., X25519 to X25519-ML-KEM-768) are marked as HIGH effort

---

## ML Engine Direct Access

The ML engine also exposes its own health endpoint:

```bash
curl http://localhost:8084/health
# {"status": "ok", "service": "qrap-ml"}
```

The ML engine endpoints (`/api/v1/score`, `/api/v1/hndl`, `/api/v1/migration-plan`) can be accessed directly on port 8084 or proxied through the Go API on port 8083.

---

## Typical API Workflow

The following diagram illustrates the complete end-to-end workflow for performing a quantum risk assessment:

```mermaid
sequenceDiagram
    participant C as Client
    participant API as Go API :8083
    participant ML as ML Engine :8084
    participant DB as PostgreSQL

    Note over C,DB: 1. Setup
    C->>API: POST /api/v1/organizations
    API->>DB: INSERT organization
    API-->>C: 201 Created {org_id}

    Note over C,DB: 2. Create Assessment
    C->>API: POST /api/v1/assessments<br/>{name, org_id, target_assets}
    API->>DB: INSERT assessment (DRAFT)
    API-->>C: 201 Created {assessment_id}

    Note over C,DB: 3. Execute Assessment
    C->>API: POST /api/v1/assessments/{id}/run
    API->>DB: UPDATE status → IN_PROGRESS
    API->>DB: INSERT findings (batch)
    API->>ML: POST /api/v1/score {findings}
    ML-->>API: {risk_score, pqc_readiness}
    API->>DB: UPDATE assessment → COMPLETED
    API-->>C: 200 OK {assessment}

    Note over C,DB: 4. Review Results
    C->>API: GET /api/v1/assessments/{id}
    API-->>C: 200 OK {assessment + summary}
    C->>API: GET /api/v1/findings?assessment_id={id}
    API-->>C: 200 OK {findings list}

    Note over C,DB: 5. Analyze & Plan
    C->>API: POST /api/v1/hndl<br/>{algorithm, data_shelf_life}
    API->>ML: Forward to ML Engine
    ML-->>API: {risk_window, urgency}
    API-->>C: 200 OK {HNDL analysis}
    C->>API: POST /api/v1/migration-plan<br/>{assets}
    API->>ML: Forward to ML Engine
    ML-->>API: {prioritized steps}
    API-->>C: 200 OK {migration plan}
```
