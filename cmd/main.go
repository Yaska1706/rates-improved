package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Yaska1706/rates-improved/cmd/config"
	"github.com/Yaska1706/rates-improved/pkg/database"
	"github.com/Yaska1706/rates-improved/pkg/handlers"
	"github.com/Yaska1706/rates-improved/pkg/service"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func main() {

	config.LoadConfig()
	logger := config.SetupLogger()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := config.SetupDatabase()
	if err != nil {
		logger.ErrorContext(ctx, "db setup failed", "error", err.Error())
		return
	}

	rateRepo := database.NewRateRepo(ctx, db, logger)
	rateService := service.NewRateService(ctx, logger, rateRepo)
	rateHandler := handlers.NewRateHandler(rateService)
	router := mux.NewRouter()

	go runRatesJobLoop(ctx, logger, rateRepo, rateService)

	router.HandleFunc("/rates/latest", rateHandler.GetLatestRates).Methods("GET")
	router.HandleFunc("/rates/analyze", rateHandler.GetCurrencyStats).Methods("GET")
	router.HandleFunc("/rates/{date}", rateHandler.GetCurrencyByDate).Methods("GET")

	http.ListenAndServe(":8080", router)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.InfoContext(ctx, "Shutting down application")

	cancel()

	// Allow some time for cleanup
	time.Sleep(1 * time.Second)
	logger.InfoContext(ctx, "Application shutdown complete")

}

func runRatesJobLoop(ctx context.Context, logger *slog.Logger, rateRepo *database.RateRepo, rateService *service.RateService) {
	if err := fetchAndInsertRates(ctx, logger, rateRepo, rateService); err != nil {
		logger.ErrorContext(ctx, "Initial rates job failed", "error", err)
	}

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.InfoContext(ctx, "Stopping rates job loop")
			return
		case <-ticker.C:
			if err := fetchAndInsertRates(ctx, logger, rateRepo, rateService); err != nil {
				logger.ErrorContext(ctx, "Rates job failed", "error", err)
				continue
			}
		}
	}
}

func fetchAndInsertRates(ctx context.Context, logger *slog.Logger, rateRepo *database.RateRepo, rateService *service.RateService) error {
	logger.InfoContext(ctx, "Starting rates job")

	currencyRates, err := rateService.FetchRates(viper.GetString(config.XmlUrl))
	if err != nil {
		logger.ErrorContext(ctx, "Failed to fetch rates", "error", err)
		return err
	}

	if err := rateRepo.Insert(currencyRates); err != nil {
		logger.ErrorContext(ctx, "Failed to insert rates", "error", err)
		return err
	}

	logger.InfoContext(ctx, "Rates job completed successfully")
	return nil
}
