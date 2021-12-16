package apiv1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/willie68/cel-service/pkg/model"

	"github.com/willie68/cel-service/internal/api"
	"github.com/willie68/cel-service/internal/celproc"
	log "github.com/willie68/cel-service/internal/logging"
	"github.com/willie68/cel-service/internal/serror"
	"github.com/willie68/cel-service/internal/utils/httputils"
)

// TenantHeader in this header thr right tenant should be inserted
const timeout = 1 * time.Minute

//APIKey the apikey of this service
var APIKey string

var (
	postEvalCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "cel_service_post_config_total",
		Help: "The total number of post eval requests",
	})
)

/*
EvalRoutes getting all routes for the config endpoint
*/
func EvalRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/", PostEvalEndpoint)
	return router
}

/*
PostEvalEndpoint create a new store for a tenant
because of the automatic store creation, this method will always return 201
*/
func PostEvalEndpoint(response http.ResponseWriter, request *http.Request) {
	//	tenant := getTenant(request)
	//	if tenant == "" {
	//		msg := fmt.Sprintf("tenant header %s missing", api.TenantHeaderKey)
	//		httputils.Err(response, request, serror.BadRequest(nil, "missing-tenant", msg))
	//		return
	//	}
	postEvalCounter.Inc()
	var celModel model.CelModel
	err := httputils.Decode(request, &celModel)
	if err != nil {
		log.Logger.Errorf("error decoding context: %v", err)
		msg := fmt.Sprintf("error decoding context: %v", err)
		httputils.Err(response, request, serror.BadRequest(nil, "server-error", msg))
		return
	}

	res, err := celproc.ProcCel(celModel)
	log.Logger.Infof("req: %v, res: %v", celModel, res)

	if err != nil {
		log.Logger.Errorf("failed to listen: %v", err)
		msg := "failed to listen"
		httputils.Err(response, request, serror.BadRequest(nil, "server-error", msg))
		return
	}
	render.Status(request, http.StatusCreated)
	render.JSON(response, request, res)
}

/*
getTenant getting the tenant from the request
*/
func getTenant(req *http.Request) string {
	return req.Header.Get(api.TenantHeaderKey)
}
