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

package db_service

import (
	"colligendis/cmd/common"
	"colligendis/internal/common/structs"
	"colligendis/internal/common/text"
	"colligendis/internal/date_service"
	"colligendis/internal/db_service/domain"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

func SaveToDB(articles []structs.HabrArticle, dateOfStats time.Time, flags *common.ColligendisFlags, db *gorm.DB) bool {
	createDBIfNotExists()

	err := db.AutoMigrate(&domain.HabrArticle{}, &domain.HabrHub{}, &domain.HabrAuthor{}, &domain.HabrStats{})
	if err != nil {
		log.Fatal("Error migrating DB")
		return false
	}
	for i := 0; i < len(articles); i++ {
		var inDBArt domain.HabrArticle
		result := db.Where("habr_number = ?", articles[i].HabrNumber).First(&inDBArt)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			if flags.ViewMode {
				log.Println("Create article: ", articles[i].Name)
			}
			newArticle := createNewArticleEntity(&articles[i], db, flags)
			db.Create(&newArticle)
			inDBArt = newArticle
			if flags.ViewMode {
				log.Printf("Article created: %s \n", newArticle.Name)
			}
		}

		var st domain.HabrStats
		result = db.Where("date_of_stats = ?", dateOfStats).
			Where("habr_article_id = ?", inDBArt.ID).First(&st)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			st.Saves = articles[i].Saves
			st.Views = articles[i].Views
			st.Comments = articles[i].Comments
			st.HabrArticle = inDBArt
			st.LikesAll = articles[i].LikesAll
			st.Likes = articles[i].Likes
			st.LikesDown = articles[i].LikesDown
			st.LikesUp = articles[i].LikesUp
			st.DateOfStats = dateOfStats
			db.Save(&st)
		} else {
			if flags.ViewMode {
				log.Printf("The statistics of %s for %s have already been saved.",
					strconv.Itoa(inDBArt.HabrNumber), dateOfStats)
			}
		}
	}

	return true
}

func createDBIfNotExists() {
	filename := "colligendis.db"
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}
}

func createNewArticleEntity(a *structs.HabrArticle, db *gorm.DB, flags *common.ColligendisFlags) domain.HabrArticle {
	var article = domain.HabrArticle{}

	article.Name = a.Name
	article.Url = a.Url
	article.HabrNumber = a.HabrNumber
	article.DateOfCreation = a.DateOfCreation
	article.TimeOfCreation = a.TimeOfCreation
	article.DateOfCreation = a.DateOfCreation
	article.DateOfPublication = a.DateOfPublication
	article.Author = getAuthor(db, a.Author, flags)
	article.Hubs = getHubs(db, a.Hubs, flags)

	return article
}

func getAuthor(db *gorm.DB, name string, flags *common.ColligendisFlags) domain.HabrAuthor {
	var author domain.HabrAuthor
	result := db.Where("name = ?", name).First(&author)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		var a domain.HabrAuthor
		a.Name = name
		db.Save(&a)
		if flags.ViewMode {
			log.Printf("Author created: %s", a.Name)
		}
		return a
	} else {
		return author
	}
}

func getAuthorByID(id uint, db *gorm.DB) domain.HabrAuthor {
	var author domain.HabrAuthor

	result := db.Where("id = ?", id).First(&author)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Println("There is no such author")
	}

	return author
}

func getHubs(db *gorm.DB, hubs []string, flags *common.ColligendisFlags) []domain.HabrHub {
	createHubsIfNotExists(db, hubs, flags)
	var hubsForDB []domain.HabrHub
	for i := 0; i < len(hubs); i++ {
		var hub domain.HabrHub
		hubs[i] = strings.TrimSpace(hubs[i])
		result := db.Where("name = ?", hubs[i]).First(&hub)
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			hubsForDB = append(hubsForDB, hub)
		}
	}
	return hubsForDB
}

func createHubsIfNotExists(db *gorm.DB, hubs []string, flags *common.ColligendisFlags) {
	for i := 0; i < len(hubs); i++ {
		var hub domain.HabrHub
		hubs[i] = strings.TrimSpace(hubs[i])
		result := db.Where("name = ?", hubs[i]).First(&hub)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			hub.Name = hubs[i]
			db.Save(&hub)
			if flags.ViewMode {
				log.Printf("Hub created: %s \n", hub.Name)
			}
		}
	}
}

