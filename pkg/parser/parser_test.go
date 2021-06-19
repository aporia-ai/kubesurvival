package parser_test

import (
	"bufio"
	"io"
	"strings"
	"testing"

	"github.com/aporia-ai/kubesurvival/v2/pkg/lexer"
	"github.com/aporia-ai/kubesurvival/v2/pkg/parser"

	"github.com/stretchr/testify/assert"
)

func TestEmptyPod(t *testing.T) {
	p := newParserNoPositions(strings.NewReader("pod()"))
	expression := p.ParseExpression()

	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.PodExpression{}, expression)
}

func TestPodArgsIntegers(t *testing.T) {
	p := newParserNoPositions(strings.NewReader("pod(cpu: 1, memory: 2, gpu: 5)"))
	expression := p.ParseExpression()

	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.PodExpression{
		CPU:    &parser.IntLiteral{Value: 1},
		Memory: &parser.IntLiteral{Value: 2},
		GPU:    &parser.IntLiteral{Value: 5},
	}, expression)
}

func TestPodArgsStrings(t *testing.T) {
	p := newParserNoPositions(strings.NewReader("pod(cpu: \"100m\", memory: \"10Gi\", gpu: \"123\")"))
	expression := p.ParseExpression()

	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.PodExpression{
		CPU:    &parser.StringLiteral{Value: "100m"},
		Memory: &parser.StringLiteral{Value: "10Gi"},
		GPU:    &parser.StringLiteral{Value: "123"},
	}, expression)
}

func TestPodArgsCombined(t *testing.T) {
	p := newParserNoPositions(strings.NewReader("pod(cpu: \"100m\", memory: \"10Gi\", gpu: 5)"))
	expression := p.ParseExpression()

	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.PodExpression{
		CPU:    &parser.StringLiteral{Value: "100m"},
		Memory: &parser.StringLiteral{Value: "10Gi"},
		GPU:    &parser.IntLiteral{Value: 5},
	}, expression)
}

func TestAdd2Pods(t *testing.T) {
	p := newParserNoPositions(strings.NewReader(`
		pod(cpu: "100m", memory: "10Gi", gpu: 5) +
		pod(cpu: 4, memory: "32Gi", gpu: 3)
	`))

	expression := p.ParseExpression()

	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.ArithmeticExpression{
		Operator: parser.Add,
		LHS: &parser.PodExpression{
			CPU:    &parser.StringLiteral{Value: "100m"},
			Memory: &parser.StringLiteral{Value: "10Gi"},
			GPU:    &parser.IntLiteral{Value: 5},
		},
		RHS: &parser.PodExpression{
			CPU:    &parser.IntLiteral{Value: 4},
			Memory: &parser.StringLiteral{Value: "32Gi"},
			GPU:    &parser.IntLiteral{Value: 3},
		},
	}, expression)
}

func TestAdd3Pods(t *testing.T) {
	p := newParserNoPositions(strings.NewReader(`
		pod(cpu: "100m", memory: "10Gi", gpu: 5) +
		pod(cpu: 4, memory: "32Gi", gpu: 3) + 
		pod(cpu: 0, memory: 0, gpu: 0)
	`))

	expression := p.ParseExpression()

	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.ArithmeticExpression{
		Operator: parser.Add,
		LHS: &parser.ArithmeticExpression{
			Operator: parser.Add,
			LHS: &parser.PodExpression{
				CPU:    &parser.StringLiteral{Value: "100m"},
				Memory: &parser.StringLiteral{Value: "10Gi"},
				GPU:    &parser.IntLiteral{Value: 5},
			},
			RHS: &parser.PodExpression{
				CPU:    &parser.IntLiteral{Value: 4},
				Memory: &parser.StringLiteral{Value: "32Gi"},
				GPU:    &parser.IntLiteral{Value: 3},
			},
		},
		RHS: &parser.PodExpression{
			CPU:    &parser.IntLiteral{Value: 0},
			Memory: &parser.IntLiteral{Value: 0},
			GPU:    &parser.IntLiteral{Value: 0},
		},
	}, expression)
}

