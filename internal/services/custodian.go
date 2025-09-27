package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/quantfidential/trading-ecosystem/custodian-simulator-go/internal/config"
)

type CustodianService struct {
	config    *config.Config
	logger    *logrus.Logger
	startTime time.Time

	// Service state
	mu       sync.RWMutex
	accounts map[string]*Account
	balances map[string]map[string]float64 // accountID -> assetID -> balance
}

type Account struct {
	ID          string            `json:"id"`
	Type        string            `json:"type"`
	Status      string            `json:"status"`
	Balances    map[string]float64 `json:"balances"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type Settlement struct {
	ID            string    `json:"id"`
	FromAccount   string    `json:"from_account"`
	ToAccount     string    `json:"to_account"`
	AssetID       string    `json:"asset_id"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	SettlementDate time.Time `json:"settlement_date"`
	CreatedAt     time.Time `json:"created_at"`
}

func NewCustodianService(cfg *config.Config, logger *logrus.Logger) *CustodianService {
	return &CustodianService{
		config:    cfg,
		logger:    logger,
		startTime: time.Now(),
		accounts:  make(map[string]*Account),
		balances:  make(map[string]map[string]float64),
	}
}

func (s *CustodianService) GetHealth(ctx context.Context) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Basic health check - verify service is operational
	if time.Since(s.startTime) < time.Second {
		return "starting", nil
	}

	return "healthy", nil
}

func (s *CustodianService) GetBalance(account string, asset string) (float64, error) {
	s.logger.WithFields(logrus.Fields{
		"account": account,
		"asset":   asset,
	}).Info("Getting balance")

	s.mu.RLock()
	defer s.mu.RUnlock()

	balances := s.balances[account]
	if balances == nil {
		return 0, nil
	}

	return balances[asset], nil
}

func (s *CustodianService) Transfer(fromAccount, toAccount, asset string, amount float64) (string, error) {
	s.logger.WithFields(logrus.Fields{
		"fromAccount": fromAccount,
		"toAccount":   toAccount,
		"asset":       asset,
		"amount":      amount,
	}).Info("Processing transfer")

	settlement := &Settlement{
		ID:            generateSettlementID(),
		FromAccount:   fromAccount,
		ToAccount:     toAccount,
		AssetID:       asset,
		Amount:        amount,
		Status:        "completed",
		SettlementDate: time.Now(),
		CreatedAt:     time.Now(),
	}

	err := s.ProcessSettlement(context.Background(), settlement)
	if err != nil {
		return "", err
	}

	return settlement.ID, nil
}

func (s *CustodianService) CreateAccount(ctx context.Context, accountType string) (*Account, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	account := &Account{
		ID:        generateAccountID(),
		Type:      accountType,
		Status:    "active",
		Balances:  make(map[string]float64),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.accounts[account.ID] = account
	s.balances[account.ID] = make(map[string]float64)

	s.logger.WithFields(logrus.Fields{
		"account_id": account.ID,
		"type":       accountType,
	}).Info("Account created")

	return account, nil
}

func (s *CustodianService) ProcessSettlement(ctx context.Context, settlement *Settlement) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Verify accounts exist
	fromAccount, exists := s.accounts[settlement.FromAccount]
	if !exists {
		return fmt.Errorf("from account %s not found", settlement.FromAccount)
	}

	toAccount, exists := s.accounts[settlement.ToAccount]
	if !exists {
		return fmt.Errorf("to account %s not found", settlement.ToAccount)
	}

	// Check balance
	fromBalances := s.balances[settlement.FromAccount]
	if fromBalances[settlement.AssetID] < settlement.Amount {
		return fmt.Errorf("insufficient balance in account %s for asset %s",
			settlement.FromAccount, settlement.AssetID)
	}

	// Process settlement
	fromBalances[settlement.AssetID] -= settlement.Amount

	toBalances := s.balances[settlement.ToAccount]
	if toBalances == nil {
		toBalances = make(map[string]float64)
		s.balances[settlement.ToAccount] = toBalances
	}
	toBalances[settlement.AssetID] += settlement.Amount

	// Update account timestamps
	now := time.Now()
	fromAccount.UpdatedAt = now
	toAccount.UpdatedAt = now

	s.logger.WithFields(logrus.Fields{
		"settlement_id": settlement.ID,
		"from_account":  settlement.FromAccount,
		"to_account":    settlement.ToAccount,
		"asset_id":      settlement.AssetID,
		"amount":        settlement.Amount,
	}).Info("Settlement processed")

	return nil
}

func (s *CustodianService) GetAccountBalance(ctx context.Context, accountID, assetID string) (float64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, exists := s.accounts[accountID]; !exists {
		return 0, fmt.Errorf("account %s not found", accountID)
	}

	balances := s.balances[accountID]
	if balances == nil {
		return 0, nil
	}

	return balances[assetID], nil
}

func generateAccountID() string {
	// Simple ID generation for simulation
	return fmt.Sprintf("ACCT_%d", time.Now().UnixNano())
}

func generateSettlementID() string {
	// Simple ID generation for simulation
	return fmt.Sprintf("SETTLE_%d", time.Now().UnixNano())
}