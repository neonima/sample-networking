package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/neonima/sample-networking/pkg/server/tcp"

	"github.com/francoispqt/onelog"
	"github.com/urfave/cli/v2"
)

var (
	Appname = "ALO"
	logger  *onelog.Logger
)

func init() {
	logger = onelog.New(
		os.Stdout,
		onelog.ERROR|onelog.FATAL|onelog.WARN,
	)
}

func main() {
	if err := RunCLI(os.Args); err != nil {
		logger.Error(err.Error())
	}
}

func RunCLI(args []string) error {
	app := cli.NewApp()
	app.Name = Appname
	app.Usage = "Alo game server"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "verbose",
			Usage:   "verbose",
			EnvVars: []string{"VERBOSE"},
			Aliases: []string{"v"},
		},
		&cli.Int64Flag{
			Name:    "port",
			Usage:   "set the port to listen to",
			EnvVars: []string{"SERVER_PORT"},
			Aliases: []string{"p"},
			Value:   8080,
		},
	}
	app.Action = CMD

	return app.Run(args)
}

func CMD(cli *cli.Context) error {
	if cli.Bool("verbose") {
		logger = onelog.New(
			os.Stdout,
			onelog.ALL,
		)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	srv := tcp.NewTCPServer(cli.Int("port")).WithLogger(logger)
	if err := srv.Start(); err != nil {
		return err
	}

	defer func() {
		if err := srv.Close(); err != nil {
			logger.Error(err.Error())
		}
	}()

	<-sig
	return nil
}
