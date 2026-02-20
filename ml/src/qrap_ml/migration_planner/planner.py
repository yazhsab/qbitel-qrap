"""PQC migration planner.

Produces a prioritised migration roadmap from classical to post-quantum
algorithms, considering algorithm mappings, risk urgency, and dependencies.
"""

from __future__ import annotations

from dataclasses import dataclass, field

# Recommended PQC replacements for classical algorithms
_MIGRATION_MAP: dict[str, str] = {
    "RSA-2048": "ML-KEM-768",
    "RSA-3072": "ML-KEM-768",
    "RSA-4096": "ML-KEM-1024",
    "ECDSA-P256": "ML-DSA-65",
    "ECDSA-P384": "ML-DSA-87",
    "Ed25519": "ML-DSA-65",
    "X25519": "X25519-ML-KEM-768",
    "DH-2048": "ML-KEM-768",
}

_PRIORITY_ORDER = {"CRITICAL": 0, "HIGH": 1, "MEDIUM": 2, "LOW": 3}


@dataclass
class MigrationStep:
    asset: str
    current_algorithm: str
    target_algorithm: str
    priority: str
    estimated_effort: str  # LOW, MEDIUM, HIGH
    notes: str = ""


@dataclass
class MigrationPlan:
    steps: list[MigrationStep] = field(default_factory=list)
    total_assets: int = 0
    critical_count: int = 0
    estimated_phases: int = 1


class MigrationPlanner:
    """Generates a PQC migration plan from a list of asset/algorithm pairs."""

    def plan(
        self,
        assets: list[dict[str, str]],
    ) -> MigrationPlan:
        """Create a migration plan.

        Args:
            assets: List of dicts with keys: asset, algorithm, urgency (optional).
        """
        steps: list[MigrationStep] = []

        for item in assets:
            asset = item["asset"]
            algo = item["algorithm"]
            urgency = item.get("urgency", "MEDIUM")

            target = _MIGRATION_MAP.get(algo)
            if target is None:
                # Already PQC or unknown
                continue

            effort = self._estimate_effort(algo, target)

            steps.append(
                MigrationStep(
                    asset=asset,
                    current_algorithm=algo,
                    target_algorithm=target,
                    priority=urgency,
                    estimated_effort=effort,
                    notes=f"Migrate from {algo} to {target}",
                )
            )

        # Sort by priority
        steps.sort(key=lambda s: _PRIORITY_ORDER.get(s.priority, 99))

        critical = sum(1 for s in steps if s.priority == "CRITICAL")

        # Estimate phases: 1 phase per 5 assets
        phases = max(1, (len(steps) + 4) // 5)

        return MigrationPlan(
            steps=steps,
            total_assets=len(assets),
            critical_count=critical,
            estimated_phases=phases,
        )

    def _estimate_effort(self, current: str, target: str) -> str:
        # Hybrid algorithms require more effort
        if "X25519" in target or "Ed25519" in target:
            return "HIGH"
        # KEM swap is usually moderate
        if "KEM" in target:
            return "MEDIUM"
        # Signature swap
        return "MEDIUM"
