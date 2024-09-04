package charsets

import (
	"colligendis/cmd/common"
	"github.com/saintfish/chardet"
	"log"
	"os"
	"os/exec"
)

func CheckCharset(filename string) (bool, string) {
	dat, err := os.ReadFile(filename)
	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest(dat)
	if err == nil {
		if result.Charset == "UTF-8" {
			return true, ""
		} else {
			return false, result.Charset
		}
	}
	return false, ""
}

func ConvertFileToUTF(pathToFile string, charsetType string, flags *common.ColligendisFlags) {
	log.Printf("File %s is not in UTF-8 encoding. Converting from %s...", pathToFile, charsetType)
	cmd := exec.Command("sh", "./convert_csv.sh", pathToFile)
	err := cmd.Run()
	if err != nil {
		return
	}
}
