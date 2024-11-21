package latex_service

import (
	"colligendis/cmd/common"
	"colligendis/cmd/version"
	"colligendis/internal/common/structs"
	"colligendis/internal/csv_service"
	"colligendis/internal/db_service"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func GenerateLaTeXFiles(tmpls []structs.TemplateStruct, flags *common.ColligendisFlags) {

	checkAndCreateTmpFolder()
	checkAndCreatePDFsFolder()
	removeOldFilesFromTmpFolder()

	for _, t := range tmpls {
		if t.Name != "convert_csv.sh" {
			path := filepath.Join("tmp", t.Name)
			err := os.WriteFile(path, t.Data, 0777)
			if err != nil {
				log.Println("Can't create file " + path)
			}
		}
	}

	templateData := getHabrData()

	err := createFileByTemplate("tmp/stats.tmpl", "tmp/stats.tex", templateData)
	if err != nil {
		log.Printf("Unable to create file: %w", err)
	}

	if !flags.DryRun {
		generatePDF(flags)
		generatePDF(flags)
		copyFile(templateData.LatestDate)
	}
}

func getHabrData() structs.TemplateData {
	var data structs.TemplateData

	data.Version = version.GetVersion()

	data.StatsInBaseCount = db_service.GetCountOfStats()
	var zeroTime time.Time
	data.AllViewsCount = db_service.GetHabrViewsCount(zeroTime)
	data.PreviousDate, data.LatestDate = getDates()
	data.ArticlesCount = db_service.GetHabrArticlesCount()
	pd, _ := time.Parse("2006-01-02", data.PreviousDate)
	data.CountOfLastArticles, data.LatestArticlesFromWeek = db_service.GetArticlesFormLastPeriod(pd, false, false)
	data.AuthorsCount = getAuthorsCount()
	_, data.AllArticlesPerWeek = db_service.GetArticlesFormLastPeriod(pd, true, false)
	data.AllArticlesPerWeek = data.AllArticlesPerWeek[0:5]
	_, data.AllArticlesGlobalWithLimit = db_service.GetArticlesFormLastPeriod(pd, true, true)
	data.AllArticlesGlobalWithLimit = data.AllArticlesGlobalWithLimit[0:5]
	_, data.AllArticlesGlobal = db_service.GetArticlesFormLastPeriod(pd, true, false)
	data.AuthorsTopGlobal = db_service.GetTopOfAuthors(false)
	data.AuthorsTopGlobal = data.AuthorsTopGlobal[0:5]
	data.Authors = db_service.GetTopOfAuthors(true)
	data.AllDates, _ = db_service.GetAllDatesOfStats()
	data.StatsForDiagram, data.WeeksCount = db_service.GetAllStatsAndDatesForDiagram()

	csv_service.PrepareCSV("tmp", "articlesCount.csv", data.StatsForDiagram)

	return data
}

func getDates() (string, string) {
	stats, state := db_service.GetTwoLatestStats()
	if state {
		switch len(stats) {
		case 0:
			return "", ""
		case 1:
			latestDate := stats[0].DateOfStats.Format("2006-01-02")
			return "", latestDate
		default:
			if len(stats) > 1 {
				previousDate := stats[1].DateOfStats.Format("2006-01-02")
				latestDate := stats[0].DateOfStats.Format("2006-01-02")
				return previousDate, latestDate
			}
		}
	}
	return "", ""
}

func getAuthorsCount() int {
	return len(db_service.GetAllHabrAutors(""))
}

func generatePDF(flags *common.ColligendisFlags) {
	cmd := exec.Command("pdflatex", "stats.tex")
	cmd.Dir = "tmp"
	if flags.ViewMode {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	err := cmd.Run()
	if err != nil {
		log.Fatalln("Error processing LaTeX file!")
	}
}

func checkAndCreateTmpFolder() {
	tmpFolder := filepath.Join("tmp")
	_, err := os.Stat(tmpFolder)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(tmpFolder, 0777)
			if err != nil {
				log.Println("Error creating tmp folder")
			}
		}
	}
}

func checkAndCreatePDFsFolder() {
	tmpFolder := filepath.Join("PDFs")
	_, err := os.Stat(tmpFolder)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(tmpFolder, 0777)
			if err != nil {
				log.Println("Error creating PDFs folder")
			}
		}
	}
}

func removeOldFilesFromTmpFolder() {
	files, err := filepath.Glob(filepath.Join("tmp/*"))
	if err != nil {
		log.Println("Can't delete tmp folder!")
	}
	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			log.Println("Can't delete tmp folder!")
		}
	}
}

func copyFile(latestDate string) {
	srcFile, err := os.Open("tmp/stats.pdf")
	if err != nil {
		panic(err)
	}
	defer srcFile.Close()

	filename := fmt.Sprintf("PDFs/flant_stats_%s.pdf", latestDate)

	destFile, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		panic(err)
	}
}

func createFileByTemplate(templatePath, filePath string, data interface{}) error {
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("unable to parse template: %w", err)
	}

	f, err := createFile(filePath, err)
	if err != nil {
		return fmt.Errorf("unable to create file: %w", err)
	}
	defer f.Close()

	err = t.Execute(f, data)
	if err != nil {
		return fmt.Errorf("unable to execute template: %w", err)
	}

	return nil
}

func createFile(filePath string, err error) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(filePath), 0o777); err != nil {
		return nil, fmt.Errorf("unable to create directory: %w", err)
	}

	f, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to create file: %w", err)
	}
	return f, nil
}
