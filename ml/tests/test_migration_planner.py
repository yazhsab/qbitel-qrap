from qrap_ml.migration_planner import MigrationPlanner


def test_basic_plan():
    planner = MigrationPlanner()
    assets = [
        {"asset": "api-server", "algorithm": "RSA-2048", "urgency": "CRITICAL"},
        {"asset": "db-conn", "algorithm": "ECDSA-P256", "urgency": "HIGH"},
    ]
    plan = planner.plan(assets)
    assert len(plan.steps) == 2
    assert plan.steps[0].priority == "CRITICAL"
    assert plan.steps[0].target_algorithm == "ML-KEM-768"
    assert plan.steps[1].target_algorithm == "ML-DSA-65"


def test_pqc_not_migrated():
    planner = MigrationPlanner()
    assets = [
        {"asset": "pqc-service", "algorithm": "ML-KEM-768"},
    ]
    plan = planner.plan(assets)
    assert len(plan.steps) == 0


def test_priority_ordering():
    planner = MigrationPlanner()
    assets = [
        {"asset": "low", "algorithm": "RSA-2048", "urgency": "LOW"},
        {"asset": "high", "algorithm": "RSA-2048", "urgency": "HIGH"},
        {"asset": "crit", "algorithm": "RSA-2048", "urgency": "CRITICAL"},
    ]
    plan = planner.plan(assets)
    assert plan.steps[0].priority == "CRITICAL"
    assert plan.steps[1].priority == "HIGH"
    assert plan.steps[2].priority == "LOW"
    assert plan.critical_count == 1
