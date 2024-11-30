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
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func SaveToDB(articles []structs.HabrArticle, dateOfStats time.Time, flags *common.ColligendisFlags) bool {
	createDBIfNotExists()
	db, err := gorm.Open(sqlite.Open("colligendis.db"),
		&gorm.Config{Logger: logger.Default.LogMode(GetLogger())})
	if err != nil {
		log.Fatal("Error opening db!")
	} else {
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

func getAuthorByID(id uint) domain.HabrAuthor {
	var author domain.HabrAuthor

	db, err := gorm.Open(sqlite.Open("colligendis.db"),
		&gorm.Config{Logger: logger.Default.LogMode(GetLogger())})
	if err != nil {
		log.Fatal("Error opening db!")
	} else {
		result := db.Where("id = ?", id).First(&author)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Println("There is no such author")
		}
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

func GetAllHabrArticles(sort string) []domain.HabrArticle {
	var articles []domain.HabrArticle
	db, err := gorm.Open(sqlite.Open("colligendis.db"),
		&gorm.Config{Logger: logger.Default.LogMode(GetLogger())})
	if err != nil {
		log.Fatal("Error opening db!")
	} else {
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
	}

	return articles
}

func GetHabrArticlesCount() int64 {
	var count int64

	db, err := gorm.Open(sqlite.Open("colligendis.db"),
		&gorm.Config{Logger: logger.Default.LogMode(GetLogger())})
	if err != nil {
		log.Fatal("Error opening db!")
	} else {
		db.Model(&domain.HabrArticle{}).Count(&count)
	}

	return count
}

func GetAllHabrAutors(sort string) []domain.HabrAuthor {
	var authors []domain.HabrAuthor

	db, err := gorm.Open(sqlite.Open("colligendis.db"),
		&gorm.Config{Logger: logger.Default.LogMode(GetLogger())})
	if err != nil {
		log.Fatal("Error opening db!")
	} else {
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
	}

	return authors
}

func GetCountOfArticlesByAuthor(authorID uint) int {
	db, err := gorm.Open(sqlite.Open("colligendis.db"),
		&gorm.Config{Logger: logger.Default.LogMode(GetLogger())})
	if err != nil {
		log.Fatal("Error opening db!")
	} else {
		var articles []domain.HabrArticle
		result := db.Where("author_id = ?", authorID).Find(&articles)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Fatalln("No records found")
		}
		return len(articles)
	}
	return 0
}

func GetLogger() logger.LogLevel {
	return logger.Silent
}

func GetLatestArticles() []domain.HabrStats {
	var stats []domain.HabrStats

	db, err := gorm.Open(sqlite.Open("colligendis.db"),
		&gorm.Config{Logger: logger.Default.LogMode(GetLogger())})
	if err != nil {
		log.Fatal("Error opening db!")
	} else {
		db.Order("date_of_stats").Find(&stats)
	}

	return stats
}

func GetLatestStatsFromArticle(articleID uint, sinceDate time.Time) ([]domain.HabrStats, bool) {
	var stats []domain.HabrStats
	state := false

	db, err := gorm.Open(sqlite.Open("colligendis.db"),
		&gorm.Config{Logger: logger.Default.LogMode(GetLogger())})
	if err != nil {
		log.Fatal("Error opening db!")
	} else {
		state = true
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
	}

	if len(stats) > 1 {
		var newStats []domain.HabrStats
		newStats = append(newStats, stats[1])
		newStats = append(newStats, stats[0])
		return newStats, state
	}

	return stats, state
}

func GetHabrViewsCount(sinceDate time.Time) int {
	articles := GetAllHabrArticles("")
	count := 0

	for i := 0; i < len(articles); i++ {
		stats, state := GetLatestStatsFromArticle(articles[i].ID, sinceDate)
		if state {
			if len(stats) > 1 {
				diff := stats[1].Views - stats[0].Views
				count = count + diff
			} else if len(stats) == 1 {
				count = count + stats[0].Views
			}
		} else {
			log.Println("There are no stats in database!")
		}
	}

	return count
}

func GetTwoLatestStats() ([]domain.HabrStats, bool) {
	var stats []domain.HabrStats
	state := false

	db, err := gorm.Open(sqlite.Open("colligendis.db"),
		&gorm.Config{Logger: logger.Default.LogMode(GetLogger())})
	if err != nil {
		log.Fatal("Error opening db!")
	} else {
		state = true
		db.
			Group("date_of_stats").
			Order("date_of_stats DESC").
			Limit(2).
			Find(&stats)
	}

	return stats, state
}

func GetArticlesFormLastPeriod(dt time.Time, getAll bool, global bool) (int, []structs.StatsArticle) {
	dt = dt.AddDate(0, 0, -1)
	count := 0
	var latestArts []structs.StatsArticle
	articles := GetAllHabrArticles("")
	for i := 0; i < len(articles); i++ {
		var zeroTime time.Time
		stats, state := GetLatestStatsFromArticle(articles[i].ID, zeroTime)
		var stat structs.StatsArticle
		if state {
			stat.Id = i
			stat.Name = text.CleanText(articles[i].Name)
			stat.Date = articles[i].DateOfPublication
			stat.Author = getAuthorByID(articles[i].Author.ID)
			stat.Author.Name = text.CleanText(stat.Author.Name)
			stat.DayBefore = date_service.GetDaysBefore(articles[i].DateOfPublication, time.Now())
			if len(stats) > 1 {
				stat.Views = stats[1].Views
				stat.Growth = stats[1].Views - stats[0].Views
			} else if len(stats) == 1 {
				stat.Views = stats[0].Views
				stat.Growth = stats[0].Views
			}
		} else {
			log.Println("There are no stats in database!")
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

func GetTopOfAuthors(sortName bool) []structs.AuthorsTop {
	var top []structs.AuthorsTop
	var authors []domain.HabrAuthor

	if sortName {
		authors = GetAllHabrAutors("name")
	} else {
		authors = GetAllHabrAutors("")
	}

	for i := 0; i < len(authors); i++ {
		var t structs.AuthorsTop
		t.Name = text.CleanText(authors[i].Name)
		t.ArticlesCount = getCountOfAuthorArticles(authors[i].ID)
		top = append(top, t)
	}

	if !sortName {
		sort.Slice(top, func(i, j int) (less bool) {
			return top[i].ArticlesCount > top[j].ArticlesCount
		})
	}

	return top
}

func getCountOfAuthorArticles(id uint) int64 {
	var count int64

	db, err := gorm.Open(sqlite.Open("colligendis.db"),
		&gorm.Config{Logger: logger.Default.LogMode(GetLogger())})
	if err != nil {
		log.Fatal("Error opening db!")
	} else {
		db.Model(&domain.HabrArticle{}).Where("author_id = ?", id).Count(&count)
	}
	return count
}

func GetCountOfStats() int64 {
	var count int64

	db, err := gorm.Open(sqlite.Open("colligendis.db"),
		&gorm.Config{Logger: logger.Default.LogMode(GetLogger())})
	if err != nil {
		log.Fatal("Error opening db!")
	} else {
		db.Model(&domain.HabrStats{}).Distinct("date_of_stats").Count(&count)
	}
	return count
}

func GetAllDatesOfStats() ([]string, []time.Time) {
	var dates []string

	var stats []domain.HabrStats
	var timeDates []time.Time

	db, err := gorm.Open(sqlite.Open("colligendis.db"),
		&gorm.Config{Logger: logger.Default.LogMode(GetLogger())})
	if err != nil {
		log.Fatal("Error opening db!")
	} else {
		db.
			Group("date_of_stats").
			Order("date_of_stats DESC").
			Find(&stats)
		for i := 0; i < len(stats); i++ {
			timeDates = append(timeDates, stats[i].DateOfStats)
			dates = append(dates, stats[i].DateOfStats.Format("January 2, 2006"))
		}
	}

	return dates, timeDates
}

func GetAllStatsAndDatesForDiagram() ([]structs.StatsForDiagram, float64) {
	var statsForDiagram []structs.StatsForDiagram

	_, dates := GetAllDatesOfStats()
	articles := GetAllHabrArticles("")

	for i := 0; i < len(dates); i++ {
		var st structs.StatsForDiagram
		st.Date = dates[i].Format("2006-01-02")

		st.Count = 0
		for y := 0; y < len(articles); y++ {
			stats, state := GetLatestStatsFromArticle(articles[y].ID, dates[i])
			if state {
				if len(stats) > 1 {
					diff := stats[1].Views - stats[0].Views
					st.Count = st.Count + diff
				} else if len(stats) == 1 {
					st.Count = st.Count + stats[0].Views
				}
			} else {
				log.Println("There are no stats in database!")
			}
		}

		statsForDiagram = append(statsForDiagram, st)
	}

	weeks := dates[0].Sub(dates[len(dates)-1]).Hours() / 168

	return statsForDiagram, weeks
}
