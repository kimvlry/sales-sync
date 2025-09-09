package pkg

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
)

type Config struct {
	App struct {
		Name string `yaml:"name"`
		Env  string `yaml:"env"`
		Port int    `yaml:"port"`
	} `yaml:"app"`

	DB struct {
		URL string `yaml:"url"`
	} `yaml:"db"`

	Redis struct {
		Addr string `yaml:"addr"`
		DB   int    `yaml:"db"`
	} `yaml:"redis"`

	Kafka struct {
		Brokers []string `yaml:"brokers"`
	} `yaml:"kafka"`
}

func replaceEnvVars(content []byte) []byte {
	re := regexp.MustCompile(`\$\{(\w+)\}`)
	return re.ReplaceAllFunc(content, func(b []byte) []byte {
		match := re.FindSubmatch(b)
		if len(match) == 2 {
			envVar := string(match[1])
			if val, ok := os.LookupEnv(envVar); ok {
				return []byte(val)
			}
			return []byte(fmt.Sprintf("MISSING_ENV_%s", envVar))
		}
		return b
	})
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	data = replaceEnvVars(data)

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
