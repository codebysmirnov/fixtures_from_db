package config

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

type RuleSet struct {
	Nullable               bool           `yaml:"nullable"`
	Length                 *RuleSetLength `yaml:"length,omitempty"`
	Prefix                 *string        `yaml:"prefix,omitempty"`
	Postfix                *string        `yaml:"postfix,omitempty"`
	NoDuplicateCharacters  bool           `yaml:"no_duplicate_characters,omitempty"`
	NoSequentialCharacters uint           `yaml:"no_sequential_characters,omitempty"`
	IncludeCharset         *string        `yaml:"include_charset,omitempty"`
	ExcludeCharset         *string        `yaml:"exclude_charset,omitempty"`
	Values                 []string       `yaml:"values,omitempty"`
	Period                 *RuleSetPeriod `yaml:"period,omitempty"`
}

type ColumnDescription struct {
	Type    string  `yaml:"type"`
	RuleSet RuleSet `yaml:"ruleSet"`
}

type TableDescription map[ColumnName]ColumnDescription

type Config map[TableName]TableDescription
