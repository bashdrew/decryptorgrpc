package decryptor

import (
	fa "bashdrew/bsscodingassignment/freqanalysis"
	"strings"
)

// DecryptInfo ... decrypt info
type DecryptInfo struct {
	FreqFunc func(string, int, []string) fa.LetterFreqList
	LtrCnt   int
	EntTop   int
}

// DecryptInfoList ... decrypt info list
type DecryptInfoList []DecryptInfo

func getLetterFreq(inTxt string, ltrCnt int, ltrExclList []string, getFunc func(string, int, []string) fa.LetterFreqList) (result fa.LetterFreqList) {
	result = getFunc(inTxt, ltrCnt, ltrExclList)

	return result
}

// Decrypt ... Decrypts a string based on decryptIfno values
func (diLst DecryptInfoList) Decrypt(encTxt, dicTxt string) (plnTxt []byte, cipherKey string) {
	var encLtrsToExcl, plnLtrsToExcl []string
	var lflEnc, lflDic fa.LetterFreqList

	ltrMap := make(fa.LetterMappingList, 0)

	for _, dInfo := range diLst {
		lflEnc = getLetterFreq(encTxt, dInfo.LtrCnt, encLtrsToExcl, dInfo.FreqFunc)
		lflDic = getLetterFreq(dicTxt, dInfo.LtrCnt, plnLtrsToExcl, dInfo.FreqFunc)
		ltrMapLtr := fa.GenEncPlainMapping(lflEnc, lflDic, dInfo.EntTop)

		encExcl, plnExcl := ltrMapLtr.GetLettersList()
		encLtrsToExcl = append(encLtrsToExcl, encExcl...)
		plnLtrsToExcl = append(plnLtrsToExcl, plnExcl...)
		ltrMap = ltrMap.Append(ltrMapLtr)
	}
	plnTxt = []byte(string(fa.SubstEncByPlain(encTxt, ltrMap)))

	cipherKeyLst := ltrMap.GetCipherKey()
	cipherKey = strings.Join(cipherKeyLst, "")

	return plnTxt, cipherKey
}
