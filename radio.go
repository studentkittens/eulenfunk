package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/studentkittens/eulenfunk/ambilight"
	"github.com/studentkittens/eulenfunk/automount"
	"github.com/studentkittens/eulenfunk/display"
	"github.com/studentkittens/eulenfunk/lightd"
	"github.com/studentkittens/eulenfunk/ui"
	"github.com/studentkittens/eulenfunk/ui/mpd"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
)

const (
	// DefaultWidth of our LCD display in runes
	DefaultWidth = 20
	// DefaultHeight of our LCD display in lines
	DefaultHeight = 4
)

//////////////////////////

func handleInfo(ctx *cli.Context, dropout context.Context) error {
	client, err := mpd.NewClient(&mpd.Config{
		MPDHost:     ctx.String("mpd-host"),
		MPDPort:     ctx.Int("mpd-port"),
		DisplayHost: ctx.String("display-host"),
		DisplayPort: ctx.Int("display-port"),
	}, dropout)

	if err != nil {
		log.Printf("Failed to create mpd client: %v", err)
		return err
	}

	client.Run()
	return nil
}

func handleUI(ctx *cli.Context, dropout context.Context) error {
	return ui.Run(&ui.Config{
		Width:         ctx.GlobalInt("width"),
		Height:        ctx.GlobalInt("height"),
		DisplayHost:   ctx.String("display-host"),
		DisplayPort:   ctx.Int("display-port"),
		MPDHost:       ctx.String("mpd-host"),
		MPDPort:       ctx.Int("mpd-port"),
		AmbilightHost: ctx.String("ambi-host"),
		AmbilightPort: ctx.Int("ambi-port"),
		LightdHost:    ctx.String("ambi-host"),
		LightdPort:    ctx.Int("ambi-port"),
	}, dropout)
}

func handleLightd(ctx *cli.Context, dropout context.Context) error {
	cfg := &lightd.Config{
		Host:         ctx.String("lightd-host"),
		Port:         ctx.Int("lightd-port"),
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

		if ctx.Bool("lock") {
			if err := locker.Lock(); err != nil {
				log.Printf("lightd-lock failed: %v", err)
			}
		} else {
			if err := locker.Unlock(); err != nil {
				log.Printf("lightd-unlock failed: %v", err)
			}
		}

		return locker.Close()
	}

	return nil
}

func handleDisplayClient(ctx *cli.Context, dropout context.Context) error {
	cfg := &display.Config{
		Host:   ctx.String("display-host"),
		Port:   ctx.Int("display-port"),
		Width:  ctx.GlobalInt("width"),
		Height: ctx.GlobalInt("height"),
	}

	if ctx.Bool("dump") {
		return display.DumpClient(
			cfg, dropout,
			ctx.String("window"),
			ctx.Bool("update"),
		)
	}

	return display.InputClient(
		cfg, dropout,
		ctx.Bool("quit"),
		ctx.String("window"),
	)
}

func handleDisplayServer(ctx *cli.Context, dropout context.Context) error {
	return display.Run(&display.Config{
		Host:         ctx.Parent().String("display-host"),
		Port:         ctx.Parent().Int("display-port"),
		Width:        ctx.GlobalInt("width"),
		Height:       ctx.GlobalInt("height"),
		NoEncoding:   ctx.Bool("no-encoding"),
		DriverBinary: ctx.String("driver"),
	}, dropout)
}

func handleAmbilightCommand(ctx *cli.Context, cfg *ambilight.Config) (bool, error) {
	on, off, quit, state := ctx.Bool("on"), ctx.Bool("off"), ctx.Bool("quit"), ctx.Bool("state")
	if !on && !off && !quit && !state {
		return false, nil
	}

	client, err := ambilight.NewClient(cfg)
	if err != nil {
		log.Printf("Failed to connect to ambilightd: %v", err)
		return true, err
	}

	switch {
	case on, off:
		return true, client.Enable(on)
	case state:
		enabled, err := client.Enabled()
		if err != nil {
			log.Printf("Failed to get state: %v", err)
			return true, err
		}

		fmt.Printf("%t", enabled)
		return true, nil
	case quit:
		return true, client.Quit()
	}

	return true, nil
}

