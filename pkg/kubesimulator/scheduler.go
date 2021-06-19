package kubesimulator

import (
	"github.com/pfnet-research/k8s-cluster-simulator/pkg/scheduler"
	"k8s.io/kubernetes/pkg/scheduler/algorithm/predicates"
	"k8s.io/kubernetes/pkg/scheduler/algorithm/priorities"
)

func buildScheduler() scheduler.Scheduler {
	// 1. Create a generic scheduler that mimics a kube-scheduler.
	sched := scheduler.NewGenericScheduler( /* preemption enabled */ true)

	// 2. Register extender(s)
	sched.AddExtender(
		scheduler.Extender{
			Name:             "MyExtender",
			Filter:           filterExtender,
			Prioritize:       prioritizeExtender,
			Weight:           1,
			NodeCacheCapable: true,
		},
	)

	// 2. Register plugin(s)
	// Predicate
	sched.AddPredicate("GeneralPredicates", predicates.GeneralPredicates)
	// Prioritizer
	sched.AddPrioritizer(priorities.PriorityConfig{
		Name:   "BalancedResourceAllocation",
		Map:    priorities.BalancedResourceAllocationMap,
		Reduce: nil,
		Weight: 1,
	})
	sched.AddPrioritizer(priorities.PriorityConfig{
		Name:   "LeastRequested",
		Map:    priorities.LeastRequestedPriorityMap,
		Reduce: nil,
		Weight: 1,
	})

	return &sched
}
