package knn

import (
	"strconv"
	"strings"
)

func Parse(data string) *KNN {
	rows := strings.Split(data, "\n")
	var x []int
	var y []int
	var c []string

	for _, a := range strings.Split(rows[0], ",") {
		_x, _ := strconv.Atoi(a)
		x = append(x, _x)
	}

	for _, a := range strings.Split(rows[1], ",") {
		_y, _ := strconv.Atoi(a)
		y = append(y, _y)
	}

	for _, a := range strings.Split(rows[2], ",") {
		c = append(c, a)
	}

	return &KNN{x, y, c}
}
