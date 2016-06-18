package ui

import (
	"log"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/studentkittens/eulenfunk/ambilight"
	"github.com/studentkittens/eulenfunk/display"
	"github.com/studentkittens/eulenfunk/ui/mpd"
)

func sysCommand(name string, args ...string) func() error {
	return func() error {
		return exec.Command(name, args...).Run()
	}
}

func boolToGlyph(b bool) string {
	if b {
		return "✓"
	} else {
		return "×"
	}
}

func ambilightChangeState(cfg *Config, enable bool) error {
	host, port := cfg.AmbilightHost, cfg.AmbilightPort
	return ambilight.WithClient(host, port, func(client *ambilight.Client) error {
		return client.Enable(enable)
	})
}

func ambilightIsEnabled(cfg *Config) (enabled bool, err error) {
	host, port := cfg.AmbilightHost, cfg.AmbilightPort
	err = ambilight.WithClient(host, port, func(client *ambilight.Client) error {
		enabled, err = client.Enabled()
		return err
	})

	return enabled, err
}

func createPartyModeEntry(cfg *Config, mgr *MenuManager) (*ToggleEntry, error) {
	partyModeEntry := &ToggleEntry{
		Text:  "Party!",
		Order: []string{"✓", "×"},
		Actions: map[string]Action{
			"✓": func() error {
				return ambilightChangeState(cfg, true)
			},
			"×": func() error {
				return ambilightChangeState(cfg, false)
			},
		},
	}

	go func() {
		for {
			enabled, err := ambilightIsEnabled(cfg)
			if err != nil {
				log.Printf("Failed to query state of ambilight: %v", err)
				continue
			} else {
				partyModeEntry.SetState(boolToGlyph(enabled), false)
				mgr.Display()
			}

			// Check every 20 seconds
			time.Sleep(20 * time.Second)
		}
	}()

	return partyModeEntry, nil
}

func createOutputEntry(mgr *MenuManager, MPD *mpd.Client) (*ToggleEntry, error) {
	outputs, err := MPD.Outputs()
	if err != nil {
		return nil, err
	}

	actionMap := map[string]Action{}
	for _, output := range outputs {
		// Stupid closure trick so we bind the right loop var:
		actionMap[output] = func(name string) func() error {
			return func() error {
				return MPD.SwitchToOutput(name)
			}
		}(output)
	}

	outputEntry := &ToggleEntry{
		Text:    "Output",
		Actions: actionMap,
		Order:   outputs,
	}

	MPD.Register("output", func() {
		active, err := MPD.ActiveOutput()
		log.Printf("Output changed to %v from extern", active)
		if err != nil {
			log.Printf("Failed to get active output: %v", err)
			return
		}

		outputEntry.SetState(active, false)
		mgr.Display()
	})

	return outputEntry, nil
}

func createPlaybackEntry(mgr *MenuManager, MPD *mpd.Client) (*ToggleEntry, error) {
	playbackEntry := &ToggleEntry{
		Text:  "Playback",
		Order: []string{"▶", "⏸", "⏹"},
		Actions: map[string]Action{
			"▶": func() error {
				return MPD.Play()
			},
			"⏸": func() error {
				return MPD.Pause()
			},
			"⏹": func() error {
				return MPD.Stop()
			},
		},
	}

	MPD.Register("player", func() {
		newState := mpd.StateToUnicode(MPD.CurrentState())
		log.Printf("State changed to %v from extern", newState)
		if err := playbackEntry.SetState(newState, false); err != nil {
			log.Printf("Failed to set new state: %v", err)
		}

		mgr.Display()
	})

	return playbackEntry, nil
}

func createRandomEntry(mgr *MenuManager, MPD *mpd.Client) (*ToggleEntry, error) {
	randomEntry := &ToggleEntry{
		Text:  "Random",
		Order: []string{"✓", "×"},
		Actions: map[string]Action{
			"✓": func() error {
				return MPD.EnableRandom(true)
			},
			"×": func() error {
				return MPD.EnableRandom(false)
			},
		},
	}

	MPD.Register("options", func() {
		log.Printf("random changed to %v from extern", MPD.IsRandom())
		randomEntry.SetState(boolToGlyph(MPD.IsRandom()), false)
		mgr.Display()
	})

	return randomEntry, nil
}

/////////////////////////

type Config struct {
	Width  int
	Height int

	DisplayHost string
	DisplayPort int

	MPDHost string
	MPDPort int

	AmbilightHost string
	AmbilightPort int
}

/////////////////////////
// MENU MAINLOOP LOGIC //
/////////////////////////

