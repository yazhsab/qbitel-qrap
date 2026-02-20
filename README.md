<p align="center">
  <h1 align="center">QRAP - Quantum Risk Assessment Platform</h1>
  <p align="center">
    <strong>Qualys/Tenable for the Quantum Era</strong>
  </p>
  <p align="center">
    Identify vulnerable cryptographic algorithms. Evaluate Harvest Now, Decrypt Later exposure. Generate prioritized PQC migration plans.
  </p>
</p>

<p align="center">
  <a href="https://github.com/quantun-opensource/qrap/actions/workflows/ci.yml"><img src="https://github.com/quantun-opensource/qrap/actions/workflows/ci.yml/badge.svg" alt="CI"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-Apache_2.0-blue.svg" alt="License: Apache 2.0"></a>
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.23-00ADD8?logo=go&logoColor=white" alt="Go 1.23"></a>
  <a href="https://python.org/"><img src="https://img.shields.io/badge/Python-3.11+-3776AB?logo=python&logoColor=white" alt="Python 3.11+"></a>
  <a href="https://react.dev/"><img src="https://img.shields.io/badge/React-19-61DAFB?logo=react&logoColor=black" alt="React 19"></a>
  <a href="https://www.postgresql.org/"><img src="https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql&logoColor=white" alt="PostgreSQL 16"></a>
</p>

---

## Overview

QRAP is an open-source platform for assessing and managing quantum computing risks to your organization's cryptographic infrastructure. As quantum computing advances, classical cryptographic algorithms like RSA and ECDSA face existential threats. QRAP helps you get ahead of this transition by:

- **Scanning** your cryptographic assets to identify quantum-vulnerable algorithms
- **Quantifying** the risk using composite scoring (0--100) weighted by severity, algorithm weakness, and exposure surface
- **Evaluating** Harvest Now, Decrypt Later (HNDL) exposure using the Mosca inequality
- **Planning** a prioritized migration roadmap from classical to post-quantum cryptography (ML-KEM, ML-DSA per FIPS 203/204)

## Key Features

### Cryptographic Risk Assessment
Scan your infrastructure to discover cryptographic assets and identify vulnerabilities. QRAP categorizes findings across six risk categories: weak algorithms, short key lengths, deprecated protocols, missing PQC, certificate expiry, and HNDL exposure.

### HNDL Risk Analysis
Calculate Harvest Now, Decrypt Later risk windows using conservative CRQC (Cryptographically Relevant Quantum Computer) timeline estimates. Understand which encrypted data captured today could be decrypted by future quantum computers.

### PQC Migration Planning
Generate prioritized migration plans mapping classical algorithms to their NIST-standardized post-quantum replacements (RSA to ML-KEM, ECDSA to ML-DSA). Plans include effort estimates and phased rollout schedules.

### Composite Risk Scoring
ML-powered risk scoring engine produces a 0--100 composite score, accounting for finding severity, category multipliers (HNDL findings weighted 1.5x), and PQC readiness percentages.

### REST API with Enterprise Auth
Full-featured Go API with JWT Bearer token and API key authentication, rate limiting (100 req/min), security headers (HSTS, CSP, X-Frame-Options), and CORS support.

### Interactive Web Dashboard
React-based dashboard for visualizing assessments, browsing findings by risk level, and running HNDL analysis interactively.

## Architecture

```
                         +-------------------+
                         |   Web Dashboard   |
                         | React 19 + Vite 6 |
                         |   localhost:3002   |
                         +---------+---------+
                                   |
                                   | HTTP/REST
                                   v
+-------------------+    +-------------------+    +-------------------+
|                   |    |                   |    |                   |
|   PostgreSQL 16   |<-->|    Go REST API    |--->|  Python ML Engine |
|                   |    |  Chi v5 Router    |    | FastAPI + uvicorn |
|  - organizations  |    |  localhost:8083   |    |  localhost:8084   |
|  - assessments    |    |                   |    |                   |
|  - findings       |    | Auth | RateLimit  |    | - Risk Scorer     |
|  - audit_log      |    | CORS | Security   |    | - HNDL Calculator |
|                   |    |                   |    | - Migration Plan  |
+-------------------+    +-------------------+    +-------------------+
```

## Quick Start

### Prerequisites

| Tool       | Version | Purpose              |
|------------|---------|----------------------|
| Go         | 1.23+   | API server           |
| Python     | 3.11+   | ML engine            |
| Node.js    | 22+     | Web dashboard        |
| PostgreSQL | 16+     | Database             |
| Docker     | 24+     | Containerized deploy |

