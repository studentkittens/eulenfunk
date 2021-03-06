package ui

import (
	"log"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/studentkittens/eulenfunk/ambilight"
	"github.com/studentkittens/eulenfunk/display"
	"github.com/studentkittens/eulenfunk/lightd"
	"github.com/studentkittens/eulenfunk/ui/mpd"
	"github.com/studentkittens/eulenfunk/util"
)

func runBinary(name string, args ...string) {
	if err := exec.Command(name, args...).Run(); err != nil {
		log.Printf("Failed to execute `%s`: %v", name, err)
	}
}

func boolToGlyph(b bool) string {
	if b {
		return "✓"
	}

	return "×"
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
				log.Printf("(Waiting 5 seconds before retrying)")
				time.Sleep(5 * time.Second)
				continue
			} else {
				partyModeEntry.SetState(boolToGlyph(enabled))
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
				if err := MPD.SwitchToOutput(name); err != nil {
					return err
				}

				// Playback might be paused when switching outputs:
				return MPD.Play()
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
		if err != nil {
			log.Printf("Failed to get active output: %v", err)
			return
		}

		outputEntry.SetState(active)
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
		playbackEntry.SetState(newState)

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
		randomEntry.SetState(boolToGlyph(MPD.IsRandom()))
		mgr.Display()
	})

	return randomEntry, nil
}

/////////////////////////

func releaseAction(mgr *MenuManager, MPD *mpd.Client) error {
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
}

func rotateAction(mgr *MenuManager, MPD *mpd.Client) error {
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
		if err := MPD.Prev(); err != nil {
			log.Printf("Failed to skip to prev: %v", err)
		}
	}

	return nil
}

func schuhuAction() error {
	go func() {
		for i := 0; i < 10; i++ {
			runBinary("aplay", "/root/hoot.wav")
		}
	}()
	return nil
}

func shutdown(cfg *Config, lw *display.LineWriter, mode, effect string) error {
	lightdCfg := &lightd.Config{
		Host: cfg.LightdHost,
		Port: cfg.LightdPort,
	}

	switchToStatic(lw, "shutdown")

	if err := ambilightChangeState(cfg, false); err != nil {
		log.Printf("NOTE: Failed to disable ambilight: %v", err)
	}

	if err := lightd.Send(lightdCfg, effect); err != nil {
		log.Printf("NOTE: Failed to send flash effect to lightd: %v", err)
	}

	// Give it a bit of time to work:
	time.Sleep(1 * time.Second)
	runBinary("systemctl", mode)
	return nil
}

func poweroffAction(cfg *Config, lw *display.LineWriter) error {
	return shutdown(cfg, lw, "poweroff", "flash{200ms|{255,0,0}|5}")
}

func rebootAction(cfg *Config, lw *display.LineWriter) error {
	return shutdown(cfg, lw, "reboot", "flash{200ms|{255,150,0}|5}")
}

/////////////////////////

func createMainMenu(mgr *MenuManager, MPD *mpd.Client) error {
	outputEntry, err := createOutputEntry(mgr, MPD)
	if err != nil {
		log.Printf("Failed to create output entry: %v", err)
		return err
	}

	partyModeEntry, err := createPartyModeEntry(mgr.Config, mgr)
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

	mainMenu := []Entry{
		&Separator{"MODES"},
		&ClickEntry{
			Text:       "Now Playing",
			ActionFunc: switcher(mgr, "mpd"),
		},
		&ClickEntry{
			Text: "Playlists",
			ActionFunc: func() error {
				entries := createPlaylistEntries(mgr, MPD)

				// Add an exit button:
				entries = append(entries, &ClickEntry{
					Text:       "(Exit)",
					ActionFunc: switcher(mgr, "menu-main"),
				})

				if err := mgr.AddMenu("menu-playlists", entries); err != nil {
					return err
				}

				return switcher(mgr, "menu-playlists")()
			},
		},
		&ClickEntry{
			Text:       "Clock",
			ActionFunc: switcher(mgr, "clock"),
		},
		&ClickEntry{
			Text:       "Weather",
			ActionFunc: switcher(mgr, "weather"),
		},
		&ClickEntry{
			Text:       "Statistics",
			ActionFunc: switcher(mgr, "stats"),
		},
		&Separator{"OPTIONS"},
		partyModeEntry,
		outputEntry,
		playbackEntry,
		randomEntry,
		&Separator{"SYSTEM"},
		&ClickEntry{
			Text:       "Powermenu",
			ActionFunc: switcher(mgr, "menu-power"),
		},
		&ClickEntry{
			Text:       "System info",
			ActionFunc: switcher(mgr, "sysinfo"),
		},
		&ClickEntry{
			Text: "About",
			ActionFunc: func() error {
				if err := schuhuAction(); err != nil {
					log.Printf("No schuhu: %v", err)
				}

				return mgr.SwitchTo("about")
			},
		},
	}

	if err := mgr.AddMenu("menu-main", mainMenu); err != nil {
		log.Printf("Add main-menu failed: %v", err)
		return err
	}

	return nil
}

