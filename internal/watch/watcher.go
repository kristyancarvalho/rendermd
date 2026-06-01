package watch

import (
	"fmt"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	Events chan struct{}
	Errors chan error
	path   string
	fw     *fsnotify.Watcher
	done   chan struct{}
}

func New(path string, debounce time.Duration) (*Watcher, error) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	if err := fw.Add(path); err != nil {
		fw.Close()
		return nil, err
	}
	w := &Watcher{
		Events: make(chan struct{}, 1),
		Errors: make(chan error, 1),
		path:   path,
		fw:     fw,
		done:   make(chan struct{}),
	}
	go w.loop(debounce)
	return w, nil
}

func (w *Watcher) Close() error {
	close(w.done)
	return w.fw.Close()
}

func (w *Watcher) loop(debounce time.Duration) {
	var timer *time.Timer
	var deletedAt time.Time
	deleted := false

	fire := func() {
		select {
		case w.Events <- struct{}{}:
		default:
		}
	}

	sendErr := func(err error) {
		select {
		case w.Errors <- err:
		default:
			fmt.Fprintln(os.Stderr, "mdp: watch error:", err)
		}
	}

	for {
		select {
		case <-w.done:
			return
		case ev, ok := <-w.fw.Events:
			if !ok {
				return
			}
			if ev.Has(fsnotify.Remove) || ev.Has(fsnotify.Rename) {
				deleted = true
				deletedAt = time.Now()
				continue
			}
			if deleted {
				if ev.Has(fsnotify.Create) || ev.Has(fsnotify.Write) {
					_ = w.fw.Add(w.path)
					deleted = false
				}
			}
			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(debounce, fire)
		case err, ok := <-w.fw.Errors:
			if !ok {
				return
			}
			sendErr(err)
		case <-time.After(100 * time.Millisecond):
			if deleted && time.Since(deletedAt) > 2*time.Second {
				deleted = false
				sendErr(fmt.Errorf("file removed: %s", w.path))
			}
		}
	}
}
