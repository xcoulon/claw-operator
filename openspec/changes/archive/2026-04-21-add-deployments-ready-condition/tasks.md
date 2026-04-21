## 1. API Changes

- [x] 1.1 Rename `ConditionTypeReady` to `ConditionTypeDeploymentsReady = "DeploymentsReady"` in `api/v1alpha1/claw_types.go`
- [x] 1.2 Rename `ConditionReasonReady` to `ConditionReasonPodsRunning = "PodsRunning"` in `api/v1alpha1/claw_types.go`
- [x] 1.3 Update kubebuilder printcolumn marker to reference `DeploymentsReady` condition type instead of `Ready`

## 2. Controller Implementation

- [x] 2.1 Rename `setReadyCondition()` method to `setDeploymentsReadyCondition()` in `internal/controller/claw_resource_controller.go`
- [x] 2.2 Update `setDeploymentsReadyCondition()` to use `ConditionTypeDeploymentsReady` instead of `ConditionTypeReady`
- [x] 2.3 Update `setDeploymentsReadyCondition()` to use `ConditionReasonPodsRunning` instead of `ConditionReasonReady` when deployments are available
- [x] 2.4 Update `updateStatus()` method to call renamed `setDeploymentsReadyCondition()`

## 3. Tests

- [x] 3.1 Update all test assertions checking for `Ready` condition to check for `DeploymentsReady` in `internal/controller/claw_status_test.go`
- [x] 3.2 Update test assertions checking for `reason="Ready"` to check for `reason="PodsRunning"` in `internal/controller/claw_status_test.go`
- [x] 3.3 Verify tests still validate deployment readiness check logic

## 4. Documentation

- [x] 4.1 Update CLAUDE.md to reference `DeploymentsReady` condition instead of `Ready`
- [x] 4.2 Update CLAUDE.md condition type constants section to reflect the rename
- [x] 4.3 Update CLAUDE.md printcolumn documentation to show `DeploymentsReady` JSONPath

## 5. Generate Manifests

- [x] 5.1 Run `make manifests` to regenerate CRD YAML with updated printcolumn and condition type
- [x] 5.2 Run `make generate` to regenerate DeepCopy methods if needed
