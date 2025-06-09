package utils

import (
	"fmt"
	"reflect"
	"regexp"
)

func getQuotedAddressRegexExpressionReplacer(contractName string) string {
	const regexPattern = `"(?:(?:\./|(?:\.\./)+)?(?:[a-zA-Z0-9_\-]+/)*%s(\.cdc)?)"`
	return fmt.Sprintf(regexPattern, contractName)
}

func getRegexExpressionReplacer(contractName string) string {
	const regexPattern = `%s`
	return fmt.Sprintf(regexPattern, contractName)
}

// ParseCadenceTemplateV3 parses the Cadence template and replaces placeholders
func ParseCadenceTemplateV3(template []byte, data ...interface{}) ([]byte, error) {
	if err := validateStruct(data); err != nil {
		return nil, err
	}
	updatedTemplate, err := replacePlaceholders(string(template), data)
	if err != nil {
		return nil, err
	}

	return []byte(updatedTemplate), nil
}

// validateStruct ensures the provided data is a struct
func validateStruct(data []interface{}) error {
	for _, item := range data {
		if item == nil {
			continue
		}
		if reflect.ValueOf(item).Kind() != reflect.Struct {
			return fmt.Errorf("data must be a struct")
		}
	}
	return nil
}

// replacePlaceholders replaces the placeholders in the template with actual values
func replacePlaceholders(template string, data []interface{}) (string, error) {
	for index, item := range data {
		if item == nil {
			continue
		}
		val := reflect.ValueOf(item)
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			fieldName := field.Name
			fieldValue := val.Field(i).String()
			if index == 0 {
				replacer := regexp.MustCompile(getQuotedAddressRegexExpressionReplacer(fieldName))
				template = replacer.ReplaceAllString(template, "0x"+fieldValue)
			} else {
				replacer := regexp.MustCompile(getRegexExpressionReplacer(fieldName))
				template = replacer.ReplaceAllString(template, fieldValue)
			}
		}

	}

	return template, nil
}
