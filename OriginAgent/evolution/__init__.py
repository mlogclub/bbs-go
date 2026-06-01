"""Local OriginAgent evolution runtime primitives."""

from OriginAgent.evolution.activation import EvolutionActivationResult, EvolutionModuleActivator
from OriginAgent.evolution.capability_gate import EvolutionCapabilityGate, EvolutionCapabilityResult
from OriginAgent.evolution.events import EventType, EvolutionEvent
from OriginAgent.evolution.identity import EvolutionIdentityStore
from OriginAgent.evolution.ledger import EvolutionLedger, LedgerStatus
from OriginAgent.evolution.manager import (
    EvolutionModuleManager,
    EvolutionStageResult,
    EvolutionVerificationResult,
)
from OriginAgent.evolution.manifest import EvolutionManifest, validate_manifest
from OriginAgent.evolution.memory_vault import (
    MemoryVaultError,
    MemoryVaultImportResult,
    export_memory_vault,
    import_memory_vault,
    inspect_memory_vault,
    read_memory_vault,
    verify_memory_vault,
)
from OriginAgent.evolution.recovery import EvolutionRecoveryManager, EvolutionRecoveryResult
from OriginAgent.evolution.state_branch import (
    EvolutionMergeConflict,
    EvolutionMergePreview,
    EvolutionStateBranchResult,
    EvolutionStateBranchStore,
)
from OriginAgent.evolution.telemetry import (
    EvolutionProofBundleResult,
    EvolutionTelemetryRecorder,
    EvolutionTelemetryResult,
    EvolutionTokenBudgetResult,
)
from OriginAgent.evolution.verifier import EvolutionModuleVerifier, EvolutionVerificationReport

__all__ = [
    "EventType",
    "EvolutionActivationResult",
    "EvolutionCapabilityGate",
    "EvolutionCapabilityResult",
    "EvolutionEvent",
    "EvolutionIdentityStore",
    "EvolutionLedger",
    "EvolutionModuleActivator",
    "EvolutionModuleManager",
    "EvolutionModuleVerifier",
    "EvolutionManifest",
    "EvolutionMergeConflict",
    "EvolutionMergePreview",
    "EvolutionProofBundleResult",
    "EvolutionRecoveryManager",
    "EvolutionRecoveryResult",
    "EvolutionStageResult",
    "EvolutionStateBranchResult",
    "EvolutionStateBranchStore",
    "EvolutionTelemetryRecorder",
    "EvolutionTelemetryResult",
    "EvolutionTokenBudgetResult",
    "EvolutionVerificationReport",
    "EvolutionVerificationResult",
    "LedgerStatus",
    "MemoryVaultError",
    "MemoryVaultImportResult",
    "export_memory_vault",
    "import_memory_vault",
    "inspect_memory_vault",
    "read_memory_vault",
    "validate_manifest",
    "verify_memory_vault",
]
