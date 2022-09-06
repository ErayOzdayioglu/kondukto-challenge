package service

import (
	"log"

	"kondukto.com/challenge/domain"
	"kondukto.com/challenge/repository"
)

type DefaultScanResultsService struct {
	Repository repository.ScanResultsRepository
}

type ScanService interface {
	Insert(domain.ScanResults) (string, error)
	Find(string) (domain.ScanResults, int)
}

func (s DefaultScanResultsService) Insert(request domain.ScanResults) (string, error) {

	result, err := s.Repository.Insert(request)

	if err != nil {
		log.Fatal("error in InsertService")
	}

	return result, nil
}

func (s DefaultScanResultsService) Find(id string) (domain.ScanResults, int) {
	result, err := s.Repository.Find(id)

	if err == -1 {
		return result, -1
	}

	return result, 0

}

func NewScanResultService(Repository repository.ScanResultsRepository) DefaultScanResultsService {
	return DefaultScanResultsService{Repository: Repository}
}