### Option 1: Docker Compose (Recommended)

Start the entire stack with a single command:

```bash
# Clone the repository
git clone https://github.com/quantun-opensource/qrap.git
cd qrap

# Start all services (API, ML engine, web dashboard, PostgreSQL)
docker compose -f infra/docker/docker-compose.yml up -d

# Run database migrations
export QRAP_DATABASE_URL="postgres://quantun:quantun_dev@localhost:5432/qrap?sslmode=disable"
make migrate
```

Services will be available at:

| Service       | URL                    |
|---------------|------------------------|
| API           | http://localhost:8083   |
| ML Engine     | http://localhost:8084   |
| Web Dashboard | http://localhost:3002   |

### Option 2: Local Development

```bash
# Install all dependencies (Go, Python venv, Node modules)
make setup

# Start PostgreSQL only
make docker-deps

# Run database migrations
export QRAP_DATABASE_URL="postgres://quantun:quantun_dev@localhost:5432/qrap?sslmode=disable"
make migrate

# Terminal 1: Start the API server
cd api && go run ./cmd/server

# Terminal 2: Start the ML engine
cd ml && .venv/bin/uvicorn qrap_ml.api.app:app --port 8084 --reload

# Terminal 3: Start the web dev server
cd web && npm run dev
```

### Verify Installation

```bash
# Health check -- API
curl http://localhost:8083/health
# {"status":"ok","service":"qrap-api"}

# Health check -- ML Engine
curl http://localhost:8084/health
# {"status":"ok","service":"qrap-ml"}
```

## API Reference

All API endpoints (except `/health`) require authentication when configured. See [docs/API.md](docs/API.md) for the complete API reference.

### Quick Examples

**Create an organization:**
```bash
curl -X POST http://localhost:8083/api/v1/organizations \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"name": "Acme Corp", "description": "Financial services"}'
```

**Create and run an assessment:**
```bash
# Create
curl -X POST http://localhost:8083/api/v1/assessments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "Q1 2026 Crypto Audit",
    "organization_id": "<org-uuid>",
    "target_assets": ["api-gateway", "payment-service", "auth-service"]
  }'

# Run
curl -X POST http://localhost:8083/api/v1/assessments/<assessment-uuid>/run \
  -H "Authorization: Bearer <token>"
```

**Calculate HNDL risk:**
```bash
curl -X POST http://localhost:8083/api/v1/hndl \
  -H "Content-Type: application/json" \
  -d '{"algorithm": "RSA-2048", "data_shelf_life_years": 15}'
```

### Endpoint Summary

| Method | Path                               | Description                    |
|--------|-------------------------------------|-------------------------------|
| GET    | `/health`                           | Health check (no auth)        |
| POST   | `/api/v1/organizations`             | Create organization           |
| GET    | `/api/v1/organizations`             | List organizations            |
| GET    | `/api/v1/organizations/{id}`        | Get organization              |
| POST   | `/api/v1/assessments`               | Create assessment             |
| GET    | `/api/v1/assessments`               | List assessments              |
| GET    | `/api/v1/assessments/{id}`          | Get assessment with summary   |
| POST   | `/api/v1/assessments/{id}/run`      | Run assessment                |
| GET    | `/api/v1/findings`                  | List findings (assessment_id) |
| GET    | `/api/v1/findings/{id}`             | Get finding                   |
| POST   | `/api/v1/score`                     | Calculate risk score          |
| POST   | `/api/v1/hndl`                      | Calculate HNDL risk           |
| POST   | `/api/v1/migration-plan`            | Generate migration plan       |

## Configuration

All configuration is via environment variables. Copy `.env.example` to `.env` to get started.

| Variable              | Default                    | Description                                |
|-----------------------|----------------------------|--------------------------------------------|
| `QRAP_PORT`           | `8083`                     | API server port                            |
| `QRAP_DATABASE_URL`   | *(required)*               | PostgreSQL connection string               |
| `QRAP_ML_ENGINE_URL`  | `http://127.0.0.1:8084`   | ML engine URL                              |
| `QRAP_LOG_LEVEL`      | `info`                     | Log level (debug, info, warn, error)       |
| `QUANTUN_JWT_SECRET`  | *(empty -- auth disabled)* | HMAC-SHA256 secret for JWT validation      |
| `QUANTUN_JWT_ISSUER`  | `quantun`                  | Expected JWT `iss` claim                   |
| `QUANTUN_API_KEYS`    | *(empty)*                  | Comma-separated `key:subject:role` entries |
| `QUANTUN_CORS_ORIGINS`| *(empty)*                  | Comma-separated allowed CORS origins       |

