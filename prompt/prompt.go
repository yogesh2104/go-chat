package prompt

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type PromptManager struct {
	Default string
	Modes   map[string]string
}

func LoadSystemPrompts(filePath string) (*PromptManager, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config struct {
		Default string            `yaml:"default"`
		Modes   map[string]string `yaml:"modes"`
	}

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &PromptManager{
		Default: config.Default,
		Modes:   config.Modes,
	}, nil
}

func (pm *PromptManager) BuildPrompt(mode string, userPrompt string) string {
	sysPrompt, exists := pm.Modes[mode]
	if !exists {
		sysPrompt = pm.Modes[pm.Default]
	}

	return fmt.Sprintf("%s", sysPrompt)
}
