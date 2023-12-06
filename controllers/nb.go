package controllers

import (
	"encoding/json"
	"fmt"
	"main/bayes"
	"net/http"
	"strconv"
	"strings"
)

func (idb *InDb) NB(w http.ResponseWriter, r *http.Request) {
	var cmd CMD
	_ = json.NewDecoder(r.Body).Decode(&cmd)
	if cmd.Command == "parse" {
		idb.nb = bayes.ParseString(cmd.Data)
		idb.nb.ClassifyLastHeader()
	}

	if cmd.Command == "showTables" {
		idb.nb.ShowTables(w)
	}

	if cmd.Command == "smooth" {
		num, _ := strconv.Atoi(cmd.Data)
		idb.nb.Smooth(num)
		fmt.Println(num)
	}
	if cmd.Command == "tree" {
		idb.nb.DecisionTree(w, true)
	}
	if cmd.Command == "root" {
		idb.nb.ShowClassEntrophy(w)
		root, _ := idb.nb.ShowGains(w)
		_, _ = w.Write([]byte(root))
	}
	if cmd.Command == "predict" {
		idb.nb.Predict(w, strings.Split(cmd.Data, ","))
	}
	if cmd.Command == "append" {
		idb.nb.PredictAppend(w, strings.Split(cmd.Data, ","))
	}
}
