package cmd

import (
	"path/filepath"

	"github.com/powerpuffpenguin/webpc/cmd/internal/slave"
	"github.com/powerpuffpenguin/webpc/utils"

	"github.com/spf13/cobra"
)

func init() {
	var (
		filename    string
		debug, test bool
		basePath    = utils.BasePath()

		addr string
	)

	cmd := &cobra.Command{
		Use:   `slave`,
		Short: `run as slave`,
		Run: func(cmd *cobra.Command, args []string) {
			slave.Run()
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&filename, `config`,
		`c`,
		utils.Abs(basePath, filepath.Join(`etc`, `slave.jsonnet`)),
		`configure file`,
	)
	flags.StringVarP(&addr, `addr`,
		`a`,
		``,
		`connect address`,
	)

	flags.BoolVarP(&debug, `debug`,
		`d`,
		false,
		`run as debug`,
	)
	flags.BoolVarP(&test, `test`,
		`t`,
		false,
		`test configure`,
	)
	rootCmd.AddCommand(cmd)
}
