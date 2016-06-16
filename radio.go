package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/studentkittens/eulenfunk/ambilight"
	"github.com/studentkittens/eulenfunk/display"
	"github.com/studentkittens/eulenfunk/lightd"
	"github.com/studentkittens/eulenfunk/ui"
	"github.com/studentkittens/eulenfunk/ui/mpdinfo"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
)

func main() {
	killCtx, cancel := context.WithCancel(context.Background())

	go func() {
		// Handle Interrupt:
		signals := make(chan os.Signal)
		signal.Notify(signals, os.Interrupt)

		<-signals
		cancel()
	}()

	app := cli.NewApp()
	app.Author = "Waldsoft"
	app.Email = "sahib@online.de"
	app.Usage = "Control the higher level eulenfunk services"

	app.Version = "0.0.1"
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
		Flags: []cli.Flag{}, // TODO
		Action: func(ctx *cli.Context) error {
			// TODO: check for running mpd (also for ambilight/mpdinfo)
			return mpdinfo.Run(&mpdinfo.Config{
				"localhost", 6600, // MPD Config
				"localhost", 7778, // Display server config
			}, killCtx, nil)
		},
	}, {
		Name:  "ui",
		Usage: "Handle window rendering and input control",
		Flags: []cli.Flag{}, // TODO
		Action: func(ctx *cli.Context) error {
			log.Printf("Starting ui...")
			return menu.Run(killCtx)
		},
	}, {
		Name:  "lightd",
		Usage: "Utility server to lock the led ownage and enable nice atomic effects",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "driver,d",
				Value:  "catlight",
				Usage:  "Which driver binary to use to send colors to",
				EnvVar: "LIGHTD_DRIVER",
			},
			cli.StringFlag{
				Name:   "host,H",
				Value:  "localhost",
				Usage:  "lightd-host to connect to",
				EnvVar: "LIGHTD_HOST",
			},
			cli.IntFlag{
				Name:   "port,p",
				Value:  3333,
				Usage:  "lightd port to connect to",
				EnvVar: "LIGHTD_PORT",
			},
			cli.StringFlag{
				Name:  "send,s",
				Usage: "Send an effect",
				Value: "",
			},
			cli.BoolFlag{
				Name:  "lock,l",
				Usage: "Lock the light",
			},
			cli.BoolFlag{
				Name:  "unlock,u",
				Usage: "Unlock the lights",
			},
		},
		Action: func(ctx *cli.Context) error {
			cfg := &lightd.Config{
				Host:         ctx.String("host"),
				Port:         ctx.Int("port"),
				DriverBinary: ctx.String("driver"),
			}

			if effect := ctx.String("send"); effect != "" {
				return lightd.Send(cfg, effect)
			}

			if ctx.Bool("lock") || ctx.Bool("unlock") {
				locker, err := lightd.NewLocker(cfg)
				if err != nil {
					return err
				}

				defer locker.Close()

				if ctx.Bool("lock") {
					locker.Lock()
				} else {
					locker.Unlock()
				}

				return nil
			}

			return lightd.Run(cfg)
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
				}, killCtx)
			},
		},
		},
	}, {
		Name:  "ambilight",
		Usage: "Control the ambilight feature",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "host,H",
				Value:  "localhost",
				Usage:  "Host of the internal control server",
				EnvVar: "AMBI_HOST",
			},
			cli.IntFlag{
				Name:   "port,p",
				Value:  4444,
				Usage:  "Port of the internal control server",
				EnvVar: "AMBI_PORT",
			},
			cli.StringFlag{
				Name:   "music-dir,m",
				Value:  "",
				Usage:  "Root directory where mpd thins the music is",
				EnvVar: "AMBI_MUSIC_DIR",
			},
			cli.StringFlag{
				Name:   "mood-dir,i",
				Value:  "",
				Usage:  "Where the mood files are stored",
				EnvVar: "AMBI_MOOD_DIR",
			},
			cli.StringFlag{
				Name:   "driver,b",
				Value:  "catlight",
				Usage:  "Which driver to output the RGB values on",
				EnvVar: "AMBI_DRIVER",
			},
			cli.BoolFlag{
				Name:  "update-mood-db,u",
				Usage: "Update the mood database and exit afterwards",
			},
			cli.BoolFlag{
				Name:  "on",
				Usage: "Enable the ambilight if it runs elsewhere",
			},
			cli.BoolFlag{
				Name:  "off",
				Usage: "Disable the ambilight temporarily if it runs elsewhere",
			},
			cli.BoolFlag{
				Name:  "quit",
				Usage: "Quit the ambilight daemon",
			},
		},
		Action: func(ctx *cli.Context) error {
			musicDir := ctx.String("music-dir")
			moodyDir := ctx.String("mood-dir")

			cfg := &ambilight.Config{
				Host:               ctx.String("host"),
				Port:               ctx.Int("port"),
				MPDHost:            ctx.GlobalString("mpd-host"),
				MPDPort:            ctx.GlobalInt("mpd-port"),
				MusicDir:           musicDir,
				MoodDir:            moodyDir,
				UpdateMoodDatabase: ctx.Bool("update-mood-db"),
				BinaryName:         ctx.String("driver"),
			}

			on, off, quit := ctx.Bool("on"), ctx.Bool("off"), ctx.Bool("quit")
			if on || off || quit {
				client, err := ambilight.NewClient(cfg)
				if err != nil {
					log.Printf("Failed to connect to ambilightd: %v", err)
					return err
				}

				switch {
				case on:
					return client.Enable(true)
				case off:
					return client.Enable(false)
				case quit:
					log.Printf("do quit")
					return client.Quit()
				}
			}

			if musicDir == "" || moodyDir == "" {
				log.Printf("Need both --music-dir and --mood-dir")
				return nil
			}

			return ambilight.RunDaemon(cfg, killCtx)
		},
	},
	}

	app.Run(os.Args)
}
