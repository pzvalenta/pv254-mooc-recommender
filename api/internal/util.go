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
	return map[string]string{"&": " ", "  ": " ", "-": " ", ".": " ", ":": " ", "(": " ",
		")": " ", "\n": " ", ",": " ", " a ": " ", " about ": " ", " above ": " ",
		" across ": " ", " after ": " ", " afterwards ": " ", " again ": " ", " against ": " ",
		" all ": " ", " almost ": " ", " alone ": " ", " along ": " ", " already ": " ",
		" also ": " ", " although ": " ", " always ": " ", " am ": " ", " among ": " ",
		" amongst ": " ", " amoungst ": " ", " amount ": " ", " an ": " ", " and ": " ",
		" another ": " ", " any ": " ", " anyhow ": " ", " anyone ": " ", " anything ": " ",
		" anyway ": " ", " anywhere ": " ", " are ": " ", " they're": " ", " they’re": " ", " around ": " ", " as ": " ", " at ": " ",
		" back ": " ", " be ": " ", " became ": " ", " because ": " ", " become ": " ", " becomes ": " ",
		" becoming ": " ", " been ": " ", " before ": " ", " beforehand ": " ", " behind ": " ", " being ": " ",
		" below ": " ", " beside ": " ", " besides ": " ", " between ": " ", " beyond ": " ", " bill ": " ",
		" both ": " ", " bottom ": " ", " but ": " ", " by ": " ", " call ": " ", " can ": " ", " cannot ": " ",
		" cant ": " ", " co ": " ", " computer ": " ", " con ": " ", " could ": " ", " couldnt ": " ", " cry ": " ",
		" de ": " ", " describe ": " ", " detail ": " ", " do ": " ", " done ": " ", " down ": " ", " due ": " ",
		" during ": " ", " each ": " ", " eg ": " ", " eight ": " ", " either ": " ", " eleven ": " ", " else ": " ",
		" elsewhere ": " ", " empty ": " ", " enough ": " ", " etc ": " ", " even ": " ", " ever ": " ", " every ": " ",
		" everyone ": " ", " everything ": " ", " everywhere ": " ", " except ": " ", " few ": " ", " fifteen ": " ",
		" fify ": " ", " fill ": " ", " find ": " ", " fire ": " ", " first ": " ", " five ": " ", " for ": " ", " former ": " ",
		" formerly ": " ", " forty ": " ", " found ": " ", " four ": " ", " from ": " ", " front ": " ", " full ": " ", " further ": " ",
		" get ": " ", " give ": " ", " go ": " ", " had ": " ", " has ": " ", " hasnt ": " ", " have ": " ", " he ": " ", " hence ": " ",
		" her ": " ", " here ": " ", " hereafter ": " ", " hereby ": " ", " herein ": " ", " hereupon ": " ", " hers ": " ",
		" herse\" ": " ", " him ": " ", " himse\" ": " ", " his ": " ", " how ": " ", " however ": " ", " hundred ": " ",
		" i ": " ", " ie ": " ", " if ": " ", " in ": " ", " inc ": " ", " indeed ": " ", " interest ": " ", " into ": " ",
		" is ": " ", " it ": " ", " its ": " ", " itse\" ": " ", " keep ": " ", " last ": " ", " latter ": " ", " latterly ": " ",
		" least ": " ", " less ": " ", " ltd ": " ", " made ": " ", " many ": " ", " may ": " ", " me ": " ",
		" meanwhile ": " ", " might ": " ", " mill ": " ", " mine ": " ", " more ": " ", " moreover ": " ", " most ": " ",
		" mostly ": " ", " move ": " ", " much ": " ", " must ": " ", " my ": " ", " myse\" ": " ", " name ": " ",
		" namely ": " ", " neither ": " ", " never ": " ", " nevertheless ": " ", " next ": " ", " nine ": " ",
		" no ": " ", " nobody ": " ", " none ": " ", " noone ": " ", " nor ": " ", " not ": " ", " nothing ": " ",
		" now ": " ", " nowhere ": " ", " of ": " ", " off ": " ", " often ": " ", " on ": " ", " once ": " ",
		" one ": " ", " only ": " ", " onto ": " ", " or ": " ", " other ": " ", " others ": " ", " otherwise ": " ",
		" our ": " ", " ours ": " ", " ourselves ": " ", " out ": " ", " over ": " ", " own ": " ",
		" part ": " ", " per ": " ", " perhaps ": " ", " please ": " ", " put ": " ", " rather ": " ",
		" re ": " ", " same ": " ", " see ": " ", " seem ": " ", " seemed ": " ", " seeming ": " ", " seems ": " ",
		" serious ": " ", " several ": " ", " she ": " ", " should ": " ", " show ": " ", " side ": " ", " since ": " ",
		" sincere ": " ", " six ": " ", " sixty ": " ", " so ": " ", " some ": " ", " somehow ": " ", " someone ": " ",
		" something ": " ", " sometime ": " ", " sometimes ": " ", " somewhere ": " ", " still ": " ", " such ": " ",
		" system ": " ", " take ": " ", " ten ": " ", " than ": " ", " that ": " ", " the ": " ", " their ": " ",
		" them ": " ", " themselves ": " ", " then ": " ", " thence ": " ", " there ": " ", " thereafter ": " ",
		" thereby ": " ", " therefore ": " ", " therein ": " ", " thereupon ": " ", " these ": " ", " they ": " ",
		" thick ": " ", " thin ": " ", " third ": " ", " this ": " ", " those ": " ", " though ": " ", " three ": " ",
		" through ": " ", " throughout ": " ", " thru ": " ", " thus ": " ", " to ": " ", " together ": " ", " too ": " ",
		" top ": " ", " toward ": " ", " towards ": " ", " twelve ": " ", " twenty ": " ", " two ": " ", " un ": " ", " under ": " ",
		" until ": " ", " up ": " ", " upon ": " ", " us ": " ", " very ": " ", " via ": " ", " was ": " ", " we ": " ", " well ": " ",
		" were ": " ", " what ": " ", " whatever ": " ", " when ": " ", " whence ": " ", " whenever ": " ", " where ": " ", " whereafter ": " ",
		" whereas ": " ", " whereby ": " ", " wherein ": " ", " whereupon ": " ", " wherever ": " ", " whether ": " ", " which ": " ", " while ": " ",
		" whither ": " ", " who ": " ", " whoever ": " ", " whole ": " ", " whom ": " ", " whose ": " ", " why ": " ", " will ": " ", " with ": " ",
		" within ": " ", " without ": " ", " would ": " ", " yet ": " ", " you ": " ", " your ": " ", " yours ": " ", " yourself ": " ", " yourselves ": " "}
}

