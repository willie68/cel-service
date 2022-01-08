package celproc

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"google.golang.org/protobuf/types/known/structpb"
	"gopkg.in/yaml.v3"

	"github.com/stretchr/testify/assert"
	"github.com/willie68/cel-service/pkg/model"
	"github.com/willie68/cel-service/pkg/protofiles"
)

func TestYaml(t *testing.T) {
	ast := assert.New(t)

	celModels := readYaml("../../test/data/data1.yaml", t)

	for _, cm := range celModels {

		result, err := ProcCel(cm.Request)
		ast.Nil(err)
		ast.NotNil(result)

		ast.Equal(cm.Result, result.Result)
	}
}

func TestJson(t *testing.T) {
	ast := assert.New(t)

	celModels := readJson("../../test/data/data1.json", t)

	for _, cm := range celModels {
		cm.Request.Context = convertJson2Map(cm.Request.Context)
		result, err := ProcCel(cm.Request)
		ast.Nil(err)
		ast.NotNil(result)

		ast.Equal(cm.Result, result.Result)
	}
}

func TestJsonMany(t *testing.T) {
	ast := assert.New(t)

	celModels := readJson("../../test/data/data1.json", t)

	for _, cm := range celModels {
		cm.Request.Context = convertJson2Map(cm.Request.Context)
		result, err := ProcCel(cm.Request)
		ast.Nil(err)
		ast.NotNil(result)

		ast.Equal(cm.Result, result.Result)
	}
}

func TestGRPCJson(t *testing.T) {
	ast := assert.New(t)

	celModels := readJson("../../test/data/data1.json", t)

	for _, cm := range celModels {
		context := convertJson2Map(cm.Request.Context)
		grpcContext, err := structpb.NewStruct(context)
		ast.Nil(err)
		celRequest := protofiles.CelRequest{
			Context:    grpcContext,
			Expression: cm.Request.Expression,
		}

		result, err := GRPCProcCel(&celRequest)
		ast.Nil(err)
		ast.NotNil(result)

		ast.Equal(cm.Result, result.Result)
	}
}

func readYaml(filename string, t *testing.T) []model.TestCelModel {
	ast := assert.New(t)
	ya, err := ioutil.ReadFile(filename)
	ast.Nil(err)
	var celModels []model.TestCelModel
	err = yaml.Unmarshal(ya, &celModels)
	ast.Nil(err)
	return celModels
}

func readJson(filename string, t *testing.T) []model.TestCelModel {
	ast := assert.New(t)
	ya, err := ioutil.ReadFile(filename)
	ast.Nil(err)
	var celModels []model.TestCelModel
	decoder := json.NewDecoder(bytes.NewReader(ya))
	decoder.UseNumber()
	err = decoder.Decode(&celModels)
	//err = json.Unmarshal(ya, &celModels)
	ast.Nil(err)
	return celModels
}
