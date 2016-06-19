package mpd

import (
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/fhs/gompd/mpd"
)

// ReMPD is a mpd connection that automatically reconnects itself.
// Individual actions might still fail, but the next call is supposed
// to work again (after a possibly long re-connect dance).
type ReMPD struct {
	sync.Mutex

	host   string
	port   int
	ctx    context.Context
	client *mpd.Client
}

// NewReMPD returns a new reconnector watching over `host` and `port`.
// It will stop re-connecting if ctx was canceled.
func NewReMPD(host string, port int, ctx context.Context) *ReMPD {
	return &ReMPD{host: host, port: port, ctx: ctx}
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

// Client returns the currently valid client.  If there is none, a new
// connection is established and the function will block until this happens.
func (rc *ReMPD) Client() *mpd.Client {
	rc.Lock()
	defer rc.Unlock()

	for {
		select {
		case <-rc.ctx.Done():
			log.Fatalf("Interrupted while waiting for connection")
			return nil
		default:
		}

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

// ReWatcher is like ReMPD, but re-connects a gompd.Watcher instance.
type ReWatcher struct {
	sync.Mutex

	host     string
	port     int
	watcher  *mpd.Watcher
	listenOn []string
	ctx      context.Context
	cancel   context.CancelFunc

	Events chan string
}

// NewReWatcher returns a new ReWatcher on `host` and `port`. It will listen on
// all events in `listenOn`.  It will stop watching when `ctx` is canceled.
func NewReWatcher(host string, port int, ctx context.Context, listenOn ...string) *ReWatcher {
	subCtx, cancel := context.WithCancel(ctx)

	rw := &ReWatcher{
		host:     host,
		port:     port,
		listenOn: listenOn,
		ctx:      subCtx,
		cancel:   cancel,
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
		select {
		case <-rw.ctx.Done():
			return
		default:
		}

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

// Close shutsdown the watcher. No events will be delievered afterwards.
func (rw *ReWatcher) Close() error {
	rw.Lock()
	defer rw.Unlock()

	rw.cancel()

	if rw.watcher == nil {
		return nil
	}

	return rw.watcher.Close()
}
