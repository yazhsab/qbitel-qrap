# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-02-20

### Added
- Initial release of QRAP as a standalone open-source product
- Go REST API with Chi v5 router and PostgreSQL backend
- Cryptographic risk assessment engine with six finding categories
- HNDL (Harvest Now, Decrypt Later) risk calculator using Mosca inequality
- PQC migration planner mapping classical algorithms to NIST post-quantum standards
- Composite risk scoring engine (0-100 scale) with ML-based weighting
- Python ML engine with FastAPI serving risk scoring, HNDL calculation, and migration planning
- React 19 web dashboard with assessment visualization and HNDL analysis views
- JWT (HMAC-SHA256) and API key authentication with constant-time comparison
- Per-IP rate limiting middleware
- Security headers middleware (HSTS, CSP, X-Frame-Options)
- PostgreSQL database schema with organizations, assessments, findings, and audit log
- Docker Compose configuration for local development and deployment
- GitHub Actions CI pipeline (Go, Python, TypeScript, Docker)
- Comprehensive REST API documentation
- Apache 2.0 license
