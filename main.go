package main

import (
	"github.com/gofiber/fiber/v2"
	"kondukto.com/challenge/app"
	"kondukto.com/challenge/config"
	"kondukto.com/challenge/repository"
	"kondukto.com/challenge/service"
)

func main() {

	appRoute := fiber.New()
	config.ConnectDB()
	dbClient := config.GetCollection(config.DB, "scanResults")

	ScanResultsRepositoryDB := repository.NewScanResultsRepositoryDB(dbClient)

	sh := app.ScanResultHandler{ScanService: service.NewScanResultService(ScanResultsRepositoryDB),
		AnalyzeService: service.NewAnalyzeService(ScanResultsRepositoryDB)}

	appRoute.Post("/api/v1/newscan", sh.CreateScanResult)
	appRoute.Get("api/v1/scan/:scan_id", sh.GetScanResult)
	appRoute.Listen(":8080")
}
