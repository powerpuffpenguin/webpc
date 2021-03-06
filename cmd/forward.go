package cmd

import (
	"strings"

	"github.com/powerpuffpenguin/webpc/cmd/internal/forward"
	"github.com/powerpuffpenguin/webpc/logger"
	"github.com/spf13/cobra"
)

func init() {
	var (
		insecure                                 bool
		url, listen, remote, user, password, log string
		heart                                    int
	)

	cmd := &cobra.Command{
		Use:   `forward`,
		Short: `port forwarding`,
		Run: func(cmd *cobra.Command, args []string) {
			logger.InitConsole(strings.ToLower(strings.TrimSpace(log)))
			forward.Run(insecure,
				strings.TrimSpace(url),
				strings.TrimSpace(listen), strings.TrimSpace(remote),
				user, password,
				heart, false,
			)
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&url, `websocket`,
		`w`,
		``,
		`connect websocket url`,
	)
	flags.BoolVarP(&insecure, `insecure`,
		`k`,
		false,
		`allow insecure server connections when using SSL`,
	)
	flags.IntVar(&heart, `heart`,
		0,
		`heart interval in seconds, if < 1 not send heart`,
	)
	flags.StringVarP(&listen, `listen`,
		`l`,
		`:10000`,
		`local listen address`,
	)
	flags.StringVarP(&remote, `remote`,
		`r`,
		``,
		`remote connect address`,
	)
	flags.StringVarP(&user, `user`,
		`u`,
		``,
		`user name`,
	)
	flags.StringVarP(&password, `password`,
		`p`,
		``,
		`user password`,
	)
	flags.StringVar(&log, `log`,
		`info`,
		`log level [debug info warn error dpanic panic fatal]`,
	)

	rootCmd.AddCommand(cmd)
}
