package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"

	"github.com/aporia-ai/kubesurvival/v2/pkg/kubesimulator"
	"github.com/aporia-ai/kubesurvival/v2/pkg/nodesource"
	"github.com/aporia-ai/kubesurvival/v2/pkg/parser"
	"github.com/aporia-ai/kubesurvival/v2/pkg/podgen"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type Config struct {
	Nodes struct {
		AWS struct {
			Region        string   `yaml:"region"`
			InstanceTypes []string `yaml:"instanceTypes"`
		} `yaml:"aws"`
	} `yaml:"nodes"`
	Pods string `yaml:"pods"`
}

type Result struct {
	InstanceType       string
	NodeCount          int
	TotalPricePerMonth float64
}

func main() {
	// Read argument
	if len(os.Args) != 2 {
		fmt.Println("USAGE: ./kubesurvival <YAML_CONFIG_PATH>")
		os.Exit(1)
	}

	// Read config file
	configFile, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("[!] Could not read config file: %s\n", err)
		return
	}

	// Parse config file
	config := &Config{}
	err = yaml.Unmarshal(configFile, config)
	if err != nil {
		fmt.Printf("[!] Could not deserialize config file: %s\n", err)
		return
	}

	// Parse & generate pods
	exp, parseErrors := parser.Parse(config.Pods)
	if len(parseErrors) > 0 {
		for _, parseError := range parseErrors {
			fmt.Printf("[!] Parse error: %s\n", parseError.Error())
		}
		return
	}

	pods, podgenErrors := podgen.Podgen(exp)
	if len(podgenErrors) > 0 {
		for _, podgenError := range podgenErrors {
			fmt.Printf("[!] PodGen error: %s\n", podgenError.Error())
		}
		return
	}

	// Generate nodes
	ns := &nodesource.AWSNodeSource{
		AWSRegion:     config.Nodes.AWS.Region,
		InstanceTypes: config.Nodes.AWS.InstanceTypes,
	}

	nodeTypes, err := ns.GetNodes()
	if err != nil {
		fmt.Printf("Could not get node types: %s\n", err)
		return
	}

	// Remove node types if there's a pod with more resources than it
	filteredNodeTypes := filterNodeTypes(nodeTypes, pods)
	if len(filteredNodeTypes) == 0 {
		fmt.Printf("[!] No nodes are available for simulation.\n")
		return
	}

	// Main loop
	var result *Result
	for _, nodeType := range filteredNodeTypes {
		// We never want a cluster with only 1 node
		nodeCount := 2

		for {
			// Calculate total price per month
			totalPricePerMonth := float64(nodeCount) * nodeType.GetHourlyPrice() * 24 * 31

			// Do we even need to simulate?
			if result != nil && totalPricePerMonth > result.TotalPricePerMonth {
				break
			}

			// Generate a list of nodes from this type
			nodes := []nodesource.Node{}
			for i := 0; i < nodeCount; i++ {
				nodes = append(nodes, nodeType)
			}

			// Simulate cluster
			simulator := &kubesimulator.KubernetesSimulator{}
			isSimulationSuccessful, err := simulator.Simulate(pods, nodes)
			if err != nil {
				fmt.Printf("[!] Failed to simulate a Kubernetes cluster: %s\n", err)
				return
			}

			if isSimulationSuccessful {
				result = &Result{
					InstanceType:       nodeType.InstanceType,
					NodeCount:          nodeCount,
					TotalPricePerMonth: totalPricePerMonth,
				}

				break
			}

			// Simple heuristic as an alternative to nodeCount++ to make convergence faster.
			nodeCount += int(math.Max(float64(nodeCount)/15, 1))
		}
	}

	if result != nil {
		fmt.Printf("Instance type: %s\n", result.InstanceType)
		fmt.Printf("Node count: %d\n", result.NodeCount)
		fmt.Printf("Total Price per Month: USD $%.2f\n", result.TotalPricePerMonth)
	} else {
		fmt.Printf("[!] Could not converge to a solution.\n")
	}
}

func filterNodeTypes(nodeTypes []*nodesource.AWSNode, pods []*v1.Pod) []*nodesource.AWSNode {
	result := []*nodesource.AWSNode{}
	for _, nodeType := range nodeTypes {
		nodeHasEnoughResources := true

		for _, pod := range pods {
			// Is Pod CPU > Node CPU?
			nodeCpu := resource.MustParse(nodeType.GetNodeConfig("node").Status.Allocatable["cpu"])
			podCpu := pod.Spec.Containers[0].Resources.Requests.Cpu()
			if podCpu.Cmp(nodeCpu) > 0 {
				fmt.Printf("WARNING: Ignoring node type %s with %s CPU because there's a pod with more CPU: %s\n",
					nodeType.InstanceType, nodeCpu.String(), podCpu.String())
				nodeHasEnoughResources = false
				break
			}

			// Is Pod Memory > Node Memory?
			nodeMemory := resource.MustParse(nodeType.GetNodeConfig("node").Status.Allocatable["memory"])
			podMemory := pod.Spec.Containers[0].Resources.Requests.Memory()
			if podMemory.Cmp(nodeMemory) > 0 {
				fmt.Printf("WARNING: Ignoring node type %s with %s memory because there's a pod with more memory: %s\n",
					nodeType.InstanceType, nodeMemory.String(), podMemory.String())
				nodeHasEnoughResources = false
				break
			}

			// Is Pod Memory > Node Memory?
			nodeGpu := resource.MustParse(nodeType.GetNodeConfig("node").Status.Allocatable["nvidia.com/gpu"])
			podGpu := pod.Spec.Containers[0].Resources.Requests["nvidia.com/gpu"]
			if podGpu.Cmp(nodeGpu) > 0 {
				fmt.Printf("WARNING: Ignoring node type %s with %s GPU because there's a pod with more GPU: %s\n",
					nodeType.InstanceType, nodeGpu.String(), podGpu.String())
				nodeHasEnoughResources = false
				break
			}
		}

		if nodeHasEnoughResources {
			result = append(result, nodeType)
		}
	}

	return result
}
