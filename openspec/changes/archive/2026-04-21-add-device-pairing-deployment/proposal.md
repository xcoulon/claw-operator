## Why

The Claw operator needs to deploy the claw-device-pairing application as part of the Claw instance. This application handles device pairing requests, enabling users to pair their devices with the Claw instance for authorized access. Currently, the operator has a NodePairingRequestApproval CRD but lacks the actual service implementation that processes pairing requests.

## What Changes

- Add `device-pairing-deployment.yaml` to `internal/assets/manifests/` with deployment manifest for claw-device-pairing application
- Add `device-pairing-service.yaml` to `internal/assets/manifests/` to expose the device pairing service
- Update `internal/assets/manifests/kustomization.yaml` to include the new device-pairing manifests in the resource list
- Configure deployment with image `quay.io/xcoulon/claw-device-pairing:latest`
- Apply security context hardening (non-root, seccomp, dropped capabilities) consistent with existing deployments

## Capabilities

### New Capabilities
- `device-pairing-deployment`: Deployment manifest for the claw-device-pairing application in the embedded manifests directory

### Modified Capabilities
- `unified-kustomize-controller`: Update kustomization.yaml resource list to include device-pairing manifests

## Impact

- **Embedded manifests**: New YAML files in `internal/assets/manifests/` directory
  - `device-pairing-deployment.yaml` — Deployment for device pairing service
  - `device-pairing-service.yaml` — Service to expose device pairing API
- **Kustomize configuration**: Updated `kustomization.yaml` to include new resources
- **Controller**: No code changes needed — controller already applies all resources from kustomization.yaml via server-side apply
- **Deployment**: Device pairing service will be deployed automatically when Claw instance reconciles
- **Backward compatibility**: Fully backward compatible — new resources are additive