func GetAllHabrArticles(sort string, db *gorm.DB) []domain.HabrArticle {
	var articles []domain.HabrArticle

	switch sort {
	case "":
		result := db.Preload(clause.Associations).Find(&articles)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Fatalln("No records found")
		}
	case "name":
		result := db.Preload(clause.Associations).Order(sort).Find(&articles)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Fatalln("No records found")
		}
	case "date":
		result := db.Preload(clause.Associations).Order("date_of_publication DESC").Find(&articles)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Fatalln("No records found")
		}
	}

	return articles
}

func GetHabrArticlesCount(db *gorm.DB) int64 {
	var count int64

	db.Model(&domain.HabrArticle{}).Count(&count)

	return count
}

func GetAllHabrAutors(sort string, db *gorm.DB) []domain.HabrAuthor {
	var authors []domain.HabrAuthor

	switch sort {
	case "":
		result := db.Preload(clause.Associations).Find(&authors)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Fatalln("No records found")
		}
	case "name":
		result := db.Preload(clause.Associations).Order(sort).Find(&authors)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Fatalln("No records found")
		}
	}

	return authors
}

func GetCountOfArticlesByAuthor(authorID uint, db *gorm.DB) int {

	var articles []domain.HabrArticle
	result := db.Where("author_id = ?", authorID).Find(&articles)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Fatalln("No records found")
	}
	return len(articles)

	return 0
}

func GetLogger() logger.LogLevel {
	return logger.Silent
}

func GetLatestArticles(db *gorm.DB) []domain.HabrStats {
	var stats []domain.HabrStats

	db.Order("date_of_stats").Find(&stats)

	return stats
}

func GetLatestStatsFromArticle(articleID uint, sinceDate time.Time, db *gorm.DB) []domain.HabrStats {
	var stats []domain.HabrStats

	if !sinceDate.IsZero() {
		db.
			Where("habr_article_id = ?", articleID).
			Where("date_of_stats <= ?", sinceDate).
			Order("date_of_stats DESC").
			Find(&stats).
			Limit(2)
	} else {
		db.
			Where("habr_article_id = ?", articleID).
			Order("date_of_stats DESC").
			Find(&stats).
			Limit(2)
	}

	if len(stats) > 1 {
		var newStats []domain.HabrStats
		newStats = append(newStats, stats[1])
		newStats = append(newStats, stats[0])
		return newStats
	}

	return stats
}

func GetHabrViewsCount(sinceDate time.Time, db *gorm.DB) int {
	articles := GetAllHabrArticles("", db)
	count := 0

	for i := 0; i < len(articles); i++ {
		stats := GetLatestStatsFromArticle(articles[i].ID, sinceDate, db)

		if len(stats) > 1 {
			diff := stats[1].Views - stats[0].Views
			count = count + diff
		} else if len(stats) == 1 {
			count = count + stats[0].Views
		}
	}

	return count
}

func GetTwoLatestStats(db *gorm.DB) ([]domain.HabrStats, bool) {
	var stats []domain.HabrStats

	db.
		Group("date_of_stats").
		Order("date_of_stats DESC").
		Limit(2).
		Find(&stats)

	return stats, true
}

func GetArticlesFormLastPeriod(dt time.Time, getAll bool, global bool, db *gorm.DB) (int, []structs.StatsArticle) {
	dt = dt.AddDate(0, 0, -1)
	count := 0
	var latestArts []structs.StatsArticle
	articles := GetAllHabrArticles("", db)
	for i := 0; i < len(articles); i++ {
		var zeroTime time.Time
		stats := GetLatestStatsFromArticle(articles[i].ID, zeroTime, db)
		var stat structs.StatsArticle
		stat.Id = i
		stat.Name = text.CleanText(articles[i].Name)
		stat.HabrNumber = articles[i].HabrNumber
		stat.Date = articles[i].DateOfPublication
		stat.Author = getAuthorByID(articles[i].Author.ID, db)
		stat.Author.Name = text.CleanText(stat.Author.Name)
		stat.DayBefore = date_service.GetDaysBefore(articles[i].DateOfPublication, time.Now())
		if len(stats) > 1 {
			stat.Views = stats[1].Views
			stat.Growth = stats[1].Views - stats[0].Views
		} else if len(stats) == 1 {
			stat.Views = stats[0].Views
			stat.Growth = stats[0].Views
		}
		if !getAll {
			if articles[i].DateOfPublication.After(dt) {
				count++
				latestArts = append(latestArts, stat)
			}
		} else {
			count++
			latestArts = append(latestArts, stat)
		}
	}
	if !global {
		sort.Slice(latestArts, func(i, j int) (less bool) {
			return latestArts[i].Growth > latestArts[j].Growth
		})
	} else {
		sort.Slice(latestArts, func(i, j int) (less bool) {
			return latestArts[i].Views > latestArts[j].Views
		})
	}
	return count, latestArts
}

