## Context

The operator currently deploys a single OpenShift Route that directs all traffic to the main `claw` service (gateway on port 18789). The device pairing service (`claw-device-pairing`) exists as a separate deployment and service but is not accessible externally.

OpenShift Routes support path-based routing through multiple Route objects with `spec.path` set. When Routes share the same `spec.host`, the router uses path-based matching to direct traffic to the appropriate backend service.

## Goals / Non-Goals

**Goals:**
- Enable external access to the device pairing service via `/pair-device` path
- Maintain backward compatibility for all existing gateway routes
- Use OpenShift-native Route features (no additional ingress controllers)
- Preserve existing Route TLS and timeout settings

**Non-Goals:**
- Complex routing logic beyond simple path prefix matching
- Load balancing between multiple backends (single service per path)
- Kubernetes Ingress support (Route is OpenShift-specific)
- Authentication/authorization at the Route level (handled by services)

## Decisions

### Decision 1: Use two Route objects with shared hostname

**Rationale**: OpenShift Routes support path-based routing when multiple Routes share the same `spec.host` value. The router performs path-based matching, directing traffic to the appropriate backend service.

**Alternatives considered:**
- **Single Route with alternateBackends**: The `spec.alternateBackends` field is designed for weighted traffic distribution (A/B testing), not path-based routing.
- **Kubernetes Ingress**: More portable but adds complexity. OpenShift-native Routes are simpler for OpenShift-only deployments.

**Approach**: Create two Route manifests sharing the same hostname:
- `device-pairing-route.yaml`: `path: /pair-device`, backend: `claw-device-pairing` service, `host: <shared-hostname>`
- `route.yaml`: `path: /`, backend: `claw` service, `host: <shared-hostname>`
- The controller applies the main Route first, extracts the assigned hostname from status, then applies the device pairing Route with the same `spec.host` value
- OpenShift router uses longest-prefix matching: `/pair-device` takes precedence over `/`

### Decision 2: Two-phase Route application with hostname injection

**Rationale**: To ensure both Routes share the same hostname, the controller must:
1. Apply the main Route first and wait for OpenShift router to assign a hostname
2. Extract the hostname from the main Route's status
3. Inject that hostname into the device pairing Route manifest before applying it

**Controller changes needed**:
1. **Extend `applyRouteOnly()` to handle hostname injection**:
   - Apply main Route (`claw`) first
   - Call `getRouteURL()` to extract the assigned hostname
   - Find the device pairing Route in the manifest objects
   - Set `spec.host` on the device pairing Route to match the main Route's hostname
   - Apply the device pairing Route

2. **Add new manifest**: Add `device-pairing-route.yaml` to the manifests directory

**Alternative approach**: Apply both Routes without explicit `spec.host`, accept that they get different hostnames, and expose both URLs in the Claw status. Rejected because it complicates the user experience.

## Risks / Trade-offs

**[Risk: Hostname injection requires two-phase application]**
→ **Mitigation**: The reconciliation flow already has phases. Add hostname extraction after main Route application and injection before device pairing Route application. Requeue if main Route status not yet populated.

**[Risk: Main Route hostname changes]**
→ **Mitigation**: Route hostnames are stable once assigned. If the main Route is deleted and recreated, the controller will detect the new hostname and update the device pairing Route accordingly.

**[Risk: Vanilla Kubernetes incompatibility]**
→ **Mitigation**: Route creation already skips gracefully on non-OpenShift clusters (NoMatchError). Device pairing will only be accessible via port-forward on vanilla Kubernetes, consistent with current behavior.

**[Trade-off: OpenShift-specific solution]**
→ **Accepted**: This operator is designed for OpenShift. Using Routes instead of Ingress keeps the solution simple and OpenShift-native.

## Migration Plan

No migration required. This is a Route manifest change that adds a new path mapping. Existing deployments will pick up the change on next reconciliation. Rollback is a simple revert of the Route manifest.

## Open Questions

None. Implementation is straightforward Route manifest update.
