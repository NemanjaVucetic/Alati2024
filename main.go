package main

import (
	_ "alati/docs"
	"alati/handler"
	"alati/model"
	"alati/repo"
	"alati/service"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"golang.org/x/time/rate"
)

// @title			Configuration API
// @version		1.0
// @description	This is a sample server for a configuration service.
// @termsOfService	http://swagger.io/terms/

// @contact.name	API Support
// @contact.url	http://www.swagger.io/support
// @contact.email	support@swagger.io

// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:8080
// @BasePath	/

func main() {

	cfg := GetConfiguration()

	// Initialize OpenTelemetry
	ctx := context.Background()
	exp, err := newExporter(cfg.JaegerEndpoint)
	if err != nil {
		log.Fatalf("failed to initialize exporter: %v", err)
	}

	tp := newTraceProvider(exp)
	defer func() { _ = tp.Shutdown(ctx) }()
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	tracer := tp.Tracer("config-service")

	logger := log.New(os.Stdout, "[config-api] ", log.LstdFlags)

	repoC, err := repo.NewConfigRepo(logger, tracer)
	if err != nil {
		logger.Fatal(err)
	}

	serviceC := service.NewConfigService(repoC, logger, tracer)

	repoG, err := repo.NewConfigGroupRepo(logger, tracer)
	if err != nil {
		logger.Fatal(err)
	}

	serviceG := service.NewConfigGroupService(repoG, logger, tracer)

	h := handler.NewConfigHandler(serviceC, logger, tracer)
	hG := handler.NewConfigGroupHandler(serviceG, serviceC, logger, tracer)

	params := make(map[string]string)
	params["param1"] = "param1"
	params["param2"] = "param2"

	labels := make(map[string]string)
	labels["l1"] = "v1"
	labels["l2"] = "v2"

	labels2 := make(map[string]string)
	labels2["l1"] = "v1"

	config := model.Config{
		Name:    "db_config",
		Version: 2,
		Params:  params,
		Labels:  labels,
	}

	config2 := model.Config{
		Name:    "db_config2",
		Version: 3,
		Params:  params,
		Labels:  labels2,
	}

	configMap := make(map[string]model.Config)
	configMap["configGroups/db_cg/2/config/l1:v1/l2:v2/db_config/2"] = config
	configMap["configGroups/db_cg/2/config/l1:v1/db_config2/3"] = config2

	group := model.ConfigGroup{
		Name:    "db_cg",
		Version: 2,
		Configs: configMap,
	}


	serviceC.Add(&config2, "npclord1", ctx)
	serviceC.Add(&config, "npclord2", ctx)
	serviceG.Add(&group, "npcgod1", ctx)


	// }

	router := mux.NewRouter()
	limiter := rate.NewLimiter(0.167, 10)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	router.Handle("/configs/{name}/{version}", count(handler.RateLimit(limiter, http.HandlerFunc(h.Get)))).Methods(http.MethodGet)
	router.Handle("/configs/", count(handler.RateLimit(limiter, http.HandlerFunc(h.Add)))).Methods(http.MethodPost)
	router.Handle("/configs/", count(handler.RateLimit(limiter, http.HandlerFunc(h.GetAll)))).Methods(http.MethodGet)
	router.Handle("/configs/{name}/{version}", count(handler.RateLimit(limiter, http.HandlerFunc(h.Delete)))).Methods(http.MethodDelete)

	router.Handle("/configGroups/", count(handler.RateLimit(limiter, http.HandlerFunc(hG.GetAll)))).Methods(http.MethodGet)
	router.Handle("/configGroups/{name}/{version}", count(handler.RateLimit(limiter, http.HandlerFunc(hG.Get)))).Methods(http.MethodGet)
	router.Handle("/configGroups/", count(handler.RateLimit(limiter, http.HandlerFunc(hG.Add)))).Methods(http.MethodPost)
	router.Handle("/configGroups/{name}/{version}", count(handler.RateLimit(limiter, http.HandlerFunc(hG.Delete)))).Methods(http.MethodDelete)
	router.Handle("/configGroups/{nameG}/{versionG}/configs/{nameC}/{versionC}", count(handler.RateLimit(limiter, http.HandlerFunc(hG.AddConfToGroup)))).Methods(http.MethodPut)
	router.Handle("/configGroups/{nameG}/{versionG}/{nameC}/{versionC}", count(handler.RateLimit(limiter, http.HandlerFunc(hG.RemoveConfFromGroup)))).Methods(http.MethodPut)

	router.Handle("/configGroups/{nameG}/{versionG}/config/{labels}", count(handler.RateLimit(limiter, http.HandlerFunc(hG.GetConfigsByLabels)))).Methods(http.MethodGet)
	router.Handle("/configGroups/{nameG}/{versionG}/config/{labels}/{nameC}", count(handler.RateLimit(limiter, http.HandlerFunc(hG.GetConfigsByLabels)))).Methods(http.MethodGet)
	router.Handle("/configGroups/{nameG}/{versionG}/config/{labels}/{nameC}/{versionC}", count(handler.RateLimit(limiter, http.HandlerFunc(hG.GetConfigsByLabels)))).Methods(http.MethodGet)

	router.Handle("/configGroups/{nameG}/{versionG}/config/{labels}", count(handler.RateLimit(limiter, http.HandlerFunc(hG.DeleteConfigsByLabels)))).Methods(http.MethodPatch)
	router.Handle("/configGroups/{nameG}/{versionG}/config/{labels}/{nameC}", count(handler.RateLimit(limiter, http.HandlerFunc(hG.DeleteConfigsByLabels)))).Methods(http.MethodPatch)
	router.Handle("/configGroups/{nameG}/{versionG}/config/{labels}/{nameC}/{versionC}", count(handler.RateLimit(limiter, http.HandlerFunc(hG.DeleteConfigsByLabels)))).Methods(http.MethodPatch)

	// Metrics endpoint
	router.Handle("/metrics", metricsHandler()).Methods(http.MethodGet)

	// Swagger documentation route
	// http://localhost:8080/swagger/index.html#/
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	server := http.Server{
		Addr:         ":" + port,        // Addr optionally specifies the TCP address for the server to listen on, in the form "host:port". If empty, ":http" (port 80) is used.
		Handler:      router,            // handler to invoke, http.DefaultServeMux if nil
		IdleTimeout:  120 * time.Second, // IdleTimeout is the maximum amount of time to wait for the next request when keep-alives are enabled.
		ReadTimeout:  1 * time.Second,   // ReadTimeout is the maximum duration for reading the entire request, including the body. A zero or negative value means there will be no timeout.
		WriteTimeout: 1 * time.Second,   // WriteTimeout is the maximum duration before timing out writes of the response.
	}

	go func() {
		log.Println("server_starting")
		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-quit

	log.Println("service_shutting_down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("server stopped")
}

func newExporter(address string) (sdktrace.SpanExporter, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(address)))
	if err != nil {
		return nil, err
	}
	return exp, nil
}

func newTraceProvider(exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	r := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("config-service"),
	)

	return sdktrace.NewTracerProvider(
		sdktrace.WithSyncer(exp),
		sdktrace.WithResource(r),
	)
}
