"""Quantum risk scoring engine.

Calculates a composite risk score (0-100) for an organisation's cryptographic
posture by weighting finding severity, algorithm weakness, and exposure surface.
"""

from __future__ import annotations

from dataclasses import dataclass, field

# Risk weights per severity level
_SEVERITY_WEIGHTS: dict[str, float] = {
    "CRITICAL": 10.0,
    "HIGH": 5.0,
    "MEDIUM": 2.0,
    "LOW": 1.0,
    "INFO": 0.0,
}

# Category multipliers -- some categories are more urgent than others
_CATEGORY_MULTIPLIERS: dict[str, float] = {
    "HARVEST_NOW_DECRYPT_LATER": 1.5,
    "MISSING_PQC": 1.3,
    "WEAK_ALGORITHM": 1.2,
    "SHORT_KEY_LENGTH": 1.1,
    "DEPRECATED_PROTOCOL": 1.0,
    "CERTIFICATE_EXPIRY": 0.8,
}


@dataclass
class Finding:
    category: str
    risk_level: str
    affected_asset: str
    current_algorithm: str | None = None
    recommended_algorithm: str | None = None


@dataclass
class RiskResult:
    risk_score: float
    overall_risk: str
    pqc_readiness: float
    finding_breakdown: dict[str, int] = field(default_factory=dict)


class RiskScorer:
    """Stateless scorer that evaluates a list of findings."""

    def score(self, findings: list[Finding], total_assets: int) -> RiskResult:
        if not findings:
            return RiskResult(
                risk_score=0.0,
                overall_risk="LOW",
                pqc_readiness=100.0,
                finding_breakdown={},
            )

        weighted_sum = 0.0
        breakdown: dict[str, int] = {}
        missing_pqc_assets: set[str] = set()

        for f in findings:
            severity = _SEVERITY_WEIGHTS.get(f.risk_level, 0.0)
            multiplier = _CATEGORY_MULTIPLIERS.get(f.category, 1.0)
            weighted_sum += severity * multiplier

            breakdown[f.risk_level] = breakdown.get(f.risk_level, 0) + 1

            if f.category == "MISSING_PQC":
                missing_pqc_assets.add(f.affected_asset)

        # Normalize to 0-100 scale
        max_possible = len(findings) * 10.0 * 1.5  # worst case: all CRITICAL + HNDL
        risk_score = min((weighted_sum / max_possible) * 100, 100.0)

        # Overall risk level
        if risk_score >= 80:
            overall_risk = "CRITICAL"
        elif risk_score >= 60:
            overall_risk = "HIGH"
        elif risk_score >= 30:
            overall_risk = "MEDIUM"
        else:
            overall_risk = "LOW"

        # PQC readiness
        if total_assets > 0:
            pqc_ready = total_assets - len(missing_pqc_assets)
            pqc_readiness = (pqc_ready / total_assets) * 100
        else:
            pqc_readiness = 100.0

        return RiskResult(
            risk_score=round(risk_score, 2),
            overall_risk=overall_risk,
            pqc_readiness=round(pqc_readiness, 2),
            finding_breakdown=breakdown,
        )
