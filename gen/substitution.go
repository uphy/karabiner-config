package gen

import (
	"errors"
	"strings"
)

type Substitution struct {
	Name  string
	Value interface{}
}

func parseSubstitution(expression string) (*Substitution, error) {
	eqIndex := strings.Index(expression, "=")
	if eqIndex < 0 {
		return nil, errors.New("Substitution requires '=': " + expression)
	}
	name := strings.Trim(expression[:eqIndex], " ")
	rawValue := strings.Trim(expression[eqIndex+1:], " ")
	value := parseValue(rawValue)
	return &Substitution{name, value}, nil
}

func (s *Substitution) ToJSON() JSONObject {
	o := make(JSONObject)
	o["name"] = s.Name
	o["value"] = s.Value
	return o
}
