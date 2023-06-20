package runtime

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/micro/micro/v3/service/logger"
)

type CallbackFunc func() error

type watcher struct {
	root        string
	watchDelay  time.Duration
	fileWatcher *fsnotify.Watcher

	eventsChan   chan string
	callbackFunc CallbackFunc
	stopChan     chan struct{}
}

func NewWatcher(root string, delay time.Duration, fn CallbackFunc) (*watcher, error) {
	fileWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &watcher{
		root:         root,
		watchDelay:   delay,
		fileWatcher:  fileWatcher,
		eventsChan:   make(chan string, 1000),
		callbackFunc: fn,
		stopChan:     make(chan struct{}),
	}, nil
}

// Watch the file changes in specific directories
func (w *watcher) Watch() error {
	err := filepath.WalkDir(w.root, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info != nil && !info.IsDir() {
			return nil
		}

		return w.watchDirectory(path)
	})

	if err != nil {
		return err
	}

	w.start()

	return nil
}

// start the watching process
func (w *watcher) start() {
	for {
		select {
		case <-w.stopChan:
			logger.Infof("Watcher is exiting...")
			return
		case <-w.eventsChan:

			time.Sleep(w.watchDelay)
			w.flushEvents()

			if err := w.callbackFunc(); err != nil {
				logger.Errorf("Watcher callback function execute error: %v", err)
				break
			}

		}
	}
}

// flushEvents flush the events in buffer channel
func (w *watcher) flushEvents() {
	for {
		select {
		case ev := <-w.eventsChan:
			logger.Debugf("Watcher flush event: %v", ev)
		default:
			return
		}
	}
}

// validExtension checks the extension of file is valid to be watched
func (w *watcher) validExtension(filepath string) bool {
	if strings.HasSuffix(filepath, ".go") || strings.HasSuffix(filepath, ".proto") {
		return true
	}

	return false
}

// watchDirectory watch all the files in the dir, recurse all the subdirectories
func (w *watcher) watchDirectory(dir string) error {
	logger.Infof("Watcher is watching path: %+v", dir)

	err := w.fileWatcher.Add(dir)
	if err != nil {
		logger.Fatal(err)
	}

	go func() {
		for {
			select {
			case <-w.stopChan:
				return
			case event, ok := <-w.fileWatcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write != fsnotify.Write {
					break
				}

				eventPath := event.Name

				if !w.validExtension(eventPath) {
					break
				}

				f, err := os.Stat(eventPath)
				if err != nil {
					logger.Errorf("File get file info: %v", err)
					break
				}

				if f.IsDir() {
					err := w.watchDirectory(eventPath)
					if err != nil {
						logger.Errorf("Watching dir error: %v", err)
					}
					break
				}

				logger.Infof("%v has changed", eventPath)

				w.eventsChan <- event.Name

			case err, ok := <-w.fileWatcher.Errors:
				if !ok {
					return
				}
				fmt.Println("error:", err)
			}
		}
	}()

	return nil
}

// Stop the watching process
func (w *watcher) Stop() {
	close(w.stopChan)
}
