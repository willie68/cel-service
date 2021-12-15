package protofiles

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestYamlProto(t *testing.T) {
	ast := assert.New(t)
	jsonContext := map[string]interface{}{
		"data": map[string]interface{}{
			"value": 1,
		},
	}
	structValue, err := structpb.NewStruct(jsonContext)
	ast.Nil(err)
	celRequest := CelRequest{
		Context:    structValue,
		Expression: "data.value == 1",
	}

	ast.NotNil(celRequest)
	out, err := proto.Marshal(&celRequest)
	ast.Nil(err)
	ast.NotNil(out)

	celReq2 := &CelRequest{}

	err = proto.Unmarshal(out, celReq2)

	ast.Nil(err)
	ast.NotNil(celReq2)

	ast.Equal(celRequest.Expression, celReq2.Expression)
}
