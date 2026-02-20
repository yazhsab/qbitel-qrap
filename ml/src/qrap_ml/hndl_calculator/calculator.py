"""Harvest Now, Decrypt Later (HNDL) risk calculator.

Estimates the time window during which captured ciphertext remains at risk,
based on the algorithm in use, estimated quantum computing timeline, and
the secrecy shelf-life of the data.
"""

from __future__ import annotations

from dataclasses import dataclass
from datetime import datetime

# Estimated years until a cryptographically relevant quantum computer (CRQC)
# can break a given algorithm.  Conservative estimates (Mosca inequality).
_ALGORITHM_BREAK_YEAR: dict[str, int] = {
    "RSA-2048": 2030,
    "RSA-3072": 2032,
    "RSA-4096": 2035,
    "ECDSA-P256": 2030,
    "ECDSA-P384": 2032,
    "Ed25519": 2030,
    "X25519": 2030,
    "DH-2048": 2030,
    "AES-128": 2040,  # Grover halves effective key length
    "AES-256": 2060,
    # PQC algorithms -- effectively safe
    "ML-KEM-512": 2080,
    "ML-KEM-768": 2080,
    "ML-KEM-1024": 2080,
    "ML-DSA-44": 2080,
    "ML-DSA-65": 2080,
    "ML-DSA-87": 2080,
}


@dataclass
class HndlResult:
    algorithm: str
    estimated_break_year: int
    data_shelf_life_years: int
    risk_window_years: int
    is_at_risk: bool
    urgency: str  # CRITICAL, HIGH, MEDIUM, LOW


class HndlCalculator:
    """Evaluates HNDL risk for a given algorithm + data shelf life."""

    def calculate(
        self,
        algorithm: str,
        data_shelf_life_years: int = 10,
        reference_year: int | None = None,
    ) -> HndlResult:
        if reference_year is None:
            reference_year = datetime.now().year

        break_year = _ALGORITHM_BREAK_YEAR.get(algorithm, 2035)

        # Mosca inequality: risk if (migration_time + shelf_life) > time_to_CRQC
        years_until_break = break_year - reference_year
        risk_window = data_shelf_life_years - years_until_break
        is_at_risk = risk_window > 0

        if risk_window >= 10:
            urgency = "CRITICAL"
        elif risk_window >= 5:
            urgency = "HIGH"
        elif risk_window > 0:
            urgency = "MEDIUM"
        else:
            urgency = "LOW"

        return HndlResult(
            algorithm=algorithm,
            estimated_break_year=break_year,
            data_shelf_life_years=data_shelf_life_years,
            risk_window_years=max(risk_window, 0),
            is_at_risk=is_at_risk,
            urgency=urgency,
        )
