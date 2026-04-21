## MODIFIED Requirements

### Requirement: Controller builds Kustomize manifests in-memory
The controller SHALL use the Kustomize API to build manifests in-memory from the embedded manifests directory.

#### Scenario: Kustomize build from embedded filesystem
- **WHEN** reconciling an Claw named 'instance'
- **THEN** the controller loads the embedded manifests filesystem
- **THEN** the controller invokes kustomize.Run() to build the resource map

#### Scenario: Kustomization file specifies resources
- **WHEN** the Kustomize build executes
- **THEN** it SHALL process the kustomization.yaml file in internal/assets/manifests/
- **THEN** the kustomization SHALL reference all resource YAML files (configmap.yaml, pvc.yaml, deployment.yaml, device-pairing-deployment.yaml, device-pairing-service.yaml)

#### Scenario: Common labels applied via Kustomize
- **WHEN** the Kustomize build executes
- **THEN** all resources SHALL have the label `app.kubernetes.io/name: claw` applied via commonLabels in kustomization.yaml