func TestAdd3PodsWithParen(t *testing.T) {
	p := newParserNoPositions(strings.NewReader(`
		pod(cpu: "100m", memory: "10Gi", gpu: 5) +
			(pod(cpu: 4, memory: "32Gi", gpu: 3) + pod(cpu: 0, memory: 0, gpu: 0))
	`))

	expression := p.ParseExpression()

	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.ArithmeticExpression{
		Operator: parser.Add,
		LHS: &parser.PodExpression{
			CPU:    &parser.StringLiteral{Value: "100m"},
			Memory: &parser.StringLiteral{Value: "10Gi"},
			GPU:    &parser.IntLiteral{Value: 5},
		},
		RHS: &parser.ArithmeticExpression{
			Operator: parser.Add,
			LHS: &parser.PodExpression{
				CPU:    &parser.IntLiteral{Value: 4},
				Memory: &parser.StringLiteral{Value: "32Gi"},
				GPU:    &parser.IntLiteral{Value: 3},
			},
			RHS: &parser.PodExpression{
				CPU:    &parser.IntLiteral{Value: 0},
				Memory: &parser.IntLiteral{Value: 0},
				GPU:    &parser.IntLiteral{Value: 0},
			},
		},
	}, expression)
}

func TestAddPodAndNumber(t *testing.T) {
	p := newParserNoPositions(strings.NewReader(`
		pod(cpu: "100m", memory: "10Gi", gpu: 5) + 5
	`))

	p.ParseExpression()

	assert.NotEmpty(t, p.Errors)
}

func TestAddPodAndString(t *testing.T) {
	p := newParserNoPositions(strings.NewReader(`
		pod(cpu: "100m", memory: "10Gi", gpu: 5) + "asdf"
	`))

	p.ParseExpression()

	assert.NotEmpty(t, p.Errors)
}

func TestAddNumberAndPod(t *testing.T) {
	p := newParserNoPositions(strings.NewReader(`
		5 + pod(cpu: "100m", memory: "10Gi", gpu: 5)
	`))

	p.ParseExpression()

	assert.NotEmpty(t, p.Errors)
}

func TestMulPodByNumber(t *testing.T) {
	p := newParserNoPositions(strings.NewReader(`
		pod(cpu: "100m", memory: "10Gi", gpu: 5) * 5
	`))

	expression := p.ParseExpression()

	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.ArithmeticExpression{
		Operator: parser.Multiply,
		LHS: &parser.PodExpression{
			CPU:    &parser.StringLiteral{Value: "100m"},
			Memory: &parser.StringLiteral{Value: "10Gi"},
			GPU:    &parser.IntLiteral{Value: 5},
		},
		RHS: &parser.IntLiteral{Value: 5},
	}, expression)
}

func TestMulNumberByPod(t *testing.T) {
	p := newParserNoPositions(strings.NewReader(`
		5 * pod(cpu: "100m", memory: "10Gi", gpu: 5)
	`))

	expression := p.ParseExpression()

	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.ArithmeticExpression{
		Operator: parser.Multiply,
		LHS:      &parser.IntLiteral{Value: 5},
		RHS: &parser.PodExpression{
			CPU:    &parser.StringLiteral{Value: "100m"},
			Memory: &parser.StringLiteral{Value: "10Gi"},
			GPU:    &parser.IntLiteral{Value: 5},
		},
	}, expression)
}

func TestMulPodByPod(t *testing.T) {
	p := newParserNoPositions(strings.NewReader(`
		pod(cpu: "100m", memory: "10Gi", gpu: 5) * pod(cpu: "100m", memory: "10Gi", gpu: 5)
	`))

	p.ParseExpression()
	assert.NotEmpty(t, p.Errors)
}

func TestMulThenAdd(t *testing.T) {
	p := newParserNoPositions(strings.NewReader(`
		pod(cpu: "100m", memory: "10Gi", gpu: 5) * 5 + pod(cpu: "300m", memory: "32Gi", gpu: 10)
	`))

	expression := p.ParseExpression()

	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.ArithmeticExpression{
		Operator: parser.Add,
		LHS: &parser.ArithmeticExpression{
			Operator: parser.Multiply,
			LHS: &parser.PodExpression{
				CPU:    &parser.StringLiteral{Value: "100m"},
				Memory: &parser.StringLiteral{Value: "10Gi"},
				GPU:    &parser.IntLiteral{Value: 5},
			},
			RHS: &parser.IntLiteral{Value: 5},
		},
		RHS: &parser.PodExpression{
			CPU:    &parser.StringLiteral{Value: "300m"},
			Memory: &parser.StringLiteral{Value: "32Gi"},
			GPU:    &parser.IntLiteral{Value: 10},
		},
	}, expression)
}

