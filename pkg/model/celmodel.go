package model

type CelModel struct {
	Id         string                 `yaml:"id" json:"id"`
	Context    map[string]interface{} `yaml:"context" json:"context"`
	Expression string                 `yaml:"expression" json:"expression"`
	Identifier string                 `yaml:"identifier" json:"identifier"`
}

type CelResult struct {
	Id      string `yaml:"id" json:"id"`
	Error   string `yaml:"error" json:"error"`
	Message string `yaml:"message" json:"message"`
	Result  bool   `yaml:"result" json:"result"`
}

type TestCelModel struct {
	Request CelModel `yaml:"request" json:"request"`
	Result  bool     `yaml:"result" json:"result"`
}
