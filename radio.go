package main

import (
	"fmt"
	"os"

	"github.com/studentkittens/eulenfunk/ambilight"
	"github.com/studentkittens/eulenfunk/display"
	"github.com/studentkittens/eulenfunk/mpdinfo"
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
		Name:  "mpdinfo",
		Usage: "Send mpd infos to the display server",
		Flags: []cli.Flag{},
		Action: func(ctx *cli.Context) error {
			return mpdinfo.Run(&mpdinfo.Config{
				"localhost", 6600,
				"localhost", 7778,
			})
		},
	}, {
		Name:  "display",
		Usage: "Display server",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "host,H",
				Value:  "localhost",
				Usage:  "Display server hostname",
				EnvVar: "DISPLAY_HOST",
			},
			cli.IntFlag{
				Name:   "port,p",
				Value:  7778,
				Usage:  "Display server port",
				EnvVar: "DISPLAY_PORT",
			},
			cli.BoolFlag{
				Name:  "dump,d",
				Usage: "Dumpy display output to terminal",
			},
			cli.BoolFlag{
				Name:  "quit,q",
				Usage: "Quit the display server",
			},
			cli.BoolFlag{
				Name:  "update,u",
				Usage: "For --dump; updates output when given",
			},
			cli.StringFlag{
				Name:  "window,w",
				Value: "1",
				Usage: "Which window to show/modify",
			},
			cli.IntFlag{
				Name:   "width",
				Value:  20,
				Usage:  "Width of each LCD display line",
				EnvVar: "DISPLAY_LCD_WITH",
			},
			cli.IntFlag{
				Name:   "height",
				Value:  4,
				Usage:  "Height of the LCD display in lines",
				EnvVar: "DISPLAY_LCD_HEIGHT",
			},
		},
		Action: func(ctx *cli.Context) error {
			cfg := &display.Config{
				Host:   ctx.String("host"),
				Port:   ctx.Int("port"),
				Width:  ctx.Int("width"),
				Height: ctx.Int("height"),
			}

			if ctx.Bool("dump") {
				return display.RunDumpClient(cfg, ctx.String("window"), ctx.Bool("update"))
			}

			return display.RunInputClient(cfg, ctx.Bool("quit"), ctx.String("window"))
		},
		Subcommands: []cli.Command{{
			Name:  "server",
			Usage: "Start the display server",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "driver",
					Value:  "cat",
					Usage:  "Driver program that takes the display output",
					EnvVar: "DISPLAY_DRIVER",
				},
			},
			Action: func(ctx *cli.Context) error {
				return display.RunDaemon(&display.Config{
					Host:         ctx.Parent().String("host"),
					Port:         ctx.Parent().Int("port"),
					Width:        ctx.Parent().Int("width"),
					Height:       ctx.Parent().Int("height"),
					DriverBinary: ctx.String("driver"),
				})
			},
		},
		},
	}, {
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
			musicDir := ctx.String("music-dir")
			moodyDir := ctx.String("mood-dir")

			if musicDir == "" || moodyDir == "" {
				return fmt.Errorf("Need both --music-dir and --mood-dir")
			}

			return ambilight.RunDaemon(&ambilight.Config{
				Host:               ctx.GlobalString("mpd-host"),
				Port:               ctx.GlobalInt("mpd-port"),
				MusicDir:           musicDir,
				MoodDir:            moodyDir,
				UpdateMoodDatabase: ctx.Bool("update-mood-db"),
				BinaryName:         ctx.String("binary"),
			})
		},
	},
	}

	app.Run(os.Args)
}