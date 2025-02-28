package yaml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateYAMLTemplate(t *testing.T) {
	type Config struct {
		Host    string   `mapstructure:"host" default:"localhost" help:"The hostname"`
		Port    int      `mapstructure:"port" default:"8080" help:"The port number"`
		Enabled bool     `mapstructure:"enabled" default:"true" help:"Enable the feature"`
		Options []string `mapstructure:"options" default:"1,2,3" help:"List of options"`
		Meta    struct {
			Version string `mapstructure:"version" default:"1.0" help:"App version"`
		} `mapstructure:"meta"`
		MapField map[string]string `mapstructure:"map_field" help:"Example map field"`
	}
	cfg := Config{}
	yamlTemplate := GenerateYAMLTemplate(cfg, true)

	expected := `host: "localhost" # The hostname
port: 8080        # The port number
enabled: true     # Enable the feature
options:          # List of options
  - 1
  - 2
  - 3
meta:
  version: "1.0"  # App version
map_field:        # Example map field
  key: value      # Map example
`

	assert.Equal(t, expected, yamlTemplate)
}

// Test basic YAML generation with primitive types.
func TestGenerateYAMLTemplate_Basic(t *testing.T) {
	cfg := struct {
		Host string `yaml:"host" default:"localhost"`
		Port int    `yaml:"port" default:"8080"`
	}{}
	yamlTemplate := GenerateYAMLTemplate(cfg, true)

	expected := `host: "localhost"
port: 8080
`

	assert.Equal(t, expected, yamlTemplate)
}

// Test YAML generation with array of primitives.
func TestGenerateYAMLTemplate_ArrayOfPrimitives(t *testing.T) {
	tests := []struct {
		name     string
		cfg      interface{}
		expected string
	}{
		{
			name: "OptionsWithDefault",
			cfg: struct {
				OptionsWithDefault []string `yaml:"options" default:"value1" help:"Array of options"`
			}{},
			expected: `options:   # Array of options
  - value1
`,
		},
		{
			name: "OptionsWithDefaults",
			cfg: struct {
				OptionsWithDefaults []string `yaml:"options" default:"1,2,3" help:"Array of options"`
			}{},
			expected: `options: # Array of options
  - 1
  - 2
  - 3
`,
		},
		{
			name: "OptionsWithoutDefaults",
			cfg: struct {
				OptionsWithoutDefaults []string `yaml:"options" help:"Array of options"`
			}{},
			expected: `options:    # Array of options
  - example
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			yamlTemplate := GenerateYAMLTemplate(tt.cfg, true)
			assert.Equal(t, tt.expected, yamlTemplate)
		})
	}
}

// Test YAML generation with array of structs.
func TestGenerateYAMLTemplate_ArrayOfStructs(t *testing.T) {
	type Item struct {
		Name  string `yaml:"name" default:"item1" help:"Item name"`
		Value int    `yaml:"value"`
	}
	cfg := struct {
		Items []Item `yaml:"items" help:"Array of items"`
	}{}
	yamlTemplate := GenerateYAMLTemplate(cfg, true)

	expected := `items:            # Array of items
  -
    name: "item1" # Item name
    value: null
`

	assert.Equal(t, expected, yamlTemplate)
}

// Test YAML generation with maps.
func TestGenerateYAMLTemplate_Map(t *testing.T) {
	cfg := struct {
		Settings map[string]string `yaml:"settings" help:"Map of settings"`
	}{}
	yamlTemplate := GenerateYAMLTemplate(cfg, true)

	expected := `settings:    # Map of settings
  key: value # Map example
`

	assert.Equal(t, expected, yamlTemplate)
}

// Test YAML generation with nested anonymous structs.
func TestGenerateYAMLTemplate_AnonymousStruct(t *testing.T) {
	cfg := struct {
		Meta struct {
			Version string `yaml:"version" default:"1.0" help:"Version"`
		} `yaml:"meta"`
	}{}
	yamlTemplate := GenerateYAMLTemplate(cfg, true)

	expected := `meta:
  version: "1.0" # Version
`

	assert.Equal(t, expected, yamlTemplate)
}

// Test YAML generation with ignored fields.
func TestGenerateYAMLTemplate_IgnoredFields(t *testing.T) {
	cfg := struct {
		Visible string `yaml:"visible" default:"shown"`
		Hidden  string `yaml:"-" default:"hidden"`
	}{}
	yamlTemplate := GenerateYAMLTemplate(cfg, true)

	expected := `visible: "shown"
`

	assert.Equal(t, expected, yamlTemplate)
}

// Test YAML generation with different tag priorities.
func TestGenerateYAMLTemplate_TagPriority(t *testing.T) {
	cfg := struct {
		Field string `yaml:"yaml_tag" mapstructure:"mapstructure_tag" default:"yaml_value" help:"YAML tag priority test"`
	}{}
	yamlTemplate := GenerateYAMLTemplate(cfg, true)

	expected := `yaml_tag: "yaml_value" # YAML tag priority test
`

	assert.Equal(t, expected, yamlTemplate)
}

// Test YAML generation with and default.
func TestGenerateYAMLTemplate_NullWithoutDefault(t *testing.T) {
	cfg := struct {
		Username string `yaml:"username" default:"default_username" help:"User login name"`
		Nickname string `yaml:"nickname" help:"User nickname"`
	}{}
	yamlTemplate := GenerateYAMLTemplate(cfg, true)

	expected := `username: "default_username" # User login name
nickname: null               # User nickname
`

	assert.Equal(t, expected, yamlTemplate)
}
