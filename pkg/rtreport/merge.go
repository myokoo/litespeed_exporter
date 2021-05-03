package rtreport

func mergeSingleMap(a, b map[string]float64) {
	for key, value := range b {
		if _, exist := a[key]; !exist {
			a[key] = value
		} else {
			a[key] = a[key] + value
		}
	}
}

func mergeDoubleMap(a, b map[string]map[string]float64) {
	for key, value := range b {
		if _, exist := a[key]; !exist {
			a[key] = value
		} else {
			mergeSingleMap(a[key], b[key])
		}
	}
}

func mergeTripleMap(a, b map[string]map[string]map[string]float64) {
	for key, value := range b {
		if _, exist := a[key]; !exist {
			a[key] = value
		} else {
			mergeDoubleMap(a[key], b[key])
		}
	}
}

func mergeQuadrupleMap(a, b map[string]map[string]map[string]map[string]float64) {
	for key, value := range b {
		if _, exist := a[key]; !exist {
			a[key] = value
		} else {
			mergeTripleMap(a[key], b[key])
		}
	}
}
