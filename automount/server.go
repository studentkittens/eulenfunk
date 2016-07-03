package automount

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fhs/gompd/mpd"
	"github.com/studentkittens/eulenfunk/util"
	"golang.org/x/net/context"
)

const (
	mountSubDir = "mounts"
)

// Config gives the user of automount some adjustment screws.
// See the fields for the available options.
type Config struct {
	AutomountHost string
	AutomountPort int
	MPDHost       string
	MPDPort       int
	MusicDir      string
}

type server struct {
	Config  *Config
	Cancel  context.CancelFunc
	Context context.Context
}

func runBinary(name string, args ...string) error {
	if stdout, err := exec.Command(name, args...).Output(); err != nil {
		log.Printf("Failed to execute `%s`: %v", name, err)

		log.Printf("Stdout output was: %s", stdout)
		if exitErr, ok := err.(*exec.ExitError); ok {
			log.Printf("Stderr output was: %s", exitErr.Stderr)
		}
		return err
	}

	return nil
}

func (srv *server) playlistFromDir(client *mpd.Client, label string) error {
	allUris, err := client.GetFiles()
	if err != nil {
		return err
	}

	// NOTE: There is the "searchaddpl" command which
	//       would do exactly the same as this function, but more efficient.
	//       Sadly, it's not supported by this library.
	//       (Also "find base <label>" would be easier than .GetFiles()
	//        but the library really failed here too)
	uris := []string{}
	prefix := filepath.Join(mountSubDir, label)
	for _, uri := range allUris {
		if strings.HasPrefix(uri, prefix) {
			uris = append(uris, uri)
		}
	}

	cmdlist := client.BeginCommandList()

	playlistName := playlistNameFromLabel(label)

	for _, uri := range uris {
		log.Printf("Adding `%s`", uri)
		cmdlist.PlaylistAdd(playlistName, uri)
	}

	return cmdlist.End()
}

func (srv *server) getUpdateID(client *mpd.Client) (int, error) {
	status, err := client.Status()
	if err != nil {
		return 0, err
	}

	// Assume an (valid) ID of 0 when no such attr is known.
	// That means that no update is currently active.
	idStr, ok := status["updating_db"]
	if !ok {
		return 0, nil
	}

	return strconv.Atoi(idStr)
}

func (srv *server) waitFor(timeout time.Duration, events ...string) error {
	addr := fmt.Sprintf("%s:%d", srv.Config.MPDHost, srv.Config.MPDPort)
	log.Printf("Waiting for %v of %s", events, addr)

	watcher, err := mpd.NewWatcher("tcp", addr, "", events...)
	if err != nil {
		return err
	}

	defer util.Closer(watcher)

	timer := time.NewTimer(timeout)

	select {
	case <-watcher.Event:
	case <-timer.C:
	}

	return nil
}

func (srv *server) updateDatabase(client *mpd.Client, label string) error {
	subDir := filepath.Join(mountSubDir, label)

	lastID, err := client.Update(subDir)
	log.Printf("Sending update")
	if err != nil {
		log.Printf("Sending update failed: %v", err)
		return err
	}

	// lastID, err := srv.getUpdateID(client)
	// if err != nil {
	// 	log.Printf("Getting initial status failed: %v", err)
	// 	return err
	// }

	log.Printf("Updating MPD database with new songs (job: %d; dir: %v)", lastID, subDir)

	for {
		if err := srv.waitFor(10 * time.Second, "update"); err != nil {
			return err
		}

		log.Printf("WaitFor done")
		currID, err := srv.getUpdateID(client)
		if err != nil {
			return err
		}

		log.Printf("Got an event; current ID is %d", currID)
		if currID == 0 || currID > lastID {
			break
		}

		log.Printf("Continueing...")
	}

	// "Security" wait:
	time.Sleep(5 * time.Second)

	return nil
}

func playlistNameFromLabel(label string) string {
	return "stick-" + label
}

