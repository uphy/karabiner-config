package gen

import (
	"errors"
	"strings"
)

type (
	Key struct {
		KeyCode   string     `json:"key_code"`
		Modifiers *Modifiers `json:"modifiers"`
	}
	Modifiers struct {
		Mandatory []string `json:"mandatory,omitempty"`
		Optional  []string `json:"optional,omitempty"`
	}
)

func parseKey(key string) (*Key, error) {
	splitted := strings.Split(key, "-")
	if len(splitted) == 0 {
		return nil, errors.New("empty key")
	}
	keyCode, err := keyCodeOf(splitted[len(splitted)-1])
	if err != nil {
		return nil, err
	}
	modifiers := splitted[:len(splitted)-1]
	mandatory := make([]string, 0)
	optional := make([]string, 0)
	for _, modifier := range modifiers {
		replacedModifier, err := modifierOf(modifier)
		if err != nil {
			return nil, err
		}
		mandatory = append(mandatory, replacedModifier)
	}
	return &Key{keyCode, &Modifiers{mandatory, optional}}, nil
}
