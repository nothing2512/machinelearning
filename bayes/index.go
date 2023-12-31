package bayes

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"slices"
	"strings"
)

type NaiveBayes struct {
	loc           string
	class         map[string]int
	headers       []string
	body          [][]string
	classifier    string
	cKey          int
	tables        map[string]map[string]int
	tableDividers map[string]map[string]int
	built         bool
	Train         *NaiveBayes
	Test          *NaiveBayes
	classEntrophy float64
}

func (b *NaiveBayes) SplitTrainData(idx []int) {
	var trainBody [][]string
	var testBody [][]string
	for k, v := range b.body {
		isTrain := false
		for _, i := range idx {
			if i-1 == k {
				isTrain = true
				break
			}
		}
		if isTrain {
			trainBody = append(trainBody, v)
		} else {
			testBody = append(testBody, v)
		}
	}

	b.Train = &NaiveBayes{
		body:    trainBody,
		headers: b.headers,
		loc:     "train",
	}
	b.Test = &NaiveBayes{
		body:    testBody,
		headers: b.headers,
		loc:     "test",
	}
	b.Train.SetClassifier(b.classifier)
	b.Test.SetClassifier(b.classifier)

	b.Train.Build()
	b.Test.Build()
}

func (b *NaiveBayes) ClassifyLastHeader() {
	b.SetClassifier(b.headers[len(b.headers)-1])
}

func (b *NaiveBayes) SetClassifier(c string) {
	b.classifier = c
	for k, h := range b.headers {
		if h == c {
			b.cKey = k
			break
		}
	}
	b.class = make(map[string]int)
	b.tables = make(map[string]map[string]int)
	b.tableDividers = make(map[string]map[string]int)
	for _, row := range b.body {
		if _, exist := b.class[row[b.cKey]]; exist {
			b.class[row[b.cKey]] += 1
		} else {
			b.class[row[b.cKey]] = 1
		}
	}
}

func (b *NaiveBayes) ShowTables(w io.Writer) {
	if !b.built {
		b.Build()
	}
	var tables [][][]string
	for _, h := range b.headers {
		if h != b.classifier {
			tables = append(tables, b.ShowTable(h))
		}
	}
	_ = json.NewEncoder(w).Encode(tables)
}

func (b *NaiveBayes) ShowTable(model string) [][]string {
	if !b.built {
		b.Build()
	}
	var headers []string
	data := [][]string{}
	data = append(data, []string{"<b>" + model + "</b>"})
	for c := range b.class {
		data[0] = append(data[0], c)
		headers = append(headers, c)
	}
	for _, c := range headers {
		data[0] = append(data[0], fmt.Sprintf("P(%v)", c))
	}
	var rowAxis []string
	for x := range b.tables[model] {
		axis := strings.Split(x, "-")[0]
		has := false
		for _, a := range rowAxis {
			if a == axis {
				has = true
				break
			}
		}
		if !has {
			rowAxis = append(rowAxis, axis)
		}
	}
	for _, axis := range rowAxis {
		rowData := []string{axis}
		for _, c := range headers {
			rowData = append(rowData, fmt.Sprintf("%v", b.tables[model][fmt.Sprintf("%v-%v", axis, c)]))
		}
		for _, c := range headers {
			rowData = append(rowData, fmt.Sprintf("%v/%v", b.tables[model][fmt.Sprintf("%v-%v", axis, c)], b.tableDividers[model][c]))
		}
		data = append(data, rowData)
	}
	footer := []string{"Total"}
	for _, c := range headers {
		footer = append(footer, fmt.Sprintf("%v", b.tableDividers[model][c]))
	}
	for range headers {
		footer = append(footer, "100%")
	}

	return data
}

func (b *NaiveBayes) Build() {
	if b.built {
		return
	}

	if b.loc == "" {
		b.loc = "Main"
	}

	for _, h := range b.headers {
		if h != b.classifier {
			b.build(h)
		}
	}

	b.built = true
}

