package celproc

import (
	"errors"
	"fmt"

	log "github.com/willie68/cel-service/internal/logging"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/willie68/cel-service/pkg/model"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func ProcCel(celModel model.CelModel) (model.CelResult, error) {
	var declList = make([]*exprpb.Decl, len(celModel.Context))
	x := 0
	for k := range celModel.Context {
		declList[x] = decls.NewVar(k, decls.Dyn)
		x++
	}
	env, err := cel.NewEnv(
		cel.Declarations(
			declList...,
		),
	)
	if err != nil {
		log.Logger.Errorf("env declaration error: %s", err)
	}
	ast, issues := env.Compile(celModel.Expression)
	if issues != nil && issues.Err() != nil {
		log.Logger.Errorf("type-check error: %s", issues.Err())
	}
	prg, err := env.Program(ast)
	if err != nil {
		log.Logger.Errorf("program construction error: %s", err)
	}
	out, details, err := prg.Eval(celModel.Context)
	fmt.Printf("result: %v\r\n", details)
	if err != nil {
		log.Logger.Errorf("program evaluation error: %s", err)
	}
	switch v := out.(type) {
	case types.Bool:
		return model.CelResult{
			Result: v == types.True,
		}, nil
	case *types.Err:
		return model.CelResult{
			Result: false,
		}, err
	default:
		return model.CelResult{
			Result: false,
		}, errors.New("unknown result type")
	}
}
