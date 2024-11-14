package csv_service

import (
	"colligendis/internal/common/structs"
	"encoding/csv"
	"log"
	"os"
	"path"
	"strconv"
)

func PrepareCSV(pathToFile string, fileNAme string, data []structs.StatsForDiagram) {
	file, err := os.Create(path.Join(pathToFile, fileNAme))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Comma = ','

	writer.Write([]string{"date", "count"})

	count := 1
	for _, d := range data {
		if count != len(data) {
			writer.Write([]string{d.Date, strconv.Itoa(d.Count)})
			count++
		}
	}

}
