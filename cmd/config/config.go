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

package config

import (
	"colligendis/cmd/common"
	"fmt"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func GetConfigCommand(flags *common.ColligendisFlags, db *gorm.DB) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "config",
		Short:   "Various utility settings",
		Long:    `Allows you to work with the configuration of the utility`,
		Example: `colligendis config`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Start loading...")
		},
	}
	cmd.AddCommand(GetGitHubCommand(flags))
	return cmd
}
