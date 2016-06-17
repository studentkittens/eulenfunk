package mpd

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/fhs/gompd/mpd"
)

type ReMPD struct {
	sync.Mutex

	host   string
	port   int
	client *mpd.Client
}

func NewReMPD(host string, port int) *ReMPD {
	return &ReMPD{host: host, port: port}
}

func (rc *ReMPD) reconnect() error {
	addr := fmt.Sprintf("%s:%d", rc.host, rc.port)
	mpdClient, err := mpd.Dial("tcp", addr)
	if err != nil {
		log.Printf("Failed to connect to mpd (%s): %v", addr, err)
		return err
	}

	rc.client = mpdClient
	return nil
}

func (rc *ReMPD) Client() *mpd.Client {
	rc.Lock()
	defer rc.Unlock()

	for {
		if rc.client == nil || rc.client.Ping() != nil {
			if err := rc.reconnect(); err != nil {
				log.Printf("No MPD connection; retry in 4s...")
				time.Sleep(5 * time.Second)
				continue
			}
		}

		break
	}

	return rc.client
}

type ReWatcher struct {
	sync.Mutex

	host     string
	port     int
	watcher  *mpd.Watcher
	listenOn []string

	Events chan string
}

func NewReWatcher(host string, port int, listenOn ...string) *ReWatcher {
	rw := &ReWatcher{
		host:     host,
		port:     port,
		listenOn: listenOn,
		Events:   make(chan string),
	}

	rw.retryUntilSuccesfull()

	go func() {
		for err := range rw.watcher.Error {
			log.Printf("MPD Watcher errored: %v", err)
			rw.retryUntilSuccesfull()
		}
	}()

	return rw
}

func (rw *ReWatcher) retryUntilSuccesfull() {
	for {
		if err := rw.reconnect(); err != nil {
			log.Printf("Failed to watch mpd: %v", err)
			log.Printf("Retrying in 5 seconds.")
			time.Sleep(5 * time.Second)
			continue
		}

		break
	}
}

func (rw *ReWatcher) reconnect() error {
	// NOTE: Old connection should time out or is already closed.
	addr := fmt.Sprintf("%s:%d", rw.host, rw.port)
	watcher, err := mpd.NewWatcher("tcp", addr, "", rw.listenOn...)
	if err != nil {
		return err
	}

	go func() {
		for ev := range watcher.Event {
			rw.Events <- ev
		}
	}()

	rw.Lock()
	defer rw.Unlock()

	rw.watcher = watcher
	return nil
}

func (rw *ReWatcher) Close() error {
	rw.Lock()
	defer rw.Unlock()

	if rw.watcher == nil {
		return nil
	}

	return rw.watcher.Close()
}
