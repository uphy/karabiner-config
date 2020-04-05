package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type (
	RootConfig struct {
		Title       string       `yaml:"title"`
		Maintainers []string     `yaml:"maintainers,omitempty"`
		Rules       []RuleConfig `yaml:"rules,omitempty"`
	}
	RuleConfig struct {
		Description  *string             `yaml:"description,omitempty"`
		Manipulators []ManipulatorConfig `yaml:"manipulators,omitempty"`
	}
	ManipulatorConfig struct {
		From            *FromConfig `yaml:"from,omitempty"`
		To              []ToConfig  `yaml:"to,omitempty"`
		ToIfAlone       []ToConfig  `yaml:"to_if_alone,omitempty"`
		ToIfHeldDown    []ToConfig  `yaml:"to_if_held_down,omitempty"`
		ToAfterKeyUp    []ToConfig  `yaml:"to_after_key_up,omitempty"`
		ToDelayedAction *struct {
			ToIfInvoked  []ToConfig `yaml:"to_if_invoked,omitempty"`
			ToIfCanceled []ToConfig `yaml:"to_if_canceled,omitempty"`
		} `yaml:"to_delayed_action,omitempty"`
		Conditions []ConditionConfig `yaml:"conditions,omitempty"`
		Parameters Parameters        `yaml:"parameters,omitempty"`

		Scope  *ScopeConfig                 `yaml:"scope,omitempty"`
		Switch map[string]ManipulatorConfig `yaml:"switch,omitempty"`
	}
	FromConfig struct {
		// keys: https://github.com/pqrs-org/Karabiner-Elements/blob/master/src/apps/PreferencesWindow/Resources/simple_modifications.json
		KeyCode         *string          `yaml:"key_code,omitempty"`
		ConsumerKeyCode *string          `yaml:"consumer_key_code,omitempty"`
		PointingButton  *string          `yaml:"pointing_button,omitempty"`
		Any             *string          `yaml:"any,omitempty"`
		Modifiers       *ModifiersConfig `yaml:"modifiers,omitempty"`
		Simultaneous    []struct {
			KeyCode         *string `yaml:"key_code,omitempty"`
			ConsumerKeyCode *string `yaml:"consumer_key_code,omitempty"`
			PointingButton  *string `yaml:"pointing_button,omitempty"`
			Any             *string `yaml:"any,omitempty"`
		} `yaml:"simultaneous,omitempty"`
		SimultaneousOptions *struct {
			DetectKeyDownUninterruptedly *bool      `yaml:"detect_key_down_uninterruptedly,omitempty"`
			KeyDownOrder                 *string    `yaml:"key_down_order,omitempty"`
			KeyUpOrder                   *string    `yaml:"key_up_order,omitempty"`
			KeyUpWhen                    *string    `yaml:"key_up_when,omitempty"`
			ToAfterKeyUp                 []ToConfig `yaml:"to_after_key_up,omitempty"`
		} `yaml:"simultaneous_options,omitempty"`

		Key *string `yaml:"key,omitempty"`
	}
	ToConfig struct {
		KeyCode           *string            `yaml:"key_code,omitempty"`
		ConsumerKeyCode   *string            `yaml:"consumer_key_code,omitempty"`
		PointingButton    *string            `yaml:"pointing_button,omitempty"`
		ShellCommand      *string            `yaml:"shell_command,omitempty"`
		SelectInputSource *InputSourceConfig `yaml:"select_input_source"`
		SetVariable       *struct {
			Name  *string `yaml:"name,omitempty"`
			Value *string `yaml:"value,omitempty"`
		} `yaml:"set_variable"`
		MouseKey *struct {
			X               *int `yaml:"x,omitempty"`
			Y               *int `yaml:"y,omitempty"`
			VerticalWheel   *int `yaml:"vertical_wheel,omitempty"`
			HorizontalWheel *int `yaml:"horizontal_wheel,omitempty"`
			SpeedMultiplier *int `yaml:"speed_multiplier,omitempty"`
		} `yaml:"mouse_key,omitempty"`
		Modifiers            *ModifiersConfig `yaml:"modifiers,omitempty"`
		Lazy                 *bool            `yaml:"lazy,omitempty"`
		Repeat               *bool            `yaml:"repeat,omitempty"`
		Halt                 *bool            `yaml:"halt,omitempty"`
		HoldDownMilliseconds *int             `yaml:"hold_down_milliseconds,omitempty"`

		Key *string `yaml:"key,omitempty"`
		Set *string `yaml:"set,omitempty"`
	}
	InputSourceConfig struct {
		Language      *string `yaml:"language"`
		InputSourceID *string `yaml:"input_source_id"`
		InputModeID   *string `yaml:"input_mode_id"`
	}
	// modifier names: https://karabiner-elements.pqrs.org/docs/json/complex-modifications-manipulator-definition/to/modifiers/
	ModifiersConfig struct {
		Mandatory []string `yaml:"mandatory,omitempty"`
		Optional  []string `yaml:"optional,omitempty"`
	}
	ConditionConfig struct {
		Type              *string  `yaml:"type,omitempty"`
		BundleIdentifiers []string `yaml:"bundle_identifiers,omitempty"`
		FilePaths         []string `yaml:"file_paths,omitempty"`
		Identifiers       []struct {
			VendorID    *string `yaml:"vendor_id,omitempty"`
			ProductID   *string `yaml:"product_id,omitempty"`
			Description *string `yaml:"description"`
		} `yaml:"identifiers,omitempty"`
		KeyboardTypes []string            `yaml:"keyboard_types,omitempty"`
		InputSources  []InputSourceConfig `yaml:"input_sources,omitempty"`
		Name          *string             `yaml:"name,omitempty"`
		Value         *interface{}        `yaml:"value,omitempty"`
		Description   *string             `yaml:"description,omitempty"`

		AppIf     *AppConfig `yaml:"app_if,omitempty"`
		AppUnless *AppConfig `yaml:"app_unless,omitempty"`
		VarIf     *string    `yaml:"var_if,omitempty"`
	}
	AppConfig struct {
		Identifiers []string `yaml:"identifiers,omitempty"`
		Paths       []string `yaml:"paths,omitempty"`
	}
	ScopeConfig struct {
		Description *string             `yaml:"description,omitempty"`
		Conditions  []ConditionConfig   `yaml:"conditions,omitempty"`
		Children    []ManipulatorConfig `yaml:"children,omitempty"`
	}
	Parameters map[string]interface{}
)

func Load(file string) (*RootConfig, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var c RootConfig
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *RootConfig) Save(file string) error {
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, b, 0644)
}

func (m *ManipulatorConfig) Merge(override *ManipulatorConfig) *ManipulatorConfig {
	merged := *m
	merged.Conditions = append(merged.Conditions, override.Conditions...)
	if v := override.From; v != nil {
		merged.From = v
	}
	if v := override.To; v != nil {
		merged.To = append(merged.To, v...)
	}
	if v := override.ToAfterKeyUp; v != nil {
		merged.ToAfterKeyUp = append(merged.ToAfterKeyUp, v...)
	}
	if v := override.ToDelayedAction; v != nil {
		merged.ToDelayedAction = v
	}
	if v := override.ToIfAlone; v != nil {
		merged.ToIfAlone = append(merged.ToIfAlone, v...)
	}
	if v := override.ToIfHeldDown; v != nil {
		merged.ToIfHeldDown = append(merged.ToIfHeldDown, v...)
	}
	if v := override.Parameters; v != nil {
		if merged.Parameters == nil {
			merged.Parameters = Parameters{}
		}
		for name, value := range v {
			merged.Parameters[name] = value
		}
	}
	return &merged
}
