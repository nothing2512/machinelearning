package controllers

import (
	"encoding/json"
	"main/ann"
	"net/http"
)

func (idb *InDb) ANN(w http.ResponseWriter, r *http.Request) {
	var cmd CMD
	_ = json.NewDecoder(r.Body).Decode(&cmd)

	// inputs := []float64{0.1, 0.2, 0.1}
	// targets := []float64{0.2, 0.1, 0.2}
	// learningRate := 0.1
	// epoch := 2
	// showWeight := false

	ann := ann.NewNeuralNetwork(w, cmd.Ann.Inputs, cmd.Ann.Targets, cmd.Ann.LearningRate, cmd.Ann.ShowWeight)
	ann.Calculate(cmd.Ann.Epoch)
}
