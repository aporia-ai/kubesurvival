package nodesource

import (
	"bufio"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	ec2instancesinfo "github.com/cristim/ec2-instances-info"
	"github.com/pfnet-research/k8s-cluster-simulator/pkg/config"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AWSNode struct {
	InstanceType  string   `json:"instanceType"`
	OnDemandPrice float64  `json:"onDemandPriceUSD"`
	VCPU          int      `json:"vcpu"`
	Memory        float32  `json:"memory"`
	GPU           int      `json:"gpu"`
	MaxPods       int      `json:"maxPods"`
	Arch          []string `json:"arch"`
	// TODO: Add VolumeSize ondemand price.
}

type AWSNodeSource struct {
	AWSRegion           string
	InstanceTypes       []string
	VolumeSizePerNodeGB int64 // TODO
}

type fetchPriceAsyncResult struct {
	node *AWSNode
	err  error
}

func (s *AWSNodeSource) GetNodes() ([]*AWSNode, error) {
	instances, err := ec2instancesinfo.Data()
	if err != nil {
		return nil, errors.Wrap(err, "could not get ec2 instances info")
	}

	maxPodsPerInstance, err := s.getMaxPodsPerInstance()
	if err != nil {
		return nil, errors.Wrap(err, "could not get max pods per instance")
	}

	nodes := []*AWSNode{}

	for _, instanceType := range s.InstanceTypes {
		// Find max pods for this instance
		maxPods, ok := maxPodsPerInstance[instanceType]
		if !ok {
			return nil, errors.New(fmt.Sprintf("Could not find max pods for instance: %s", instanceType))
		}

		// Find info for this instance
		found := false
		for _, instance := range *instances {
			if instanceType == instance.InstanceType {
				nodes = append(nodes, &AWSNode{
					InstanceType:  instance.InstanceType,
					OnDemandPrice: instance.Pricing[s.AWSRegion].Linux.OnDemand,
					VCPU:          instance.VCPU,
					Memory:        instance.Memory,
					GPU:           instance.GPU,
					Arch:          instance.Arch,
					MaxPods:       maxPods,
				})

				found = true
				break
			}
		}

		if !found {
			return nil, errors.New(fmt.Sprintf("Could not find instance data for %s", instanceType))
		}
	}

	return nodes, nil
}

func (s *AWSNodeSource) getMaxPodsPerInstance() (map[string]int, error) {
	response, err := http.Get("https://raw.githubusercontent.com/awslabs/amazon-eks-ami/master/files/eni-max-pods.txt")
	if err != nil {
		return nil, errors.Wrap(err, "could not fetch max pods list")
	}
	defer response.Body.Close()

	maxPodsPerInstance := make(map[string]int)

	// Check status code
	if response.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(err, "fetch max pods list did not return 200 (%d instead)", response.StatusCode)
	}

	// Read and parse response
	scanner := bufio.NewScanner(response.Body)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip comments
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Parse line
		splitted := strings.Split(line, " ")
		if len(splitted) != 2 {
			return nil, errors.Errorf("could not parse eni-max-pods.txt file, bad line: %s", line)
		}

		instanceType := splitted[0]
		maxPods, err := strconv.ParseInt(splitted[1], 10, 32)
		if err != nil {
			return nil, errors.Errorf("could not parse eni-max-pods.txt file, bad line: %s", line)
		}

		maxPodsPerInstance[instanceType] = int(maxPods)
	}

	return maxPodsPerInstance, nil
}

func (n *AWSNode) GetHourlyPrice() float64 {
	// TODO: Add storage price
	return n.OnDemandPrice
}

func (n *AWSNode) GetNodeConfig(nodeName string) *config.NodeConfig {
	return &config.NodeConfig{
		Metadata: metav1.ObjectMeta{
			Name: nodeName,
			Labels: map[string]string{
				"beta.kubernetes.io/os": "simulated",
			},
		},
		Spec: v1.NodeSpec{
			Unschedulable: false,
		},
		Status: config.NodeStatus{
			Allocatable: map[v1.ResourceName]string{
				"cpu":            fmt.Sprintf("%d", n.VCPU),
				"memory":         fmt.Sprintf("%dGi", int(n.Memory)),
				"nvidia.com/gpu": fmt.Sprintf("%d", n.GPU),
				"pods":           fmt.Sprintf("%d", n.MaxPods),
			},
		},
	}
}
