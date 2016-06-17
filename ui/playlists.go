package ui

import (
	"log"

	"github.com/studentkittens/eulenfunk/display"
	"github.com/studentkittens/eulenfunk/ui/mpd"
)

func showPlaylistWindow(lw *display.LineWriter, MPD *mpd.Client) error {
	for idx, playlist := range MPD.ListPlaylists() {
		if _, err := lw.Formatf("line playlists %d %s", idx, playlist); err != nil {
			log.Printf("Failed to format playlist line %d: %v", idx, err)
			return err
		}
	}

	_, err := lw.Formatf("switch playlists")
	return err
}
