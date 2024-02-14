package internal

import (
	"encoding/json"
	"errors"
	"os"
)

type Configurations struct {
	Owner            string `json:"owner"`
	Name             string `json:"name"`
	Auth             string `json:"auth"`
	GithubApiVersion string `json:"githubApiVersion"`
}

// LoadConfigurations read configs from filepath and returns an structure with the configs
func LoadConfigurations(filepath string) (*Configurations, error) {
	var config Configurations
	// read file
	fdata, err := os.ReadFile(filepath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// empty configuration
			return &config, nil
		}
		return nil, err
	}
	// unmarshal data
	err = json.Unmarshal(fdata, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
