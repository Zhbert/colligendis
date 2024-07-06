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
	"colligendis/cmd/common"
	"colligendis/internal/common/structs"
	"colligendis/internal/db_service"
	"colligendis/internal/latex_service"
	"github.com/spf13/cobra"
	"log"
	"strconv"
)

const (
	collectLong = "Displays statistics of changes since the last time information was added by the load command.\n\n" +
		"First, you need to load statistics into the database using the load command."
)

func GetCollectCommand(flags *common.ColligendisFlags, tmpls []structs.TemplateStruct) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "collect",
		Short:   "Сollect statistics from loaded sources.",
		Long:    collectLong,
		Example: `colligendis collect`,
		Run: func(cmd *cobra.Command, args []string) {
			if flags.Habr {
				switch flags.File {
				case true:
					latex_service.GenerateLaTeXFiles(tmpls, flags)
				case false:
					log.Printf("Total habr views: %s", strconv.Itoa(db_service.GetHabrViewsCount()))
					if flags.Full {
						getFullHabrViewsCount(flags.Limit, flags.SortType)
					}
				}
			}
		},
	}
	cmd.Flags().BoolVarP(&flags.All, "all", "", false, "Starts the collection of statistics from all resources")
	cmd.Flags().BoolVarP(&flags.Habr, "habr", "", false, "Starts the collection of statistics from Habr")
	cmd.Flags().BoolVarP(&flags.GitHub, "github", "", false, "Starts the collection of statistics from GitHub resources")
	cmd.Flags().BoolVarP(&flags.Full, "full", "f", false, "Show full statistic")
	cmd.Flags().BoolVarP(&flags.File, "file", "", false, "Generate PDF file")
	cmd.Flags().IntVarP(&flags.Limit, "limit", "l", 0, "Limit of table items")
	cmd.Flags().StringVarP(&flags.SortType, "sort", "s", "", "Output sorting filter")
	return cmd
}