> **Note:** When neither `QUANTUN_JWT_SECRET` nor `QUANTUN_API_KEYS` is set, all endpoints are accessible without authentication. This is convenient for development but should never be used in production.

## Testing

```bash
# Run the full test suite (Go + Python + TypeScript)
make test

# Individual test suites
make test-go       # Go unit tests (api/ + shared/go/)
make test-python   # Python tests (ml/)
make test-node     # TypeScript type checking (web/)

# Linting
make lint          # Go (golangci-lint) + Python (ruff)
```

## Project Structure

```
qrap/
+-- api/                             Go REST API
|   +-- cmd/
|   |   +-- server/main.go          Server entrypoint
|   |   +-- migrate/main.go         Migration CLI helper
|   +-- internal/
|   |   +-- config/                  Environment configuration
|   |   +-- handler/                 HTTP handlers (health, org, assessment, finding)
|   |   +-- model/                   Data models and response types
|   |   +-- repository/              PostgreSQL repositories (pgx)
|   |   +-- service/                 Business logic layer
|   +-- Dockerfile
|   +-- go.mod
+-- ml/                              Python ML Engine
|   +-- src/qrap_ml/
|   |   +-- api/app.py              FastAPI application
|   |   +-- risk_scorer/            Composite risk scoring engine
|   |   +-- hndl_calculator/        HNDL risk calculator (Mosca inequality)
|   |   +-- migration_planner/      PQC migration roadmap generator
|   +-- tests/                       Pytest test suite
|   +-- pyproject.toml
|   +-- Dockerfile
+-- web/                             React Dashboard
|   +-- src/
|   |   +-- App.tsx                  Router and navigation
|   |   +-- pages/                   Dashboard, Assessments, Findings, HNDL
|   +-- index.html
|   +-- package.json
|   +-- vite.config.ts
|   +-- Dockerfile
+-- shared/go/                       Shared Go Libraries
|   +-- database/pool.go            PostgreSQL connection pool (pgx)
|   +-- middleware/
|       +-- auth.go                  JWT + API key authentication
|       +-- ratelimit.go             Per-IP rate limiting
|       +-- security.go              Security headers, CORS, body size limit
|       +-- pagination.go            Query parameter pagination
+-- db/migrations/                   PostgreSQL schema migrations
+-- infra/docker/                    Docker Compose files
+-- .github/workflows/ci.yml        GitHub Actions CI pipeline
+-- .env.example                     Environment variable template
+-- Makefile                         Build, test, and deploy targets
+-- go.work                          Go workspace configuration
+-- LICENSE                          Apache License 2.0
```

## Documentation

| Document                                    | Description                                 |
|---------------------------------------------|---------------------------------------------|
| [Architecture](docs/ARCHITECTURE.md)        | System design, data flow, and tech choices  |
| [API Reference](docs/API.md)               | Complete REST API documentation             |
| [Deployment Guide](docs/DEPLOYMENT.md)     | Production deployment and operations        |
| [Development Guide](docs/DEVELOPMENT.md)   | Local setup, coding standards, debugging    |
| [Contributing](CONTRIBUTING.md)            | How to contribute to QRAP                   |
| [Security Policy](SECURITY.md)             | Vulnerability reporting and security info   |
| [Changelog](CHANGELOG.md)                  | Release history                             |

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on:

- Reporting bugs and requesting features
- Setting up the development environment
- Submitting pull requests
- Coding standards and commit message format

## Security

QRAP is a security-focused platform. If you discover a vulnerability, please follow our [responsible disclosure process](SECURITY.md) rather than opening a public issue.

Key security features:
- HMAC-SHA256 JWT validation with constant-time comparison
- API key validation resistant to timing side-channel attacks
- Rate limiting (100 requests/minute per IP)
- Security headers (HSTS, CSP, X-Frame-Options, X-Content-Type-Options)
- Request body size limits (1 MB)
- Parameterized SQL queries (pgx)
- Graceful shutdown with connection draining

## License

Copyright 2026 Quantun

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for the full license text.

---

<p align="center">
  Built with care by the <a href="https://github.com/quantun-opensource">Quantun</a> team.
</p>
