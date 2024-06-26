package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type TableName string

type ColumnName string

type RuleSetLength struct {
	Min uint `yaml:"min"`
	Max uint `yaml:"max"`
}

type RuleSetPeriod struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
}

type RuleSetRange struct {
	Min uint `yaml:"min"`
	Max uint `yaml:"max"`
}

type RuleSet struct {
	Nullable               bool           `yaml:"nullable"`
	Unique                 bool           `yaml:"unique"`
	Length                 *RuleSetLength `yaml:"length,omitempty"`
	Prefix                 *string        `yaml:"prefix,omitempty"`
	Postfix                *string        `yaml:"postfix,omitempty"`
	NoDuplicateCharacters  bool           `yaml:"no_duplicate_characters,omitempty"`
	NoSequentialCharacters uint           `yaml:"no_sequential_characters,omitempty"`
	IncludeCharset         *string        `yaml:"include_charset,omitempty"`
	ExcludeCharset         *string        `yaml:"exclude_charset,omitempty"`
	Values                 []string       `yaml:"values,omitempty"`
	Period                 *RuleSetPeriod `yaml:"period,omitempty"`
	Range                  *RuleSetRange  `yaml:"range,omitempty"`
}

type ColumnDescriptionReference struct {
	Table  TableName  `yaml:"table"`
	Column ColumnName `yaml:"column"`
}

type ColumnDescription struct {
	Type      string                      `yaml:"type"`
	RuleSet   RuleSet                     `yaml:"ruleSet"`
	Reference *ColumnDescriptionReference `yaml:"reference,omitempty"`
}

type TableDescription map[ColumnName]ColumnDescription

type TablesDescription map[TableName]TableDescription

type Config struct {
	Tables      TablesDescription `yaml:"tables"`
	TablesOrder []TableName       `yaml:"tables_order"`
}

func LoadFromFile(path string) (*Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return LoadFromBytes(content)
}

func LoadFromBytes(data []byte) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}
	return &cfg, nil
}

func SaveToBytes(cfg *Config) ([]byte, error) {
	return yaml.Marshal(cfg)
}
