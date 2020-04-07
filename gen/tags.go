package gen

import "github.com/uphy/karabiner-config/config"

type Tags struct {
	tagToManipulators map[string][]config.ManipulatorConfig
}

func newTags() *Tags {
	return &Tags{map[string][]config.ManipulatorConfig{}}
}

func (t *Tags) BindManipulators(manipulator *config.ManipulatorConfig, tags []string) {
	for _, tag := range tags {
		t.tagToManipulators[tag] = append(t.tagToManipulators[tag], *manipulator)
	}
}

func (t *Tags) ManipulatorsFor(includes []string, excludes []string) []config.ManipulatorConfig {
	manipulators := make([]config.ManipulatorConfig, 0)
l:
	for tag, m := range t.tagToManipulators {
		included := false
		for _, include := range includes {
			if tag == include {
				included = true
				break
			}
		}
		if !included {
			continue
		}
		for _, exclude := range excludes {
			if tag == exclude {
				continue l
			}
		}
		manipulators = append(manipulators, m...)
	}
	return manipulators
}