func (srv *server) mountToPlaylist(destPath, label string) error {
	addr := fmt.Sprintf("%s:%d", srv.Config.MPDHost, srv.Config.MPDPort)
	client, err := mpd.Dial("tcp", addr)
	if err != nil {
		return err
	}

	defer util.Closer(client)

	if dbErr := srv.updateDatabase(client, label); dbErr != nil {
		log.Printf("Updating MPD failed: %v", dbErr)
		return dbErr
	}

	playlists, err := client.ListPlaylists()
	if err != nil {
		return err
	}

	playlistName := playlistNameFromLabel(label)
	for _, playlist := range playlists {
		if playlist["playlist"] == playlistName {
			if err := client.PlaylistClear(playlistName); err != nil {
				return err
			}

			break
		}
	}

	return srv.playlistFromDir(client, label)
}

func (srv *server) mount(device, label string) error {
	destPath := filepath.Join(srv.Config.MusicDir, mountSubDir, label)
	if err := os.MkdirAll(destPath, 0777); err != nil {
		return err
	}

	log.Printf("Mounting `%s` to `%s`\n", device, destPath)
	if err := runBinary("mount", device, destPath); err != nil {
		return err
	}

	if err := srv.mountToPlaylist(destPath, playlistNameFromLabel(label)); err != nil {
		return err
	}

	return nil
}

func (srv *server) unmount(device, label string) error {
	addr := fmt.Sprintf("%s:%d", srv.Config.MPDHost, srv.Config.MPDPort)
	client, err := mpd.Dial("tcp", addr)
	if err != nil {
		return err
	}

	defer util.Closer(client)

	if err := client.PlaylistRemove(playlistNameFromLabel(label)); err != nil {
		return err
	}

	log.Printf("Unmounting`%s`\n", device)
	if err := runBinary("umount", "-l", device); err != nil {
		return err
	}

	// Do not remove recursively. That might go very bad if umount did not work:
	destPath := filepath.Join(srv.Config.MusicDir, label)
	if err := os.Remove(destPath); err != nil {
		log.Printf("Failed to remove mount dir: %v", err)
		return err
	}

	return nil
}

func (srv *server) handleLine(line string) bool {
	log.Printf("Received: %v", line)
	split := strings.Split(line, " ")

	switch split[0] {
	case "mount":
		if len(split) >= 3 {
			if err := srv.mount(split[1], split[2]); err != nil {
				log.Printf("Failed to mount: %v", err)
			}
		}
	case "unmount":
		if len(split) >= 3 {
			if err := srv.unmount(split[1], split[2]); err != nil {
				log.Printf("Failed to unmount: %v", err)
			}
		}
	case "close":
		return false
	case "quit":
		srv.Cancel()
		return false
	}

	return true
}

func (srv *server) handleRequests(conn io.ReadCloser) {
	defer util.Closer(conn)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		if !srv.handleLine(strings.TrimSpace(scanner.Text())) {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("line scanning failed: %v", err)
	}
}

func cancelled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

// Run creates a new automountd on the specified host and port.
func Run(cfg *Config, ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", cfg.AutomountHost, cfg.AutomountPort)
	lsn, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("Error listening: %v", err.Error())
		return err
	}

	subCtx, cancel := context.WithCancel(ctx)

	srv := &server{
		Config:  cfg,
		Context: subCtx,
		Cancel:  cancel,
	}

	defer util.Closer(lsn)
	log.Println("Listening on " + addr)

	// Manually trigger a udevadm run after a few seconds:
	go func() {
		time.Sleep(2 * time.Second)
		if err := exec.Command("udevadm", "trigger", "-c", "add").Run(); err != nil {
			log.Printf("Failed to trigger udev: %v", err)
		}
	}()

	for !cancelled(ctx) {
		if tcpLsn, ok := lsn.(*net.TCPListener); ok {
			if err := tcpLsn.SetDeadline(time.Now().Add(2 * time.Second)); err != nil {
				log.Printf("Setting deadline failed: %v", err)
				return err
			}
		}

		conn, err := lsn.Accept()
		if err, ok := err.(*net.OpError); ok && err.Timeout() {
			continue
		}

		if err != nil {
			log.Printf("Error accepting: %v", err.Error())
			return err
		}

		go srv.handleRequests(conn)
	}

	return nil
}
