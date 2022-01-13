package health

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/opentracing/opentracing-go"
	log "github.com/willie68/cel-service/internal/logging"
)

var myhealthy bool

/*
This is the healtchcheck you will have to provide.
*/
func check(tracer opentracing.Tracer) (bool, string) {
	myhealthy = true
	message := "healthy"
	return myhealthy, message
}

//##### template internal functions for processing the healthchecks #####
var healthmessage string
var healthy bool
var lastChecked time.Time
var period int

// CheckConfig configuration for the healthcheck system
type CheckConfig struct {
	Period int
}

// Msg a health message
type Msg struct {
	Message   string `json:"message"`
	LastCheck string `json:"lastCheck,omitempty"`
}

// InitHealthSystem initialise the complete health system
func InitHealthSystem(config CheckConfig, tracer opentracing.Tracer) {
	period = config.Period
	log.Logger.Infof("healthcheck starting with period: %d seconds", period)
	healthmessage = "service starting"
	healthy = false
	doCheck(tracer)
	go func() {
		background := time.NewTicker(time.Second * time.Duration(period))
		for _ = range background.C {
			doCheck(tracer)
		}
	}()
}

// docheck internal function to process the health check
func doCheck(tracer opentracing.Tracer) {
	var msg string
	healthy, msg = check(tracer)
	if !healthy {
		healthmessage = msg
	} else {
		healthmessage = ""
	}
	lastChecked = time.Now()
}

// Routes getting all routes for the health endpoint

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/livez", GetLiveness)
	router.Get("/readyz", GetReadiness)
	router.Head("/livez", GetLiveness)
	router.Head("/readyz", GetReadiness)
	return router
}

// GetLiveness liveness probe
func GetLiveness(response http.ResponseWriter, req *http.Request) {
	render.Status(req, http.StatusOK)
	render.JSON(response, req, Msg{
		Message: "service started",
	})
}

// GetReadiness is this service ready for taking requests, e.g. formaly known as health checks
func GetReadiness(response http.ResponseWriter, req *http.Request) {
	t := time.Now()
	if t.Sub(lastChecked) > (time.Second * time.Duration(2*period)) {
		healthy = false
		healthmessage = "Healthcheck not running"
	}
	if healthy {
		render.Status(req, http.StatusOK)
		render.JSON(response, req, Msg{
			Message:   "service up and running",
			LastCheck: lastChecked.String(),
		})
	} else {
		render.Status(req, http.StatusServiceUnavailable)
		render.JSON(response, req, Msg{
			Message:   fmt.Sprintf("service is unavailable: %s", healthmessage),
			LastCheck: lastChecked.String(),
		})
	}
}

// sendMessage sending a span message to tracer
func sendMessage(tracer opentracing.Tracer, message string) {
	span := tracer.StartSpan("say-hello")
	println(message)
	span.Finish()
}
