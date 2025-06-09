package types

import (
	"encoding/json"
	"fmt"
	"regexp"
)

type RegexList []*regexp.Regexp

func (r *RegexList) UnmarshalJSON(data []byte) error {
	var raw []string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for _, pattern := range raw {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("invalid regex %q: %w", pattern, err)
		}
		*r = append(*r, re)
	}

	return nil
}

func (r *RegexList) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw []string
	if err := unmarshal(&raw); err != nil {
		return err
	}

	for _, pattern := range raw {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("invalid regex %q: %w", pattern, err)
		}
		*r = append(*r, re)
	}

	return nil
}

func (r RegexList) MarshalJSON() ([]byte, error) {
	var raw []string
	for _, re := range r {
		raw = append(raw, re.String())
	}
	return json.Marshal(raw)
}

func (r RegexList) MarshalYAML() (interface{}, error) {
	var raw []string
	for _, re := range r {
		raw = append(raw, re.String())
	}
	return raw, nil
}
