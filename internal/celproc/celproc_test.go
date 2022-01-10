package celproc

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/structpb"
	"gopkg.in/yaml.v3"

	"github.com/stretchr/testify/assert"
	"github.com/willie68/cel-service/pkg/model"
	"github.com/willie68/cel-service/pkg/protofiles"
)

const MAX_TEST_COUNT = 10000

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

func TestCacheUpdate(t *testing.T) {
	ast := assert.New(t)
	context := make(map[string]interface{})

	context["user"] = "willie"
	context["number"] = 1

	celModel := model.CelModel{
		Context:    context,
		Identifier: "willie",
		Expression: "number == 1 && user == \"willie\"",
	}
	result, err := ProcCel(celModel)
	ast.Nil(err)
	ast.NotNil(result)

	ast.Equal(true, result.Result)

	context["user2"] = "wutz"

	celModel = model.CelModel{
		Context:    context,
		Identifier: "willie",
		Expression: "number == 1 && user2 == \"wutz\"",
	}

	result, err = ProcCel(celModel)
	ast.Nil(err)
	ast.NotNil(result)

	ast.Equal(true, result.Result)
}

func TestWrongID(t *testing.T) {
	ast := assert.New(t)
	context := make(map[string]interface{})

	context["user"] = "willie"
	context["number"] = 1

	celModel := model.CelModel{
		Context:    context,
		Identifier: "willie",
		Expression: "number == 1 && user == \"willie\"",
	}
	result, err := ProcCel(celModel)
	ast.Nil(err)
	ast.NotNil(result)

	ast.Equal(true, result.Result)

	celModel = model.CelModel{
		Context:    context,
		Identifier: "willie",
		Expression: "hurtz ==  \"wutz\"",
	}

	result, err = ProcCel(celModel)
	ast.NotNil(err)
	ast.NotNil(result)

	ast.Equal(false, result.Result)
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

func BenchmarkJsonManyWithoutCache(t *testing.B) {
	ast := assert.New(t)

	celModels := readJsonB("../../test/data/data1.json", t)
	stt := time.Now()
	for i := 0; i < MAX_TEST_COUNT; i++ {
		for _, cm := range celModels {
			cm.Request.Context = convertJson2Map(cm.Request.Context)
			cm.Request.Identifier = "" //fmt.Sprintf("%d_%d", i, x)
			result, err := ProcCel(cm.Request)
			ast.Nil(err)
			ast.NotNil(result)

			ast.Equal(cm.Result, result.Result)
		}
	}
	ste := time.Now()

	t.Logf("execution: %d", ste.Sub(stt).Milliseconds())
}

func BenchmarkJsonManyWithCache(t *testing.B) {
	ast := assert.New(t)

	celModels := readJsonB("../../test/data/data1.json", t)
	stt := time.Now()
	for i := 0; i < MAX_TEST_COUNT; i++ {
		for _, cm := range celModels {
			cm.Request.Context = convertJson2Map(cm.Request.Context)
			result, err := ProcCel(cm.Request)
			ast.Nil(err)
			ast.NotNil(result)

			ast.Equal(cm.Result, result.Result)
		}
	}
	ste := time.Now()

	t.Logf("execution: %d", ste.Sub(stt).Milliseconds())
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

func readJsonB(filename string, t *testing.B) []model.TestCelModel {
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
