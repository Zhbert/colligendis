package version

import (
	"colligendis/cmd/common"
	"fmt"
	"github.com/spf13/cobra"
)

var (
	Version = "dev"
)

func GetVersionCommand(flags *common.ColligendisFlags) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "version",
		Short:   "Displaying version of utility",
		Example: `colligendis version`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Version)
		},
	}
	return cmd
}