func Run(cfg *Config, ctx context.Context) error {
	log.Printf("Connecting to displayd...")
	lw, err := display.Connect(&display.Config{
		Host: cfg.DisplayHost,
		Port: cfg.DisplayPort,
	})

	if err != nil {
		return err
	}

	defer lw.Close()

	if err := drawStaticScreens(lw); err != nil {
		log.Printf("Failed to draw static screens: %v", err)
		return err
	}

	log.Printf("Creating menus...")
	mgr, err := NewMenuManager(cfg, lw, "startup")
	if err != nil {
		return err
	}

	// Wait a second to give the startup screen a bit time
	// to show before switching to mpd status.
	go func() {
		time.Sleep(1 * time.Second)
		mgr.SwitchTo("mpd")
	}()

	// Some flags to coordinate actions:
	ignoreRelease := false

	switcher := func(name string) func() error {
		return func() error {
			ignoreRelease = true
			return mgr.SwitchTo(name)
		}
	}

	// Start auxillary services:
	log.Printf("Starting background services...")
	MPD, err := mpd.NewClient(&mpd.Config{
		Host:        cfg.MPDHost,
		Port:        cfg.MPDPort,
		DisplayHost: cfg.DisplayHost,
		DisplayPort: cfg.DisplayPort,
	})

	if err != nil {
		log.Printf("Failed to create mpd client: %v", err)
		return err
	}

	go MPD.Run(ctx)
	go RunClock(lw, cfg.Width, ctx)
	go RunSysinfo(lw, cfg.Width, ctx)

	// Create some special entries with extended logic:

	outputEntry, err := createOutputEntry(mgr, MPD)
	if err != nil {
		log.Printf("Failed to create output entry: %v", err)
		return err
	}

	partyModeEntry, err := createPartyModeEntry(cfg, mgr)
	if err != nil {
		log.Printf("Failed to create party-mode entry: %v", err)
		return err
	}

	playbackEntry, err := createPlaybackEntry(mgr, MPD)
	if err != nil {
		log.Printf("Failed to create playback entry: %v", err)
		return err
	}

	randomEntry, err := createRandomEntry(mgr, MPD)
	if err != nil {
		log.Printf("Failed to create random entry: %v", err)
		return err
	}

	// Define the menu structure:

	mainMenu := []Entry{
		&Separator{"MODES"},
		&ClickEntry{
			Text:       "Music info",
			ActionFunc: switcher("mpd"),
		},
		&ClickEntry{
			Text: "Playlists",
			ActionFunc: func() error {
				entries := createPlaylistEntries(MPD)
				if err := mgr.AddMenu("playlists", entries); err != nil {
					return err
				}

				return switcher("playlists")()
			},
		},
		&ClickEntry{
			Text:       "Clock",
			ActionFunc: switcher("clock"),
		},
		&ClickEntry{
			Text:       "System info",
			ActionFunc: switcher("sysinfo"),
		},
		&ClickEntry{
			Text:       "Statistics",
			ActionFunc: switcher("stats"),
		},
		&Separator{"OPTIONS"},
		partyModeEntry,
		outputEntry,
		playbackEntry,
		randomEntry,
		&Separator{"SYSTEM"},
		&ClickEntry{
			Text:       "Powermenu",
			ActionFunc: switcher("menu-power"),
		},
		&ClickEntry{
			Text:       "About",
			ActionFunc: switcher("about"),
		},
	}

	powerMenu := []Entry{
		&ClickEntry{
			Text: "Poweroff",
			ActionFunc: func() error {
				switchToStatic(lw, "shutdown")
				return sysCommand("systemctl", "poweroff")()
			},
		},
		&ClickEntry{
			Text: "Reboot",
			ActionFunc: func() error {
				switchToStatic(lw, "shutdown")
				return sysCommand("systemctl", "reboot")()
			},
		},
		&ClickEntry{
			Text:       "Exit",
			ActionFunc: switcher("menu-main"),
		},
	}

	if err := mgr.AddMenu("menu-main", mainMenu); err != nil {
		log.Printf("Add main-menu failed: %v", err)
		return err
	}

	if err := mgr.AddMenu("menu-power", powerMenu); err != nil {
		log.Printf("Add main-power failed: %v", err)
		return err
	}

	mgr.AddTimedAction(10*time.Millisecond, func() error {
		return nil
	})

	mgr.AddTimedAction(600*time.Millisecond, func() error {
		ignoreRelease = true
		return mgr.SwitchTo("menu-main")
	})

	mgr.AddTimedAction(3*time.Second, func() error {
		ignoreRelease = true
		return mgr.SwitchTo("menu-power")
	})

	mgr.AddTimedAction(8*time.Second, func() error {
		ignoreRelease = true
		cmd := sysCommand("aplay", "/root/hoot.wav")
		go cmd()
		return nil
	})

	mgr.ReleaseAction(func() error {
		if ignoreRelease {
			ignoreRelease = false
			return nil
		}

		switch currWin := mgr.ActiveWindow(); currWin {
		case "mpd":
			if err := MPD.TogglePlayback(); err != nil {
				log.Printf("Failed to toggle playback: %v", err)
			}
		default:
			// This is a bit of a hack:
			// Enable "click to exit window" on most non-menu windows:
			if !strings.Contains(currWin, "menu") {
				return mgr.SwitchTo("menu-main")
			}
		}

		return nil
	})

	mgr.RotateAction(func() error {
		if mgr.ActiveWindow() != "mpd" {
			return nil
		}

		log.Printf("rotate action")
		switch mgr.Direction() {
		case DirectionRight:
			log.Printf("Play next")
			if err := MPD.Next(); err != nil {
				log.Printf("Failed to skip to next: %v", err)
			}
		case DirectionLeft:
			log.Printf("Play prev")
			if err := MPD.Next(); err != nil {
				log.Printf("Failed to skip to prev: %v", err)
			}
		}

		return nil
	})

	log.Printf("Waiting for a silent death...")
	<-ctx.Done()

	return mgr.Close()
}