func (b *NaiveBayes) build(model string) {

	// get model Key
	mKey := 0
	for k, h := range b.headers {
		if h == model {
			mKey = k
			break
		}
	}

	// get model values
	mValues := make(map[string]bool)
	for _, row := range b.body {
		mValues[row[mKey]] = true
	}

	// get keys
	singleTable := make(map[string]int)
	singleDivider := make(map[string]int)
	for c := range b.class {
		singleDivider[c] = 0
		for m := range mValues {
			singleTable[fmt.Sprintf("%v-%v", m, c)] = 0
		}
	}

	// get Data
	for _, row := range b.body {
		m := row[mKey]
		c := row[b.cKey]
		singleTable[fmt.Sprintf("%v-%v", m, c)] += 1
		singleDivider[c] += 1
	}

	b.tables[model] = singleTable
	b.tableDividers[model] = singleDivider
}

func (b *NaiveBayes) PredictAppend(w io.Writer, data []string) {
	if !b.built {
		b.Build()
	}
	c := b.Predict(w, data)
	for _k, header := range b.headers {
		k := _k
		if _k == b.cKey {
			continue
		}
		if _k > b.cKey {
			k -= 1
		}
		b.tables[header][fmt.Sprintf("%v-%v", data[k], c)] += 1
		b.tableDividers[header][c] += 1
	}
}

func (b *NaiveBayes) Predict(w io.Writer, data []string) string {
	if !b.built {
		b.Build()
	}
	predictions := make(map[string]float64)
	for c := range b.class {
		fmt.Fprint(w, "\n")
		predictions[c] = b.predict(w, c, data)
		fmt.Fprint(w, "\n")
	}
	concultionC := ""
	var concultionV float64 = 0
	for k, v := range predictions {
		if concultionV < v {
			concultionV = v
			concultionC = k
		}
	}
	fmt.Fprint(w, fmt.Sprintf("\nConclution : %v : %.3f", concultionC, concultionV))
	return concultionC
}

func (b *NaiveBayes) predict(w io.Writer, c string, data []string) float64 {
	fmt.Fprint(w, "<b>"+c+"</b>\n")
	steps := [][]string{{}, {}}
	values := 1
	dividers := 1
	for headerKey, header := range b.headers {
		k := headerKey
		if headerKey == b.cKey {
			continue
		}
		if headerKey > b.cKey {
			k -= 1
		}
		steps[0] = append(steps[0], fmt.Sprintf("P(%v|%v)", data[k], c))
		divKey := fmt.Sprintf("%v-%v", data[k], c)
		value := b.tables[header][divKey]
		divider := b.tableDividers[header][c]
		steps[1] = append(steps[1], fmt.Sprintf("<sup>%v</sup>/<sub>%v</sub>", value, divider))
		values *= value
		dividers *= divider
	}

	cDiv := 0
	for _, v := range b.class {
		cDiv += v
	}
	steps[0] = append(steps[0], fmt.Sprintf("P(%v)", c))
	steps[1] = append(steps[1], fmt.Sprintf("<sup>%v</sup>/<sub>%v</sub>", b.class[c], cDiv))

	values *= b.class[c]
	dividers *= cDiv

	fmt.Fprint(w, " = ", strings.Join(steps[0], " * "), "\n")
	fmt.Fprint(w, " = ", strings.Join(steps[1], " * "), "\n")
	fmt.Fprint(w, " = ", fmt.Sprintf("%v / %v", values, dividers), "\n")
	fmt.Fprint(w, " = ", fmt.Sprintf("%.3f", float64(values)/float64(dividers)), "\n")
	fmt.Fprint(w, "\n\n")
	if values == 0 {
		return 0
	}
	return float64(values) / float64(dividers)
}

func (b *NaiveBayes) Smooth(num int) {
	if !b.built {
		b.Build()
	}
	for k := range b.class {
		b.class[k] += num
	}
	for tk, v := range b.tableDividers {
		for k := range v {
			b.tableDividers[tk][k] += num
		}
	}
	for tk, v := range b.tables {
		for k := range v {
			b.tables[tk][k] += num
		}
	}
}

