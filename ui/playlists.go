package ui

import "github.com/studentkittens/eulenfunk/ui/mpd"

func createPlaylistEntries(MPD *mpd.Client) []Entry {
	entries := []Entry{&Separator{"Playlists"}}

	for _, name := range MPD.ListPlaylists() {
		entries = append(entries, &ClickEntry{
			Text: name,

			// Closure trick so we don't get the last loop var:
			ActionFunc: func(name string) func() error {
				return func() error {
					return MPD.LoadAndPlayPlaylist(name)
				}
			}(name),
		})
	}

	return entries
}
