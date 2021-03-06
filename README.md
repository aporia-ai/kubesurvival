# 💰 KubeSurvival
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/aporia-ai/kubesurvival?sort=semver&style=flat)](https://github.com/aporia-ai/kubesurvival/releases/latest)
[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/aporia-ai/kubesurvival/Test?label=build%20%26%20tests&style=flat)](https://github.com/aporia-ai/kubesurvival/actions?workflow=test)
[![Maintainability](https://api.codeclimate.com/v1/badges/e301f215e966335dc6bd/maintainability)](https://codeclimate.com/github/aporia-ai/kubesurvival/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/e301f215e966335dc6bd/test_coverage)](https://codeclimate.com/github/aporia-ai/kubesurvival/test_coverage)

KubeSurvival allows you to **significantly reduce your Kubernetes compute costs** by finding the cheapest machine types that can run your workloads successfully.

If you have a multi-tenant environment, ML training jobs, a large number of ML model servers, etc, this tool can help you optimize your K8s compute costs.

To easily define workloads, KubeSurvival uses a very simple DSL:

```python
  (
    # Some microservice
    pod(cpu: 1, memory: "1Gi") + 

    # Another microservice - with 3 replicas
    pod(cpu: "500m", memory: "2Gi") * 3 +

    # More microservices!
    (
      pod(cpu: 1, memory: "1Gi") +
      pod(cpu: "250m", memory: "1Gi")
    ) * 3
  ) * 2  # Production, Staging
```

This will give you a result such as:

    Instance type: t3.medium
    Node count: 11
    Total Price per Month: USD $340.45
    
<p align="center">Made with :heart: by <a href="https://www.aporia.com?utm_source=github&utm_medium=github&utm_campaign=kubesurvival" target="_blank">Aporia</a></p>

## Installation

Download a precompiled binary for your operating system from the [Releases](https://github.com/aporia-ai/kubesurvival/releases) page.

Alternatively, if you have Go installed, you can run:

```console
$ go install github.com/aporia-ai/kubesurvival/v2
```

## Usage

To run KubeSurvival:

    ./kubesurvival config.yaml

See the [examples](examples/) directory for example config files.

## How does it work?

KubeSurvival uses [k8s-cluster-simulator](https://github.com/pfnet-research/k8s-cluster-simulator) to simulate Kubernetes pod scheduling, without running on the actual underlying machines. It iterates over all possible instance types and node counts, simulates a K8s cluster with your workload, and checks if there are any pending pods. 

For each simulation it calculates the on-demand cost per month using the [ec2-instances-info](https://github.com/cristim/ec2-instances-info) library. Additionally, it queries the [eni-max-pods.txt](https://github.com/awslabs/amazon-eks-ami/blob/master/files/eni-max-pods.txt) file to determine what's the maximum number of pods in each instance type.

When simulating a cluster, KubeSurvival always makes sure you have 10% free CPU and Memory on each node.

Finally, KubeSurvival selects the cheapest configuration without pending pods.

## What's missing from this?

Well... a lot actually. Here's a partial list:

* Support for AKS and GKE
* Support for calculating costs of EBS storages
* Support for different node groups (e.g 2 machines with GPU + 4 machines without GPU)
* and probably much more!

We would love your help! ❤️
