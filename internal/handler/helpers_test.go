package handler_test

import (
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/GulzhanKarakul/payment-service/internal/domain"
	"github.com/GulzhanKarakul/payment-service/internal/handler"
	"github.com/GulzhanKarakul/payment-service/internal/service"
)

// handler test types
type errorResponse struct {
	Error string `json:"error"`
}

type messageResponse struct {
	Message string `json:"message"`
}

// handler tests constants
const (
	testClientID = "11111111-1111-1111-1111-111111111111"
	testBusinessID = "22222222-2222-2222-2222-222222222222"
	testTransactionID = "33333333-3333-3333-3333-333333333333"
	testPhone = "+77771156580"
	testName = "Gulzhan Karakul"
	testAmount int64 = 500_000
	testBalance int64 = 10_000_00
	testDescription = "оплата кофе"
)

// logger
func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// handler factory
func newHandler(
	clientSvc service.ClientService,
	bizSvc service.BusinessService,
	txSvc service.TransactionService,
	bsSvc service.BonusSettingsService,
) http.Handler {
	h := handler.NewHandler(clientSvc, bizSvc, txSvc, bsSvc, testLogger())
	return h.Routes()
}

// domain builders
func testClient() domain.Client {
	return domain.Client{
		ID: testClientID,
		Phone: testPhone,
		Name: testName,
		BonusBalance: 0,
		IsActive: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func testBusiness() domain.Business {
	return domain.Business{
		ID: testBusinessID,
		Name: "Test Business",
		OwnerPhone: testPhone,
		BonusBalance: 10_000_00,
		IsActive: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func testBonusSettings() domain.BonusSettings {
	return domain.BonusSettings{
		ID: "44444444-4444-4444-4444-444444444444",
		BusinessID: testBusinessID,
		BonusPercent: 5.0,
		IsActive: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func testInactiveBonusSettings() domain.BonusSettings {
	return domain.BonusSettings{
		ID: "55555555-5555-5555-5555-555555555555",
		BusinessID: testBusinessID,
		BonusPercent: 5.0,
		IsActive: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func testTransaction() domain.Transaction {
	return domain.Transaction{
		ID: testTransactionID,
		ClientID: testClientID,
		BusinessID: testBusinessID,
		Amount: testAmount,
		BonusAccrued: 25_000,
		Status: domain.StatusCompleted,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Description: strPtr(testDescription),
	}
}

func strPtr(s string) *string {
	return &s
}

// makeTransactions - to create n tx slice
func makeTransactions(n int) []domain.Transaction {
	txs := make([]domain.Transaction, 0, n)
	for i:=0; i<n; i++ {
		txs = append(txs, testTransaction())
	}
	return txs
}