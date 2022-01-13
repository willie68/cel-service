package apiv1

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/willie68/cel-service/pkg/model"

	"github.com/willie68/cel-service/internal/celproc"
	log "github.com/willie68/cel-service/internal/logging"
	"github.com/willie68/cel-service/internal/serror"
	"github.com/willie68/cel-service/internal/utils/httputils"
)

var (
	postEvalCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "cel_service_post_eval_total",
		Help: "The total number of post eval requests",
	})
	postEvalManyCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "cel_service_post_eval_many_total",
		Help: "The total number of post eval many requests",
	})
)

/*
EvalRoutes getting all routes for the config endpoint
*/
func EvalRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/evaluate", PostEval)
	router.Post("/evaluatemany", PostEvalMany)
	return router
}

// PostEval Evaluates the given context from payload against the specified CEL expression
// @Summary Post Evaluation
// @Description Evaluates the given context from payload against the specified CEL expression
// @Tags evaluation
// @Accept  json
// @Produce  json
// @Security apikey
// @Param payload body model.CelModel true "Context and expression"
// @Success 201 {object} model.CelResult "Evaluation result"
// @Failure 400 {object} serror.Serr "client error information as json"
// @Failure 500 {object} serror.Serr "server error information as json"
// @Router /evaluate [post]
func PostEval(response http.ResponseWriter, request *http.Request) {
	postEvalCounter.Inc()
	var celModel model.CelModel
	err := decode(request, &celModel)
	if err != nil {
		log.Logger.Errorf("error decoding context: %v", err)
		msg := fmt.Sprintf("error decoding context: %v", err)
		httputils.Err(response, request, serror.BadRequest(nil, "server-error", msg))
		return
	}
	if celModel.Expression == "" {
		httputils.Err(response, request, serror.BadRequest(nil, "empty expression not allowed"))
		return
	}
	res, err := celproc.ProcCel(celModel)
	log.Logger.Infof("req: %v, res: %v", celModel, res)

	if err != nil {
		log.Logger.Errorf("processing error: %v", err)
		render.Status(request, http.StatusBadRequest)
		render.JSON(response, request, res)
		return
	}
	render.Status(request, http.StatusOK)
	render.JSON(response, request, res)
}

// PostEvalMany Evaluates a list of context/expression from payload 
// @Summary Post Evaluation Many
// @Description Evaluates a list of given context from payload against the CEL expression
// @Tags evaluation
// @Accept  json
// @Produce  json
// @Security apikey
// @Param payload body []model.CelModel true "Context and expression"
// @Success 201 {object} []model.CelResult "Evaluation result"
// @Failure 400 {object} serror.Serr "client error information as json"
// @Failure 500 {object} serror.Serr "server error information as json"
// @Router /evaluatemany [post]
func PostEvalMany(response http.ResponseWriter, request *http.Request) {
	postEvalManyCounter.Inc()
	var celModels []model.CelModel
	err := decode(request, &celModels)
	if err != nil {
		log.Logger.Errorf("error decoding context: %v", err)
		msg := fmt.Sprintf("error decoding context: %v", err)
		httputils.Err(response, request, serror.BadRequest(nil, "server-error", msg))
		return
	}
	res, err := celproc.ProcCelMany(celModels)
	log.Logger.Infof("req: %v, res: %v", celModels, res)
	if err != nil {
		log.Logger.Errorf("processing error: %v", err)
		render.Status(request, http.StatusBadRequest)
		render.JSON(response, request, res)
		return
	}

	render.Status(request, http.StatusOK)
	render.JSON(response, request, res)
}

// Validate validator
var Validate *validator.Validate = validator.New()

// Decode decodes and validates an object
func decode(r *http.Request, v interface{}) error {
	err := defaultDecoder(r, v)
	if err != nil {
		return serror.BadRequest(err, "decode-body", "could not decode body")
	}
	if err := Validate.Struct(v); err != nil {
		return serror.BadRequest(err, "validate-body", "body invalid")
	}
	return nil
}

func defaultDecoder(r *http.Request, v interface{}) error {
	var err error

	switch render.GetRequestContentType(r) {
	case render.ContentTypeJSON:
		err = decodeJSON(r.Body, v)
	case render.ContentTypeXML:
		err = render.DecodeXML(r.Body, v)
	// case ContentTypeForm: // TODO
	default:
		err = errors.New("render: unable to automatically decode the request content type")
	}

	return err
}

func decodeJSON(r io.Reader, v interface{}) error {
	defer io.Copy(ioutil.Discard, r)
	d := json.NewDecoder(r)
	d.UseNumber()
	return d.Decode(v)
}
