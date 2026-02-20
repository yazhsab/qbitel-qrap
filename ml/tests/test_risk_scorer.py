from qrap_ml.risk_scorer import RiskScorer
from qrap_ml.risk_scorer.scorer import Finding


def test_empty_findings():
    scorer = RiskScorer()
    result = scorer.score([], 10)
    assert result.risk_score == 0.0
    assert result.overall_risk == "LOW"
    assert result.pqc_readiness == 100.0


def test_critical_findings():
    scorer = RiskScorer()
    findings = [
        Finding(
            category="HARVEST_NOW_DECRYPT_LATER",
            risk_level="CRITICAL",
            affected_asset="api-server",
            current_algorithm="RSA-2048",
        ),
        Finding(
            category="MISSING_PQC",
            risk_level="HIGH",
            affected_asset="api-server",
            current_algorithm="RSA-2048",
        ),
    ]
    result = scorer.score(findings, 1)
    assert result.risk_score > 0
    assert result.overall_risk in ("CRITICAL", "HIGH")
    assert result.pqc_readiness == 0.0


def test_pqc_readiness():
    scorer = RiskScorer()
    findings = [
        Finding(
            category="MISSING_PQC",
            risk_level="HIGH",
            affected_asset="service-a",
        ),
    ]
    result = scorer.score(findings, 3)
    # 2 out of 3 assets are PQC-ready
    assert result.pqc_readiness > 60
