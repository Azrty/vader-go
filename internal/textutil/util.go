package textutil

import (
	"math"
	"regexp"
	"strings"
)

// Go port of vader sentiment analysis tool, source:
// Hutto, C.J. & Gilbert, E.E. (2014). VADER`: A Parsimonious Rule-based Model for
// Sentiment Analysis of Social Media Text. Eighth International Conference on
// Weblogs and Social Media (ICWSM-14). Ann Arbor, MI, June 2014.

const (
	// BoosterIncrease empirically derived mean sentiment intensity rating increase for booster words
	BoosterIncrease float64 = 0.293
	// BoosterDecrease empirically derived mean sentiment intensity rating increase for booster words
	BoosterDecrease float64 = -0.293
	// CapsIncrease empirically derived mean sentiment intensity rating increase for using ALLCAPs to emphasize a word
	CapsIncrease float64 = 0.733
	// NScalar .
	NScalar float64 = -0.74
)

// Punc simple regex to remove punctuation
var Punc *regexp.Regexp

// Spaces simple regex to remove spaces
var Spaces *regexp.Regexp

func init() {
	Punc = regexp.MustCompile(`[!"#$%&'()*+,-./:;<=>?@[\\\]^_` + "`" + `{|}~…]`)
	Spaces = regexp.MustCompile("[^\\s]+")
}

// AllCapsDifferential Check whether just some words in the input are ALL CAPS
func AllCapsDifferential(words []string) bool {
	var allCapsWords int
	for _, word := range words {
		if strings.ToUpper(word) == word {
			allCapsWords++
		}
	}
	capDiff := len(words) - allCapsWords
	//only true if words are partially caps.
	return capDiff > 0 && capDiff < len(words)
}

// CleanExtraPunc removes contiguous puncs
func CleanExtraPunc(text string) string {
	lastPunc := ' '
	out := strings.Builder{}
	for _, char := range text {
		if Punc.MatchString(string(char)) {
			//special case for ... and …
			if (lastPunc == char) || (lastPunc == '…' && char == '.') || (lastPunc == '.' && char == '…') {
				//we ignore this one because it's a duplicate letter
				continue
			}
		}
		lastPunc = char
		out.WriteRune(char)
	}
	return out.String()
}

//Normalize the score to be between -1 and 1 using an alpha that approximates the max expected value
func Normalize(score float64, alpha float64) float64 {
	norm := score / math.Sqrt((score*score)+alpha)
	switch {
	case norm < -1.0:
		return -1.0
	case norm > 1.0:
		return 1.0
	default:
		return norm
	}
}