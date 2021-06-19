package kubesimulator

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/scheduler/api"
)

func filterExtender(args api.ExtenderArgs) api.ExtenderFilterResult {
	// Filters out no nodes.
	return api.ExtenderFilterResult{
		Nodes:       &v1.NodeList{},
		NodeNames:   args.NodeNames,
		FailedNodes: api.FailedNodesMap{},
		Error:       "",
	}
}

func prioritizeExtender(args api.ExtenderArgs) api.HostPriorityList {
	// Ranks all nodes equally.
	priorities := make(api.HostPriorityList, 0, len(*args.NodeNames))
	for _, name := range *args.NodeNames {
		priorities = append(priorities, api.HostPriority{Host: name, Score: 1})
	}

	return priorities
}
