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
	Command string `json:"command"`
	Data    string `json:"data"`
	Knn     struct {
		Data string `json:"data"`
		X    int    `json:"x"`
		Y    int    `json:"y"`
	} `json:"knn"`
}
