package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/powerpuffpenguin/webpc/cmd/internal/master"
	"github.com/powerpuffpenguin/webpc/configure"
	"github.com/powerpuffpenguin/webpc/db/manipulator"
	"github.com/powerpuffpenguin/webpc/logger"
	"github.com/powerpuffpenguin/webpc/sessions"
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
		Use:   `master`,
		Short: `run as master`,
		Run: func(cmd *cobra.Command, args []string) {
			// load configure
			cnf := configure.DefaultConfigure()
			e := cnf.Load(filename)
			if e != nil {
				log.Fatalln(e)
			}
			if addr != `` {
				cnf.HTTP.Addr = addr
			}
			if test {
				fmt.Println(cnf)
				return
			}
			// init logger
			e = logger.Init(basePath, &cnf.Logger)
			if e != nil {
				log.Fatalln(e)
			}

			// init db
			manipulator.Init(&cnf.DB)
			sessions.Init(&cnf.Session)

			master.Run(&cnf.HTTP, debug)
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&filename, `config`,
		`c`,
		utils.Abs(basePath, filepath.Join(`etc`, `master.jsonnet`)),
		`configure file`,
	)
	flags.StringVarP(&addr, `addr`,
		`a`,
		``,
		`listen address`,
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
