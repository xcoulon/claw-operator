## REMOVED Requirements

### Requirement: Ready condition indicates overall readiness
**Reason**: Renamed to DeploymentsReady for clearer semantics that explicitly indicate deployment pod state tracking

**Migration**: Replace all references to the `Ready` condition type with `DeploymentsReady`. Update automation and monitoring to check for `DeploymentsReady` instead of `Ready`.

### Requirement: Ready condition reasons are well-defined
**Reason**: Replaced by DeploymentsReady condition with updated reason constants

**Migration**: Use `ConditionReasonPodsRunning` instead of `ConditionReasonReady` when checking for deployment availability. `ConditionReasonProvisioning` remains unchanged.

### Requirement: Condition type constants defined in API
**Reason**: `ConditionTypeReady` and `ConditionReasonReady` constants renamed to better reflect their purpose

**Migration**: Update code references from `ConditionTypeReady` to `ConditionTypeDeploymentsReady` and from `ConditionReasonReady` to `ConditionReasonPodsRunning`.

## ADDED Requirements

### Requirement: DeploymentsReady condition indicates deployment readiness
The controller SHALL maintain a DeploymentsReady condition type to indicate whether the Claw deployment pods are running.

#### Scenario: DeploymentsReady condition set to True when pods running
- **WHEN** both claw and claw-proxy Deployments have Available=True status
- **THEN** the controller SHALL set DeploymentsReady condition with status=True, reason=PodsRunning, message confirming both deployments are available

#### Scenario: DeploymentsReady condition set to False when pods not ready
- **WHEN** either claw or claw-proxy Deployment has Available condition not equal to True
- **THEN** the controller SHALL set DeploymentsReady condition with status=False, reason=Provisioning, message indicating which deployments are pending

#### Scenario: DeploymentsReady uses same deployment check logic
- **WHEN** evaluating DeploymentsReady condition status
- **THEN** the controller SHALL fetch both claw and claw-proxy Deployments and check their Available conditions
- **THEN** the controller SHALL use the same deployment readiness check that was previously used for the Ready condition

### Requirement: DeploymentsReady condition type and reason constants defined
The API package SHALL define constants for the DeploymentsReady condition type and its reasons.

#### Scenario: DeploymentsReady condition type constant exists
- **WHEN** examining api/v1alpha1/claw_types.go
- **THEN** it SHALL define ConditionTypeDeploymentsReady = "DeploymentsReady"

#### Scenario: PodsRunning reason constant exists
- **WHEN** examining api/v1alpha1/claw_types.go
- **THEN** it SHALL define ConditionReasonPodsRunning = "PodsRunning"
- **THEN** this reason SHALL be used when both deployments are available

#### Scenario: Provisioning reason constant remains
- **WHEN** examining api/v1alpha1/claw_types.go
- **THEN** it SHALL still define ConditionReasonProvisioning = "Provisioning"
- **THEN** this reason SHALL be used when one or both deployments are not ready

## MODIFIED Requirements

### Requirement: Claw CRD includes printcolumn for condition status
The Claw CRD SHALL include a printcolumn that displays the DeploymentsReady condition status.

#### Scenario: Printcolumn shows DeploymentsReady status
- **WHEN** examining the Claw CRD kubebuilder markers in api/v1alpha1/claw_types.go
- **THEN** the printcolumn SHALL be named "Ready" for display purposes
- **THEN** the printcolumn SHALL use JSONPath `.status.conditions[?(@.type=="DeploymentsReady")].status`

#### Scenario: Printcolumn shows DeploymentsReady reason
- **WHEN** examining the Claw CRD kubebuilder markers in api/v1alpha1/claw_types.go
- **THEN** the Reason printcolumn SHALL use JSONPath `.status.conditions[?(@.type=="DeploymentsReady")].reason`
