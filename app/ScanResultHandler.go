package app

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"kondukto.com/challenge/domain"
	"kondukto.com/challenge/service"
)

type ScanResultHandler struct {
	ScanService    service.ScanService
	AnalyzeService service.AnalyzeService
}

func (h ScanResultHandler) CreateScanResult(c *fiber.Ctx) error {
	var request domain.NewScanRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	var scanResults domain.ScanResults
	scanResults.GithubUrl = request.Url

	result, err := h.ScanService.Insert(scanResults)

	if err != nil {
		log.Fatalln("error")
	}

	go h.AnalyzeService.Analyze(result, request.Url)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"id": result,
	})
}

func (h ScanResultHandler) GetScanResult(c *fiber.Ctx) error {

	id := c.Params("scan_id")
	result, err := h.ScanService.Find(id)

	if err == -1 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "There is no scan results with id : " + id,
		})
	}

	return c.Status(http.StatusOK).JSON(result)

}
