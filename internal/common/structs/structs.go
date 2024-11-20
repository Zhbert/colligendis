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

package structs

import (
	"colligendis/internal/db_service/domain"
	"time"
)

type HabrArticle struct {
	Id                int
	HabrNumber        int
	Name              string
	Url               string
	Author            string
	Hubs              []string
	DateOfCreation    time.Time
	TimeOfCreation    time.Time
	DateOfPublication time.Time
	TimeOfPublication time.Time
	LikesAll          int
	LikesUp           int
	LikesDown         int
	Likes             int
	Comments          int
	Saves             int
	Views             int
}

type StatsArticle struct {
	Id        int
	Name      string
	Date      time.Time
	Views     int
	Growth    int
	Author    domain.HabrAuthor
	DayBefore int
}

type TemplateStruct struct {
	Name string
	Data []byte
}

type TemplateData struct {
	Version                    string
	StatsInBaseCount           int64
	AllViewsCount              int
	PreviousDate               string
	LatestDate                 string
	ArticlesCount              int64
	CountOfLastArticles        int
	LatestArticlesFromWeek     []StatsArticle
	AllArticlesPerWeek         []StatsArticle
	AllArticlesGlobal          []StatsArticle
	AllArticlesGlobalWithLimit []StatsArticle
	AuthorsCount               int
	AuthorsTopGlobal           []AuthorsTop
	Authors                    []AuthorsTop
	AllDates                   []string
	StatsForDiagram            []StatsForDiagram
	WeeksCount                 float64
}

type AuthorsTop struct {
	Name          string
	ID            int
	ArticlesCount int64
}

type StatsForDiagram struct {
	Date  string
	Count int
}
