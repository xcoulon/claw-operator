## Context

The Claw operator uses Kustomize to manage all Kubernetes resources for a Claw instance. Manifests are embedded in `internal/assets/manifests/` as a complete directory and applied via server-side apply during reconciliation. The controller already handles:
- Loading manifests from embedded filesystem
- Building resources using the Kustomize API
- Setting namespace and owner references dynamically
- Applying all resources via server-side apply

The device-pairing application needs to be deployed alongside the existing claw and claw-proxy deployments. This requires adding new manifest files to the embedded directory and updating the kustomization.yaml resource list.

## Goals / Non-Goals

**Goals:**
- Add deployment and service manifests for claw-device-pairing to embedded manifests directory
- Follow existing security hardening patterns (non-root, seccomp, dropped capabilities)
- Ensure device-pairing resources are applied automatically during Claw reconciliation
- Use consistent naming and labeling with existing resources

**Non-Goals:**
- Modifying the controller code (no changes needed - controller already applies all resources from kustomization.yaml)
- Adding new configuration options or CRD fields
- Implementing the device-pairing application logic itself
- Adding NetworkPolicies or Routes for device-pairing (can be added later if needed)

## Decisions

### Decision 1: Manifest location in internal/assets/manifests/
**Approach:** Add `device-pairing-deployment.yaml` and `device-pairing-service.yaml` directly to `internal/assets/manifests/` alongside existing manifests.

**Rationale:** This is the established pattern for all Claw resources. The controller's `//go:embed manifests` directive will automatically include new files in the embedded filesystem.

**Alternative considered:** Create a subdirectory for device-pairing manifests → Rejected because it would require changes to the Kustomize build process and complicate the structure unnecessarily.

### Decision 2: Security context matches existing deployments
**Approach:** Use identical security hardening as the claw deployment: `runAsNonRoot: true`, `allowPrivilegeEscalation: false`, `readOnlyRootFilesystem: true`, drop all capabilities, seccomp RuntimeDefault.

**Rationale:** Consistency across deployments, proven security posture, follows Kubernetes security best practices.

**Alternative considered:** Use looser security context → Rejected because there's no justification for relaxing security hardening.

### Decision 3: Single replica with Recreate strategy
**Approach:** Run a single replica with `Recreate` update strategy (same as claw deployment).

**Rationale:** Simplicity for initial implementation. Device pairing likely doesn't need high availability at this stage. Recreate strategy avoids having two instances during updates.

**Alternative considered:** Multiple replicas with RollingUpdate → Deferred until there's a demonstrated need for HA.

### Decision 4: Use latest tag for now
**Approach:** Use `quay.io/xcoulon/claw-device-pairing:latest` as specified by the user.

**Rationale:** Follows user's explicit requirement. Allows rapid iteration during development.

**Trade-off:** Latest tag is not recommended for production. Should be replaced with semantic versioning tags (e.g., `v1.0.0`) when the application stabilizes.

### Decision 5: ClusterIP service for internal access
**Approach:** Expose device-pairing via a ClusterIP Service (not exposed externally).

**Rationale:** Device pairing is an internal service. External access can be added later via Route or Ingress if needed.

**Alternative considered:** Add Route immediately → Deferred because requirements for external access are unclear.

## Risks / Trade-offs

**Risk:** Using `latest` tag makes deployments non-deterministic.
→ **Mitigation:** Document that this should be replaced with versioned tags for production. Consider adding a TODO comment in the manifest.

**Risk:** Unknown port/configuration requirements for device-pairing application.
→ **Mitigation:** Use placeholder port (8080) in service manifest. Update during implementation if needed.

**Risk:** Resource limits may be too low or too high.
→ **Mitigation:** Use conservative defaults (128Mi/256Mi memory, 100m/500m CPU similar to other services). Monitor and adjust based on actual usage.

**Trade-off:** Single replica means no high availability.
→ **Acceptance rationale:** Acceptable for initial implementation. Can be increased later if needed.

## Migration Plan

- **Deployment:** No migration needed. New manifests are additive.
- **Rollback:** If device-pairing causes issues, remove the manifests from kustomization.yaml and redeploy the operator.
- **User impact:** None. Device pairing service is automatically deployed with new Claw instances. Existing instances will get the device-pairing deployment on next reconciliation.
