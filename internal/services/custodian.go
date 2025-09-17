package services

import (
	"github.com/sirupsen/logrus"
)

type CustodianService struct {
	logger *logrus.Logger
}

func NewCustodianService(logger *logrus.Logger) *CustodianService {
	return &CustodianService{
		logger: logger,
	}
}

func (s *CustodianService) GetBalance(account string, asset string) (float64, error) {
	s.logger.WithFields(logrus.Fields{
		"account": account,
		"asset":   asset,
	}).Info("Getting balance")
	return 1000.0, nil
}

func (s *CustodianService) Transfer(fromAccount, toAccount, asset string, amount float64) (string, error) {
	s.logger.WithFields(logrus.Fields{
		"fromAccount": fromAccount,
		"toAccount":   toAccount,
		"asset":       asset,
		"amount":      amount,
	}).Info("Processing transfer")
	return "transfer-123", nil
}