func (b *NaiveBayes) ShowEntrophy(w io.Writer, model string) (ent float64, sVal int, sDiv int) {
	if !b.built {
		b.Build()
	}
	fmt.Fprint(w, "\t", model)

	steps := [][]string{{}, {}, {}}
	var values float64
	var div int
	for k, v := range b.tables[model] {
		fmt.Fprint(w, k, v)
		div += v
	}
	for k, v := range b.tables[model] {
		vDiv := float64(v) / float64(div)
		val := (vDiv) * math.Log2(vDiv)
		values += val
		steps[0] = append(steps[0], fmt.Sprintf("Ent(%v)", strings.ReplaceAll(k, "-", "|")))
		steps[1] = append(steps[1], fmt.Sprintf("(%v/%v)log2(%v/%v)", v, div, v, div))
		steps[2] = append(steps[2], fmt.Sprintf("%v", val))
	}
	fmt.Fprint(w, fmt.Sprintf("Ent(%v) = - (%v)", model, strings.Join(steps[0], " + ")))
	fmt.Fprint(w, fmt.Sprintf("Ent(%v) = - (%v)", model, strings.Join(steps[1], " + ")))
	fmt.Fprint(w, fmt.Sprintf("Ent(%v) = - (%v)", model, strings.Join(steps[2], " + ")))
	fmt.Fprint(w, fmt.Sprintf("Ent(%v) = %v", model, -values))

	return -values, 0, div
}

func (b *NaiveBayes) ShowClassEntrophy(w io.Writer) {
	if !b.built {
		b.Build()
	}
	fmt.Fprint(w, "<b>Classifier</b>\n")
	fmt.Fprint(w, "\n")

	steps := [][]string{{}, {}}
	var values float64
	var div int
	for _, v := range b.class {
		div += v
	}
	for k, v := range b.class {
		vDiv := float64(v) / float64(div)
		val := (vDiv) * math.Log2(vDiv)
		values += val
		steps[0] = append(steps[0], fmt.Sprintf("Ent(%v)", strings.ReplaceAll(k, "-", "|")))
		steps[1] = append(steps[1], fmt.Sprintf("%v", val))
	}
	fmt.Fprint(w, fmt.Sprintf("Ent(D) = - (%v)\n", strings.Join(steps[0], " + ")))
	fmt.Fprint(w, fmt.Sprintf("Ent(D) = - (%v)\n", strings.Join(steps[1], " + ")))
	fmt.Fprint(w, fmt.Sprintf("Ent(D) = %.3f\n", -values))
	fmt.Fprintf(w, "\n")

	b.classEntrophy = -values
}

