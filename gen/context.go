package gen

import (
	"github.com/uphy/karabiner-config/config"
)

type ScopeContext struct {
	parent      *ScopeContext
	scopeConfig *config.ScopeConfig
}

func newRootScopeContext() *ScopeContext {
	return newScopeContext(nil, nil)
}

func newScopeContext(parent *ScopeContext, scopeConfig *config.ScopeConfig) *ScopeContext {
	return &ScopeContext{parent, scopeConfig}
}

func (c *ScopeContext) ChildScope(childScopeConfig *config.ScopeConfig) *ScopeContext {
	return newScopeContext(c, childScopeConfig)
}

func (c *ScopeContext) MergeManipulatorConfig(manipulator *config.ManipulatorConfig) *config.ManipulatorConfig {
	merged := manipulator
	if c.parent != nil {
		merged = c.parent.MergeManipulatorConfig(merged)
	}
	if c.scopeConfig != nil {
		for _, condition := range c.scopeConfig.Conditions {
			merged.Conditions = append(merged.Conditions, condition)
		}
	}
	return merged
}
