package model

type CelModel struct {
	Context    map[string]interface{} `yaml:"context" json:"context"`
	Expression string                 `yaml:"expression" json:"expression"`
}

type CelResult struct {
	Error   error  `yaml:"error" json:"error"`
	Message string `yaml:"message" json:"message"`
	Result  bool   `yaml:"result" json:"result"`
}

type TestCelModel struct {
	Request CelModel `yaml:"request" json:"request"`
	Result  bool     `yaml:"result" json:"result"`
}
