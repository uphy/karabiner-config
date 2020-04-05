package gen

import (
	"errors"
	"strings"
)

var modifierMap = map[string]string{
	"c": "control",
	"a": "option",
	"o": "option",
	"s": "shift",
}

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
	keyCode := splitted[len(splitted)-1]
	modifiers := splitted[:len(splitted)-1]
	mandatory := make([]string, 0)
	optional := make([]string, 0)
	for _, modifier := range modifiers {
		replacedModifier, exist := modifierMap[strings.ToLower(modifier)]
		if exist {
			modifier = replacedModifier
		}
		mandatory = append(mandatory, modifier)
	}
	return &Key{keyCode, &Modifiers{mandatory, optional}}, nil
}