func handleAmbilight(ctx *cli.Context, dropout context.Context) error {
	musicDir := ctx.String("music-dir")
	moodyDir := ctx.String("mood-dir")

	cfg := &ambilight.Config{
		AmbiHost:           ctx.String("ambi-host"),
		AmbiPort:           ctx.Int("ambi-port"),
		MPDHost:            ctx.String("mpd-host"),
		MPDPort:            ctx.Int("mpd-port"),
		LightdHost:         ctx.String("lightd-host"),
		LightdPort:         ctx.Int("lightd-port"),
		UpdateMoodDatabase: ctx.Bool("update-mood-db"),
		BinaryName:         ctx.String("driver"),
		MusicDir:           musicDir,
		MoodDir:            moodyDir,
	}

	handled, err := handleAmbilightCommand(ctx, cfg)
	if err != nil {
		return err
	}

	if handled {
		return nil
	}

	if musicDir == "" || moodyDir == "" {
		log.Printf("Need both --music-dir and --mood-dir")
		return nil
	}

	return ambilight.Run(cfg, dropout)
}

func handleAutomount(ctx *cli.Context, dropout context.Context) error {
	cfg := &automount.Config{
		AutomountHost: ctx.String("automount-host"),
		AutomountPort: ctx.Int("automount-port"),
		MPDHost:       ctx.String("mpd-host"),
		MPDPort:       ctx.Int("mpd-port"),
		MusicDir:      ctx.String("music-dir"),
	}

	if ctx.Bool("quit") {
		return automount.WithClient(cfg, func(cl *automount.Client) error {
			return cl.Quit()
		})
	}

	device := ctx.String("device")
	label := ctx.String("label")
	unmount := ctx.Bool("unmount")

	if device != "" || label != "" {
		if device == "" {
			log.Printf("Need --device for mounting and unmounting")
			return nil
		}

		if unmount {
			return automount.WithClient(cfg, func(cl *automount.Client) error {
				return cl.Unmount(device)
			})
		}

		if label == "" {
			log.Printf("Need --label for mounting")
			return nil
		}

		return automount.WithClient(cfg, func(cl *automount.Client) error {
			return cl.Mount(device, label)
		})
	}

	return automount.Run(cfg, dropout)
}

//////////////////////////

type handlerFunc func(ctx *cli.Context, dropout context.Context) error

func withCancelCtx(dropout context.Context, fn handlerFunc) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		return fn(ctx, dropout)
	}
}

func concat(packs ...[]cli.Flag) []cli.Flag {
	result := []cli.Flag{}
	for _, pack := range packs {
		result = append(result, pack...)
	}

	return result
}