func (b *NaiveBayes) ShowGains(w io.Writer) (col string, vals []string) {
	if !b.built {
		b.Build()
	}

	result := make(map[string]struct {
		gain      float64
		gainRatio float64
		iv        float64
	})

	totalData := len(b.body)

	for tk, v := range b.tables {
		fmt.Fprint(w, "\n")
		fmt.Fprint(w, "<b>", tk, "</b>\n") // color

		dividers := make(map[string]int)
		for k, _v := range v {
			if _, exist := dividers[strings.Split(k, "-")[0]]; !exist {
				dividers[strings.Split(k, "-")[0]] = _v
			} else {
				dividers[strings.Split(k, "-")[0]] += _v
			}
		}

		var gainSteps []string
		var ivSteps []string
		gainSub := 0.0
		iv := 0.0

		// get model entropies
		for k, d := range dividers { // k=dark, d=/10
			value := 0.0
			var steps []string
			for c := range b.class { // c=true, cVal=10
				toDivide := float64(v[k+"-"+c])
				steps = append(steps, fmt.Sprintf("(<sup>%v</sup>/<sub>%v</sub>)log₂(<sup>%v</sup>/<sub>%v</sub>)", toDivide, d, toDivide, d))
				val := toDivide / float64(d)
				value += val * math.Log2(val)
			}
			value = -value
			if math.IsNaN(value) {
				value = 0
			}
			fmt.Fprint(w, fmt.Sprintf("Ent(%v) = -(%v) = %.3f", k, strings.Join(steps, " + "), value))
			gainSteps = append(gainSteps, fmt.Sprintf("%.3f(<sup>%v</sup>/<sub>%v</sub>)", value, d, totalData))
			ivSteps = append(steps, fmt.Sprintf("(<sup>%v</sup>/<sub>%v</sub>)log₂(<sup>%v</sup>/<sub>%v</sub>)", d, totalData, d, totalData))
			gainSub += float64(d) * value / float64(totalData)
			iv += (float64(d) / float64(totalData)) * math.Log2(float64(d)/float64(totalData))
			fmt.Fprint(w, "\n")
		}

		iv = -iv
		gain := b.classEntrophy - gainSub
		fmt.Fprint(w, fmt.Sprintf("Gain(D, %v) = %.3f - (%v) = %.3f\n", tk, b.classEntrophy, strings.Join(gainSteps, " + "), gain))
		fmt.Fprint(w, fmt.Sprintf("IV(%v) = -(%v) = %.3f\n", tk, strings.Join(ivSteps, " + "), iv))
		fmt.Fprint(w, fmt.Sprintf("GainRatio(%v) = Gain(D, %v) / IV(%v) = %.3f / %.3f = %.3f\n", tk, tk, tk, gain, iv, gain/iv))
		fmt.Fprint(w, "\n")

		result[tk] = struct {
			gain      float64
			gainRatio float64
			iv        float64
		}{gain, gain / iv, iv}
	}

	fmt.Fprint(w, "\n")
	rootValue := 0.0
	root := ""

	var _headers []string
	for _, h := range b.headers {
		_headers = append(_headers, h)
	}
	slices.Reverse(_headers)

	fmt.Fprint(w, "<b>Gains</b>\n")
	for _, h := range _headers {
		if h == b.classifier {
			continue
		}
		v := result[h]
		if v.gain >= rootValue {
			rootValue = v.gain
			root = h
		}
		fmt.Fprint(w, fmt.Sprintf("k(%v) gr(%.3f) g(%.3f) iv(%.3f)\n", h, v.gainRatio, v.gain, v.iv))
	}

	fmt.Fprint(w, "\n")
	fmt.Fprint(w, "<b>Root</b> : ", root)
	col = root

	for kRoot := range b.tables[root] {
		k := strings.Split(kRoot, "-")[0]
		has := false
		for _, v := range vals {
			if v == k {
				has = true
				break
			}
		}
		if !has {
			vals = append(vals, k)
		}
	}

	return
}

func (b *NaiveBayes) GetClassEntrophy() {
	if !b.built {
		b.Build()
	}

	var values float64
	var div int
	for _, v := range b.class {
		div += v
	}
	for _, v := range b.class {
		vDiv := float64(v) / float64(div)
		val := (vDiv) * math.Log2(vDiv)
		values += val
	}

	b.classEntrophy = -values
}

