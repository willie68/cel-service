package celproc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/willie68/cel-service/pkg/model"
)

func TestToken(t *testing.T) {
	ast := assert.New(t)

	result, err := ProcCel(model.CelModel{
		Context: map[string]interface{}{
			"data": map[string]interface{}{
				"value": 2,
			},
		},
		Expression: "data.v1.value == 1",
	})
	ast.Nil(err)
	ast.NotNil(result)

	ast.True(result.Result)
}
