package convert

import (
	"colligendis/cmd/common"
	"colligendis/internal/charsets"
	"colligendis/internal/common/structs"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

func GetConvertCommand(flags *common.ColligendisFlags, tmpls []structs.TemplateStruct) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "convert",
		Short:   "Command for automatic conversion of CSV files",
		Long:    "Automatically converts CSV files of Habr statistics to UTF-8 encoding",
		Example: `colligendis convert`,
		Run: func(cmd *cobra.Command, args []string) {
			//create shell script
			for _, t := range tmpls {
				if t.Name == "convert_csv.sh" {
					path := filepath.Join(t.Name)
					err := os.WriteFile(path, t.Data, 0777)
					if err != nil {
						log.Println("Can't create file " + path)
					}
				}
			}
			// process csv files
			files, err := filepath.Glob("*.csv")
			if err != nil {
				fmt.Println(err)
				return
			}

			log.Printf("Found CSV files: %d\n", len(files))

			for _, pathToFile := range files {
				charsetOK, charsetType := charsets.CheckCharset(pathToFile)
				if !charsetOK {
					charsets.ConvertFileToUTF(pathToFile, charsetType, flags)
				}
			}
		},
	}
	return cmd
}