func TestAddThenMul(t *testing.T) {
	p := newParserNoPositions(strings.NewReader(`
		pod(cpu: "300m", memory: "32Gi", gpu: 10) + pod(cpu: "100m", memory: "10Gi", gpu: 5) * 5
	`))

	expression := p.ParseExpression()

	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.ArithmeticExpression{
		Operator: parser.Add,
		LHS: &parser.PodExpression{
			CPU:    &parser.StringLiteral{Value: "300m"},
			Memory: &parser.StringLiteral{Value: "32Gi"},
			GPU:    &parser.IntLiteral{Value: 10},
		},
		RHS: &parser.ArithmeticExpression{
			Operator: parser.Multiply,
			LHS: &parser.PodExpression{
				CPU:    &parser.StringLiteral{Value: "100m"},
				Memory: &parser.StringLiteral{Value: "10Gi"},
				GPU:    &parser.IntLiteral{Value: 5},
			},
			RHS: &parser.IntLiteral{Value: 5},
		},
	}, expression)
}

func TestAddThenMulWithParen(t *testing.T) {
	p := newParserNoPositions(strings.NewReader(`
		(pod(cpu: "300m", memory: "32Gi", gpu: 10) + pod(cpu: "100m", memory: "10Gi", gpu: 5)) * 6
	`))

	expression := p.ParseExpression()

	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.ArithmeticExpression{
		Operator: parser.Multiply,
		LHS: &parser.ArithmeticExpression{
			Operator: parser.Add,
			LHS: &parser.PodExpression{
				CPU:    &parser.StringLiteral{Value: "300m"},
				Memory: &parser.StringLiteral{Value: "32Gi"},
				GPU:    &parser.IntLiteral{Value: 10},
			},
			RHS: &parser.PodExpression{
				CPU:    &parser.StringLiteral{Value: "100m"},
				Memory: &parser.StringLiteral{Value: "10Gi"},
				GPU:    &parser.IntLiteral{Value: 5},
			},
		},
		RHS: &parser.IntLiteral{Value: 6},
	}, expression)
}

func TestAddNumbers(t *testing.T) {
	p := newParserNoPositions(strings.NewReader(`
		pod(cpu: 5 + 10, memory: "32Gi", gpu: 10)
	`))

	p.ParseExpression()
	assert.NotEmpty(t, p.Errors)
}

