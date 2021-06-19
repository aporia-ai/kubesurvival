package kubesimulator

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/scheduler/algorithm"

	"github.com/pfnet-research/k8s-cluster-simulator/pkg/clock"
	"github.com/pfnet-research/k8s-cluster-simulator/pkg/metrics"
	"github.com/pfnet-research/k8s-cluster-simulator/pkg/submitter"
)

type Submitter struct {
	pods []*v1.Pod
}

func newSubmitter(pods []*v1.Pod) *Submitter {
	return &Submitter{
		pods: pods,
	}
}

func (s *Submitter) Submit(clock clock.Clock, _ algorithm.NodeLister, met metrics.Metrics) ([]submitter.Event, error) {
	events := []submitter.Event{}

	for _, pod := range s.pods {
		if pod.ObjectMeta.Namespace == "" {
			pod.ObjectMeta.Namespace = "default"
		}

		events = append(events, &submitter.SubmitEvent{Pod: pod})
	}

	events = append(events, &submitter.TerminateSubmitterEvent{})

	return events, nil
}
