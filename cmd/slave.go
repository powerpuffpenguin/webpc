package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/powerpuffpenguin/webpc/cmd/internal/slave"
	"github.com/powerpuffpenguin/webpc/configure"
	"github.com/powerpuffpenguin/webpc/db/manipulator"
	"github.com/powerpuffpenguin/webpc/logger"
	"github.com/powerpuffpenguin/webpc/single/logger/db"
	"github.com/powerpuffpenguin/webpc/single/mount"
	"github.com/powerpuffpenguin/webpc/utils"

	"github.com/spf13/cobra"
)

func init() {
	var (
		filename              string
		insecure, debug, test bool
		basePath              = utils.BasePath()

		url string
	)

	cmd := &cobra.Command{
		Use:   `slave`,
		Short: `run as slave`,
		Run: func(cmd *cobra.Command, args []string) {
			// load configure
			cnf := configure.DefaultSlave()
			e := cnf.Load(filename)
			if e != nil {
				log.Fatalln(e)
			}
			if url != `` {
				cnf.Connect.URL = url
			}
			if insecure {
				cnf.Connect.Insecure = insecure
			}
			if test {
				fmt.Println(cnf)
				return
			}
			// init logger
			e = logger.Init(basePath, false, &cnf.Logger)
			if e != nil {
				log.Fatalln(e)
			}
			db.Init(cnf.Logger.Filename)

			// init db
			manipulator.Init(&cnf.DB)
			if cnf.System.Enable {
				mount.Init(cnf.System.Mount)
			}
			slave.Run(&cnf.Connect, &cnf.System, debug)
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&filename, `config`,
		`c`,
		utils.Abs(basePath, filepath.Join(`etc`, `slave.jsonnet`)),
		`configure file`,
	)
	flags.StringVarP(&url, `url`,
		`u`,
		``,
		`connect websocket url`,
	)
	flags.BoolVarP(&insecure, `insecure`,
		`k`,
		false,
		`allow insecure server connections when using SSL`,
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
