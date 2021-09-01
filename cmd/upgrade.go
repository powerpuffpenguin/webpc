package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/powerpuffpenguin/webpc/logger"
	"github.com/powerpuffpenguin/webpc/single/upgrade"
	"github.com/powerpuffpenguin/webpc/version"
	"github.com/spf13/cobra"
)

func init() {
	var (
		yes   bool
		loglv string
	)

	cmd := &cobra.Command{
		Use:   `upgrade`,
		Short: `upgrade to the latest version`,
		Run: func(cmd *cobra.Command, args []string) {
			logger.InitConsole(strings.ToLower(strings.TrimSpace(loglv)))
			upgraded, ver, e := upgrade.DefaultUpgrade().Do(yes)
			if e == nil {
				if upgraded {
					fmt.Println(`upgrade success:`, version.Version, `->`, ver)
				} else {
					fmt.Println(`already the latest version:`, version.Version)
				}
			} else {
				log.Fatalln(e)
			}
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&loglv, `log`,
		`info`,
		`log level [debug info warn error dpanic panic fatal]`,
	)
	flags.BoolVarP(&yes, `yes`,
		`y`,
		false,
		`automatic yes to prompts`,
	)
	rootCmd.AddCommand(cmd)
}
