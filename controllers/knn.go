package controllers

import (
	"encoding/json"
	"main/knn"
	"net/http"
)

func (idb *InDb) KNN(w http.ResponseWriter, r *http.Request) {
	var cmd CMD
	_ = json.NewDecoder(r.Body).Decode(&cmd)
	idb.k = knn.Parse(cmd.Knn.Data)
	idb.k.Count(w, cmd.Knn.X, cmd.Knn.Y)
}
