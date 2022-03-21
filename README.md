[![Build](https://github.com/ardoq/k8s-ardoq-bridge/actions/workflows/build.yml/badge.svg)](https://github.com/ardoq/k8s-ardoq-bridge/actions/workflows/build.yml)
[![Release](https://github.com/ardoq/k8s-ardoq-bridge/actions/workflows/release.yml/badge.svg)](https://github.com/ardoq/k8s-ardoq-bridge/actions/workflows/release.yml)

# K8s-Ardoq Bridge

## What is this?

The project captures data on application resources and nodes running in a Kubernetes cluster and continuously syncs the current state and select information into Ardoq as part of the application hosting Model. The project runs as a lightweight operator in your cluster, watching the resources and ensuring the data in your Ardoq workspace is constantly updated.

## Setup and Installation
### Setup
- Log into Ardoq
- Create an empty workspace based on the "Blank workspace" template
- The operator shall bootstrap the required model and custom field types on its first run, no need to worry about that.

### Installation
**Remote Helm repository**
Add and update the helm repository
```shell
helm repo add ardoq https://ardoq.github.io/k8s-ardoq-bridge/
helm repo update
```
Ensure repository is now available in your repo list
```shell
helm repo list
```
Deploy the helm chart
```shell
helm upgrade --install k8s-ardoq-bridge ardoq/k8s-ardoq-bridge --set "ardoq.baseUri='https://{your_custom_domain_here}/api/',ardoq.org='{your_org_label_here}',ardoq.workspaceId='{your_workspace_id_here}',ardoq.apiKey={your_api_key_here},ardoq.cluster='{your_cluster_name_here}'"
```
**Local Helm repository**
```shell
helm upgrade --install k8s-ardoq-bridge ./helm/chart --set "ardoq.baseUri='https://{your_custom_domain_here}/api/',ardoq.org='{your_org_label_here}',ardoq.workspaceId='{your_workspace_id_here}',ardoq.apiKey={your_api_key_here},ardoq.cluster='{your_cluster_name_here}'"
```
As an alternative to the --set option of the helm command, you can also edit the values.yaml file.

## How do I onboard an Application Resource?
On either the namespace or the specific Resource, you simply add:
```yaml
sync-to-ardoq: "enabled"
```
Additionally, you can enrich the resources as it is stored in Ardoq. This can further/later be used to assign ownership of the project. Options:
```yaml
ardoq/stack: "nginx"
ardoq/team: DevOps
ardoq/project: "TestProject"
```

### How can I exclude a resource in an enabled namespace?
You can disable single/multiple resource(s) in a monitored namespace, by:
```yaml
sync-to-ardoq: "disabled"
```

## What are Application resources?

These are the types of Kubernetes resources that we monitor. Currently, we only watch StatefulSets and Deployments.

## How does it works?

The operator uses the Kubernetes watch interface to consistently get updates on the resources and post them to your
workspace. Each resource type monitored,i.e Deployments,StatefulSets and Nodes, gets its own goroutine (lightweight thread) and performs the
sync serially ensuring "no resource gets left behind". It performs a single syncing pass on initialisation capturing all
the labelled Application Resources and only performs subsequent syncs if the data stored in-memory differs from resource
details being updated.

### Can I get a bit more detail?

Sure,you are welcome to dig into the code or join us in the discussion board (https://github.com/ardoq/k8s-ardoq-bridge/discussions), we are very happy to help.

## What resources does it monitor?

We are currently only monitoring:

- Application Resources:
    - Deployments
    - StatefulSets
- Nodes


## What information does it capture?

### Application Resources:

- Name
- ResourceType
- Namespace
- Replicas
- Image:  List of container images running
- CreationTimestamp
- Stack: Custom Field based on a unique label; "ardoq/stack"
- Team: Custom Field based on a unique label; "ardoq/team"
- Project: Custom Field based on a unique label; "ardoq/project"

### Nodes
- Name
- Architecture
- Capacity: CPU, Memory, Storage and Number of Pods
- Allocatable: CPU, Memory, Storage and Number of Pods
- ContainerRuntime
- KernelVersion
- KubeletVersion
- KubeProxyVersion
- OperatingSystem
- OSImage
- Provider
- CreationTimestamp
- Region: Based on the label; "failure-domain.beta.kubernetes.io/region"
- Zone: Based on the label; "failure-domain.beta.kubernetes.io/zone"

## What events does it monitor?
The operator only captures data on: Addition, Modification or Deletion of a resource
## Who has access to this data?
This syncs directly to a workspace you have created in your organisation. Access to this workspace is based on the access matrix you have defined for the given workspace. More details: https://help.ardoq.com/en/articles/1812349-user-roles-and-workspace-permissions

## Where can I find more details of how it interfaces with Ardoq?
The interface is based on the Ardoq Rest API as documented in the Swagger docs.

## Can I contribute?
Sure thing, all contributions are very welcome; improvements and optimisations even more.

### How?
We have a contribution guideline available under the docs [Contribution Guideline](./docs/CONTRIBUTING.md) and [Development](./docs/DEVELOPMENT.md)

## Any more details?
Yes. A lot more available in: [Documentation](./docs)

# Shoutouts(Built on top of previous work by):
- https://github.com/mories76/ardoq-client-go
- https://github.com/AlexsJones/KubeOps