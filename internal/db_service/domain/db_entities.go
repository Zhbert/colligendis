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

package domain

import (
	"gorm.io/gorm"
	"time"
)

type HabrArticle struct {
	gorm.Model
	ID                uint `gorm:"primary_key;auto_increment;not_null"`
	HabrNumber        int
	Name              string
	Url               string
	AuthorID          int
	Author            HabrAuthor
	Hubs              []HabrHub `gorm:"many2many:habr_article_hubs;"`
	DateOfCreation    time.Time
	TimeOfCreation    time.Time
	DateOfPublication time.Time
	TimeOfPublication time.Time
}

type HabrHub struct {
	gorm.Model
	ID   uint `gorm:"primary_key;auto_increment;not_null"`
	Name string
}

type HabrAuthor struct {
	gorm.Model
	ID   uint `gorm:"primary_key;auto_increment;not_null"`
	Name string
}

type HabrStats struct {
	gorm.Model
	ID            uint `gorm:"primary_key;auto_increment;not_null"`
	HabrArticleID uint
	HabrArticle   HabrArticle
	DateOfStats   time.Time
	LikesAll      int
	LikesUp       int
	LikesDown     int
	Likes         int
	Comments      int
	Saves         int
	Views         int
}
