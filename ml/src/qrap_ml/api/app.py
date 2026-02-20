"""FastAPI application exposing the QRAP ML engine."""

from __future__ import annotations

from fastapi import FastAPI
from pydantic import BaseModel

from qrap_ml.hndl_calculator import HndlCalculator
from qrap_ml.migration_planner import MigrationPlanner
from qrap_ml.risk_scorer import RiskScorer
from qrap_ml.risk_scorer.scorer import Finding

app = FastAPI(title="QRAP ML Engine", version="0.1.0")

_scorer = RiskScorer()
_hndl = HndlCalculator()
_planner = MigrationPlanner()


# ---------- Health ----------


@app.get("/health")
def health() -> dict[str, str]:
    return {"status": "ok", "service": "qrap-ml"}


# ---------- Risk Scoring ----------


class FindingInput(BaseModel):
    category: str
    risk_level: str
    affected_asset: str
    current_algorithm: str | None = None
    recommended_algorithm: str | None = None


class ScoreRequest(BaseModel):
    findings: list[FindingInput]
    total_assets: int


class ScoreResponse(BaseModel):
    risk_score: float
    overall_risk: str
    pqc_readiness: float
    finding_breakdown: dict[str, int]


@app.post("/api/v1/score", response_model=ScoreResponse)
def score_risk(req: ScoreRequest) -> ScoreResponse:
    findings = [
        Finding(
            category=f.category,
            risk_level=f.risk_level,
            affected_asset=f.affected_asset,
            current_algorithm=f.current_algorithm,
            recommended_algorithm=f.recommended_algorithm,
        )
        for f in req.findings
    ]
    result = _scorer.score(findings, req.total_assets)
    return ScoreResponse(
        risk_score=result.risk_score,
        overall_risk=result.overall_risk,
        pqc_readiness=result.pqc_readiness,
        finding_breakdown=result.finding_breakdown,
    )


# ---------- HNDL ----------


class HndlRequest(BaseModel):
    algorithm: str
    data_shelf_life_years: int = 10


class HndlResponse(BaseModel):
    algorithm: str
    estimated_break_year: int
    data_shelf_life_years: int
    risk_window_years: int
    is_at_risk: bool
    urgency: str


@app.post("/api/v1/hndl", response_model=HndlResponse)
def calculate_hndl(req: HndlRequest) -> HndlResponse:
    result = _hndl.calculate(req.algorithm, req.data_shelf_life_years)
    return HndlResponse(
        algorithm=result.algorithm,
        estimated_break_year=result.estimated_break_year,
        data_shelf_life_years=result.data_shelf_life_years,
        risk_window_years=result.risk_window_years,
        is_at_risk=result.is_at_risk,
        urgency=result.urgency,
    )


# ---------- Migration Planning ----------


class AssetInput(BaseModel):
    asset: str
    algorithm: str
    urgency: str = "MEDIUM"


class MigrationRequest(BaseModel):
    assets: list[AssetInput]


class MigrationStepResponse(BaseModel):
    asset: str
    current_algorithm: str
    target_algorithm: str
    priority: str
    estimated_effort: str
    notes: str


class MigrationResponse(BaseModel):
    steps: list[MigrationStepResponse]
    total_assets: int
    critical_count: int
    estimated_phases: int


@app.post("/api/v1/migration-plan", response_model=MigrationResponse)
def create_migration_plan(req: MigrationRequest) -> MigrationResponse:
    assets = [a.model_dump() for a in req.assets]
    plan = _planner.plan(assets)
    return MigrationResponse(
        steps=[
            MigrationStepResponse(
                asset=s.asset,
                current_algorithm=s.current_algorithm,
                target_algorithm=s.target_algorithm,
                priority=s.priority,
                estimated_effort=s.estimated_effort,
                notes=s.notes,
            )
            for s in plan.steps
        ],
        total_assets=plan.total_assets,
        critical_count=plan.critical_count,
        estimated_phases=plan.estimated_phases,
    )
