package gen

import (
	"github.com/uphy/karabiner-config/config"
)

func generateCondition(ctx *ScopeContext, conditionConfig *config.ConditionConfig) (JSONObject, error) {
	condition := make(JSONObject)

	// app_if
	if found := generateAppCondition("frontmost_application_if", conditionConfig.AppIf, condition); found {
		return condition, nil
	}
	// app_unless
	if found := generateAppCondition("frontmost_application_unless", conditionConfig.AppUnless, condition); found {
		return condition, nil
	}
	// var_if
	if v := conditionConfig.VarIf; v != nil {
		cond, err := parseVariableCondition(*v)
		if err != nil {
			return nil, err
		}
		return cond.ToJSON(), nil
	}

	if v := conditionConfig.Type; v != nil {
		condition["type"] = *v
	}
	if v := conditionConfig.BundleIdentifiers; v != nil {
		condition["bundle_identifiers"] = v
	}
	if v := conditionConfig.FilePaths; v != nil {
		condition["file_paths"] = v
	}
	if v := conditionConfig.Identifiers; v != nil {
		identifiers := make(JSONArray, 0)
		for _, identifierConfig := range v {
			identifier := make(JSONObject)
			if identifierConfig.VendorID != nil {
				identifier["vendor_id"] = *identifierConfig.VendorID
			}
			if identifierConfig.ProductID != nil {
				identifier["product_id"] = *identifierConfig.ProductID
			}
			if identifierConfig.Description != nil {
				identifier["description"] = *identifierConfig.Description
			}
			identifiers = append(identifiers, identifier)
		}
		condition["identifiers"] = identifiers
	}
	if v := conditionConfig.KeyboardTypes; v != nil {
		condition["keyboard_types"] = v
	}
	if v := conditionConfig.InputSources; v != nil {
		inputSources := make(JSONArray, 0)
		for _, inputSourceConfig := range v {
			inputSources = append(inputSources, generateInputSource(&inputSourceConfig))
		}
		condition["input_sources"] = inputSources
	}
	if v := conditionConfig.Name; v != nil {
		condition["name"] = *v
	}
	if v := conditionConfig.Value; v != nil {
		condition["value"] = v
	}
	if v := conditionConfig.Description; v != nil {
		condition["description"] = v
	}
	return condition, nil
}

func generateAppCondition(typ string, appConfig *config.AppConfig, dest JSONObject) bool {
	if appConfig == nil {
		return false
	}
	dest["type"] = typ
	if appConfig.Identifiers != nil {
		dest["bundle_identifiers"] = appConfig.Identifiers
	}
	if appConfig.Paths != nil {
		dest["file_paths"] = appConfig.Paths
	}
	return true
}
