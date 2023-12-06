package controllers

import (
	"encoding/json"
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
		switch cmd.Context {
		case "train":
			idb.nb.Train.ShowTables(w)
			break
		case "test":
			idb.nb.Test.ShowTables(w)
			break
		default:
			idb.nb.ShowTables(w)
		}
	}

	if cmd.Command == "smooth" {
		num, _ := strconv.Atoi(cmd.Data)
		switch cmd.Context {
		case "train":
			idb.nb.Train.Smooth(num)
			break
		case "test":
			idb.nb.Test.Smooth(num)
			break
		default:
			idb.nb.Smooth(num)
		}
	}
	if cmd.Command == "tree" {
		switch cmd.Context {
		case "train":
			idb.nb.Train.DecisionTree(w, true)
			break
		case "test":
			idb.nb.Test.DecisionTree(w, true)
			break
		default:
			idb.nb.DecisionTree(w, true)
		}
	}
	if cmd.Command == "root" {
		switch cmd.Context {
		case "train":
			idb.nb.Train.ShowClassEntrophy(w)
			root, _ := idb.nb.Train.ShowGains(w)
			_, _ = w.Write([]byte(root))
			break
		case "test":
			idb.nb.Test.ShowClassEntrophy(w)
			root, _ := idb.nb.Test.ShowGains(w)
			_, _ = w.Write([]byte(root))
			break
		default:
			idb.nb.ShowClassEntrophy(w)
			root, _ := idb.nb.ShowGains(w)
			_, _ = w.Write([]byte(root))
		}
	}
	if cmd.Command == "predict" {
		switch cmd.Context {
		case "train":
			idb.nb.Train.Predict(w, strings.Split(cmd.Data, ","))
			break
		case "test":
			idb.nb.Test.Predict(w, strings.Split(cmd.Data, ","))
			break
		default:
			idb.nb.Predict(w, strings.Split(cmd.Data, ","))
		}
	}
	if cmd.Command == "append" {
		switch cmd.Context {
		case "train":
			idb.nb.Train.PredictAppend(w, strings.Split(cmd.Data, ","))
			break
		case "test":
			idb.nb.Test.PredictAppend(w, strings.Split(cmd.Data, ","))
			break
		default:
			idb.nb.PredictAppend(w, strings.Split(cmd.Data, ","))
		}
	}
	if cmd.Command == "splitTrainData" {
		var idx []int
		for _, x := range strings.Split(cmd.Data, ",") {
			_x, _ := strconv.Atoi(x)
			idx = append(idx, _x)
		}
		idb.nb.SplitTrainData(idx)
	}
	if cmd.Command == "raw" {
		switch cmd.Context {
		case "train":
			idb.nb.Train.PrintRaw(w)
			break
		case "test":
			idb.nb.Test.PrintRaw(w)
			break
		default:
			idb.nb.PrintRaw(w)
			break
		}
	}
}
