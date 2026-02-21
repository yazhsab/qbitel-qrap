# Security Policy

The QRAP team takes security seriously. This document describes how to report vulnerabilities, our response process, and the security features built into the platform.

## Table of Contents

- [Supported Versions](#supported-versions)
- [Reporting a Vulnerability](#reporting-a-vulnerability)
- [What to Include in Reports](#what-to-include-in-reports)
- [Response Timeline](#response-timeline)
- [Responsible Disclosure](#responsible-disclosure)
- [Security Features](#security-features)
- [Out of Scope](#out-of-scope)

---

## Supported Versions

| Version | Supported |
|---|---|
| 0.1.x | Yes |

Only the latest release within a supported version line receives security updates. We recommend always running the most recent patch release.

---

## Reporting a Vulnerability

**Do not report security vulnerabilities through public GitHub issues.**

Instead, please send a detailed report to:

**Email**: [security@qbitel.dev](mailto:security@qbitel.dev)

If you would like to encrypt your report, our PGP key is available at:

- [https://qbitel.dev/.well-known/pgp-key.txt](https://qbitel.dev/.well-known/pgp-key.txt)
- Key fingerprint published on the project website

Please use encrypted email for reports involving sensitive details such as proof-of-concept exploits or credentials.

---

## What to Include in Reports

To help us triage and resolve the issue quickly, please include:

1. **Description**: A clear, concise description of the vulnerability.
2. **Affected component**: Which part of QRAP is affected (API, ML engine, web dashboard, database, authentication, etc.).
3. **QRAP version**: The version or commit hash where the vulnerability was observed.
4. **Steps to reproduce**: A minimal set of steps or a proof-of-concept that demonstrates the vulnerability.
5. **Impact assessment**: Your assessment of the severity and potential impact (e.g., data exposure, authentication bypass, denial of service).
6. **Environment**: Operating system, browser, and any relevant configuration details.
7. **Suggested fix** (optional): If you have ideas about how to address the issue, we welcome them.

---

## Response Timeline

| Stage | Target |
|---|---|
| **Acknowledgment** | Within 48 hours of receiving the report |
| **Triage and assessment** | Within 1 week |
| **Fix development** | Target within 30 days (severity-dependent) |
| **Public disclosure** | Coordinated with the reporter after the fix is released |

For critical vulnerabilities that are actively being exploited, we will prioritize an expedited fix and release cycle.

If you do not receive an acknowledgment within 48 hours, please follow up by email to ensure your report was received.

---

## Responsible Disclosure

We follow a coordinated disclosure process:

1. **Report received**: We acknowledge receipt and begin investigation.
2. **Assessment**: We evaluate the severity, impact, and scope of the vulnerability.
3. **Fix development**: We develop and test a fix in a private branch.
4. **Release**: The fix is included in a new patch release.
5. **Advisory**: We publish a security advisory on GitHub with credit to the reporter (unless anonymity is requested).
6. **Public disclosure**: Full details are disclosed after users have had reasonable time to update.

We ask that reporters:

- Allow us reasonable time to investigate and address the vulnerability before public disclosure.
- Avoid exploiting the vulnerability beyond what is necessary to demonstrate the issue.
- Do not access, modify, or delete data belonging to other users.
- Act in good faith to avoid disruption to QRAP users and infrastructure.

We will not take legal action against researchers who follow this responsible disclosure process.

---

## Security Features

QRAP incorporates the following security measures:

### Authentication

- **JWT (HMAC-SHA256)**: JSON Web Tokens are signed using HMAC-SHA256 with a configurable secret (`QUANTUN_JWT_SECRET`). Token verification uses constant-time comparison to prevent timing attacks.
- **API key authentication**: API keys are validated using timing-safe comparison functions (`crypto/subtle.ConstantTimeCompare` in Go) to prevent timing-based enumeration of valid keys.
- **Dual authentication**: The API supports both JWT bearer tokens and API key headers, allowing flexible integration patterns.

### Rate Limiting

- **Per-IP rate limiting**: The API enforces a configurable rate limit (default: 100 requests per minute per IP address) to prevent abuse and brute-force attacks.
- Rate limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`) are included in API responses.

### Security Headers

The API applies the following security headers to all responses:

| Header | Value | Purpose |
|---|---|---|
| `Strict-Transport-Security` | `max-age=31536000; includeSubDomains` | Enforce HTTPS |
| `Content-Security-Policy` | Restrictive policy | Prevent XSS and injection |
| `X-Frame-Options` | `DENY` | Prevent clickjacking |
| `X-Content-Type-Options` | `nosniff` | Prevent MIME-type sniffing |
| `X-XSS-Protection` | `1; mode=block` | Legacy XSS filter |
| `Referrer-Policy` | `strict-origin-when-cross-origin` | Control referrer leakage |

### Request Validation

- **Body size limits**: Request bodies are limited to 1 MB by default (`QRAP_MAX_BODY_SIZE`) to prevent resource exhaustion.
- **Input validation**: All API inputs are validated and sanitized before processing.

### Database Security

- **Parameterized queries**: All database queries use pgx parameterized statements, which prevents SQL injection by design. No string concatenation is used to build SQL queries.
- **Least privilege**: The database user should be granted only the permissions required by the application.

### Operational Security

- **Graceful shutdown**: The API server handles `SIGTERM` and `SIGINT` signals, draining active connections before shutting down. This prevents data corruption during deployments.
- **Structured logging**: All logs are structured JSON via the zap library. Sensitive data (tokens, passwords, API keys) is never written to logs.
- **Audit trail**: The `qrap_audit_log` table records all significant operations for forensic analysis and compliance.

---

## Out of Scope

The following are not considered vulnerabilities for the purposes of this security policy:

- **Social engineering** attacks against QRAP maintainers, contributors, or users.
- **Denial of service (DoS)** attacks against QRAP infrastructure, including volumetric attacks. Rate limiting is provided as a best-effort defense, not a guarantee against determined attackers.
- **Attacks requiring physical access** to the server or network.
- **Issues in third-party dependencies** that do not directly affect QRAP. Please report these to the upstream project. If a dependency vulnerability does affect QRAP, please include that context in your report.
- **Missing security headers on non-production endpoints** (e.g., health checks).
- **Clickjacking on pages with no state-changing actions**.
- **Content spoofing or text injection** without a demonstrated security impact.
- **Version or banner disclosure** in HTTP response headers.

---

Thank you for helping keep QRAP and its users safe. We appreciate the efforts of security researchers and the broader community in identifying and responsibly disclosing vulnerabilities.
