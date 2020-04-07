package gen

import (
	"github.com/uphy/karabiner-config/config"
)

type ScopeContext struct {
	parent     *ScopeContext
	conditions []config.ConditionConfig
	tags       *Tags
}

func newRootScopeContext() *ScopeContext {
	return newScopeContext(nil, nil)
}

func newScopeContext(parent *ScopeContext, conditions []config.ConditionConfig) *ScopeContext {
	return &ScopeContext{parent, conditions, newTags()}
}

func (c *ScopeContext) ChildScope(conditions []config.ConditionConfig) *ScopeContext {
	return newScopeContext(c, conditions)
}

func (c *ScopeContext) MergeManipulatorConfig(manipulator *config.ManipulatorConfig) *config.ManipulatorConfig {
	merged := manipulator
	if c.parent != nil {
		merged = c.parent.MergeManipulatorConfig(merged)
	}
	if c.conditions != nil {
		for _, condition := range c.conditions {
			merged.Conditions = append(merged.Conditions, condition)
		}
	}
	return merged
}

func (c *ScopeContext) ManipulatorsFor(includes []string, excludes []string) []config.ManipulatorConfig {
	manipulators := make([]config.ManipulatorConfig, 0)
	if c.parent != nil {
		manipulators = append(manipulators, c.parent.ManipulatorsFor(includes, excludes)...)
	}
	manipulators = append(manipulators, c.tags.ManipulatorsFor(includes, excludes)...)
	return manipulators
}

func (c *ScopeContext) BindManipulators(manipulator *config.ManipulatorConfig, tags []string) {
	c.tags.BindManipulators(manipulator, tags)
}