func createPowerMenu(mgr *MenuManager, lw *display.LineWriter) error {
	powerMenu := []Entry{
		&Separator{"POWER MENU"},
		&ClickEntry{
			Text: "Poweroff",
			ActionFunc: func() error {
				return poweroffAction(mgr.Config, lw)
			},
		},
		&ClickEntry{
			Text: "Reboot",
			ActionFunc: func() error {
				return rebootAction(mgr.Config, lw)
			},
		},
		&ClickEntry{
			Text:       "(Exit)",
			ActionFunc: switcher(mgr, "menu-main"),
		},
	}

	if err := mgr.AddMenu("menu-power", powerMenu); err != nil {
		log.Printf("Add main-power failed: %v", err)
		return err
	}

	return nil
}

/////////////////////////

// Config allows the user to configure to which services the ui connects.
type Config struct {
	Width  int
	Height int

	DisplayHost string
	DisplayPort int

	MPDHost string
	MPDPort int

	AmbilightHost string
	AmbilightPort int

	LightdHost string
	LightdPort int
}

/////////////////////////
// MENU MAINLOOP LOGIC //
/////////////////////////

func switcher(mgr *MenuManager, name string) func() error {
	return func() error {
		return mgr.SwitchTo(name)
	}
}

func initialSwitchToMPD(mgr *MenuManager, MPD *mpd.Client) {
	initial := true

	MPD.Register("player", func() {
		// Wait for the first "real" event.
		// First one is just emiited on startup.
		if initial {
			if err := mgr.SwitchTo("mpd"); err != nil {
				log.Printf("Initial switch to mpd failed: %v", err)
			}

			initial = false
		}
	})
}

// Run starts the UI with the settings in `cfg` and until `ctx` is canceled.
func Run(cfg *Config, ctx context.Context) error {
	log.Printf("Connecting to displayd...")
	lw, err := display.Connect(&display.Config{
		Host: cfg.DisplayHost,
		Port: cfg.DisplayPort,
	}, ctx)

	if err != nil {
		return err
	}

	defer util.Closer(lw)

	if staticErr := drawStaticScreens(lw); staticErr != nil {
		log.Printf("Failed to draw static screens: %v", staticErr)
		return staticErr
	}

	log.Printf("Creating menus...")
	mgr, err := NewMenuManager(cfg, lw, "startup")
	if err != nil {
		return err
	}

	// Start auxiliary services:
	log.Printf("Starting background services...")
	MPD, err := mpd.NewClient(&mpd.Config{
		MPDHost:     cfg.MPDHost,
		MPDPort:     cfg.MPDPort,
		DisplayHost: cfg.DisplayHost,
		DisplayPort: cfg.DisplayPort,
	}, ctx)

	if err != nil {
		log.Printf("Failed to create mpd client: %v", err)
		return err
	}

	defer util.Closer(MPD)

	initialSwitchToMPD(mgr, MPD)

	go MPD.Run()
	go RunClock(lw, cfg.Width, ctx)
	go RunSysinfo(lw, cfg.Width, ctx)
	go RunWeather(lw, cfg.Width, ctx)

	if err := createMainMenu(mgr, MPD); err != nil {
		return err
	}

	if err := createPowerMenu(mgr, lw); err != nil {
		return err
	}

	mgr.AddTimedAction(600*time.Millisecond, switcher(mgr, "menu-main"))
	mgr.AddTimedAction(2*time.Second, switcher(mgr, "menu-playlists"))
	mgr.AddTimedAction(3*time.Second, switcher(mgr, "menu-power"))
	mgr.AddTimedAction(8*time.Second, schuhuAction)

	mgr.ReleaseAction(func() error {
		return releaseAction(mgr, MPD)
	})

	mgr.RotateAction(func() error {
		return rotateAction(mgr, MPD)
	})

	log.Printf("Waiting for a silent death...")
	<-ctx.Done()

	return mgr.Close()
}
