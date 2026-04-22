## 1. Create Device Pairing Route Manifest

- [x] 1.1 Create device-pairing-route.yaml in internal/assets/manifests/
- [x] 1.2 Set metadata.name to "claw-device-pairing"
- [x] 1.3 Add app: claw label
- [x] 1.4 Set spec.path to "/pair-device"
- [x] 1.5 Configure spec.to to target claw-device-pairing service on port 8080
- [x] 1.6 Add TLS configuration (edge termination, redirect HTTP)

## 2. Update Main Route Manifest

- [x] 2.1 Add spec.path: "/" to route.yaml for explicit catch-all
- [x] 2.2 Verify existing TLS and timeout settings are preserved
- [x] 2.3 Verify backend remains "claw" service on port 18789

## 3. Update Kustomization

- [x] 3.1 Add device-pairing-route.yaml to resources list in kustomization.yaml

## 4. Update Controller Route Application Logic

- [x] 4.1 Update applyRouteOnly() to handle hostname injection
- [x] 4.2 Apply main Route first without explicit spec.host
- [x] 4.3 Extract hostname from main Route status using getRouteURL()
- [x] 4.4 Inject hostname into device pairing Route spec.host field
- [x] 4.5 Apply device pairing Route with shared hostname

## 5. Update Controller Constants

- [x] 5.1 Add ClawDevicePairingRouteName constant with value "claw-device-pairing"
- [x] 5.2 Update buildKustomizedObjects() manifest file map to include device-pairing-route.yaml

## 6. Update Controller RBAC

- [x] 6.1 Verify Route RBAC permissions are sufficient (already present)
- [x] 6.2 Run make manifests to regenerate RBAC

## 7. Verification

- [x] 7.1 Run make manifests to ensure YAML is valid
- [x] 7.2 Run make generate to update generated code
- [x] 7.3 Run make test to verify controller tests pass
- [x] 7.4 Review both Route manifests to confirm all spec requirements are met
