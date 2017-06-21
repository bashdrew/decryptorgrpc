package freqanalysis

import (
	"regexp"
	"sort"
	"strings"
	"unicode"
)

const (
	alphaCnt = 26
)

// LetterFreq ... Struct of letter and it's frequency
type LetterFreq struct {
	Letter string
	Freq   int
}

var areLetters = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString

// LetterFreqList ... List of letters and each letter's frequency
type LetterFreqList []LetterFreq

// LetterMappingList ... List of LetterMapping
type LetterMappingList map[string]string

func (lfl LetterFreqList) Len() int {
	return len(lfl)
}

func (lfl LetterFreqList) Less(i, j int) bool {
	return lfl[i].Freq < lfl[j].Freq
}

func (lfl LetterFreqList) Swap(i, j int) {
	lfl[i], lfl[j] = lfl[j], lfl[i]
}

// Swap ... Swaps the key and value in a map
func (lml LetterMappingList) Swap() (result LetterMappingList) {
	result = make(LetterMappingList, alphaCnt)
	for key, val := range lml {
		result[val] = key
	}

	return result
}

// GetLettersList ... Returns encrypted and plain letters as arrays
func (lml LetterMappingList) GetLettersList() (encLtrs []string, plnLtrs []string) {
	encLtrs = make([]string, len(lml))
	plnLtrs = make([]string, len(lml))

	idx := 0
	for key, val := range lml {
		encLtrs[idx] = key
		plnLtrs[idx] = val

		idx++
	}

	return encLtrs, plnLtrs
}

// GetCipherKey ... Returns cipherkey
func (lml LetterMappingList) GetCipherKey() (cipherkey []string) {
	ltrMap := lml.Swap()
	i := 0
	cipherkey = make([]string, len(lml))
	for charKey := 'A'; charKey <= 'Z'; charKey++ {
		strKey, _ := ltrMap[string(charKey)]
		cipherkey[i] = strKey
		i++
	}

	return cipherkey
}

// Append ... Appends LetterMappingList
func (lml LetterMappingList) Append(ltrMapList LetterMappingList) LetterMappingList {
	for key, val := range ltrMapList {
		lml[key] = val
	}

	return lml
}

func sortByLetterCount(ltrFrq map[string]int) (result LetterFreqList) {
	result = make(LetterFreqList, len(ltrFrq))

	i := 0
	for k, v := range ltrFrq {
		result[i] = LetterFreq{k, v}
		i++
	}
	sort.Sort(sort.Reverse(result))

	return result
}

// GetLetterFreq ...  returns a list of letters and frequencies in a string
func GetLetterFreq(inTxt string) (result LetterFreqList) {
	lfMap := make(map[string]int, alphaCnt)

	for _, val := range inTxt {
		if unicode.IsLetter(val) {
			valUpper := strings.ToUpper(string(val))
			lfMap[valUpper] = lfMap[valUpper] + 1
		}
	}

	result = sortByLetterCount(lfMap)
	return result
}

// GetLetterFreqMulti ...  returns a list of letters and frequencies in a string
func GetLetterFreqMulti(inTxt string, ltrCnt int, ltrExcl []string) (result LetterFreqList) {
	lfMap := make(map[string]int, alphaCnt)
	ltrExclStr := strings.Join(ltrExcl, ",")

	var valLetter string
	i := 0
	for _, val := range inTxt {
		if unicode.IsLetter(val) {
			valUpper := strings.ToUpper(string(val))

			i++
			valLetter = valLetter + valUpper

			if i >= ltrCnt {
				if !strings.Contains(ltrExclStr, valLetter) {
					lfMap[valLetter] = lfMap[valLetter] + 1
				}
				valLetter = ""
				i = 0
			}
		} else {
			i = 0
			valLetter = ""
		}
	}

	result = sortByLetterCount(lfMap)
	return result
}

// GetWordFreq ...  returns a list of letters and frequencies in a string
func GetWordFreq(inTxt string, ltrCnt int, wrdExcl []string) (result LetterFreqList) {
	lfMap := make(map[string]int, alphaCnt)
	wrdExclStr := strings.Join(wrdExcl, ",")

	for _, val := range strings.Fields(inTxt) {
		if areLetters(val) && len(val) == ltrCnt && !strings.Contains(wrdExclStr, val) {
			valUpper := strings.ToUpper(string(val))
			lfMap[valUpper] = lfMap[valUpper] + 1
		}
	}

	result = sortByLetterCount(lfMap)
	return result
}

// GenEncPlainMapping ... Generate letter mapping
func GenEncPlainMapping(lflEnc, lflPlain LetterFreqList, iTop int) (result LetterMappingList) {
	result = make(LetterMappingList, len(lflEnc))
	for idx := range lflEnc {
		if len(lflEnc[idx].Letter) > 1 {
			for i, val := range lflEnc[idx].Letter {
				valLetter := string(val)
				if _, ok := result[valLetter]; !ok {
					result[valLetter] = string(lflPlain[idx].Letter[i])
				}
			}
		} else {
			result[lflEnc[idx].Letter] = lflPlain[idx].Letter
		}

		if iTop > 0 && (iTop-1) == idx {
			break
		}
	}

	return result
}

// SubstEncByPlain ...  returns a list of letters and frequencies in a string
func SubstEncByPlain(inTxt string, ltrMap LetterMappingList) (result []rune) {
	for _, val := range inTxt {
		plainLtr := val
		if unicode.IsLetter(val) {
			valUpper := strings.ToUpper(string(val))
			if plainUpper, ok := ltrMap[valUpper]; ok {
				plainLtr = rune(plainUpper[0])
				if unicode.IsLower(val) {
					plainLtr = unicode.ToLower(rune(plainUpper[0]))
				}
			}
		}

		result = append(result, plainLtr)
	}

	return result
}
