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

func getStopChars() map[string]string {
	return map[string]string{"&": " ", "  ": " ", "-": " ", "_": " ",
		".": " ", ":": " ", "(": " ", "’": "", "'": "", ";": " ",
		")": " ", "\n": " ", ",": " ", "\"": " ",
	}
}

func getStopWords() map[string]string {
	return map[string]string{" course ": " ", " a ": " ", " about ": " ", " above ": " ", " after ": " ", " all ": " ", " also ": " ", " always ": " ", " am ": " ", " an ": " ", " and ": " ", " any ": " ", " are ": " ", " at ": " ", " be ": " ", " been ": " ", " being ": " ", " but ": " ", " by ": " ", " came ": " ", " can ": " ", " cant ": " ", " come ": " ", " could ": " ", " did ": " ", " didnt ": " ", " do ": " ", " does ": " ", " doesnt ": " ", " doing ": " ", " dont ": " ", " else ": " ", " for ": " ", " from ": " ", " get ": " ", " give ": " ", " goes ": " ", " going ": " ", " had ": " ", " happen ": " ", " has ": " ", " have ": " ", " having ": " ", " how ": " ", " i ": " ", " if ": " ", " ill ": " ", " im ": " ", " in ": " ", " into ": " ", " is ": " ", " isnt ": " ", " it ": " ", " its ": " ", " ive ": " ", " just ": " ", " keep ": " ", " let ": " ", " like ": " ", " made ": " ", " make ": " ", " many ": " ", " may ": " ", " me ": " ", " mean ": " ", " more ": " ", " most ": " ", " much ": " ", " no ": " ", " not ": " ", " now ": " ", " of ": " ", " only ": " ", " or ": " ", " our ": " ", " really ": " ", " say ": " ", " see ": " ", " some ": " ", " something ": " ", " take ": " ", " tell ": " ", " than ": " ", " that ": " ", " the ": " ", " their ": " ", " them ": " ", " then ": " ", " there ": " ", " they ": " ", " thing ": " ", " this ": " ", " to ": " ", " try ": " ", " up ": " ", " us ": " ", " use ": " ", " used ": " ", " uses ": " ", " very ": " ", " want ": " ", " was ": " ", " way ": " ", " we ": " ", " what ": " ", " when ": " ", " where ": " ", " which ": " ", " who ": " ", " why ": " ", " will ": " ", " with ": " ", " without ": " ", " wont ": " ", " you ": " ", " your ": " ", " youre ": " ", " him ": " ", " her ": " ", " again ": " ", " against ": " ", " arent ": " ", " as ": " ", " because ": " ", " before ": " ", " below ": " ", " between ": " ", " both ": " ", " cannot ": " ", " couldnt ": " ", " down ": " ", " during ": " ", " each ": " ", " few ": " ", " further ": " ", " hadnt ": " ", " hasnt ": " ", " havent ": " ", " he ": " ", " hed ": " ", " hell ": " ", " hes ": " ", " here ": " ", " heres ": " ", " hers ": " ", " herself ": " ", " himself ": " ", " his ": " ", " hows ": " ", " id ": " ", " itself ": " ", " lets ": " ", " mustnt ": " ", " my ": " ", " myself ": " ", " nor ": " ", " off ": " ", " on ": " ", " once ": " ", " other ": " ", " ought ": " ", " ours ": " ", " ourselves ": " ", " out ": " ", " over ": " ", " own ": " ", " same ": " ", " shant ": " ", " she ": " ", " shed ": " ", " shell ": " ", " shes ": " ", " should ": " ", " shouldnt ": " ", " so ": " ", " such ": " ", " thats ": " ", " theirs ": " ", " themselves ": " ", " theres ": " ", " these ": " ", " theyd ": " ", " theyll ": " ", " theyre ": " ", " theyve ": " ", " those ": " ", " through ": " ", " too ": " ", " under ": " ", " until ": " ", " wasnt ": " ", " wed ": " ", " well ": " ", " were ": " ", " weve ": " ", " werent ": " ", " whats ": " ", " whens ": " ", " wheres ": " ", " while ": " ", " whos ": " ", " whom ": " ", " whys ": " ", " would ": " ", " wouldnt ": " ", " youd ": " ", " youll ": " ", " youve ": " ", " yours ": " ", " yourself ": " ", " yourselves ": " "}
}
func getStopWordsNoSpaces() map[string]string {
	return map[string]string{"a": " ", "about": " ", "above": " ", "after": " ", "all": " ", "also": " ", "always": " ", "am": " ", "an": " ", "and": " ", "any": " ", "are": " ", "at": " ", "be": " ", "been": " ", "being": " ", "but": " ", "by": " ", "came": " ", "can": " ", "cant": " ", "come": " ", "could": " ", "did": " ", "didnt": " ", "do": " ", "does": " ", "doesnt": " ", "doing": " ", "dont": " ", "else": " ", "for": " ", "from": " ", "get": " ", "give": " ", "goes": " ", "going": " ", "had": " ", "happen": " ", "has": " ", "have": " ", "having": " ", "how": " ", "i": " ", "if": " ", "ill": " ", "im": " ", "in": " ", "into": " ", "is": " ", "isnt": " ", "it": " ", "its": " ", "ive": " ", "just": " ", "keep": " ", "let": " ", "like": " ", "made": " ", "make": " ", "many": " ", "may": " ", "me": " ", "mean": " ", "more": " ", "most": " ", "much": " ", "no": " ", "not": " ", "now": " ", "of": " ", "only": " ", "or": " ", "our": " ", "really": " ", "say": " ", "see": " ", "some": " ", "something": " ", "take": " ", "tell": " ", "than": " ", "that": " ", "the": " ", "their": " ", "them": " ", "then": " ", "there": " ", "they": " ", "thing": " ", "this": " ", "to": " ", "try": " ", "up": " ", "us": " ", "use": " ", "used": " ", "uses": " ", "very": " ", "want": " ", "was": " ", "way": " ", "we": " ", "what": " ", "when": " ", "where": " ", "which": " ", "who": " ", "why": " ", "will": " ", "with": " ", "without": " ", "wont": " ", "you": " ", "your": " ", "youre": " ", "him": " ", "her": " ", "again": " ", "against": " ", "arent": " ", "as": " ", "because": " ", "before": " ", "below": " ", "between": " ", "both": " ", "cannot": " ", "couldnt": " ", "down": " ", "during": " ", "each": " ", "few": " ", "further": " ", "hadnt": " ", "hasnt": " ", "havent": " ", "he": " ", "hed": " ", "hell": " ", "hes": " ", "here": " ", "heres": " ", "hers": " ", "herself": " ", "himself": " ", "his": " ", "hows": " ", "id": " ", "itself": " ", "lets": " ", "mustnt": " ", "my": " ", "myself": " ", "nor": " ", "off": " ", "on": " ", "once": " ", "other": " ", "ought": " ", "ours": " ", "ourselves": " ", "out": " ", "over": " ", "own": " ", "same": " ", "shant": " ", "she": " ", "shed": " ", "shell": " ", "shes": " ", "should": " ", "shouldnt": " ", "so": " ", "such": " ", "thats": " ", "theirs": " ", "themselves": " ", "theres": " ", "these": " ", "theyd": " ", "theyll": " ", "theyre": " ", "theyve": " ", "those": " ", "through": " ", "too": " ", "under": " ", "until": " ", "wasnt": " ", "wed": " ", "well": " ", "were": " ", "weve": " ", "werent": " ", "whats": " ", "whens": " ", "wheres": " ", "while": " ", "whos": " ", "whom": " ", "whys": " ", "would": " ", "wouldnt": " ", "youd": " ", "youll": " ", "youve": " ", "yours": " ", "yourself": " ", "yourselves": " "}
}

