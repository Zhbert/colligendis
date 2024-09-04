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
	"colligendis/internal/charsets"
	"colligendis/internal/common/structs"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func loadHabrFromFile(fileName string) []structs.HabrArticle {
	newFile, err := os.Open(fileName)
	if err != nil {
		log.Fatalln("Can't open result CSV-file!", err)
	}
	defer newFile.Close()

	charsetType, _ := charsets.CheckCharset(fileName)

	if charsetType {
		r := csv.NewReader(newFile)
		r.Comma = '\t'
		r.FieldsPerRecord = 16
		r.LazyQuotes = true

		// skip first line
		if _, err := r.Read(); err != nil {
			log.Fatalln("Error reading file!", err)
		}

		records, err := r.ReadAll()

		if err != nil {
			log.Fatalln("Error parsing CSV!", err)
		}

		var articles []structs.HabrArticle

		for _, record := range records {
			a := structs.HabrArticle{}
			a.HabrNumber, _ = strconv.Atoi(record[0])
			a.Url = record[1]
			a.Name = strings.TrimSpace(record[2])
			a.Author = record[3]

			hubsReader := csv.NewReader(strings.NewReader(record[4]))
			hubs, err := hubsReader.ReadAll()
			if err != nil {
				log.Fatalln("Error parsing hubs string!", err)
			}
			a.Hubs = hubs[0]

			a.DateOfCreation, err = time.Parse("2.1.2006", record[5])
			if err != nil {
				log.Fatalln("Error parsing DateOfCreation", err)
			}
			a.TimeOfCreation, err = time.Parse("15:04", record[6])
			if err != nil {
				log.Fatalln("Error parsing TimeOfCreation", err)
			}
			a.DateOfPublication, err = time.Parse("2.1.2006", record[7])
			if err != nil {
				log.Fatalln("Error parsing DateOfPublication", err)
			}
			a.TimeOfPublication, err = time.Parse("15:04", record[8])
			if err != nil {
				log.Fatalln("Error parsing TimeOfPublication", err)
			}
			a.Likes, _ = strconv.Atoi(record[9])
			a.LikesUp, _ = strconv.Atoi(record[10])
			a.LikesDown, _ = strconv.Atoi(record[11])
			a.LikesAll, _ = strconv.Atoi(record[12])
			a.Comments, _ = strconv.Atoi(record[13])
			a.Saves, _ = strconv.Atoi(record[14])
			a.Views, _ = strconv.Atoi(record[15])

			articles = append(articles, a)
		}

		return articles
	}
	return nil
}
