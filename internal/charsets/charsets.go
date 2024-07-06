package charsets

import (
	"github.com/saintfish/chardet"
	"log"
	"os"
)

func CheckCharset(filename string) bool {
	dat, err := os.ReadFile(filename)
	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest(dat)
	if err == nil {
		if result.Charset == "UTF-8" {
			return true
		} else {
			log.Printf("Charset of file %s is not UTF-8. "+
				"You need to convert it from %s", filename, result.Charset)
			return false
		}
	}
	return false
}
