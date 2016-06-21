package ui

import "github.com/studentkittens/eulenfunk/ui/mpd"

func createPlaylistEntries(mgr *MenuManager, MPD *mpd.Client) []Entry {
	entries := []Entry{&Separator{"Playlists"}}

	for _, name := range MPD.ListPlaylists() {
		entries = append(entries, &ClickEntry{
			Text: name,

			// Closure trick so we don't get the last loop var:
			ActionFunc: func(name string) func() error {
				return func() error {
					if err := MPD.LoadAndPlayPlaylist(name); err != nil {
						return err
					}

					return mgr.SwitchTo("mpd")
				}
			}(name),
		})
	}

	return entries
}
