package gen

import (
	"errors"
	"strconv"
	"strings"
)

type VariableCondition struct {
	Name  string
	Value interface{}
	Eq    bool
}

func parseVariableCondition(expression string) (*VariableCondition, error) {
	eqIndex := strings.Index(expression, "==")
	notEqIndex := strings.Index(expression, "!=")
	var name, rawValue string
	var eq bool
	if eqIndex >= 0 {
		name = expression[0:eqIndex]
		rawValue = expression[eqIndex+2:]
		eq = true
	} else if notEqIndex >= 0 {
		name = expression[0:notEqIndex]
		rawValue = expression[notEqIndex+2:]
		eq = false
	} else {
		return nil, errors.New("'==' or '!=' must be specified in 'var_if'")
	}
	name = strings.Trim(name, " ")
	rawValue = strings.Trim(rawValue, " ")
	value := parseValue(rawValue)
	return &VariableCondition{name, value, eq}, nil
}

func parseValue(v string) interface{} {
	if strings.HasPrefix(v, "'") && strings.HasSuffix(v, "'") ||
		strings.HasPrefix(v, `"`) && strings.HasSuffix(v, `"`) {
		return v[1 : len(v)-1]
	}
	switch v {
	case "true", "yes":
		return true
	case "false", "no":
		return false
	default:
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return v
			}
			return f
		}
		return n
	}
}

func (v *VariableCondition) ToJSON() JSONObject {
	condition := make(JSONObject)
	if v.Eq {
		condition["type"] = "variable_if"
	} else {
		condition["type"] = "variable_unless"
	}
	condition["name"] = v.Name
	condition["value"] = v.Value
	return condition
}
