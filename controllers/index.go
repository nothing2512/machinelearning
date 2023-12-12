package controllers

import (
	"main/bayes"
	"main/knn"
)

type InDb struct {
	nb *bayes.NaiveBayes
	k  *knn.KNN
}

type CMD struct {
	Context string `json:"context"`
	Command string `json:"command"`
	Data    string `json:"data"`
	Knn     struct {
		Data string `json:"data"`
		X    int    `json:"x"`
		Y    int    `json:"y"`
	} `json:"knn"`
	Ann struct {
		Inputs       []float64 `json:"inputs"`
		Targets      []float64 `json:"targets"`
		LearningRate float64   `json:"learningRate"`
		Epoch        int       `json:"epoch"`
		ShowWeight   bool      `json:"showWeight"`
	}
}
