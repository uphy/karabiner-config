package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/uphy/karabiner-config/config"
	"github.com/uphy/karabiner-config/gen"
)

type (
	ConfigWriter interface {
		Write(conf *config.RootConfig) error
	}
	OverwriteConfigWriter struct {
		output string
	}
	KarabinierJSONWriter struct {
		output  string
		profile string
	}
)

func (w *KarabinierJSONWriter) read() (gen.JSONObject, error) {
	f, err := os.Open(w.output)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var obj gen.JSONObject
	if err := json.Unmarshal(b, &obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func (w *KarabinierJSONWriter) ReplaceRule(original gen.JSONObject, generated gen.JSONObject) (gen.JSONObject, error) {
	profiles := original["profiles"].([]interface{})
	found := false
	for i, profile := range profiles {
		p := profile.(map[string]interface{})
		if p["name"].(string) == w.profile {
			c := p["complex_modifications"].(map[string]interface{})
			c["rules"] = generated["rules"]
			p["complex_modifications"] = c
			break
		}
		profiles[i] = p
		found = true
	}
	if !found {
		return nil, errors.New("profile not found in karabiner.json")
	}
	original["profiles"] = profiles
	return original, nil
}

func (w *KarabinierJSONWriter) Write(conf *config.RootConfig) error {
	original, err := w.read()
	if err != nil {
		return err
	}
	generated, err := gen.GenerateObject(conf)
	if err != nil {
		return err
	}
	replaced, err := w.ReplaceRule(original, generated)
	if err != nil {
		return err
	}
	f, err := os.Create(w.output)
	if err != nil {
		return err
	}
	defer f.Close()
	return replaced.Write(f)
}

func (w *OverwriteConfigWriter) Write(conf *config.RootConfig) error {
	f, err := os.Create(w.output)
	if err != nil {
		return fmt.Errorf("failed to generate while creating file: file=%s, err=%w", w.output, err)
	}
	defer f.Close()
	if err := gen.Generate(conf, f); err != nil {
		return fmt.Errorf("failed to generate: err=%w", err)
	}
	return nil
}
