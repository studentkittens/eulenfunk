package main

import (
	"os"

	"github.com/studentkittens/eulenfunk/ambilight"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "mpd-host,H",
			Value:  "localhost",
			Usage:  "MPD host to connect to",
			EnvVar: "MPD_HOST",
		},
		cli.IntFlag{
			Name:   "mpd-port,p",
			Value:  6600,
			Usage:  "MPD port to connect to",
			EnvVar: "MPD_PORT",
		},
	}

	app.Commands = []cli.Command{{
		Name:  "ambilight",
		Usage: "Control the ambilight feature",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "music-dir,m",
				Value:  "",
				Usage:  "MPD port to connect to",
				EnvVar: "AMBI_MUSIC_DIR",
			},
			cli.StringFlag{
				Name:   "mood-dir,d",
				Value:  "",
				Usage:  "MPD port to connect to",
				EnvVar: "AMBI_MOOD_DIR",
			},
			cli.StringFlag{
				Name:   "binary,b",
				Value:  "catlight",
				Usage:  "MPD port to connect to",
				EnvVar: "AMBI_BINARY",
			},
			cli.BoolFlag{
				Name:  "update-mood-db,u",
				Usage: "Update the mood database and exit afterwards",
			},
		},
		Action: func(ctx *cli.Context) error {
			return ambilight.RunDaemon(&ambilight.Config{
				Host:               ctx.GlobalString("mpd-host"),
				Port:               ctx.GlobalInt("mpd-port"),
				MusicDir:           ctx.String("music-dir"),
				MoodDir:            ctx.String("mood-dir"),
				UpdateMoodDatabase: ctx.Bool("update-mood-db"),
				BinaryName:         ctx.String("binary"),
			})
		},
	},
	}

	app.Run(os.Args)
}
