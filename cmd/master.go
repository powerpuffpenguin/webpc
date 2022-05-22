package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/powerpuffpenguin/webpc/cmd/internal/master"
	"github.com/powerpuffpenguin/webpc/configure"
	"github.com/powerpuffpenguin/webpc/db/manipulator"
	"github.com/powerpuffpenguin/webpc/logger"
	"github.com/powerpuffpenguin/webpc/sessionid"
	"github.com/powerpuffpenguin/webpc/single/logger/db"
	"github.com/powerpuffpenguin/webpc/single/mount"
	"github.com/powerpuffpenguin/webpc/single/upgrade"
	"github.com/powerpuffpenguin/webpc/utils"
	"github.com/spf13/cobra"
)

func init() {
	var (
		filename    string
		debug, test bool
		basePath    = utils.BasePath()

		addr, vnc string
		noupgrade bool
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
			if vnc != `` {
				cnf.System.VNC = vnc
			}
			if test {
				fmt.Println(cnf)
				return
			}
			// init logger
			e = logger.Init(basePath, true, &cnf.Logger)
			if e != nil {
				log.Fatalln(e)
			}
			db.Init(cnf.Logger.Filename)

			// init db
			manipulator.Init(&cnf.DB)
			sessionid.Init(&cnf.Session)
			if cnf.System.Enable {
				mount.Init(cnf.System.Mount)
			}
			if !noupgrade {
				go upgrade.DefaultUpgrade().Serve()
			}
			master.Run(&cnf.HTTP, &cnf.System, debug)
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
	flags.StringVarP(&vnc, `vnc`,
		`v`,
		``,
		`connect vnc address`,
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
	flags.BoolVar(&noupgrade, `no-upgrade`,
		false,
		`disable automatic upgrades`,
	)
	rootCmd.AddCommand(cmd)
}
