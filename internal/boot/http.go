package boot

import (
	"log"
	"net/http"
	"testing/internal/data/auth"
	"testing/pkg/httpclient"
	"testing/pkg/tracing"

	"testing/internal/config"
	jaegerLog "testing/pkg/log"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	testingData "testing/internal/data/testing"
	testingServer "testing/internal/delivery/http"
	testingHandler "testing/internal/delivery/http/testing"
	testingService "testing/internal/service/testing"
)

// HTTP will load configuration, do dependency injection and then start the HTTP server
func HTTP() error {
	var (
		s        testingServer.Server    // HTTP Server Object
		testingD testingData.Data        // Domain data layer
		testingS testingService.Service  // Domain service layer
		testingH *testingHandler.Handler // Domain handler
		cfg      *config.Config          // Configuration object
		httpc    *httpclient.Client
		authD    auth.Data

		logger *zap.Logger
	)

	err := config.Init()
	if err != nil {
		log.Fatalf("[CONFIG] Failed to initialize config: %v", err)
	}
	cfg = config.Get()
	// Open MySQL DB Connection
	db, err := sqlx.Open("mysql", cfg.Database.Master)
	if err != nil {
		log.Fatalf("[DB] Failed to initialize database connection: %v", err)
	}

	// Set logger used for jaeger
	logger, _ = zap.NewDevelopment(
		zap.AddStacktrace(zapcore.FatalLevel),
		zap.AddCallerSkip(1),
	)
	zapLogger := logger.With(zap.String("service", "testing"))
	zlogger := jaegerLog.NewFactory(zapLogger)

	// Set tracer for service
	tracer, closer := tracing.Init("testing", zlogger)
	defer closer.Close()

	httpc = httpclient.NewClient(tracer)
	authD = auth.New(httpc, cfg.API.Auth)

	// Diganti dengan domain yang anda buat
	testingD = testingData.New(db, tracer, zlogger)
	testingS = testingService.New(testingD, authD, tracer, zlogger)
	testingH = testingHandler.New(testingS, tracer, zlogger)

	s = testingServer.Server{
		Testing: testingH,
	}

	if err := s.Serve(cfg.Server.Port); err != http.ErrServerClosed {
		return err
	}

	return nil
}
