package internal

import (
	"math"
	"regexp"
	"strings"
)

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

func wordCount(a []string) *map[string]int {
	m1 := make(map[string]int)
	for _, word := range a {
		if val, ok := m1[word]; ok {
			m1[word] = val + 1
		} else {
			m1[word] = 1
		}
	}
	return &m1
}

func getStopWords() map[string]string {
	return map[string]string{" a ": " ", " and ": " ", " the ": " ", " of ": " ", " is ": " ", " are ": " ",
		" in ": " ", " to ": " ", " from ": " ", " on ": " ", ".": "", ":": "", "(": " ",
		")": " ", "\n": " ", ",": " ", "  ": " "}
}

func tokenize(text string) []string {
	text = normalize(text)
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	tks := strings.Split(text, " ")
	var cleanToks []string
	for _, val := range tks {
		if val != "" {
			cleanVal := reg.ReplaceAllString(val, "")
			cleanToks = append(cleanToks, cleanVal)
		}
	}
	return cleanToks
}

func normalize(text string) string {
	text = strings.ToLower(text)
	for word, newWord := range getStopWords() {
		text = strings.Replace(text, word, newWord, -1)
	}
	return text
}

func computeIdf(texts []string) map[string]float64 {
	wordCounts := make(map[string]float64)

	N := float64(len(texts))
	for _, text := range texts {
		toks := tokenize(text)
		words := make(map[string]bool)
		for _, word := range toks {
			if _, ok := words[word]; !ok {
				words[word] = true
			}
		}
		for word := range words {
			if _, ok := wordCounts[word]; ok {
				wordCounts[word] = wordCounts[word] + 1
			} else {
				wordCounts[word] = 1
			}
		}
	}
	wordIdf := make(map[string]float64)

	for word, count := range wordCounts {
		wordIdf[word] = math.Log(N/ count)
	}
	return wordIdf

}

func computeTf(text string) *map[string]float64 {
	tf := make(map[string]float64)
	tokens := tokenize(text)
	N := float64(len(tokens))
	wordCounts := wordCount(tokens)
	for word, count := range *wordCounts {
		tf[word] = float64(count) / N
	}

	return &tf
}
