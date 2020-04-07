package gen

import (
	"sort"

	"github.com/uphy/karabiner-config/config"
)

func generateManipulators(ctx *ScopeContext, manipulatorConfig *config.ManipulatorConfig) ([]JSONObject, error) {
	if len(manipulatorConfig.Switch) == 0 {
		return generateManipulatorsRaw(ctx, manipulatorConfig)
	}
	manipulators := make([]JSONObject, 0)

	// for testing, sort 'switch' keys
	expressions := make([]string, 0)
	for expression := range manipulatorConfig.Switch {
		expressions = append(expressions, expression)
	}
	sort.Strings(expressions)
	for _, expression := range expressions {
		m := manipulatorConfig.Switch[expression]
		merged := manipulatorConfig.Merge(&m)
		if expression != "default" {
			merged.Conditions = append(merged.Conditions, config.ConditionConfig{
				VarIf: &expression,
			})
		}
		manipulator, err := generateManipulatorsRaw(ctx, merged)
		if err != nil {
			return nil, err
		}
		manipulators = append(manipulators, manipulator...)
	}
	return manipulators, nil
}

func generateManipulatorsRaw(ctx *ScopeContext, manipulatorConfig *config.ManipulatorConfig) ([]JSONObject, error) {
	merged := ctx.MergeManipulatorConfig(manipulatorConfig)
	manipulators := make([]JSONObject, 0)
	if merged.Scope != nil {
		m, err := generateScope(ctx, merged.Scope)
		if err != nil {
			return nil, err
		}
		manipulators = append(manipulators, m...)
	}
	m := make(JSONObject)
	m["type"] = "basic"
	if merged.From != nil {
		from, err := generateFrom(ctx, merged.From)
		if err != nil {
			return nil, err
		}
		m["from"] = from
	}
	// move 'delayed' 'to' to the 'to_delayed_action'
	for i := len(merged.To) - 1; i >= 0; i-- {
		t := merged.To[i]
		if t.Delayed != nil && *t.Delayed {
			// remove 'to'
			merged.To = append(merged.To[:i], merged.To[i+1:]...)
			if merged.ToDelayedAction == nil {
				merged.ToDelayedAction = new(config.DelayedActionConfig)
			}
			// insert the 'to' to the head of existing 'to_delayed_action'
			merged.ToDelayedAction.ToIfCanceled = append([]config.ToConfig{t}, merged.ToDelayedAction.ToIfCanceled...)
			merged.ToDelayedAction.ToIfInvoked = append([]config.ToConfig{t}, merged.ToDelayedAction.ToIfInvoked...)
		}
	}
	if err := generateTos(ctx, merged.To, m, "to"); err != nil {
		return nil, err
	}
	if err := generateTos(ctx, merged.ToIfAlone, m, "to_if_alone"); err != nil {
		return nil, err
	}
	if err := generateTos(ctx, merged.ToIfHeldDown, m, "to_if_held_down"); err != nil {
		return nil, err
	}
	if err := generateTos(ctx, merged.ToAfterKeyUp, m, "to_after_key_up"); err != nil {
		return nil, err
	}
	if merged.ToDelayedAction != nil {
		toDelayedAction := make(JSONObject)
		if err := generateTos(ctx, merged.ToDelayedAction.ToIfInvoked, toDelayedAction, "to_if_invoked"); err != nil {
			return nil, err
		}
		if err := generateTos(ctx, merged.ToDelayedAction.ToIfCanceled, toDelayedAction, "to_if_canceled"); err != nil {
			return nil, err
		}
		m["to_delayed_action"] = toDelayedAction
	}
	if merged.Conditions != nil && len(merged.Conditions) > 0 {
		conditions := make(JSONArray, 0)
		for _, conditionConfig := range merged.Conditions {
			condition, err := generateCondition(ctx, &conditionConfig)
			if err != nil {
				return nil, err
			}
			conditions = append(conditions, condition)
		}
		m["conditions"] = conditions
	}
	if len(merged.Parameters) > 0 {
		m["parameters"] = merged.Parameters
	}
	if len(m) > 1 {
		manipulators = append(manipulators, m)
	}
	return manipulators, nil
}