func main() {
	dropout, cancel := context.WithCancel(context.Background())

	go func() {
		// Handle Interrupt:
		signals := make(chan os.Signal)
		signal.Notify(signals, os.Interrupt)

		<-signals
		cancel()
	}()

	//////////////////////////
	// APPLICATION METADATA //
	//////////////////////////

	app := cli.NewApp()
	app.Author = "Waldsoft"
	app.Email = "sahib@online.de"
	app.Usage = "Control the higher level eulenfunk services"

	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:   "width",
			Value:  DefaultWidth,
			Usage:  "Width of the LCD screen",
			EnvVar: "LCD_HEIGHT",
		},
		cli.IntFlag{
			Name:   "height",
			Value:  DefaultHeight,
			Usage:  "Height of the LCD screen",
			EnvVar: "LCD_HEIGHT",
		},
	}

	//////////////////////
	// CONNECTION FLAGS //
	//////////////////////

	mpdNetFlags := []cli.Flag{
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

	displaydNetFlags := []cli.Flag{
		cli.StringFlag{
			Name:   "display-host",
			Value:  "localhost",
			Usage:  "Display server hostname",
			EnvVar: "DISPLAY_HOST",
		},
		cli.IntFlag{
			Name:   "display-port",
			Value:  7777,
			Usage:  "Display server port",
			EnvVar: "DISPLAY_PORT",
		},
	}

	ambiNetFlags := []cli.Flag{
		cli.StringFlag{
			Name:   "ambi-host",
			Value:  "localhost",
			Usage:  "Host of the internal control server",
			EnvVar: "AMBI_HOST",
		},
		cli.IntFlag{
			Name:   "ambi-port",
			Value:  4444,
			Usage:  "Port of the internal control server",
			EnvVar: "AMBI_PORT",
		},
	}

	lightdNetFlags := []cli.Flag{
		cli.StringFlag{
			Name:   "lightd-host",
			Value:  "localhost",
			Usage:  "Host of the lightd server",
			EnvVar: "LIGHTD_HOST",
		},
		cli.IntFlag{
			Name:   "lightd-port",
			Value:  3333,
			Usage:  "Port of the lightd server",
			EnvVar: "LIGHTD_PORT",
		},
	}

	////////////////////////
	// ACTUAL SUBCOMMANDS //
	////////////////////////

	app.Commands = []cli.Command{{
		Name:   "info",
		Usage:  "Send mpd infos to the display server on the `mpd` window",
		Action: withCancelCtx(dropout, handleInfo),
		Flags:  concat(mpdNetFlags, displaydNetFlags),
	}, {
		Name:   "ui",
		Usage:  "Handle window rendering and input control",
		Action: withCancelCtx(dropout, handleUI),
		Flags:  concat(displaydNetFlags, mpdNetFlags, ambiNetFlags, lightdNetFlags),
	}, {
		Name:   "automount",
		Usage:  "Control the automount for usb sticks filled with music",
		Action: withCancelCtx(dropout, handleAutomount),
		Flags: concat(mpdNetFlags, []cli.Flag{
			cli.StringFlag{
				Name:   "automount-host",
				Value:  "localhost",
				Usage:  "The host on which the control daemon listens on",
				EnvVar: "AUTOMOUNT_HOST",
			},
			cli.IntFlag{
				Name:   "automount-port",
				Value:  5555,
				Usage:  "The port on which the control daemon listens on",
				EnvVar: "AUTOMOUNT_PORT",
			},
			cli.StringFlag{
				Name:  "device,d",
				Value: "",
				Usage: "The device (under /dev) to mount; absolute path.",
			},
			cli.StringFlag{
				Name:  "label,l",
				Value: "",
				Usage: "Which label the device has",
			},
			cli.StringFlag{
				Name:   "music-dir,m",
				Value:  "",
				Usage:  "The root directory of the music collection",
				EnvVar: "AUTOMOUNT_MUSIC_DIR",
			},
			cli.BoolFlag{
				Name:  "unmount,u",
				Usage: "Unmount the device",
			},
		}),
	}, {
		Name:   "lightd",
		Usage:  "Utility server to lock the led ownage and enable nice atomic effects",
		Action: withCancelCtx(dropout, handleLightd),
		Flags: concat(lightdNetFlags, []cli.Flag{
			cli.StringFlag{
				Name:   "driver,d",
				Value:  "catlight",
				Usage:  "Which driver binary to use to send colors to",
				EnvVar: "LIGHTD_DRIVER",
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
		}),
	}, {

		Name:   "display",
		Usage:  "Display manager",
		Action: withCancelCtx(dropout, handleDisplayClient),
		Flags: concat(displaydNetFlags, []cli.Flag{
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
		}),
		Subcommands: []cli.Command{{
			Name:   "server",
			Usage:  "Start the display server",
			Action: withCancelCtx(dropout, handleDisplayServer),
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "driver",
					Value:  "cat",
					Usage:  "Driver program that takes the display output",
					EnvVar: "DISPLAY_DRIVER",
				},
				cli.BoolFlag{
					Name:   "no-encoding",
					Usage:  "Disables special encoding for optimized LCD output for testing",
					EnvVar: "DISPLAY_NO_ENCODING",
				},
			},
		},
		},
	}, {
		Name:   "ambilight",
		Usage:  "Control the ambilight feature",
		Action: withCancelCtx(dropout, handleAmbilight),
		Flags: concat(mpdNetFlags, lightdNetFlags, ambiNetFlags, []cli.Flag{
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
				Name:  "state",
				Usage: "Print the current state of ambilight (on/off)",
			},
			cli.BoolFlag{
				Name:  "quit",
				Usage: "Quit the ambilight daemon",
			},
		}),
	},
	}

	if err := app.Run(os.Args); err != nil {
		log.Printf("eulenfunk failed: %v", err)
		os.Exit(1)
	}
}
