from qrap_ml.hndl_calculator import HndlCalculator


def test_rsa_at_risk():
    calc = HndlCalculator()
    result = calc.calculate("RSA-2048", data_shelf_life_years=15, reference_year=2025)
    assert result.is_at_risk
    assert result.urgency in ("CRITICAL", "HIGH")
    assert result.risk_window_years > 0


def test_pqc_not_at_risk():
    calc = HndlCalculator()
    result = calc.calculate("ML-KEM-768", data_shelf_life_years=10, reference_year=2025)
    assert not result.is_at_risk
    assert result.urgency == "LOW"
    assert result.risk_window_years == 0


def test_short_shelf_life():
    calc = HndlCalculator()
    result = calc.calculate("RSA-2048", data_shelf_life_years=2, reference_year=2025)
    assert not result.is_at_risk