func TestComplexExpression(t *testing.T) {
	p := newParserNoPositions(strings.NewReader(`
		# Global microservice 1
		pod(cpu: "250m", memory: "64Gi", gpu: 10) +

		# Global microservice 2
		pod(cpu: "2500m", memory: "31Gi", gpu: 980) +

		# Global microservice 3
		pod(cpu: "350m", memory: "564Gi", gpu: 455) +

		# Environment-specific microservices
		# We have 3 environments
		(
			# Environment microservice 1
			pod(cpu: "550m", memory: "64Gi", gpu: 10) +

			# Environment microservice 2 (with 3 replicas)
			pod(cpu: "650m", memory: "100Gi", gpu: 11) * 3 +

			# Environment microservice 3
			pod(cpu: "750m", memory: "64Gi", gpu: 10) +

			# In each environment, there are 60 instances of something with pod with 32 replicas + another pod
			(
				32 * pod(cpu: "2200m", memory: "32Gi", gpu: 10) + 
				pod(cpu: "100m", memory: "10Gi", gpu: 5)
			) * 6
		) * 3 + 

		# Global microservice 4
		pod(cpu: "2550m", memory: "64Gi", gpu: 10)
	`))

	expression := p.ParseExpression()

	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.ArithmeticExpression{
		Operator: parser.Add,

		LHS: &parser.ArithmeticExpression{
			Operator: parser.Add,
			LHS: &parser.ArithmeticExpression{
				Operator: parser.Add,
				LHS: &parser.ArithmeticExpression{
					Operator: parser.Add,

					// Global microservice 1
					// pod(cpu: "250m", memory: "64Gi", gpu: 10) +
					LHS: &parser.PodExpression{
						CPU:    &parser.StringLiteral{Value: "250m"},
						Memory: &parser.StringLiteral{Value: "64Gi"},
						GPU:    &parser.IntLiteral{Value: 10},
					},

					// Global microservice 2
					// pod(cpu: "2500m", memory: "31Gi", gpu: 980) +
					RHS: &parser.PodExpression{
						CPU:    &parser.StringLiteral{Value: "2500m"},
						Memory: &parser.StringLiteral{Value: "31Gi"},
						GPU:    &parser.IntLiteral{Value: 980},
					},
				},

				// Global microservice 3
				// pod(cpu: "350m", memory: "564Gi", gpu: 455) +
				RHS: &parser.PodExpression{
					CPU:    &parser.StringLiteral{Value: "350m"},
					Memory: &parser.StringLiteral{Value: "564Gi"},
					GPU:    &parser.IntLiteral{Value: 455},
				},
			},

			// Environment-specific microservices
			// We have 3 environments
			RHS: &parser.ArithmeticExpression{
				Operator: parser.Multiply,

				LHS: &parser.ArithmeticExpression{
					Operator: parser.Add,
					LHS: &parser.ArithmeticExpression{
						Operator: parser.Add,
						LHS: &parser.ArithmeticExpression{
							Operator: parser.Add,

							// Environment microservice 1
							// pod(cpu: "550m", memory: "64Gi", gpu: 10) +
							LHS: &parser.PodExpression{
								CPU:    &parser.StringLiteral{Value: "550m"},
								Memory: &parser.StringLiteral{Value: "64Gi"},
								GPU:    &parser.IntLiteral{Value: 10},
							},

							// Environment microservice 2 (with 3 replicas)
							// pod(cpu: "650m", memory: "100Gi", gpu: 11) * 3 +
							RHS: &parser.ArithmeticExpression{
								Operator: parser.Multiply,
								LHS: &parser.PodExpression{
									CPU:    &parser.StringLiteral{Value: "650m"},
									Memory: &parser.StringLiteral{Value: "100Gi"},
									GPU:    &parser.IntLiteral{Value: 11},
								},
								RHS: &parser.IntLiteral{Value: 3},
							},
						},

						// Environment microservice 3
						// pod(cpu: "750m", memory: "64Gi", gpu: 10)
						RHS: &parser.PodExpression{
							CPU:    &parser.StringLiteral{Value: "750m"},
							Memory: &parser.StringLiteral{Value: "64Gi"},
							GPU:    &parser.IntLiteral{Value: 10},
						},
					},

					// In each environment, there are 60 instances of something with pod with 32 replicas + another pod
					RHS: &parser.ArithmeticExpression{
						Operator: parser.Multiply,

						LHS: &parser.ArithmeticExpression{
							Operator: parser.Add,

							// 32 * pod(cpu: "2200m", memory: "32Gi", gpu: 10)
							LHS: &parser.ArithmeticExpression{
								Operator: parser.Multiply,
								LHS:      &parser.IntLiteral{Value: 32},
								RHS: &parser.PodExpression{
									CPU:    &parser.StringLiteral{Value: "2200m"},
									Memory: &parser.StringLiteral{Value: "32Gi"},
									GPU:    &parser.IntLiteral{Value: 10},
								},
							},

							// pod(cpu: "100m", memory: "10Gi", gpu: 5)
							RHS: &parser.PodExpression{
								CPU:    &parser.StringLiteral{Value: "100m"},
								Memory: &parser.StringLiteral{Value: "10Gi"},
								GPU:    &parser.IntLiteral{Value: 5},
							},
						},

						RHS: &parser.IntLiteral{Value: 6},
					},
				},

				RHS: &parser.IntLiteral{Value: 3},
			},
		},

		// Global microservice 4
		// pod(cpu: "2550m", memory: "64Gi", gpu: 10)
		RHS: &parser.PodExpression{
			CPU:    &parser.StringLiteral{Value: "2550m"},
			Memory: &parser.StringLiteral{Value: "64Gi"},
			GPU:    &parser.IntLiteral{Value: 10},
		},
	}, expression)
}

func newParserNoPositions(reader io.Reader) *parser.Parser {
	scanner := &lexer.Scanner{
		Reader:           bufio.NewReader(reader),
		DisablePositions: true,
	}

	return parser.NewParser(scanner)
}
