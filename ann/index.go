package ann

import (
	"fmt"
	"io"
	"math"
	"strings"
)

type NeuralNetwork struct {
	InputSize    int
	OutputSize   int
	Weights      [][]float64
	Bias         []float64
	LearningRate float64
	Input        []float64
	Target       []float64
	w            io.Writer
	showWeight   bool
}

// NewNeuralNetwork initializes a new neural network.
func NewNeuralNetwork(w io.Writer, inputs, targets []float64, learningRate float64, showWeight bool) *NeuralNetwork {
	inputSize := len(inputs)
	outputSize := len(targets)

	weights := make([][]float64, outputSize)
	for i := range weights {
		weights[i] = make([]float64, inputSize)
	}
	bias := make([]float64, outputSize)

	return &NeuralNetwork{
		InputSize:    inputSize,
		OutputSize:   outputSize,
		Weights:      weights,
		Bias:         bias,
		LearningRate: learningRate,
		Input:        inputs,
		Target:       targets,
		w:            w,
		showWeight:   showWeight,
	}
}

// Sigmoid activation function.
func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

// ForwardPass performs the forward pass of the neural network.
func (nn *NeuralNetwork) ForwardPass(inputs []float64) []float64 {
	outputs := make([]float64, nn.OutputSize)

	for i := 0; i < nn.OutputSize; i++ {
		// Calculate the weighted sum of inputs
		z := nn.Bias[i]
		var steps []string
		for j := 0; j < nn.InputSize; j++ {
			steps = append(steps, fmt.Sprintf("(w<sub>%v%v</sub> * I%v)", i+1, j+1, j+1))
			z += nn.Weights[i][j] * inputs[j]
		}

		fmt.Fprintf(nn.w, fmt.Sprintf("X%v = %v<br>", i+1, strings.Join(steps, " + ")))
		fmt.Fprintf(nn.w, fmt.Sprintf("X%v = %v<br>", i+1, z))
		fmt.Fprintf(nn.w, fmt.Sprintf("(sigmoid activation) X%v = 1 / ( 1 + e<sup>%v</sup> )<br>", i+1, z))
		// Apply activation function (sigmoid)
		outputs[i] = sigmoid(z)
		fmt.Fprintf(nn.w, fmt.Sprintf("(sigmoid activation) X%v = %v<br>", i+1, outputs[i]))
		fmt.Fprintf(nn.w, "<br>")
	}

	var outputString []string
	for _, x := range outputs {
		outputString = append(outputString, fmt.Sprintf("%v", x))
	}

	fmt.Fprintf(nn.w, fmt.Sprintf("<b>Outputs</b> : [%v]<br>", strings.Join(outputString, ", ")))
	fmt.Fprintf(nn.w, "<br>")

	return outputs
}

// BackwardPass performs the backward pass and updates weights and bias.
func (nn *NeuralNetwork) BackwardPass(inputs, targets []float64) {

	for i := 0; i < nn.OutputSize; i++ {
		// Calculate the error
		error := targets[i] - nn.Bias[i]
		fmt.Fprintf(nn.w, fmt.Sprintf("E%v = T%v - B%v<br>", i+1, i+1, i+1))
		fmt.Fprintf(nn.w, fmt.Sprintf("E%v = %v - %.3f<br>", i+1, targets[i], nn.Bias[i]))
		fmt.Fprintf(nn.w, fmt.Sprintf("E%v = %v<br>", i+1, error))
		fmt.Fprintf(nn.w, "<br>")

		// Update bias
		nn.Bias[i] += nn.LearningRate * error
		fmt.Fprintf(nn.w, fmt.Sprintf("B%v = learningRate * E%v<br>", i+1, i+1))
		fmt.Fprintf(nn.w, fmt.Sprintf("B%v = %v * %v<br>", i+1, nn.LearningRate, error))
		fmt.Fprintf(nn.w, fmt.Sprintf("B%v = %.3f<br>", i+1, nn.Bias[i]))
		fmt.Fprintf(nn.w, "<br>")

		for j := 0; j < nn.InputSize; j++ {
			if nn.showWeight {
				fmt.Fprintf(nn.w, fmt.Sprintf("W<sub>%v%v</sub> = W<sub>%v%v</sub> + (learningRate * E%v * I%v)<br>", i+1, j+1, i+1, j+1, i+1, j+1))
				fmt.Fprintf(nn.w, fmt.Sprintf("W<sub>%v%v</sub> = %.3f + (%v * %v * %v)<br>", i+1, j+1, nn.Weights[i][j], nn.LearningRate, error, inputs[j]))
			}

			nn.Weights[i][j] += nn.LearningRate * error * inputs[j]

			if nn.showWeight {
				fmt.Fprintf(nn.w, fmt.Sprintf("W<sub>%v%v</sub> = %.3f<br>", i+1, j+1, nn.Weights[i][j]))
				fmt.Fprintf(nn.w, "<br>")
			}
		}
	}
}

func (nn *NeuralNetwork) CalculateError(outputs []float64) float64 {
	var sumError float64
	var steps []string
	var steps2 []string
	for i := range nn.Target {
		error := nn.Target[i] - outputs[i]
		sumError += error * error
		steps = append(steps, fmt.Sprintf("(T%v - Z%v)<sup>2</sup>", i+1, i+1))
		steps2 = append(steps, fmt.Sprintf("(%v - %v)<sup>2</sup>", nn.Target[i], outputs[i]))
	}
	fmt.Fprint(nn.w, fmt.Sprintf("E = (%v) / total target\n", strings.Join(steps, " + ")))
	fmt.Fprint(nn.w, fmt.Sprintf("E = (%v) / %v\n", strings.Join(steps2, " + "), len(nn.Target)))
	fmt.Fprint(nn.w, fmt.Sprintf("E = %v / %v\n", sumError, len(nn.Target)))
	fmt.Fprint(nn.w, fmt.Sprintf("<b>E</b> = %v\n", sumError/float64(len(nn.Target))))
	return sumError / float64(len(nn.Target))
}

func (nn *NeuralNetwork) Calculate(epoch int) {
	inputs := nn.Input
	targets := nn.Target
	for _epoch := 0; _epoch < epoch; _epoch++ {
		fmt.Fprintf(nn.w, fmt.Sprintf("<b>Iteration #%v</b><br>", _epoch+1))
		fmt.Fprintf(nn.w, "<br>")
		fmt.Fprintf(nn.w, fmt.Sprintf("Forward Pass <b>#%v</b><br>", _epoch+1))
		fmt.Fprintf(nn.w, "<br>")
		outputs := nn.ForwardPass(inputs)

		fmt.Fprintf(nn.w, fmt.Sprintf("Backward Pass <b>#%v</b><br>", _epoch+1))
		fmt.Fprintf(nn.w, "<br>")
		nn.BackwardPass(inputs, targets)
		nn.CalculateError(outputs)
	}
}

func (nn *NeuralNetwork) setBias(bias float64) {
	for k := range nn.Bias {
		nn.Bias[k] = bias
	}
}

// func main() {
// 	// Training data (input and target)
// 	inputs := []float64{0.1, 0.2, 0.1}
// 	targets := []float64{0.2, 0.1, 0.2}

// 	nn := NewNeuralNetwork(inputs, targets, 0.1)
// 	nn.Calculate(2)
// }
