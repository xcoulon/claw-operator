## Context

The Claw operator currently maintains several status conditions (`Ready`, `CredentialsResolved`, `ProxyConfigured`) to track the instance state. The `Ready` condition tracks deployment readiness by checking both `claw` and `claw-proxy` Deployment statuses. This change renames the `Ready` condition to `DeploymentsReady` for clearer semantics - the new name explicitly indicates the condition tracks deployment pod state.

The existing implementation in `internal/controller/claw_resource_controller.go` includes:
- `checkDeploymentsReady()` — checks both deployments' Available conditions
- `getDeploymentAvailableStatus()` — fetches a deployment and reads its Available condition
- `setReadyCondition()` — sets the Ready condition based on deployment states (will be renamed to `setDeploymentsReadyCondition()`)
- `updateStatus()` — orchestrates status updates using the status subresource

The printcolumn in the CRD kubebuilder markers currently references the `Ready` condition type and will need to be updated to reference `DeploymentsReady`.

## Goals / Non-Goals

**Goals:**
- Rename `Ready` condition to `DeploymentsReady` for clearer semantics
- Rename `ConditionTypeReady` constant to `ConditionTypeDeploymentsReady`
- Rename `ConditionReasonReady` constant to `ConditionReasonPodsRunning`
- Update printcolumn JSONPath to reference `DeploymentsReady` instead of `Ready`
- Maintain identical deployment readiness check logic (no behavior change)

**Non-Goals:**
- Changing the deployment readiness check logic itself
- Adding new deployment status checks or monitoring
- Maintaining backward compatibility with the old `Ready` condition name (this is a breaking change)

## Decisions

### Decision 1: Complete rename instead of deprecation cycle
**Approach:** Directly replace the `Ready` condition with `DeploymentsReady` in a single change, rather than adding the new condition while keeping the old one.

**Rationale:** This is a simple rename for clarity. A deprecation cycle would add complexity (maintaining two identical conditions) for minimal benefit in an early-stage operator. Clean break is simpler.

**Alternative considered:** Add `DeploymentsReady` alongside `Ready`, then deprecate `Ready` → Rejected because it unnecessarily complicates the code and status surface for a cosmetic improvement.

### Decision 2: Use "PodsRunning" reason instead of "Ready"
**Approach:** When both deployments are available, set `reason=PodsRunning` instead of `Ready`.

**Rationale:** More descriptive and aligns with the condition name. "PodsRunning" clearly indicates the pods are in running state, while "Ready" is vague.

**Alternative considered:** Keep reason as "Ready" → Rejected because it doesn't improve clarity and wastes the opportunity to make the reason more descriptive.

### Decision 3: Rename helper method setReadyCondition to setDeploymentsReadyCondition
**Approach:** Rename the controller helper method to match the new condition name.

**Rationale:** Maintains naming consistency between the condition type and the method that sets it.

**Alternative considered:** Keep method name as `setReadyCondition()` → Rejected because it would create confusion between the method name and the condition it actually sets.

## Risks / Trade-offs

**Risk:** Breaking change requires users to update any automation checking the `Ready` condition.
→ **Mitigation:** Clear communication in release notes. The JSONPath change in printcolumn will be visible immediately in `kubectl get claw` output. Document the migration path.

**Risk:** Existing Claw instances may have `Ready` condition in their status from previous reconciliations.
→ **Mitigation:** On next reconciliation, the new controller will set `DeploymentsReady` and leave the old `Ready` condition orphaned. The orphaned condition will not be updated but is harmless (consumers should check for condition presence before reading). Consider adding a one-time migration to remove orphaned `Ready` conditions if needed.

**Trade-off:** Breaking change vs. better naming.
→ **Acceptance rationale:** Better naming improves long-term maintainability and user experience. The operator is in alpha (v1alpha1), so breaking changes are acceptable. The migration burden is low (simple find-replace in consumer code).

## Migration Plan

- **User migration:**
  1. Update any scripts, automation, or monitoring checking for `Ready` condition to check for `DeploymentsReady` instead
  2. Update JSONPath queries from `.status.conditions[?(@.type=="Ready")]` to `.status.conditions[?(@.type=="DeploymentsReady")]`
  3. Update checks for `reason="Ready"` to check for `reason="PodsRunning"` when deployments are available
  
- **Deployment:**
  1. Deploy updated operator with renamed condition
  2. Existing Claw instances will have `DeploymentsReady` set on next reconciliation
  3. Old `Ready` conditions will become orphaned (not updated, but not removed)
  
- **Rollback:**
  1. If critical issues arise, rollback to previous operator version
  2. Old controller will resume setting `Ready` condition
  3. Both conditions may exist temporarily (harmless)
  
- **User impact:** **BREAKING** - Any tooling checking for the `Ready` condition will break until updated to check for `DeploymentsReady`
