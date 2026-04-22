## ADDED Requirements

### Requirement: Device pairing Route manifest exists
The embedded manifests SHALL include a separate Route manifest for the device pairing service.

#### Scenario: Device pairing Route file exists
- **WHEN** examining the internal/assets/manifests directory
- **THEN** it SHALL contain a file named device-pairing-route.yaml

#### Scenario: Device pairing Route has correct metadata
- **WHEN** examining the device-pairing-route.yaml manifest
- **THEN** metadata.name SHALL be "claw-device-pairing"
- **THEN** it SHALL include label "app: claw" for consistency

### Requirement: Device pairing Route specifies path
The device pairing Route SHALL specify /pair-device as its path to enable path-based routing.

#### Scenario: Path is set to /pair-device
- **WHEN** examining the device-pairing-route.yaml manifest
- **THEN** spec.path SHALL be "/pair-device"

### Requirement: Device pairing Route targets device pairing service
The device pairing Route SHALL route traffic to the claw-device-pairing service.

#### Scenario: Backend service is claw-device-pairing
- **WHEN** examining the device-pairing-route.yaml manifest
- **THEN** spec.to.name SHALL be "claw-device-pairing"
- **THEN** spec.to.kind SHALL be "Service"
- **THEN** spec.to.weight SHALL be 100

#### Scenario: Port configuration for device pairing
- **WHEN** examining the device-pairing-route.yaml manifest
- **THEN** spec.port.targetPort SHALL be 8080

### Requirement: Device pairing Route has TLS configuration
The device pairing Route SHALL use edge TLS termination with HTTP redirect.

#### Scenario: TLS settings configured
- **WHEN** examining the device-pairing-route.yaml manifest
- **THEN** spec.tls.termination SHALL be "edge"
- **THEN** spec.tls.insecureEdgeTerminationPolicy SHALL be "Redirect"

### Requirement: Main Route has explicit catch-all path
The main Route SHALL explicitly specify path "/" to catch all requests not matching other Routes.

#### Scenario: Main Route path is root
- **WHEN** examining the route.yaml manifest
- **THEN** spec.path SHALL be "/"

#### Scenario: Main Route backend unchanged
- **WHEN** examining the route.yaml manifest
- **THEN** spec.to.name SHALL be "claw"
- **THEN** spec.to.kind SHALL be "Service"
- **THEN** spec.to.weight SHALL be 100

### Requirement: Main Route settings preserved
The main Route SHALL maintain existing TLS and timeout configuration.

#### Scenario: TLS settings unchanged
- **WHEN** examining the route.yaml manifest
- **THEN** spec.tls.termination SHALL be "edge"
- **THEN** spec.tls.insecureEdgeTerminationPolicy SHALL be "Redirect"

#### Scenario: Timeout annotation preserved
- **WHEN** examining the route.yaml manifest
- **THEN** annotations SHALL include "haproxy.router.openshift.io/timeout"
- **THEN** the timeout value SHALL be "3600s"

#### Scenario: Port configuration unchanged
- **WHEN** examining the route.yaml manifest
- **THEN** spec.port.targetPort SHALL be 18789

### Requirement: Kustomization includes both Routes
The kustomization.yaml SHALL include both Route manifests in its resources list.

#### Scenario: Both Routes in kustomization
- **WHEN** examining the kustomization.yaml manifest
- **THEN** resources SHALL include "route.yaml"
- **THEN** resources SHALL include "device-pairing-route.yaml"

### Requirement: Controller injects hostname into device pairing Route
The controller SHALL inject the main Route's hostname into the device pairing Route to ensure both Routes share the same hostname.

#### Scenario: applyRouteOnly handles hostname injection
- **WHEN** the controller executes applyRouteOnly
- **THEN** it SHALL apply the main Route first without explicit spec.host
- **THEN** it SHALL extract the hostname from the main Route's status.ingress[0].host
- **THEN** it SHALL inject that hostname into the device pairing Route's spec.host field
- **THEN** it SHALL apply the device pairing Route with the shared hostname

#### Scenario: Hostname extraction from Route status
- **WHEN** the controller calls getRouteURL to extract the hostname
- **THEN** it SHALL return the HTTPS URL (e.g., "https://claw-namespace.apps.cluster.com")
- **THEN** the controller SHALL extract just the hostname portion (without "https://")
- **THEN** that hostname SHALL be injected into spec.host of the device pairing Route

#### Scenario: Both Routes share same hostname
- **WHEN** both Routes are successfully applied
- **THEN** the main Route SHALL have an OpenShift-assigned hostname in status.ingress[0].host
- **THEN** the device pairing Route SHALL have the same hostname in spec.host and status.ingress[0].host
- **THEN** OpenShift router SHALL use path-based routing (longest-prefix matching)

### Requirement: Controller constants updated for device pairing Route
The controller SHALL define a constant for the device pairing Route name.

#### Scenario: ClawDevicePairingRouteName constant exists
- **WHEN** examining the controller constants
- **THEN** ClawDevicePairingRouteName SHALL be defined with value "claw-device-pairing"
