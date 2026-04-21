## 1. Create Deployment Manifest

- [x] 1.1 Create `internal/assets/manifests/device-pairing-deployment.yaml` with metadata name `claw-device-pairing`
- [x] 1.2 Set deployment spec.replicas to 1 and strategy.type to Recreate
- [x] 1.3 Configure container image to `quay.io/xcoulon/claw-device-pairing:latest` with imagePullPolicy IfNotPresent
- [x] 1.4 Add security context: runAsNonRoot=true, allowPrivilegeEscalation=false, readOnlyRootFilesystem=true, drop ALL capabilities
- [x] 1.5 Set pod-level securityContext with seccompProfile type RuntimeDefault
- [x] 1.6 Add resource requests (memory: 128Mi, cpu: 100m) and limits (memory: 256Mi, cpu: 500m)
- [x] 1.7 Configure container port 8080 for the device pairing API
- [x] 1.8 Set automountServiceAccountToken to false for security

## 2. Create Service Manifest

- [x] 2.1 Create `internal/assets/manifests/device-pairing-service.yaml` with metadata name `claw-device-pairing`
- [x] 2.2 Set service type to ClusterIP
- [x] 2.3 Configure port mapping: port 8080, targetPort 8080, protocol TCP
- [x] 2.4 Add selector labels matching the device-pairing deployment (app: claw-device-pairing)

## 3. Update Kustomization

- [x] 3.1 Add `device-pairing-deployment.yaml` to the resources list in `internal/assets/manifests/kustomization.yaml`
- [x] 3.2 Add `device-pairing-service.yaml` to the resources list in `internal/assets/manifests/kustomization.yaml`

## 4. Verify Manifests

- [x] 4.1 Run `make manifests` to ensure CRD generation still works
- [x] 4.2 Run `make test` to ensure controller tests pass with new embedded manifests
- [x] 4.3 Verify device-pairing manifests follow the same labeling pattern as existing resources (app.kubernetes.io/name: claw will be added via kustomization.yaml commonLabels)
