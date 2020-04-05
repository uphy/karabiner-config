package app

import (
	"fmt"
	"log"

	"github.com/uphy/karabiner-config/config"
	"github.com/uphy/karabiner-config/watch"
)

type (
	Runner struct {
		writer ConfigWriter
	}
)

func (r *Runner) Run(input string) error {
	conf, err := config.Load(input)
	if err != nil {
		return fmt.Errorf("failed to load: %w", err)
	}
	return r.writer.Write(conf)
}

func (r *Runner) Watch(input string) error {
	err := r.Run(input)
	if err != nil {
		return err
	}
	watch.WatchFile(input, func() {
		if err := r.Run(input); err != nil {
			log.Println(err)
		}
	})
	return nil
}
