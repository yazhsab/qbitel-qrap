# Development Guide

This guide covers setting up a local development environment for QRAP, along with conventions, tooling, and workflows for contributing to each component of the platform.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Local Setup](#local-setup)
- [Go API Development](#go-api-development)
- [Python ML Engine Development](#python-ml-engine-development)
- [Web Dashboard Development](#web-dashboard-development)
- [Database Development](#database-development)
- [Code Style](#code-style)
- [Debugging Tips](#debugging-tips)
- [IDE Setup](#ide-setup)
- [Makefile Reference](#makefile-reference)

---

## Prerequisites

Ensure the following tools are installed on your development machine:

| Tool | Version | Purpose |
|---|---|---|
| Go | 1.23+ | API server |
| Python | 3.11+ | ML engine |
| Node.js | 22+ | Web dashboard |
| PostgreSQL | 16+ | Database |
| Docker | 24+ | Containerized dependencies |
| Docker Compose | v2 | Multi-service orchestration |
| Make | any | Build automation |

---

## Local Setup

### 1. Clone the Repository

```bash
git clone https://github.com/quantun/qrap.git
cd qrap
```

### 2. Start Infrastructure Dependencies

Use Docker Compose to start PostgreSQL:

```bash
make docker-deps
```

This starts a PostgreSQL 16 container on port 5432 with the default credentials configured in the Makefile.

### 3. Run Database Migrations

```bash
make migrate
```

### 4. Set Environment Variables

Create a `.env` file in the project root (this file is git-ignored):

```bash
# Database
QRAP_DB_HOST=localhost
QRAP_DB_PORT=5432
QRAP_DB_NAME=qrap
QRAP_DB_USER=qrap
QRAP_DB_PASSWORD=qrap_dev_password
QRAP_DB_SSLMODE=disable

# ML Engine
QRAP_ML_URL=http://localhost:8084

# Auth
QUANTUN_JWT_SECRET=dev-secret-change-in-production-must-be-at-least-32-chars
QUANTUN_JWT_ISSUER=qrap
QUANTUN_API_KEYS=dev-api-key-1,dev-api-key-2

# Logging
QRAP_LOG_LEVEL=debug
```

### 5. Start All Services

In separate terminal windows, start each service:

```bash
# Terminal 1: Go API
make run-api

# Terminal 2: Python ML Engine
make run-ml

# Terminal 3: React Dashboard
make run-web
```

Or start everything with Docker:

```bash
make docker-all
```

### 6. Verify the Setup

```bash
# API health check
curl http://localhost:8083/healthz

# ML engine health check
curl http://localhost:8084/health

# Web dashboard
open http://localhost:3002
```

---

## Go API Development

The Go API lives in the `api/` directory and uses Chi v5 as the HTTP router with pgx v5 for PostgreSQL access.

### Module Structure

```
api/
  cmd/            # Application entrypoint (main.go)
  internal/
    handler/      # HTTP request handlers
    service/      # Business logic layer
    repository/   # Database access layer (pgx queries)
    model/        # Domain types and DTOs
    middleware/    # HTTP middleware (auth, rate limiting, etc.)
  go.mod
  go.sum
```

### Adding a New Endpoint

1. **Define the model** in `api/internal/model/`:

```go
type Widget struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
}
```

2. **Create the repository** in `api/internal/repository/`:

```go
func (r *WidgetRepo) Create(ctx context.Context, w *model.Widget) error {
    _, err := r.pool.Exec(ctx,
        "INSERT INTO widgets (id, name, created_at) VALUES ($1, $2, $3)",
        w.ID, w.Name, w.CreatedAt,
    )
    return err
}
```

3. **Implement the service** in `api/internal/service/`:

```go
func (s *WidgetService) Create(ctx context.Context, req CreateWidgetRequest) (*model.Widget, error) {
    // Business logic, validation, ML engine calls
}
```

4. **Write the handler** in `api/internal/handler/`:

```go
func (h *WidgetHandler) Create(w http.ResponseWriter, r *http.Request) {
    // Parse request, call service, write JSON response
}
```

5. **Register the route** in the router setup:

```go
r.Route("/api/v1/widgets", func(r chi.Router) {
    r.Post("/", widgetHandler.Create)
    r.Get("/{id}", widgetHandler.Get)
})
```

### Running Tests

```bash
# All tests
cd api && go test ./...

# Specific package with verbose output
cd api && go test -v ./internal/handler/...

# With race detection
cd api && go test -race ./...

# With coverage
cd api && go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## Python ML Engine Development

The ML engine lives in the `ml/` directory and uses FastAPI with numpy and scikit-learn.

### Project Structure

```
ml/
  app/
    main.py           # FastAPI application entrypoint
    routers/          # API route definitions
    engines/          # ML computation engines
    models/           # Pydantic request/response models
    services/         # Business logic
  tests/              # pytest test suite
  requirements.txt    # Python dependencies
  pyproject.toml      # Project metadata
```

### Setting Up the Virtual Environment

```bash
cd ml
python -m venv .venv
source .venv/bin/activate    # macOS/Linux
# .venv\Scripts\activate     # Windows

pip install -r requirements.txt
pip install -r requirements-dev.txt  # Testing and linting tools
```

### Running the ML Engine

```bash
cd ml
uvicorn app.main:app --host 0.0.0.0 --port 8084 --reload
```

### Adding a New ML Engine

1. **Create the engine** in `ml/app/engines/`:

```python
import numpy as np

class NewRiskEngine:
    """Computes risk scores for a new category."""

    def calculate(self, params: dict) -> float:
        # ML computation logic
        score = np.clip(raw_score, 0.0, 100.0)
        return float(score)
```

2. **Define request/response models** in `ml/app/models/`:

```python
from pydantic import BaseModel, Field

class NewRiskRequest(BaseModel):
    algorithm: str
    key_size: int = Field(ge=128)

class NewRiskResponse(BaseModel):
    score: float = Field(ge=0, le=100)
    confidence: float
```

3. **Create the router** in `ml/app/routers/`:

```python
from fastapi import APIRouter

router = APIRouter(prefix="/api/v1/new-risk", tags=["new-risk"])

@router.post("/", response_model=NewRiskResponse)
async def calculate_new_risk(req: NewRiskRequest):
    engine = NewRiskEngine()
    return engine.calculate(req.model_dump())
```

4. **Register the router** in `ml/app/main.py`:

```python
app.include_router(new_risk_router)
```

### Running Tests

```bash
cd ml
pytest

# With verbose output
pytest -v

# Specific test file
pytest tests/test_risk_scoring.py

# With coverage
pytest --cov=app --cov-report=html
```

---

## Web Dashboard Development

The web dashboard lives in the `web/` directory and is built with React 19, Vite 6, and TypeScript 5.7.

### Project Structure

```
web/
  src/
    components/       # Reusable UI components
    pages/            # Route-level page components
    hooks/            # Custom React hooks
    services/         # API client functions
    types/            # TypeScript type definitions
    utils/            # Utility functions
    App.tsx           # Root component with routing
    main.tsx          # Application entrypoint
  public/             # Static assets
  index.html
  vite.config.ts
  tsconfig.json
  package.json
```

### Running the Dev Server

```bash
cd web
npm install
npm run dev
```

The Vite dev server starts on `http://localhost:3002` with hot module replacement enabled.

### Adding a New Page

1. **Create the page component** in `web/src/pages/`:

```tsx
export default function NewPage() {
  return (
    <div>
      <h1>New Page</h1>
    </div>
  );
}
```

2. **Add the route** in `web/src/App.tsx`:

```tsx
<Route path="/new-page" element={<NewPage />} />
```

3. **Create API service functions** in `web/src/services/`:

```tsx
export async function fetchNewData(): Promise<NewDataResponse> {
  const response = await apiClient.get('/api/v1/new-data');
  return response.data;
}
```

### TypeScript Conventions

- Use explicit return types on exported functions.
- Prefer `interface` for object shapes and `type` for unions and intersections.
- Use `unknown` instead of `any` wherever possible.
- Define API response types in `web/src/types/`.

### Building for Production

```bash
cd web
npm run build
```

The production bundle is output to `web/dist/`.

---

## Database Development

### Writing Migrations

Migrations are stored in `db/migrations/` with the naming convention:

```
YYYYMMDDHHMMSS_description.up.sql
YYYYMMDDHHMMSS_description.down.sql
```

Example migration:

```sql
-- 20260220120000_add_widgets.up.sql
CREATE TABLE widgets (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id      UUID NOT NULL REFERENCES organizations(id),
    name        TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_widgets_org_id ON widgets(org_id);
```

```sql
-- 20260220120000_add_widgets.down.sql
DROP TABLE IF EXISTS widgets;
```

### Running Migrations

```bash
# Apply all pending migrations
make migrate

# Roll back the last migration (manual)
psql -h localhost -U qrap -d qrap -f db/migrations/20260220120000_add_widgets.down.sql
```

### Best Practices

- Every `up` migration must have a corresponding `down` migration.
- Never modify a migration that has already been applied to a shared environment.
- Use `TIMESTAMPTZ` for all timestamp columns.
- Add indexes for foreign keys and commonly queried columns.
- Use `gen_random_uuid()` for primary keys.
- Always include `created_at` and `updated_at` columns on mutable tables.

---

## Code Style

### Go

- **Formatter**: `gofmt` (applied automatically on save in most editors)
- **Linter**: `golangci-lint`

```bash
# Format
gofmt -w api/

# Lint
golangci-lint run ./api/...
```

Key conventions:
- Follow the standard Go project layout.
- Use `context.Context` as the first parameter in functions that perform I/O.
- Return `error` as the last return value; never panic in library code.
- Use table-driven tests.
- Keep handler functions thin; push logic into services.

### Python

- **Linter and formatter**: `ruff`

```bash
# Check
ruff check ml/

# Fix automatically
ruff check --fix ml/

# Format
ruff format ml/
```

Key conventions:
- Use type hints on all function signatures.
- Use Pydantic models for request/response validation.
- Follow PEP 8 naming conventions.
- Write docstrings for public classes and functions.

### TypeScript

- **Formatter**: `prettier`

```bash
# Check
cd web && npx prettier --check "src/**/*.{ts,tsx}"

# Fix
cd web && npx prettier --write "src/**/*.{ts,tsx}"
```

Key conventions:
- Use functional components with hooks (no class components).
- Prefer named exports over default exports for components.
- Keep components focused; extract custom hooks for complex logic.
- Use `const` by default; use `let` only when reassignment is needed.

### Running All Linters

```bash
make lint
```

---

## Debugging Tips

### Structured Logging

Set `QRAP_LOG_LEVEL=debug` for verbose output. The API logs all incoming requests with method, path, status code, and duration:

```json
{
  "level": "debug",
  "ts": "2026-02-20T10:15:30.123Z",
  "msg": "request completed",
  "method": "POST",
  "path": "/api/v1/assessments",
  "status": 201,
  "duration_ms": 45.2,
  "remote_addr": "127.0.0.1"
}
```

### Curl Examples

```bash
# Authenticate and get a JWT token
curl -s http://localhost:8083/api/v1/auth/token \
  -H "Content-Type: application/json" \
  -d '{"api_key": "dev-api-key-1"}' | jq .

# Create an assessment (with JWT)
curl -s http://localhost:8083/api/v1/assessments \
  -H "Authorization: Bearer <your-jwt-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "org_id": "org_abc123",
    "name": "Production Infrastructure Scan"
  }' | jq .

# Use API key authentication directly
curl -s http://localhost:8083/api/v1/assessments \
  -H "X-API-Key: dev-api-key-1" | jq .

# Check ML engine risk scoring
curl -s http://localhost:8084/api/v1/risk-score \
  -H "Content-Type: application/json" \
  -d '{
    "algorithm": "RSA",
    "key_size": 2048,
    "data_sensitivity": "high"
  }' | jq .

# HNDL calculation
curl -s http://localhost:8084/api/v1/hndl \
  -H "Content-Type: application/json" \
  -d '{
    "algorithm": "RSA-2048",
    "shelf_life_years": 10,
    "migration_time_years": 3
  }' | jq .
```

### Database Debugging

```bash
# Connect to the local database
psql -h localhost -U qrap -d qrap

# View recent audit log entries
psql -h localhost -U qrap -d qrap -c \
  "SELECT * FROM qrap_audit_log ORDER BY created_at DESC LIMIT 10;"

# Check active connections
psql -h localhost -U qrap -d qrap -c \
  "SELECT count(*) FROM pg_stat_activity WHERE datname = 'qrap';"
```

---

## IDE Setup

### Visual Studio Code

Recommended extensions:

- **Go** (`golang.go`) -- Go language support with debugging, linting, and testing
- **Python** (`ms-python.python`) -- Python language support
- **Ruff** (`charliermarsh.ruff`) -- Fast Python linting and formatting
- **ESLint** (`dbaeumer.vscode-eslint`) -- JavaScript/TypeScript linting
- **Prettier** (`esbenp.prettier-vscode`) -- Code formatting
- **Docker** (`ms-azuretools.vscode-docker`) -- Dockerfile and Compose support
- **PostgreSQL** (`ckolkman.vscode-postgres`) -- Database client

Recommended workspace settings (`.vscode/settings.json`):

```json
{
  "go.lintTool": "golangci-lint",
  "go.lintFlags": ["--fast"],
  "go.testFlags": ["-v"],
  "[go]": {
    "editor.formatOnSave": true,
    "editor.defaultFormatter": "golang.go"
  },
  "[python]": {
    "editor.formatOnSave": true,
    "editor.defaultFormatter": "charliermarsh.ruff"
  },
  "[typescript][typescriptreact]": {
    "editor.formatOnSave": true,
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  },
  "python.analysis.typeCheckingMode": "basic"
}
```

### GoLand / IntelliJ IDEA

- Enable the `golangci-lint` integration under **Settings > Tools > Go Linter**.
- Configure the database data source pointing to `localhost:5432/qrap`.
- Set the Go module root to the `api/` directory.

### PyCharm

- Set the Python interpreter to the virtual environment at `ml/.venv/bin/python`.
- Enable Ruff as the external tool under **Settings > Tools > External Tools**.
- Configure the pytest runner under **Settings > Tools > Python Integrated Tools**.

---

## Makefile Reference

| Target | Description |
|---|---|
| `make setup` | Install all dependencies (Go, Python, Node.js) |
| `make build` | Build all components (API binary, ML package, web bundle) |
| `make test` | Run all test suites (Go, Python, TypeScript) |
| `make lint` | Run all linters (golangci-lint, ruff, prettier) |
| `make migrate` | Apply database migrations |
| `make docker-deps` | Start infrastructure dependencies (PostgreSQL) |
| `make docker-all` | Build and start all services via Docker Compose |
| `make run-api` | Start the Go API in development mode |
| `make run-ml` | Start the Python ML engine in development mode |
| `make run-web` | Start the Vite dev server |
| `make clean` | Remove build artifacts, caches, and temporary files |

Run `make help` (if available) or inspect the `Makefile` directly for additional targets and options.
