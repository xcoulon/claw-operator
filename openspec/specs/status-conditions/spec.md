## ADDED Requirements

### Requirement: Claw status has Conditions field
The ClawStatus struct SHALL include a Conditions field of type []metav1.Condition to track instance state.

#### Scenario: Conditions field in CRD status
- **WHEN** the Claw CRD is defined
- **THEN** the ClawStatus struct SHALL contain a field `Conditions []metav1.Condition` with JSON tag `conditions`
- **THEN** the field SHALL have kubebuilder marker `+listType=map` and `+listMapKey=type`
- **THEN** the field SHALL be optional (conditions can be empty on creation)

#### Scenario: Generated CRD includes status subresource
- **WHEN** make manifests is run
- **THEN** the generated CRD YAML SHALL include status subresource configuration
- **THEN** the status.conditions field SHALL be present in the OpenAPI schema

### Requirement: DeploymentsReady condition indicates deployment readiness
The controller SHALL maintain a DeploymentsReady condition type to indicate whether the Claw deployment pods are running.

#### Scenario: DeploymentsReady condition set to True when pods running
- **WHEN** both claw and claw-proxy Deployments have Available=True status
- **THEN** the controller SHALL set DeploymentsReady condition with status=True, reason=PodsRunning, message confirming both deployments are available

#### Scenario: DeploymentsReady condition set to False when pods not ready
- **WHEN** either claw or claw-proxy Deployment has Available condition not equal to True
- **THEN** the controller SHALL set DeploymentsReady condition with status=False, reason=Provisioning, message indicating which deployments are pending

#### Scenario: DeploymentsReady uses same deployment check logic
- **WHEN** evaluating DeploymentsReady condition status
- **THEN** the controller SHALL fetch both claw and claw-proxy Deployments and check their Available conditions
- **THEN** the controller SHALL use the same deployment readiness check that was previously used for the Ready condition

### Requirement: Claw CRD includes printcolumn for condition status
The Claw CRD SHALL include a printcolumn that displays the DeploymentsReady condition status.

#### Scenario: Printcolumn shows DeploymentsReady status
- **WHEN** examining the Claw CRD kubebuilder markers in api/v1alpha1/claw_types.go
- **THEN** the printcolumn SHALL be named "Ready" for display purposes
- **THEN** the printcolumn SHALL use JSONPath `.status.conditions[?(@.type=="DeploymentsReady")].status`

#### Scenario: Printcolumn shows DeploymentsReady reason
- **WHEN** examining the Claw CRD kubebuilder markers in api/v1alpha1/claw_types.go
- **THEN** the Reason printcolumn SHALL use JSONPath `.status.conditions[?(@.type=="DeploymentsReady")].reason`

### Requirement: CredentialsResolved condition tracks credential validation
The controller SHALL maintain a CredentialsResolved condition type to indicate whether all credential Secrets have been validated.

#### Scenario: CredentialsResolved set to True when all Secrets exist
- **WHEN** all Secrets referenced in spec.credentials exist and are accessible
- **THEN** the controller SHALL set CredentialsResolved condition with status=True, reason=Resolved

#### Scenario: CredentialsResolved set to False when Secret missing
- **WHEN** any Secret referenced in spec.credentials does not exist
- **THEN** the controller SHALL set CredentialsResolved condition with status=False, reason=ValidationFailed, message identifying missing Secret

#### Scenario: CredentialsResolved omitted when no credentials configured
- **WHEN** spec.credentials is empty or omitted
- **THEN** the controller MAY omit the CredentialsResolved condition OR set it to True

### Requirement: ProxyConfigured condition tracks proxy configuration
The controller SHALL maintain a ProxyConfigured condition type to indicate whether the proxy has been configured successfully.

#### Scenario: ProxyConfigured set to True when proxy configured
- **WHEN** the openclaw-proxy Deployment has been configured with credentials
- **THEN** the controller SHALL set ProxyConfigured condition with status=True, reason=Configured

#### Scenario: ProxyConfigured set to False on configuration failure
- **WHEN** proxy configuration fails (e.g., cannot stamp Secret versions)
- **THEN** the controller SHALL set ProxyConfigured condition with status=False, reason=ConfigFailed, message describing the failure

### Requirement: Controller checks Deployment status conditions
The controller SHALL read the Available condition from both openclaw and openclaw-proxy Deployment status to determine readiness.

