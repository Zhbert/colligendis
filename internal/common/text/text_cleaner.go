package text

import (
	"regexp"
)

func CleanText(text string) string {
	var deleteApostrophe = regexp.MustCompile(`'`)
	var deleteSharp = regexp.MustCompile(`#`)
	var deleteSingleAmp = regexp.MustCompile(`&`)
	var deleteColon = regexp.MustCompile(`:`)
	var deleteDollar = regexp.MustCompile(`\$`)
	var deleteUnderline = regexp.MustCompile(`_`)
	var deletePlus = regexp.MustCompile(`\+`)
	result := deleteApostrophe.ReplaceAllString(text, "")
	result = deleteSharp.ReplaceAllString(result, "")
	result = deleteColon.ReplaceAllString(result, "")
	result = deleteDollar.ReplaceAllString(result, "")
	result = deleteUnderline.ReplaceAllString(result, " ")
	result = deletePlus.ReplaceAllString(result, "")
	result = deleteSingleAmp.ReplaceAllString(result, "")
	return result
}