func generateFrom(ctx *ScopeContext, fromConfig *config.FromConfig) (interface{}, error) {
	from := make(JSONObject)
	if fromConfig.Key != nil {
		key, err := parseKey(*fromConfig.Key)
		if err != nil {
			return nil, err
		}
		ToObject(key).InjectMembersTo(from)
	}
	if fromConfig.KeyCode != nil {
		from["key_code"] = *fromConfig.KeyCode
	}
	if fromConfig.ConsumerKeyCode != nil {
		from["consumer_key_code"] = *fromConfig.ConsumerKeyCode
	}
	if fromConfig.PointingButton != nil {
		from["pointing_button"] = *fromConfig.PointingButton
	}
	if fromConfig.Any != nil {
		from["any"] = *fromConfig.Any
	}
	if fromConfig.Modifiers != nil {
		modifiers := make(JSONObject)
		if m := fromConfig.Modifiers.Mandatory; m != nil {
			modifiers["mandatory"] = m
		}
		if o := fromConfig.Modifiers.Optional; o != nil {
			modifiers["optional"] = o
		}
		from["modifiers"] = modifiers
	}
	if simultaneousConfig := fromConfig.Simultaneous; simultaneousConfig != nil {
		simultaneous := make(JSONArray, 0)
		for _, s := range simultaneousConfig {
			obj := make(JSONObject)
			if s.KeyCode != nil {
				obj["key_code"] = *s.KeyCode
			}
			if s.ConsumerKeyCode != nil {
				obj["consumer_key_code"] = *s.ConsumerKeyCode
			}
			if s.PointingButton != nil {
				obj["pointing_button"] = *s.PointingButton
			}
			if s.Any != nil {
				obj["any"] = *s.Any
			}
			simultaneous = append(simultaneous, obj)
		}
		from["simultaneous"] = simultaneous
	}
	if simultaneousOptionsConfig := fromConfig.SimultaneousOptions; simultaneousOptionsConfig != nil {
		simultaneousOptions := make(JSONObject)
		if o := simultaneousOptionsConfig.DetectKeyDownUninterruptedly; o != nil {
			simultaneousOptions["detect_key_down_uninterruptedly"] = o
		}
		if o := simultaneousOptionsConfig.KeyDownOrder; o != nil {
			simultaneousOptions["key_down_order"] = o
		}
		if o := simultaneousOptionsConfig.KeyUpOrder; o != nil {
			simultaneousOptions["key_up_order"] = o
		}
		if o := simultaneousOptionsConfig.KeyUpWhen; o != nil {
			simultaneousOptions["key_up_when"] = o
		}
		if err := generateTos(ctx, simultaneousOptionsConfig.ToAfterKeyUp, simultaneousOptions, "to_after_key_up"); err != nil {
			return nil, err
		}
		from["simultaneous_options"] = simultaneousOptions
	}
	return from, nil
}

func generateTos(ctx *ScopeContext, toConfigs []config.ToConfig, dest JSONObject, name string) error {
	if toConfigs != nil && len(toConfigs) > 0 {
		to := make([]JSONObject, 0)
		for _, toConfig := range toConfigs {
			t, err := generateTo(ctx, &toConfig)
			if err != nil {
				return err
			}
			to = append(to, t)
		}
		dest[name] = to
	}
	return nil
}

func generateTo(ctx *ScopeContext, toConfig *config.ToConfig) (JSONObject, error) {
	to := make(JSONObject)
	if toConfig.Key != nil {
		k, err := parseKey(*toConfig.Key)
		if err != nil {
			return nil, err
		}
		to["key_code"] = k.KeyCode
		if len(k.Modifiers.Mandatory) > 0 {
			to["modifiers"] = k.Modifiers.Mandatory
		}
	}
	if toConfig.Set != nil {
		s, err := parseSubstitution(*toConfig.Set)
		if err != nil {
			return nil, err
		}
		to["set_variable"] = s.ToJSON()
	}
	if v := toConfig.KeyCode; v != nil {
		to["key_code"] = v
	}
	if v := toConfig.ConsumerKeyCode; v != nil {
		to["consumer_key_code"] = v
	}
	if v := toConfig.PointingButton; v != nil {
		to["pointing_button"] = v
	}
	if v := toConfig.ShellCommand; v != nil {
		to["shell_command"] = v
	}
	if v := toConfig.SelectInputSource; v != nil {
		to["select_input_source"] = generateInputSource(toConfig.SelectInputSource)
	}
	if v := toConfig.SetVariable; v != nil {
		s := &Substitution{*v.Name, parseValue(*v.Value)}
		to["set_variable"] = s.ToJSON()
	}
	if v := toConfig.MouseKey; v != nil {
		to["mouse_key"] = v
	}
	if v := toConfig.MouseKey; v != nil {
		mouseKey := make(JSONObject)
		if v.X != nil {
			mouseKey["x"] = v.X
		}
		if v.Y != nil {
			mouseKey["y"] = v.Y
		}
		if v.VerticalWheel != nil {
			mouseKey["vertical_wheel"] = v.VerticalWheel
		}
		if v.HorizontalWheel != nil {
			mouseKey["horizontal_wheel"] = v.HorizontalWheel
		}
		if v.SpeedMultiplier != nil {
			mouseKey["speed_multiplier"] = v.SpeedMultiplier
		}
		to["mouse_key"] = mouseKey
	}
	if v := toConfig.Modifiers; v != nil {
		modifiers := make(JSONObject)
		if m := v.Mandatory; m != nil {
			modifiers["mandatory"] = m
		}
		if o := v.Optional; o != nil {
			modifiers["optional"] = o
		}
		to["modifiers"] = modifiers
	}
	if v := toConfig.Lazy; v != nil {
		to["lazy"] = *v
	}
	if v := toConfig.Repeat; v != nil {
		to["repeat"] = *v
	}
	if v := toConfig.Halt; v != nil {
		to["halt"] = *v
	}
	if v := toConfig.HoldDownMilliseconds; v != nil {
		to["hold_down_milliseconds"] = *v
	}
	return to, nil
}

func generateScope(ctx *ScopeContext, scopeConfig *config.ScopeConfig) ([]JSONObject, error) {
	childCtx := ctx.ChildScope(scopeConfig.Conditions)
	manipulators := make([]JSONObject, 0)
	for _, childConfig := range scopeConfig.Children {
		m, err := generateManipulators(childCtx, &childConfig)
		if err != nil {
			return nil, err
		}
		manipulators = append(manipulators, m...)
	}
	return manipulators, nil
}

func generateInputSource(src *config.InputSourceConfig) JSONObject {
	s := make(JSONObject)
	if src.Language != nil {
		s["language"] = src.Language
	}
	if src.InputSourceID != nil {
		s["input_source_id"] = src.InputSourceID
	}
	if src.InputModeID != nil {
		s["input_mode_id"] = src.InputModeID
	}
	return s
}
