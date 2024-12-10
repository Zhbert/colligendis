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

package collect

import (
	"colligendis/internal/common/colors"
	"colligendis/internal/common/structs"
	"colligendis/internal/db_service"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"gorm.io/gorm"
	"os"
	"sort"
	"strconv"
	"time"
)

func getFullHabrViewsCount(limit int, sortType string, db *gorm.DB) []structs.StatsArticle {
	articles := db_service.GetAllHabrArticles("date", db)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Name", "Date of publication", "Views count", "Growth"})

	var rows []table.Row
	var rowStructs []structs.StatsArticle

	for i := 0; i < len(articles); i++ {
		var zeroTime time.Time
		stats := db_service.GetLatestStatsFromArticle(articles[i].ID, zeroTime, db)
		var stat structs.StatsArticle
		stat.Id = i
		stat.Name = articles[i].Name
		stat.Date = articles[i].DateOfPublication
		if len(stats) > 1 {
			stat.Views = stats[1].Views
			stat.Growth = stats[1].Views - stats[0].Views
		} else if len(stats) == 1 {
			stat.Views = stats[0].Views
			stat.Growth = stats[0].Views
		}
		rowStructs = append(rowStructs, stat)
	}

	if sortType == "views" {
		sort.Slice(rowStructs, func(i, j int) (less bool) {
			return rowStructs[i].Growth > rowStructs[j].Growth
		})
	}

	for i := 0; i < len(rowStructs); i++ {
		var r table.Row
		r = append(r, rowStructs[i].Id)
		r = append(r, rowStructs[i].Name)
		r = append(r, rowStructs[i].Date.Format("02-Jan-2006"))
		r = append(r, rowStructs[i].Views)
		r = append(r, getColorForDiff(rowStructs[i].Growth))

		if limit > 0 && len(rows) <= limit {
			rows = append(rows, r)
		} else if limit == 0 {
			rows = append(rows, r)
		}
	}

	t.AppendRows(rows)
	t.SetStyle(table.StyleLight)
	t.Render()

	return rowStructs
}

func getColorForDiff(diff int) string {
	if diff > 0 {
		return fmt.Sprintf("%s%s%s", colors.Green, strconv.Itoa(diff), colors.Reset)
	}
	return strconv.Itoa(diff)
}