func getWords(text string) []string {
	words := regexp.MustCompile("\\w+")
	return words.FindAllString(text, -1)
}

func tokenize(text string) []string {
	var cleanToks []string
	text = " " + text + " "
	text = strings.ToLower(text)
	regURL := regexp.MustCompile(`(http|ftp|https)://([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?`)
	regHTML := regexp.MustCompile("(&gt;|&lt;)")
	regAp := regexp.MustCompile("&#039;")
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")

	matches := regURL.FindAllStringSubmatch(text, -1)
	for _, url := range matches {
		cleanToks = append(cleanToks, url[0])
	}
	text = regURL.ReplaceAllString(text, "")

	text = regAp.ReplaceAllString(text, "")
	text = regHTML.ReplaceAllString(text, " ")
	for word, newWord := range getStopChars() {
		text = strings.ReplaceAll(text, word, newWord)
	}
	for word, newWord := range getStopWords() {
		text = strings.ReplaceAll(text, word, newWord)
	}

	tks := getWords(text)

	for _, val := range tks {
		// if _, ok := getStopWordsNoSpaces()[val]; ok {
		// 	continue //if some stopwords slip through... but it takes a long time so..
		// }
		if len(val) < 2 {
			continue
		}
		cleanVal := reg.ReplaceAllString(val, " ")
		cleanToks = append(cleanToks, cleanVal)
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
