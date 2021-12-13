package protofiles

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestTokenYaml(t *testing.T) {
	ast := assert.New(t)

	celRequest := CelRequest{
		Context: map[string]*ContextValue{
			"data": {
				Type: ContextValue_map,
				Vmap: map[string]*ContextValue{
					"value": {
						Type:  1,
						Value: []byte{0, 0, 0, 1},
					},
				},
			},
		},
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