func tokenize(text string) []string {
	var cleanToks []string
	text = strings.ToLower(text)
	regHTML := regexp.MustCompile("(&gt;|&lt;)")
	regAp := regexp.MustCompile("&#039;")
	text = regAp.ReplaceAllString(text, "")
	text = regHTML.ReplaceAllString(text, " ")
	regURL := regexp.MustCompile(`(http|ftp|https)://([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?`)
	matches := regURL.FindAllStringSubmatch(text, -1)
	for _, url := range matches {
		cleanToks = append(cleanToks, url[0])
	}
	text = regURL.ReplaceAllString(text, "")

	for word, newWord := range getStopWords() {
		text = strings.Replace(text, word, newWord, -1)
	}
	text = strings.Replace(text, "  ", " ", -1)
	tks := strings.Split(text, " ")
	reg, _ := regexp.Compile("[^a-zA-Z0-9’]+")

	for _, val := range tks {
		val = strings.ReplaceAll(val, " ", "")
		val = strings.ReplaceAll(val, "’", "")
		if len(val) < 2 {
			continue
		}
		if val != "" {
			cleanVal := reg.ReplaceAllString(val, " ")
			tks2 := strings.Split(cleanVal, " ")
			if len(tks2) > 1 {
				for tk := range tks2 {
					if tks2[tk] != "" && len(tks2[tk])>1 {
						cleanToks = append(cleanToks, tks2[tk])
					}
				}
			} else {
				cleanToks = append(cleanToks, cleanVal)
			}
		}
	}
	return cleanToks
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
		wordIdf[word] = math.Log10(N / count)
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
