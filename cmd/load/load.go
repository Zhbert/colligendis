/*
 * Copyright (c) 2024. Konstantin Nezhbert.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "colligendis"), to deal in
 * the Software without restriction, including without limitation the rights to use,
 * copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the
 * Software, and to permit persons to whom the Software is furnished to do so, subject
 * to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
 * INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
 * PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
 * LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
 * TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE
 * USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package load

import (
	"colligendis/cmd/common"
	"colligendis/internal/db_service"
	"fmt"
	"github.com/spf13/cobra"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

func GetLoadCommand(flags *common.ColligendisFlags) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "load",
		Short:   "Load statistics from sources",
		Long:    `Loads data from available resources`,
		Example: `colligendis load`,
		Run: func(cmd *cobra.Command, args []string) {
			if flags.Habr {
				if flags.FromCSV != "" && !flags.AllCSV {
					processFile(flags.FromCSV, flags)
				}
				if flags.AllCSV {
					currentDirectory, err := os.Getwd()
					if err != nil {
						fmt.Println(err)
					} else {
						if flags.ViewMode {
							log.Printf("The download of CSV files from the catalog starts: %s...", currentDirectory)
						}
						for _, csvFile := range findCsvFiles(currentDirectory, ".csv") {
							processFile(csvFile, flags)
						}
					}
				}
			}
		},
	}
	cmd.Flags().BoolVarP(&flags.Habr, "habr", "", false, "Downloading data from a CSV file with statistics from Habr")
	cmd.Flags().BoolVarP(&flags.GitHub, "github", "", false, "Downloading data from GitHub repositories")
	cmd.Flags().StringVarP(&flags.FromRepo, "from-repo", "r", "", "The address of the GitHub repository")
	cmd.Flags().StringVarP(&flags.FromCSV, "from-csv", "f", "", "The path to the CSV file from the Habr")
	cmd.Flags().BoolVarP(&flags.AllCSV, "all-csv", "a", false, "Upload all CSV files")
	return cmd
}

func getStatsDate(filename string) time.Time {
	re := regexp.MustCompile("[0-9][0-9][0-9][0-9]_[0-9][0-9]_[0-9][0-9]")
	match := re.FindStringSubmatch(filename)
	dateOfStats, err := time.Parse("2006_01_02", match[0])
	if err != nil {
		log.Fatalln("Error parsing DateOfPublication", err)
	}
	return dateOfStats
}

func findCsvFiles(pathToDir, ext string) []string {
	var a []string
	err := filepath.WalkDir(pathToDir, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			a = append(a, s)
		}
		return nil
	})
	if err != nil {
		return nil
	}
	return a
}

func processFile(csvFile string, flags *common.ColligendisFlags) {
	start := time.Now()
	if flags.ViewMode {
		log.Printf("Parsing file %s... \n", csvFile)
	}
	articles := loadHabrFromFile(csvFile)
	if articles != nil {
		if db_service.SaveToDB(articles, getStatsDate(csvFile), flags) {
			if flags.ViewMode {
				duration := time.Since(start)
				log.Printf("The file '%s' has been processed in %s \n", csvFile, duration)
			}
		} else {
			log.Printf("File %s processing error", csvFile)
		}
	} else {
		log.Printf("The file %s encoding is not UTF-8!", csvFile)
	}
}
