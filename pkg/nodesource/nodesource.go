package nodesource

import "github.com/pfnet-research/k8s-cluster-simulator/pkg/config"

type Node interface {
	GetHourlyPrice() float64
	GetNodeConfig(nodeName string) *config.NodeConfig
}

type NodeSource interface {
	GetNodes() ([]Node, error)
}
