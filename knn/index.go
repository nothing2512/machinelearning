package knn

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

const ROOT = "√"
const SQUARE = "²"

type KNN struct {
	X     []int
	Y     []int
	class []string
}

func (*KNN) subscript(x int) string {
	data := []string{
		"\u2080",
		"\u2081",
		"\u2082",
		"\u2083",
		"\u2084",
		"\u2085",
		"\u2086",
		"\u2087",
		"\u2088",
		"\u2089",
	}
	_data := strings.Split(fmt.Sprintf("%v", x), "")
	res := ""
	for _, a := range _data {
		_x, _ := strconv.Atoi(string(a))
		res += data[_x]
	}
	return res
}

func (k *KNN) Count(w io.Writer, x, y int) {
	var steps []string
	var sorted []string
	var values []float64
	steps = append(steps, "\nCalculate\n")
	for i := range k.X {
		val := math.Pow((float64(x-k.X[i])), 2) + math.Pow((float64(y-k.Y[i])), 2)
		step := fmt.Sprintf("ex%v = %v(%v-%v)%v + (%v-%v)%v = %v%v (%v)\n",
			i+1,
			ROOT,
			x, k.X[i], SQUARE,
			y, k.Y[i], SQUARE,
			ROOT,
			val,
			k.class[i],
		)

		steps = append(steps, step)
		values = append(values, val)
		sorted = append(sorted, step)
	}
	steps = append(steps, "\nSorted\n")
	for i := 0; i < len(values)-1; i++ {
		for j := 0; j < len(values)-i-1; j++ {
			if values[j] > values[j+1] {
				values[j], values[j+1] = values[j+1], values[j]
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}
	steps = append(steps, sorted...)
	for _, s := range steps {
		fmt.Fprint(w, s)
	}
}