#### Scenario: Fetch openclaw Deployment status
- **WHEN** updating Claw status conditions
- **THEN** the controller SHALL fetch the Deployment named 'openclaw' in the same namespace
- **THEN** the controller SHALL read the Available condition from deployment.Status.Conditions

#### Scenario: Fetch openclaw-proxy Deployment status
- **WHEN** updating Claw status conditions
- **THEN** the controller SHALL fetch the Deployment named 'openclaw-proxy' in the same namespace
- **THEN** the controller SHALL read the Available condition from deployment.Status.Conditions

#### Scenario: Handle missing Deployment
- **WHEN** a Deployment is not found during status check
- **THEN** the controller SHALL treat the deployment as not ready
- **THEN** the Ready condition SHALL remain False with reason=Provisioning

### Requirement: Status updates use status subresource
The controller SHALL update Claw status using the Kubernetes status subresource, not the main resource.

#### Scenario: Status updated via client.Status()
- **WHEN** the controller updates Claw status conditions
- **THEN** the controller SHALL use `client.Status().Update(ctx, clawInstance)` or `client.Status().Patch(ctx, clawInstance, patch)`
- **THEN** status updates SHALL NOT trigger spec reconciliation

#### Scenario: Failed status update retried
- **WHEN** a status update fails due to conflict or API error
- **THEN** the controller SHALL return an error to trigger reconciliation retry
- **THEN** the next reconciliation SHALL attempt the status update again

### Requirement: Condition transitions update LastTransitionTime
The controller SHALL update the LastTransitionTime field only when the condition status changes.

#### Scenario: Status change updates LastTransitionTime
- **WHEN** the Ready condition changes from False to True
- **THEN** the controller SHALL set LastTransitionTime to the current timestamp
- **THEN** the controller SHALL update the reason and message fields

#### Scenario: Same status preserves LastTransitionTime
- **WHEN** the Ready condition status remains the same (e.g., False to False)
- **THEN** the controller SHALL preserve the existing LastTransitionTime value
- **THEN** the controller MAY update the reason or message fields

### Requirement: Condition uses standard meta fields
Each condition SHALL include all standard metav1.Condition fields: Type, Status, ObservedGeneration, LastTransitionTime, Reason, and Message.

#### Scenario: Condition fields populated
- **WHEN** the controller sets a condition
- **THEN** the condition SHALL have Type set to the condition type string (e.g., "Ready", "CredentialsResolved", "ProxyConfigured")
- **THEN** the condition SHALL have Status set to "True", "False", or "Unknown"
- **THEN** the condition SHALL have Reason set to a CamelCase single-word or hyphenated reason
- **THEN** the condition SHALL have Message set to a human-readable description
- **THEN** the condition SHALL have ObservedGeneration set to the Claw resource's metadata.generation
- **THEN** the condition SHALL have LastTransitionTime set to the time of the status change


### Requirement: DeploymentsReady condition type and reason constants defined
The API package SHALL define constants for the DeploymentsReady condition type and its reasons.

#### Scenario: DeploymentsReady condition type constant exists
- **WHEN** examining api/v1alpha1/claw_types.go
- **THEN** it SHALL define ConditionTypeDeploymentsReady = "DeploymentsReady"

#### Scenario: PodsRunning reason constant exists
- **WHEN** examining api/v1alpha1/claw_types.go
- **THEN** it SHALL define ConditionReasonPodsRunning = "PodsRunning"
- **THEN** this reason SHALL be used when both deployments are available

#### Scenario: Provisioning reason constant remains
- **WHEN** examining api/v1alpha1/claw_types.go
- **THEN** it SHALL still define ConditionReasonProvisioning = "Provisioning"
- **THEN** this reason SHALL be used when one or both deployments are not ready

### Requirement: Condition type constants defined in API
The API package SHALL define constants for condition types and reasons.

#### Scenario: Condition type constants exist
- **WHEN** examining api/v1alpha1/claw_types.go
- **THEN** it SHALL define ConditionTypeCredentialsResolved = "CredentialsResolved"
- **THEN** it SHALL define ConditionTypeProxyConfigured = "ProxyConfigured"

#### Scenario: Condition reason constants exist
- **WHEN** examining api/v1alpha1/claw_types.go
- **THEN** it SHALL define ConditionReasonResolved = "Resolved"
- **THEN** it SHALL define ConditionReasonValidationFailed = "ValidationFailed"
- **THEN** it SHALL define ConditionReasonConfigured = "Configured"
- **THEN** it SHALL define ConditionReasonConfigFailed = "ConfigFailed"
