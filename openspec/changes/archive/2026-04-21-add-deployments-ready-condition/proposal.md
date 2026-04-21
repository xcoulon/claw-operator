## Why

The Claw operator currently uses the `Ready` condition to signal deployment readiness. Renaming it to `DeploymentsReady` provides clearer semantics - the condition name explicitly indicates it tracks deployment pod state. This improves observability and makes the condition's purpose immediately clear to users and automation tools.

## What Changes

- **BREAKING**: Replace `Ready` condition type with `DeploymentsReady` condition type in `Claw.Status.Conditions`
- Set `DeploymentsReady` to `True` with `Reason="PodsRunning"` when both `claw` and `claw-proxy` Deployments report `Available=True`
- Set `DeploymentsReady` to `False` with `Reason="Provisioning"` when one or both Deployments are not ready
- Rename condition type constant from `ConditionTypeReady` to `ConditionTypeDeploymentsReady` in API package
- Add new reason constant `ConditionReasonPodsRunning` to replace `ConditionReasonReady` when deployments are available
- Update controller to set `DeploymentsReady` instead of `Ready`
- Update printcolumn in CRD to show `DeploymentsReady` status instead of `Ready`

## Capabilities

### New Capabilities
<!-- No new capabilities - this is a modification to existing status condition system -->

### Modified Capabilities
- `status-conditions`: Replace `Ready` condition type with `DeploymentsReady`, rename constants, update controller logic and CRD printcolumn

## Impact

- **API changes**: 
  - Rename `ConditionTypeReady` to `ConditionTypeDeploymentsReady` in `api/v1alpha1/claw_types.go`
  - Rename `ConditionReasonReady` to `ConditionReasonPodsRunning` in `api/v1alpha1/claw_types.go`
  - Keep `ConditionReasonProvisioning` unchanged (still used when deployments not ready)
- **Controller changes**: Update `internal/controller/claw_resource_controller.go` to set `DeploymentsReady` condition instead of `Ready`
- **CRD manifest**: Update printcolumn from `Ready` to `DeploymentsReady` in `api/v1alpha1/claw_types.go` kubebuilder markers
- **Tests**: Update all status condition tests in `internal/controller/claw_status_test.go` to use `DeploymentsReady` instead of `Ready`
- **CLAUDE.md**: Update documentation to reflect the condition rename
- **Backward compatibility**: **BREAKING CHANGE** - any automation or tooling checking for the `Ready` condition will need to update to check for `DeploymentsReady` instead
