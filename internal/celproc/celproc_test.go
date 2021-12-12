package celproc

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/stretchr/testify/assert"
	"github.com/willie68/cel-service/pkg/model"
)

func TestTokenYaml(t *testing.T) {
	ast := assert.New(t)

	ya, err := ioutil.ReadFile("../../test/data/data1.yaml")
	ast.Nil(err)
	var celModels []model.CelModel
	err = yaml.Unmarshal(ya, &celModels)
	ast.Nil(err)

	for _, cm := range celModels {

		result, err := ProcCel(cm)
		ast.Nil(err)
		ast.NotNil(result)

		ast.True(result.Result)
	}
}

func TestTokenJson(t *testing.T) {
	ast := assert.New(t)

	ya, err := ioutil.ReadFile("../../test/data/data1.json")
	ast.Nil(err)
	var celModels []model.CelModel
	err = json.Unmarshal(ya, &celModels)
	ast.Nil(err)

	for _, cm := range celModels {

		result, err := ProcCel(cm)
		ast.Nil(err)
		ast.NotNil(result)

		ast.True(result.Result)
	}
}
