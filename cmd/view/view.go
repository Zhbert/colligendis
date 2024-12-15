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
	"colligendis/cmd/common"
	"colligendis/internal/db_service"
	"fmt"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"log"
)

func GetViewCommand(flags *common.ColligendisFlags, db *gorm.DB) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "view",
		Short:   "Displaying various information in the terminal",
		Long:    `Allows you to view various information and lists by articles, authors, and so on`,
		Example: `colligendis view`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Select the display type: habr, github")
		},
	}
	cmd.AddCommand(getViewHabrArticlesCommand(flags, db))
	return cmd
}

func getViewHabrArticlesCommand(flags *common.ColligendisFlags, db *gorm.DB) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "habr",
		Short:   "Displaying Habr information in the terminal",
		Long:    `Allows you to view various information and lists by articles, authors, and so on`,
		Example: `colligendis view habr`,
		Run: func(cmd *cobra.Command, args []string) {
			if flags.ViewHabrArticles {
				switch flags.SortType {
				case "":
					viewAllHabrArticles("", db)
				case "name":
					viewAllHabrArticles("name", db)
				case "date":
					viewAllHabrArticles("date", db)
				default:
					log.Fatal("You need to specify the sort type: name, date")
				}
			} else if flags.ViewHabrAuthors {
				switch flags.SortType {
				case "":
					viewAllHabrAuthors("", db)
				case "name":
					viewAllHabrAuthors("name", db)
				default:
					log.Fatal("You need to specify the sort type: name, date")
				}
			} else if flags.ViewStatsDates {
				dates, _ := db_service.GetAllDatesOfStats(db)
				for _, st := range dates {
					fmt.Println(st)
				}
			} else {
				log.Println("Select the display type: articles, authors")
			}
		},
	}
	cmd.Flags().BoolVarP(&flags.ViewHabrArticles, "articles", "", false, "Show all articles")
	cmd.Flags().BoolVarP(&flags.ViewHabrAuthors, "authors", "", false, "Show all authors")
	cmd.Flags().BoolVarP(&flags.ViewStatsDates, "dates-of-stats", "", false, "Show all dates from BD with stats")
	cmd.Flags().StringVarP(&flags.SortType, "sort", "s", "", "Output sorting filter")
	return cmd
}
