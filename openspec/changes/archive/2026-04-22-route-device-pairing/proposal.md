## Why

The device pairing functionality requires its own service endpoint at `/pair-device` to handle node pairing requests separately from the main OpenClaw gateway. Currently, the Route only supports routing to a single service, preventing access to the device pairing service.

## What Changes

- Update the OpenShift Route manifest to support path-based routing
- Configure `/pair-device` path to route to `claw-device-pairing` service
- Maintain all other paths routing to the main `claw` service
- Preserve existing Route annotations (timeout, TLS settings)

## Capabilities

### New Capabilities
- `route-path-based-routing`: Path-based routing configuration for the OpenShift Route to support multiple backend services based on request path

### Modified Capabilities
<!-- No existing spec requirements are changing -->

## Impact

- `internal/assets/manifests/route.yaml`: Route manifest will be updated to include path-based routing rules
- `internal/controller/claw_resource_controller.go`: May need verification that Route reconciliation works with multiple backends
- External access: Users will be able to access device pairing at `https://<route-host>/pair-device`
