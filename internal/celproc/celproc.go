package celproc

import (
	"encoding/json"
	"errors"
	"fmt"

	log "github.com/willie68/cel-service/internal/logging"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/willie68/cel-service/pkg/model"
	"github.com/willie68/cel-service/pkg/protofiles"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

var lrucache LRUList

func init() {
	lrucache = LRUList{
		MaxCount: int(10000),
	}
	lrucache.Init()
}

func GRPCProcCel(celRequest *protofiles.CelRequest) (*protofiles.CelResponse, error) {
	context := convertJson2Map(celRequest.Context.AsMap())
	celModel := model.CelModel{
		Context:    context,
		Expression: celRequest.Expression,
	}

	rep, err := ProcCel(celModel)
	celResponse := protofiles.CelResponse{
		Error:   rep.Error,
		Message: rep.Message,
		Result:  rep.Result,
	}
	return &celResponse, err
}

func convertJson2Map(src map[string]interface{}) (dst map[string]interface{}) {
	if src == nil {
		return nil
	}
	dst = make(map[string]interface{})
	for key, value := range src {
		switch v := value.(type) {
		case json.Number:
			iv, err := v.Int64()
			if err == nil {
				dst[key] = iv
			} else {
				fv, err := v.Float64()
				if err == nil {
					dst[key] = fv
				} else {
					dst[key] = v.String()
				}
			}
		case map[string]interface{}:
			dst[key] = convertJson2Map(v)
		default:
			dst[key] = value
		}
	}
	return
}

func ProcCel(celModel model.CelModel) (model.CelResult, error) {
	context := convertJson2Map(celModel.Context)
	ok := false
	var prg cel.Program
	var expression string
	id := celModel.Identifier
	if id != "" {
		var entry LRUEntry
		entry, ok = lrucache.Get(id)
		if ok {
			prg = entry.Program
			expression = entry.Expression
		}
		// Check if we have to update the cache
		if expression != celModel.Expression {
			ok = false
		}
	}
	if !ok {
		var declList = make([]*exprpb.Decl, len(context))
		x := 0
		for k := range context {
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
			log.Logger.Errorf("type-check error: %v", issues.Err())
			return model.CelResult{
				Error:   fmt.Sprintf("%v", issues.Err()),
				Message: issues.Err().Error(),
			}, issues.Err()
		}
		prg, err = env.Program(ast)
		if err != nil {
			log.Logger.Errorf("program construction error: %v", err)
			return model.CelResult{
				Error:   fmt.Sprintf("%v", err),
				Message: fmt.Sprintf("program construction error: %s", err.Error()),
			}, err
		}
		if id != "" {
			entry := LRUEntry{
				ID:         id,
				Expression: celModel.Expression,
				Program:    prg,
			}
			if lrucache.Has(celModel.Identifier) {
				lrucache.Update(entry)
			} else {
				lrucache.Add(entry)
				go lrucache.HandleContrains()
			}
		}
	}
	out, details, err := prg.Eval(context)
	//fmt.Printf("result: %v\ndetails: %v\nerror: %v\n", out, details, err)

	if err != nil {
		log.Logger.Errorf("program evaluation error: %v", err)

		return model.CelResult{
			Error:   fmt.Sprintf("%v", err),
			Message: fmt.Sprintf("program evaluation error: %s\r\ndetails: %s", err.Error(), details),
		}, err
	}
	switch v := out.(type) {
	case types.Bool:
		return model.CelResult{
			Message: fmt.Sprintf("result ok: %v", v),
			Result:  v == types.True,
		}, nil
	case *types.Err:
		return model.CelResult{
			Error:   fmt.Sprintf("%v", err),
			Message: fmt.Sprintf("unknown cel engine error: %v", err),
			Result:  false,
		}, err
	default:
		return model.CelResult{
			Message: "unknown result type",
			Result:  false,
		}, errors.New("unknown result type")
	}
}
