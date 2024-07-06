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

package view

import (
	"colligendis/internal/db_service"
	"colligendis/internal/db_service/domain"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
)

func viewAllHabrArticles(sort string) {
	var articles []domain.HabrArticle

	articles = db_service.GetAllHabrArticles(sort)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Name", "Habr number", "Author", "Time of publication"})

	var rows []table.Row
	for i, a := range articles {
		var r table.Row
		r = append(r, i+1)
		r = append(r, a.Name)
		r = append(r, a.HabrNumber)
		r = append(r, a.Author.Name)
		r = append(r, a.DateOfPublication.Format("02-Jan-2006"))
		rows = append(rows, r)
	}

	t.AppendRows(rows)
	t.SetStyle(table.StyleLight)
	t.Render()
}

func viewAllHabrAuthors(sort string) {
	var authors []domain.HabrAuthor

	authors = db_service.GetAllHabrAutors(sort)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Name", "Articles count"})

	var rows []table.Row
	for i, a := range authors {
		var r table.Row
		r = append(r, i+1)
		r = append(r, a.Name)
		r = append(r, db_service.GetCountOfArticlesByAuthor(a.ID))
		rows = append(rows, r)
	}

	t.AppendRows(rows)
	t.SetStyle(table.StyleLight)
	t.Render()
}