func GetTopOfAuthors(sortName bool, db *gorm.DB) []structs.AuthorsTop {
	var top []structs.AuthorsTop
	var authors []domain.HabrAuthor

	if sortName {
		authors = GetAllHabrAutors("name", db)
	} else {
		authors = GetAllHabrAutors("", db)
	}

	for i := 0; i < len(authors); i++ {
		var t structs.AuthorsTop
		t.Name = text.CleanText(authors[i].Name)
		t.ArticlesCount = getCountOfAuthorArticles(authors[i].ID, db)
		top = append(top, t)
	}

	if !sortName {
		sort.Slice(top, func(i, j int) (less bool) {
			return top[i].ArticlesCount > top[j].ArticlesCount
		})
	}

	return top
}

func getCountOfAuthorArticles(id uint, db *gorm.DB) int64 {
	var count int64

	db.Model(&domain.HabrArticle{}).Where("author_id = ?", id).Count(&count)

	return count
}

func GetCountOfStats(db *gorm.DB) int64 {
	var count int64
	db.Model(&domain.HabrStats{}).Distinct("date_of_stats").Count(&count)
	return count
}

func GetAllDatesOfStats(db *gorm.DB) ([]string, []time.Time) {
	var dates []string

	var stats []domain.HabrStats
	var timeDates []time.Time

	db.
		Group("date_of_stats").
		Order("date_of_stats DESC").
		Find(&stats)
	for i := 0; i < len(stats); i++ {
		timeDates = append(timeDates, stats[i].DateOfStats)
		dates = append(dates, stats[i].DateOfStats.Format("January 2, 2006"))
	}

	return dates, timeDates
}

func GetAllStatsAndDatesForDiagram(db *gorm.DB) ([]structs.StatsForDiagram, float64) {
	var statsForDiagram []structs.StatsForDiagram

	_, dates := GetAllDatesOfStats(db)
	articles := GetAllHabrArticles("", db)

	for i := 0; i < len(dates); i++ {
		var st structs.StatsForDiagram
		st.Date = dates[i].Format("2006-01-02")

		st.Count = 0
		for y := 0; y < len(articles); y++ {
			stats := GetLatestStatsFromArticle(articles[y].ID, dates[i], db)
			if len(stats) > 1 {
				diff := stats[1].Views - stats[0].Views
				st.Count = st.Count + diff
			} else if len(stats) == 1 {
				st.Count = st.Count + stats[0].Views
			}
		}

		statsForDiagram = append(statsForDiagram, st)
	}

	weeks := dates[0].Sub(dates[len(dates)-1]).Hours() / 168

	return statsForDiagram, weeks
}

func GetEachArticleStats(db *gorm.DB) []structs.EachArticleStat {
	articles := GetAllHabrArticles("", db)
	_, dates := GetAllDatesOfStats(db)
	var eachArticleStats []structs.EachArticleStat

	for i := 0; i < len(articles); i++ {
		var eas structs.EachArticleStat
		eas.Name = text.CleanText(articles[i].Name)
		eas.HabrNumber = articles[i].HabrNumber
		for y := 0; y < len(dates); y++ {
			var st structs.StatsForDiagram
			st.Date = dates[y].Format("2006-01-02")
			stats := GetLatestStatsFromArticle(articles[i].ID, dates[y], db)
			if len(stats) > 1 {
				st.Count = stats[1].Views - stats[0].Views
			} else if len(stats) == 1 {
				st.Count = stats[0].Views
			}
			eas.Stats = append(eas.Stats, st)
		}
		eachArticleStats = append(eachArticleStats, eas)
		slices.Reverse(eachArticleStats)
	}
	return eachArticleStats
}
