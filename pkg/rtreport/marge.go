package rtreport

func margeSingleMap(a, b map[string]float64) {
	for key, value := range b {
		if _, exist := a[key]; !exist {
			a[key] = value
		} else {
			a[key] = a[key] + value
		}
	}
}

func margeDoubleMap(a, b map[string]map[string]float64) {
	for key, value := range b {
		if _, exist := a[key]; !exist {
			a[key] = value
		} else {
			margeSingleMap(a[key], b[key])
		}
	}
}

func margeTripleMap(a, b map[string]map[string]map[string]float64) {
	for key, value := range b {
		if _, exist := a[key]; !exist {
			a[key] = value
		} else {
			margeDoubleMap(a[key], b[key])
		}
	}
}

func margeQuadrupleMap(a, b map[string]map[string]map[string]map[string]float64) {
	for key, value := range b {
		if _, exist := a[key]; !exist {
			a[key] = value
		} else {
			margeTripleMap(a[key], b[key])
		}
	}
}
