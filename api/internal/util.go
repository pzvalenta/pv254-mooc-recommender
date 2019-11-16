package internal

//Min ...
func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func intersection(a, b []string) []string {
	// interacting on the smallest list first can potentially be faster...but not by much, worse case is the same
	low, high := a, b
	if len(a) > len(b) {
		low = b
		high = a
	}

	var inter []string
	done := false
	for i, l := range low {
		for j, h := range high {
			// get future index values
			f1 := i + 1
			f2 := j + 1
			if l == h {
				inter = append(inter, h)
				if f1 < len(low) && f2 < len(high) {
					// if the future values aren't the same then that's the end of the intersection
					if low[f1] != high[f2] {
						done = true
					}
				}
				// we don't want to interate on the entire list everytime, so remove the parts we already looped on will make it faster each pass
				high = high[:j+copy(high[j:], high[j+1:])]
				break
			}
		}
		// nothing in the future so we are done
		if done {
			break
		}
	}
	return inter
}

func wordCount(a []string) *map[string]float64 {
	m1 := make(map[string]float64)
	for _, word := range a {
		if val, ok := m1[word]; ok {
			m1[word] = val + 1
		} else {
			m1[word] = 1
		}

	}

	return &m1
}
