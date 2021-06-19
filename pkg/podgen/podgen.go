package podgen

import (
	"fmt"

	"github.com/aporia-ai/kubesurvival/v2/pkg/parser"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/api/resource"
)

type PodGenerator struct {
	errors          []Error
	pods            []*corev1.Pod
	currentPodIndex int64
}

// Podgen generates a list of pods from an expression
func Podgen(expression parser.Expression) ([]*corev1.Pod, []Error) {
	c := &PodGenerator{
		errors: []Error{},
		pods:   []*corev1.Pod{},
	}

	c.PodgenExpression(expression)

	return c.pods, c.errors
}

// PodgenExpression generates pods for an expression.
func (c *PodGenerator) PodgenExpression(node parser.Expression) {
	switch s := node.(type) {
	case *parser.PodExpression:
		c.PodgenPodExpression(s)
	case *parser.ArithmeticExpression:
		c.PodgenArithmeticExpression(s)
	}
}

// PodgenPodExpression generates pods for a pod expression.
func (c *PodGenerator) PodgenPodExpression(node *parser.PodExpression) {
	resources := corev1.ResourceList{}

	cpu := c.ParseQuantity(node.CPU)
	if cpu != nil {
		resources["cpu"] = *cpu
	}

	memory := c.ParseQuantity(node.Memory)
	if memory != nil {
		resources["memory"] = *memory
	}

	gpu := c.ParseQuantity(node.GPU)
	if gpu != nil {
		resources["nvidia.com/gpu"] = *gpu
	}

	c.pods = append(c.pods, &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("pod-%d", c.currentPodIndex),
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "container",
					Image: "container",
					Resources: corev1.ResourceRequirements{
						Requests: resources,
					},
				},
			},
		},
	})

	c.currentPodIndex++
}

func (c *PodGenerator) ParseQuantity(node parser.Expression) *resource.Quantity {
	switch q := node.(type) {
	case *parser.IntLiteral:
		return resource.NewScaledQuantity(q.Value, 0)

	case *parser.StringLiteral:
		result, err := resource.ParseQuantity(q.Value)
		if err != nil {
			c.errors = append(c.errors, Error{
				Message: err.Error(),
				Pos:     q.Position,
			})

			return nil
		}

		return &result

	default:
		return nil
	}
}

// PodgenArithmeticExpression generates pods for an arithmetic expression.
func (c *PodGenerator) PodgenArithmeticExpression(node *parser.ArithmeticExpression) {

	switch node.Operator {
	case parser.Add:
		c.PodgenExpression(node.LHS)
		c.PodgenExpression(node.RHS)

	case parser.Multiply:
		// One of LHS or RHS must be an integer.

		// Try to parse LHS as integer.
		multiplier, isLHSInteger := node.LHS.(*parser.IntLiteral)
		if !isLHSInteger {
			// If it didn't work, then RHS must be an integer.
			var ok bool
			if multiplier, ok = node.RHS.(*parser.IntLiteral); !ok {
				c.errors = append(c.errors, Error{
					Message: "one of [lhs, rhs] must be an integer in a multiply expression",
					Pos:     node.Position,
				})
				return
			}
		}

		var i int64
		for i = 0; i < multiplier.Value; i++ {
			var exp parser.Expression
			if isLHSInteger {
				exp = node.RHS
			} else {
				exp = node.LHS
			}

			c.PodgenExpression(exp)
		}
	}
}
