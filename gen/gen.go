package gen

import (
	"io"

	"github.com/uphy/karabiner-config/config"
)

func GenerateObject(c *config.RootConfig) (JSONObject, error) {
	ctx := newRootScopeContext()
	root := make(JSONObject)
	rules := make(JSONArray, 0)
	for _, manipulatorConfig := range c.Manipulators {
		tags := manipulatorConfig.Tags
		if len(tags) == 0 {
			tags = append(tags, "default")
		}
		ctx.BindManipulators(&manipulatorConfig, tags)
	}
	for _, ruleConfig := range c.Rules {
		rule, err := generateRule(ctx, &ruleConfig)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	root["title"] = c.Title
	root["maintainers"] = c.Maintainers
	root["rules"] = rules
	return root, nil
}

func Generate(c *config.RootConfig, writer io.Writer) error {
	obj, err := GenerateObject(c)
	if err != nil {
		return err
	}
	return obj.Write(writer)
}

func generateRule(ctx *ScopeContext, ruleConfig *config.RuleConfig) (JSONObject, error) {
	rule := make(JSONObject, 0)
	manipulators := make(JSONArray, 0)
	ruleContext := ctx.ChildScope(ruleConfig.Conditions)

	// load manipulators for this rule
	manipulatorConfigs := append(ruleConfig.Manipulators, ruleContext.ManipulatorsFor(ruleConfig.Includes, ruleConfig.Excludes)...)
	for _, manipulatorConfig := range manipulatorConfigs {
		m, err := generateManipulators(ruleContext, &manipulatorConfig)
		if err != nil {
			return nil, err
		}
		for _, mm := range m {
			manipulators = append(manipulators, mm)
		}
	}
	rule["manipulators"] = manipulators
	rule["description"] = ruleConfig.Description
	return rule, nil
}