func (b *NaiveBayes) GetRoot() (col string) {
	if !b.built {
		b.Build()
	}

	b.GetClassEntrophy()

	result := make(map[string]struct {
		gain      float64
		gainRatio float64
		iv        float64
	})

	totalData := len(b.body)

	for tk, v := range b.tables {

		dividers := make(map[string]int)
		for k, _v := range v {
			if _, exist := dividers[strings.Split(k, "-")[0]]; !exist {
				dividers[strings.Split(k, "-")[0]] = _v
			} else {
				dividers[strings.Split(k, "-")[0]] += _v
			}
		}

		gainSub := 0.0
		iv := 0.0

		// get model entropies
		for k, d := range dividers { // k=dark, d=/10
			value := 0.0
			var steps []string
			for c := range b.class { // c=true, cVal=10
				toDivide := float64(v[k+"-"+c])
				steps = append(steps, fmt.Sprintf("(<sup>%v</sup>/<sub>%v</sub>)log₂(<sup>%v</sup>/<sub>%v</sub>)", toDivide, d, toDivide, d))
				val := toDivide / float64(d)
				value += val * math.Log2(val)
			}
			value = -value
			if math.IsNaN(value) {
				value = 0
			}
			gainSub += float64(d) * value / float64(totalData)
			iv += (float64(d) / float64(totalData)) * math.Log2(float64(d)/float64(totalData))
		}

		iv = -iv
		gain := b.classEntrophy - gainSub

		result[tk] = struct {
			gain      float64
			gainRatio float64
			iv        float64
		}{gain, gain / iv, iv}
	}

	rootValue := 0.0
	root := ""

	var _headers []string
	for _, h := range b.headers {
		_headers = append(_headers, h)
	}
	slices.Reverse(_headers)

	for _, h := range _headers {
		if h == b.classifier {
			continue
		}
		v := result[h]
		if v.gain >= rootValue {
			rootValue = v.gain
			root = h
		}
	}

	col = root

	return
}

func (b *NaiveBayes) Filter(col, val string) *NaiveBayes {
	k := 0
	for i, h := range b.headers {
		if h == col {
			k = i
			break
		}
	}
	var body [][]string
	for _, bd := range b.body {
		if bd[k] == val {
			body = append(body, bd)
		}
	}
	nb := &NaiveBayes{
		body:    body,
		headers: b.headers,
		loc:     "main",
	}
	nb.SetClassifier(b.classifier)
	return nb
}

type Json map[string]interface{}

func (b *NaiveBayes) DecisionTree(w io.Writer, show bool) Json {
	if !b.built {
		b.Build()
	}
	models := make(map[string][]string)
	for k, t := range b.tables {
		if _, exist := models[k]; !exist {
			models[k] = []string{}
		}
		for m := range t {
			tm := strings.Split(m, "-")[0]
			has := false
			for _, v := range models[k] {
				if v == tm {
					has = true
				}
			}
			if !has {
				models[k] = append(models[k], tm)
			}
		}
	}

	tree := b.getTree(models, b)

	if show {
		by, _ := json.MarshalIndent(tree, "", "    ")
		fmt.Fprint(w, string(by))
	}

	return tree
}

func (*NaiveBayes) getTree(models map[string][]string, b *NaiveBayes) Json {
	root := b.GetRoot()
	data := Json{}
	for _, branch := range models[root] {
		branchData := b.Filter(root, branch)

		if len(branchData.body) == 0 {
			_c := ""
			for c := range b.class {
				_c = fmt.Sprintf("%v=%v", b.classifier, c)
				break
			}
			data[branch] = _c
			continue
		}

		bodySame := true
		for k, v := range branchData.body {
			for kk := range v {
				breaked := false
				if b.body[k][kk] != v[kk] {
					bodySame = false
					breaked = true
					break
				}
				if breaked {
					break
				}
			}
		}
		if bodySame {
			var bodySameData []string
			for c := range b.class {
				bodySameData = append(bodySameData, c)
			}
			data[branch] = fmt.Sprintf("%v=%v", b.classifier, strings.Join(bodySameData, "|"))
			continue
		}

		cData := branchData.body[0][b.cKey]
		isSame := true
		for _, row := range branchData.body {
			if row[b.cKey] != cData {
				isSame = false
			}
		}
		if isSame {
			data[branch] = fmt.Sprintf("%v=%v", b.classifier, cData)
			continue
		}

		data[branch] = b.getTree(models, branchData)
	}
	return Json{fmt.Sprintf("%v=?", root): data}
}

func (b *NaiveBayes) PrintRaw(w io.Writer) {
	fmt.Fprint(w, strings.Join(b.headers, ","), "\n")
	for _, body := range b.body {
		fmt.Fprint(w, strings.Join(body, ","), "\n")
	}
}
