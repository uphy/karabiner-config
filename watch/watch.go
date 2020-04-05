package watch

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

type WatchFunc func()

func WatchFile(file string, onChange WatchFunc) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					onChange()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	watcher.Add(file)
	onChange()
	if err != nil {
		log.Fatal(err)
	}
	<-done
	return nil
}
