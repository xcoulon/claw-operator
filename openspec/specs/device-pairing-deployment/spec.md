## ADDED Requirements

### Requirement: Device pairing deployment manifest exists
The embedded manifests SHALL include a deployment manifest for the claw-device-pairing application.

#### Scenario: Deployment manifest file exists
- **WHEN** examining the internal/assets/manifests directory
- **THEN** it SHALL contain a file named device-pairing-deployment.yaml

#### Scenario: Deployment uses correct image
- **WHEN** examining the device-pairing-deployment.yaml manifest
- **THEN** it SHALL specify image quay.io/xcoulon/claw-device-pairing:latest
- **THEN** it SHALL use imagePullPolicy IfNotPresent

#### Scenario: Deployment has security hardening
- **WHEN** examining the device-pairing-deployment.yaml manifest
- **THEN** it SHALL set runAsNonRoot to true
- **THEN** it SHALL set allowPrivilegeEscalation to false
- **THEN** it SHALL set readOnlyRootFilesystem to true
- **THEN** it SHALL drop ALL capabilities
- **THEN** it SHALL set seccompProfile type to RuntimeDefault

#### Scenario: Deployment has resource limits
- **WHEN** examining the device-pairing-deployment.yaml manifest
- **THEN** it SHALL define resource requests for memory and CPU
- **THEN** it SHALL define resource limits for memory and CPU

### Requirement: Device pairing service manifest exists
The embedded manifests SHALL include a service manifest to expose the claw-device-pairing application.

#### Scenario: Service manifest file exists
- **WHEN** examining the internal/assets/manifests directory
- **THEN** it SHALL contain a file named device-pairing-service.yaml

#### Scenario: Service exposes device pairing port
- **WHEN** examining the device-pairing-service.yaml manifest
- **THEN** it SHALL be of type ClusterIP
- **THEN** it SHALL define a port mapping for the device pairing API
- **THEN** it SHALL use selector labels matching the device-pairing deployment

### Requirement: Deployment name is claw-device-pairing
The device pairing deployment SHALL be named claw-device-pairing for consistency with naming conventions.

#### Scenario: Deployment metadata name
- **WHEN** examining the device-pairing-deployment.yaml manifest
- **THEN** the metadata.name SHALL be "claw-device-pairing"

#### Scenario: Service metadata name
- **WHEN** examining the device-pairing-service.yaml manifest
- **THEN** the metadata.name SHALL be "claw-device-pairing"

### Requirement: Deployment has single replica
The device pairing deployment SHALL run a single replica for simplicity.

#### Scenario: Replica count is one
- **WHEN** examining the device-pairing-deployment.yaml manifest
- **THEN** spec.replicas SHALL be 1

### Requirement: Deployment uses Recreate strategy
The device pairing deployment SHALL use Recreate strategy to avoid multiple instances during updates.

#### Scenario: Update strategy is Recreate
- **WHEN** examining the device-pairing-deployment.yaml manifest
- **THEN** spec.strategy.type SHALL be "Recreate"
