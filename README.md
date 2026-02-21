<div align="center">

<!-- HERO BANNER -->
<br/>

<img src="https://img.shields.io/badge/Q-R_A_P-000000?style=for-the-badge&labelColor=6C3AED&color=000000" alt="QRAP" height="60"/>

# Quantum Risk Assessment Platform

### *Qualys/Tenable for the Quantum Era*

<br/>

**Identify vulnerable cryptography. Quantify harvest-now-decrypt-later exposure.<br/>Generate prioritized post-quantum migration plans&mdash;before it's too late.**

<br/>

<!-- BADGE ROW 1: Build & Status -->
[![CI](https://github.com/yazhsab/qbitel-qrap/actions/workflows/ci.yml/badge.svg)](https://github.com/yazhsab/qbitel-qrap/actions/workflows/ci.yml)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache_2.0-blue?style=flat-square)](LICENSE)
[![GitHub Stars](https://img.shields.io/github/stars/yazhsab/qbitel-qrap?style=flat-square&logo=github&label=Stars)](https://github.com/yazhsab/qbitel-qrap)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen?style=flat-square)](CONTRIBUTING.md)

<!-- BADGE ROW 2: Tech Stack -->
[![Go 1.23](https://img.shields.io/badge/Go-1.23-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev/)
[![Python 3.11+](https://img.shields.io/badge/Python-3.11+-3776AB?style=flat-square&logo=python&logoColor=white)](https://python.org/)
[![React 19](https://img.shields.io/badge/React-19-61DAFB?style=flat-square&logo=react&logoColor=black)](https://react.dev/)
[![PostgreSQL 16](https://img.shields.io/badge/PostgreSQL-16-4169E1?style=flat-square&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat-square&logo=docker&logoColor=white)](https://docker.com/)

<!-- BADGE ROW 3: Qbitel -->
[![Made with love by Qbitel](https://img.shields.io/badge/Made%20with%20%E2%9D%A4%EF%B8%8F%20by-Qbitel-6C3AED?style=flat-square)](https://github.com/yazhsab)

<br/>

[Getting Started](#-quick-start) &bull; [Features](#-features-showcase) &bull; [API Reference](#-api-at-a-glance) &bull; [Documentation](#-documentation) &bull; [Contributing](#-community--contributing)

<br/>

</div>

---

<br/>

## The Quantum Threat Is Not Theoretical&mdash;It's Inevitable

> **"By 2030, cryptographically relevant quantum computers could break RSA-2048 in hours.**
> **Your TLS certificates, VPN tunnels, and encrypted databases are on borrowed time."**

Every day, adversaries execute **Harvest Now, Decrypt Later (HNDL)** attacks&mdash;capturing encrypted traffic today to decrypt it with tomorrow's quantum computers. Financial records, medical data, trade secrets, and national security communications that you transmit now may already be stored, waiting for Q-Day.

**The question is not *if* your cryptography will break. It's *when*&mdash;and whether you'll be ready.**

QRAP gives you the answer. It scans your cryptographic infrastructure, quantifies your quantum risk exposure down to a single composite score, calculates exactly when your data becomes vulnerable using the **Mosca inequality**, and generates a step-by-step migration plan to NIST-standardized post-quantum algorithms.

<br/>

---

<br/>

## Why QRAP?

<div align="center">

<table>
<tr>
<td align="center" width="33%">

**First Open-Source<br/>Quantum Risk Platform**

The only production-grade, full-stack open-source solution for quantum cryptographic risk assessment.

</td>
<td align="center" width="33%">

**ML-Powered<br/>Risk Scoring**

Composite 0&ndash;100 risk scores weighted by severity, algorithm vulnerability, HNDL exposure, and PQC readiness.

</td>
<td align="center" width="33%">

**HNDL Exposure<br/>Calculator**

Precise harvest-now-decrypt-later risk windows using the Mosca inequality and conservative CRQC timelines.

</td>
</tr>
<tr>
<td align="center" width="33%">

**Automated PQC<br/>Migration Planning**

One-click migration roadmaps: RSA to ML-KEM, ECDSA to ML-DSA, aligned with FIPS 203/204 standards.

</td>
<td align="center" width="33%">

**Enterprise-Grade<br/>Security**

JWT + API key auth, rate limiting, HSTS, CSP, constant-time comparisons, and parameterized queries.

</td>
<td align="center" width="33%">

**Full-Stack<br/>Architecture**

Go API + Python ML engine + React dashboard + PostgreSQL. Deploy with a single `docker compose up`.

</td>
</tr>
</table>

</div>

<br/>

---

<br/>

## Features Showcase

### Discovery & Assessment

- Scan infrastructure to discover cryptographic assets across your entire organization
- Identify quantum-vulnerable algorithms (RSA, ECDSA, DH, DSA) with precise version and key-length detection
- Categorize findings across six risk categories: weak algorithms, short key lengths, deprecated protocols, missing PQC, certificate expiry, and HNDL exposure
- Multi-asset assessment targeting (API gateways, payment services, auth services, and more)

### ML-Powered Analysis

- Composite risk scoring engine producing a single 0&ndash;100 quantum risk score
- Severity-weighted scoring with category-specific multipliers (HNDL findings weighted 1.5x)
- PQC readiness percentage calculation per organization
- Intelligent prioritization of remediation actions by impact and effort

### Risk Scoring & Visualization

- Interactive React dashboard for real-time risk visualization
- Drill-down from organization-level scores to individual findings
- Risk trend analysis across multiple assessment cycles
- Exportable assessment summaries and compliance-ready reports

### HNDL Protection

- Mosca inequality-based risk window calculation
- Conservative CRQC timeline estimates (configurable per threat model)
- Data shelf-life analysis: understand which encrypted data is at risk and when
- Priority flagging for long-lived secrets (financial records, medical data, PII)

### Migration Planning

- Automated mapping from classical to NIST-standardized PQC algorithms
- RSA to ML-KEM (FIPS 203) and ECDSA to ML-DSA (FIPS 204) migration paths
- Phased rollout schedules with effort estimates (hours, complexity, dependencies)
- Priority-ordered remediation plans based on composite risk scores

### Enterprise Security

- HMAC-SHA256 JWT validation with constant-time comparison
- API key authentication resistant to timing side-channel attacks
- Per-IP rate limiting (100 requests/minute)
- Security headers: HSTS, CSP, X-Frame-Options, X-Content-Type-Options
- Request body size limits (1 MB) and parameterized SQL queries (pgx)
- Graceful shutdown with connection draining

<br/>

---

<br/>

## Architecture

```mermaid
flowchart LR
    subgraph Web["Web Dashboard"]
        direction TB
        WD["React 19 + Vite 6<br/>localhost:3002"]
    end

    subgraph API["Go REST API"]
        direction TB
        GA["Chi v5 Router<br/>localhost:8083"]
        MW["Auth | CORS<br/>RateLimit | Security"]
    end

    subgraph ML["Python ML Engine"]
        direction TB
        MLE["FastAPI + uvicorn<br/>localhost:8084"]
        RS["Risk Scorer"]
        HC["HNDL Calculator"]
        MP["Migration Planner"]
    end

    subgraph DB["PostgreSQL 16"]
        direction TB
        PG[("organizations<br/>assessments<br/>findings<br/>audit_log")]
    end

    Web -- "REST / HTTP" --> API
    API -- "gRPC / HTTP" --> ML
    API -- "SQL" --> DB
```

<br/>

---

<br/>

## Quick Start

### Docker Compose &mdash; Up and Running in 60 Seconds

```bash
# 1. Clone the repository
git clone https://github.com/yazhsab/qbitel-qrap.git && cd qbitel-qrap

# 2. Launch the full stack (API, ML engine, dashboard, PostgreSQL)
docker compose -f infra/docker/docker-compose.yml up -d

# 3. Run database migrations
export QRAP_DATABASE_URL="postgres://quantun:quantun_dev@localhost:5432/qrap?sslmode=disable"
make migrate
```

**That's it.** Open your browser:

| Service            | URL                          | Description                |
|:-------------------|:-----------------------------|:---------------------------|
| **Web Dashboard**  | http://localhost:3002         | Interactive risk dashboard |
| **API Server**     | http://localhost:8083         | Go REST API                |
| **ML Engine**      | http://localhost:8084         | Python ML scoring engine   |

Verify everything is running:

```bash
curl http://localhost:8083/health   # {"status":"ok","service":"qrap-api"}
curl http://localhost:8084/health   # {"status":"ok","service":"qrap-ml"}
```

<details>
<summary><strong>Local Development Setup (without Docker)</strong></summary>

<br/>

**Prerequisites:**

| Tool       | Version | Purpose              |
|:-----------|:--------|:---------------------|
| Go         | 1.23+   | API server           |
| Python     | 3.11+   | ML engine            |
| Node.js    | 22+     | Web dashboard        |
| PostgreSQL | 16+     | Database             |

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

</details>

<br/>

---

<br/>

## Screenshots & Demo

<div align="center">

<table>
<tr>
<td align="center" width="33%">

**Dashboard Overview**

<!-- <img src="docs/assets/screenshots/dashboard.png" alt="QRAP Dashboard" width="300"/> -->

*Organization-level quantum risk scores, PQC readiness metrics, and real-time threat indicators at a glance.*

</td>
<td align="center" width="33%">

**Risk Assessment**

<!-- <img src="docs/assets/screenshots/assessment.png" alt="Risk Assessment" width="300"/> -->

*Drill into individual assessments with finding-level detail, severity breakdowns, and remediation priorities.*

</td>
<td align="center" width="33%">

**HNDL Analysis**

<!-- <img src="docs/assets/screenshots/hndl.png" alt="HNDL Analysis" width="300"/> -->

*Visualize harvest-now-decrypt-later exposure windows using Mosca inequality calculations.*

</td>
</tr>
</table>

> Screenshots coming soon. Star the repo to be notified when the first release ships.

</div>

<br/>

---

<br/>

## API at a Glance

All API endpoints (except `/health`) support JWT Bearer token and API key authentication.

| Method | Endpoint | Description |
|:------:|:---------|:------------|
| `GET`  | `/health` | Health check (no auth required) |
| `POST` | `/api/v1/organizations` | Create a new organization |
| `GET`  | `/api/v1/organizations` | List all organizations |
| `GET`  | `/api/v1/organizations/{id}` | Get organization details |
| `POST` | `/api/v1/assessments` | Create a new assessment |
| `GET`  | `/api/v1/assessments` | List assessments |
| `GET`  | `/api/v1/assessments/{id}` | Get assessment with risk summary |
| `POST` | `/api/v1/assessments/{id}/run` | Execute an assessment |
| `GET`  | `/api/v1/findings` | List findings (filter by `assessment_id`) |
| `GET`  | `/api/v1/findings/{id}` | Get finding details |
| `POST` | `/api/v1/score` | Calculate composite risk score |
| `POST` | `/api/v1/hndl` | Calculate HNDL risk window |
| `POST` | `/api/v1/migration-plan` | Generate PQC migration plan |

<details>
<summary><strong>Example: Full Assessment Workflow</strong></summary>

<br/>

```bash
# Create an organization
curl -X POST http://localhost:8083/api/v1/organizations \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"name": "Acme Corp", "description": "Financial services"}'

# Create an assessment
curl -X POST http://localhost:8083/api/v1/assessments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "Q1 2026 Crypto Audit",
    "organization_id": "<org-uuid>",
    "target_assets": ["api-gateway", "payment-service", "auth-service"]
  }'

# Run the assessment
curl -X POST http://localhost:8083/api/v1/assessments/<assessment-uuid>/run \
  -H "Authorization: Bearer <token>"

# Calculate HNDL risk for a specific algorithm
curl -X POST http://localhost:8083/api/v1/hndl \
  -H "Content-Type: application/json" \
  -d '{"algorithm": "RSA-2048", "data_shelf_life_years": 15}'

# Generate a migration plan
curl -X POST http://localhost:8083/api/v1/migration-plan \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"assessment_id": "<assessment-uuid>"}'
```

</details>

Full API documentation: **[docs/API.md](docs/API.md)**

<br/>

---

<br/>

## Tech Stack

<div align="center">

[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![Python](https://img.shields.io/badge/Python-3776AB?style=for-the-badge&logo=python&logoColor=white)](https://python.org/)
[![React](https://img.shields.io/badge/React-20232A?style=for-the-badge&logo=react&logoColor=61DAFB)](https://react.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://docker.com/)
[![GitHub Actions](https://img.shields.io/badge/GitHub_Actions-2088FF?style=for-the-badge&logo=github-actions&logoColor=white)](https://github.com/features/actions)
[![Vite](https://img.shields.io/badge/Vite-646CFF?style=for-the-badge&logo=vite&logoColor=white)](https://vite.dev/)
[![FastAPI](https://img.shields.io/badge/FastAPI-009688?style=for-the-badge&logo=fastapi&logoColor=white)](https://fastapi.tiangolo.com/)

</div>

| Layer | Technology | Role |
|:------|:-----------|:-----|
| **API Server** | Go 1.23 &bull; Chi v5 | High-performance REST API with middleware pipeline |
| **ML Engine** | Python 3.11+ &bull; FastAPI &bull; uvicorn | Risk scoring, HNDL calculations, migration planning |
| **Dashboard** | React 19 &bull; Vite 6 &bull; TypeScript | Interactive risk visualization and management UI |
| **Database** | PostgreSQL 16 &bull; pgx | Persistent storage with parameterized queries |
| **Infrastructure** | Docker &bull; Docker Compose | One-command deployment for the full stack |
| **CI/CD** | GitHub Actions | Automated testing, linting, and build verification |
| **Auth** | JWT (HMAC-SHA256) &bull; API Keys | Dual authentication strategy for enterprise use |

<br/>

---

<br/>

## QRAP vs. Traditional Tools

<div align="center">

| Capability | QRAP | Qualys | Tenable | CryptoSense |
|:-----------|:----:|:------:|:-------:|:-----------:|
| Quantum risk assessment | **Yes** | No | No | Partial |
| HNDL exposure calculation | **Yes** | No | No | No |
| Mosca inequality analysis | **Yes** | No | No | No |
| PQC migration planning | **Yes** | No | No | Partial |
| ML-powered risk scoring | **Yes** | Yes | Yes | No |
| FIPS 203/204 alignment | **Yes** | No | No | Partial |
| Open source | **Yes** | No | No | No |
| Self-hosted option | **Yes** | No | No | Yes |
| REST API | **Yes** | Yes | Yes | Yes |
| Interactive dashboard | **Yes** | Yes | Yes | Yes |
| Enterprise auth (JWT + API keys) | **Yes** | Yes | Yes | Yes |
| Cost | **Free** | $$$ | $$$ | $$ |

</div>

<br/>

---

<br/>

## Roadmap

| Quarter | Milestone | Description |
|:--------|:----------|:------------|
| **Q2 2026** | Certificate Scanner | Automated TLS/X.509 certificate discovery and quantum vulnerability scanning across network endpoints |
| **Q2 2026** | SBOM Integration | Software Bill of Materials ingestion for cryptographic dependency analysis |
| **Q3 2026** | Compliance Frameworks | Mapping to NIST SP 1800-38, CNSA 2.0, and EU Quantum-Safe recommendations |
| **Q3 2026** | Agent-Based Collection | Lightweight agents for continuous cryptographic asset discovery on hosts and containers |
| **Q4 2026** | Hybrid PQC Testing | Validate hybrid classical/PQC deployments (e.g., X25519+ML-KEM) and measure performance impact |
| **Q4 2026** | Multi-Tenant SaaS Mode | Cloud-hosted offering with tenant isolation, SSO, and usage-based billing |

<br/>

---

<br/>

## Documentation

| Document | Description |
|:---------|:------------|
| **[Architecture](docs/ARCHITECTURE.md)** | System design, data flow, component interactions, and technology decisions |
| **[API Reference](docs/API.md)** | Complete REST API documentation with request/response examples |
| **[Deployment Guide](docs/DEPLOYMENT.md)** | Production deployment, scaling, monitoring, and operational best practices |
| **[Development Guide](docs/DEVELOPMENT.md)** | Local setup, coding standards, debugging, and contribution workflow |
| **[Contributing](CONTRIBUTING.md)** | How to contribute: issues, PRs, coding standards, and commit conventions |
| **[Security Policy](SECURITY.md)** | Vulnerability reporting process and security architecture details |
| **[Changelog](CHANGELOG.md)** | Release history and version notes |

<br/>

---

<br/>

## Configuration

All configuration is via environment variables. Copy `.env.example` to `.env` to get started.

| Variable | Default | Description |
|:---------|:--------|:------------|
| `QRAP_PORT` | `8083` | API server port |
| `QRAP_DATABASE_URL` | *(required)* | PostgreSQL connection string |
| `QRAP_ML_ENGINE_URL` | `http://127.0.0.1:8084` | ML engine URL |
| `QRAP_LOG_LEVEL` | `info` | Log level (debug, info, warn, error) |
| `QUANTUN_JWT_SECRET` | *(empty &mdash; auth disabled)* | HMAC-SHA256 secret for JWT validation |
| `QUANTUN_JWT_ISSUER` | `quantun` | Expected JWT `iss` claim |
| `QUANTUN_API_KEYS` | *(empty)* | Comma-separated `key:subject:role` entries |
| `QUANTUN_CORS_ORIGINS` | *(empty)* | Comma-separated allowed CORS origins |

> **Warning:** When neither `QUANTUN_JWT_SECRET` nor `QUANTUN_API_KEYS` is set, all endpoints are accessible without authentication. This is convenient for development but must never be used in production.

<br/>

---

<br/>

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

<br/>

---

<br/>

## Project Structure

```
qbitel-qrap/
├── api/                             # Go REST API
│   ├── cmd/
│   │   ├── server/main.go          # Server entrypoint
│   │   └── migrate/main.go         # Migration CLI helper
│   ├── internal/
│   │   ├── config/                  # Environment configuration
│   │   ├── handler/                 # HTTP handlers (health, org, assessment, finding)
│   │   ├── model/                   # Data models and response types
│   │   ├── repository/              # PostgreSQL repositories (pgx)
│   │   └── service/                 # Business logic layer
│   └── Dockerfile
├── ml/                              # Python ML Engine
│   ├── src/qrap_ml/
│   │   ├── api/app.py              # FastAPI application
│   │   ├── risk_scorer/            # Composite risk scoring engine
│   │   ├── hndl_calculator/        # HNDL risk calculator (Mosca inequality)
│   │   └── migration_planner/      # PQC migration roadmap generator
│   ├── tests/                       # Pytest test suite
│   └── Dockerfile
├── web/                             # React Dashboard
│   ├── src/
│   │   ├── App.tsx                  # Router and navigation
│   │   └── pages/                   # Dashboard, Assessments, Findings, HNDL
│   └── Dockerfile
├── shared/go/                       # Shared Go Libraries
│   └── middleware/
│       ├── auth.go                  # JWT + API key authentication
│       ├── ratelimit.go             # Per-IP rate limiting
│       ├── security.go              # Security headers, CORS, body limits
│       └── pagination.go            # Query parameter pagination
├── db/migrations/                   # PostgreSQL schema migrations
├── infra/docker/                    # Docker Compose files
├── .github/workflows/ci.yml        # GitHub Actions CI pipeline
├── Makefile                         # Build, test, and deploy targets
└── go.work                          # Go workspace configuration
```

<br/>

---

<br/>

## Community & Contributing

QRAP is built in the open. We welcome contributions of all kinds&mdash;from bug fixes and documentation improvements to new features and integrations.

**Getting started is easy:**

1. Read the **[Contributing Guide](CONTRIBUTING.md)** for coding standards and workflow
2. Browse [open issues](https://github.com/yazhsab/qbitel-qrap/issues) or create a new one
3. Fork the repo, create a branch, and submit a pull request
4. Join the discussion in [GitHub Discussions](https://github.com/yazhsab/qbitel-qrap/discussions)

Whether you're a cryptographer, a Go/Python/React developer, or a security professional concerned about quantum threats&mdash;there's a place for you here.

<br/>

---

<br/>

## Security

QRAP is a security-focused platform. If you discover a vulnerability, please follow our **[responsible disclosure process](SECURITY.md)** rather than opening a public issue.

Security highlights:

- HMAC-SHA256 JWT validation with constant-time comparison
- API key validation resistant to timing side-channel attacks
- Per-IP rate limiting (100 requests/minute)
- Security headers (HSTS, CSP, X-Frame-Options, X-Content-Type-Options)
- Request body size limits (1 MB)
- Parameterized SQL queries via pgx (zero SQL injection surface)
- Graceful shutdown with connection draining

<br/>

---

<br/>

## License

```
Copyright 2026 Qbitel

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

See **[LICENSE](LICENSE)** for the full text.

<br/>

---

<br/>

<div align="center">

**Built with ❤️ by the [Qbitel](https://github.com/yazhsab) team**

*Securing today's cryptography against tomorrow's quantum computers*

<br/>

[![GitHub Stars](https://img.shields.io/github/stars/yazhsab/qbitel-qrap?style=social)](https://github.com/yazhsab/qbitel-qrap)
&nbsp;&nbsp;
[![GitHub Forks](https://img.shields.io/github/forks/yazhsab/qbitel-qrap?style=social)](https://github.com/yazhsab/qbitel-qrap/fork)
&nbsp;&nbsp;
[![GitHub Watchers](https://img.shields.io/github/watchers/yazhsab/qbitel-qrap?style=social)](https://github.com/yazhsab/qbitel-qrap)

<br/>

**More from Qbitel:** &nbsp; [PQC-SDK](https://github.com/yazhsab) &bull; [QSGW](https://github.com/yazhsab)

<br/>

</div>
