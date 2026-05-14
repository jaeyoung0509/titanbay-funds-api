package integration

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	postgresadapter "github.com/jaeyoung0509/titanbay-funds-api/internal/adapter/postgres"
	appserver "github.com/jaeyoung0509/titanbay-funds-api/internal/app"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/domain/enum"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/domain/vo"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/platform"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/request"
	apiresponse "github.com/jaeyoung0509/titanbay-funds-api/internal/response"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	testDB        *sql.DB
	testApp       *fiber.App
	testContainer *postgres.PostgresContainer
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("titanbay"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "docker") {
			fmt.Fprintln(os.Stderr, "integration tests skipped:", err)
			os.Exit(0)
		}
		panic(err)
	}
	testContainer = container

	connString, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}

	db, err := platform.Open(connString)
	if err != nil {
		panic(err)
	}
	testDB = db

	migrationsPath := filepath.Join("..", "..", "migrations")
	if err := platform.Migrate(testDB, migrationsPath); err != nil {
		panic(err)
	}

	swaggerPath := filepath.Join("..", "..", "docs", "swagger", "openapi.yaml")
	testApp = appserver.New(appserver.Dependencies{
		Repo:            postgresadapter.New(testDB),
		SwaggerFilePath: swaggerPath,
	})

	code := m.Run()

	_ = testDB.Close()
	_ = testContainer.Terminate(ctx)
	os.Exit(code)
}

func resetDB(t *testing.T) {
	t.Helper()
	_, err := testDB.Exec(`TRUNCATE TABLE investments, investors, funds RESTART IDENTITY CASCADE`)
	if err != nil {
		t.Fatalf("truncate: %v", err)
	}
}

func requestJSON(t *testing.T, method, path string, body any) (*http.Response, []byte) {
	t.Helper()

	var payload []byte
	var err error
	if body != nil {
		payload, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal request: %v", err)
		}
	}

	var reader io.Reader
	if len(payload) > 0 {
		reader = bytes.NewReader(payload)
	}

	req := httptest.NewRequest(method, path, reader)
	if len(payload) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := testApp.Test(req, fiber.TestConfig{Timeout: 0})
	if err != nil {
		t.Fatalf("request %s %s: %v", method, path, err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response: %v", err)
	}
	_ = resp.Body.Close()

	return resp, respBody
}

func createFund(t *testing.T, name string) apiresponse.Fund {
	t.Helper()
	resp, body := requestJSON(t, http.MethodPost, "/funds", request.CreateFundRequest{
		Name:          name,
		VintageYear:   2024,
		TargetSizeUSD: mustMoney(t, "250000000.00"),
		Status:        enum.FundStatusFundraising.String(),
	})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", resp.StatusCode, string(body))
	}

	var fund apiresponse.Fund
	if err := json.Unmarshal(body, &fund); err != nil {
		t.Fatalf("decode fund: %v", err)
	}
	return fund
}

func createInvestor(t *testing.T, name string, email string) apiresponse.Investor {
	t.Helper()
	resp, body := requestJSON(t, http.MethodPost, "/investors", request.CreateInvestorRequest{
		Name:         name,
		InvestorType: enum.InvestorTypeInstitution.String(),
		Email:        email,
	})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", resp.StatusCode, string(body))
	}

	var investor apiresponse.Investor
	if err := json.Unmarshal(body, &investor); err != nil {
		t.Fatalf("decode investor: %v", err)
	}
	return investor
}

func createInvestment(t *testing.T, fundID, investorID string, amount string, date string) apiresponse.Investment {
	t.Helper()

	resp, body := requestJSON(t, http.MethodPost, "/funds/"+fundID+"/investments", request.CreateInvestmentRequest{
		InvestorID:     investorID,
		AmountUSD:      mustMoney(t, amount),
		InvestmentDate: mustDate(t, date),
	})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", resp.StatusCode, string(body))
	}

	var investment apiresponse.Investment
	if err := json.Unmarshal(body, &investment); err != nil {
		t.Fatalf("decode investment: %v", err)
	}
	return investment
}

func TestSwaggerRoute(t *testing.T) {
	resetDB(t)

	resp, _ := requestJSON(t, http.MethodGet, "/swagger", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}
}

func TestHealthRoute(t *testing.T) {
	resetDB(t)

	resp, body := requestJSON(t, http.MethodGet, "/", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, body = %s", resp.StatusCode, string(body))
	}

	var health apiresponse.Health
	if err := json.Unmarshal(body, &health); err != nil {
		t.Fatalf("decode health: %v", err)
	}
	if health.Status != "ok" {
		t.Fatalf("status = %q", health.Status)
	}
}

func TestCreateFundAndValidation(t *testing.T) {
	resetDB(t)

	resp, body := requestJSON(t, http.MethodPost, "/funds", request.CreateFundRequest{
		Name:          "Titanbay Growth Fund II",
		VintageYear:   2025,
		TargetSizeUSD: mustMoney(t, "500000000.00"),
		Status:        enum.FundStatusFundraising.String(),
	})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", resp.StatusCode, string(body))
	}

	var created apiresponse.Fund
	if err := json.Unmarshal(body, &created); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if created.Name != "Titanbay Growth Fund II" {
		t.Fatalf("name = %q", created.Name)
	}
	if created.TargetSizeUSD.String() != "500000000.00" {
		t.Fatalf("target_size_usd = %s", created.TargetSizeUSD.String())
	}
	if created.CreatedAt.String() == "" {
		t.Fatal("expected created_at")
	}

	resp, body = requestJSON(t, http.MethodPost, "/funds", request.CreateFundRequest{
		Name:          "Broken Fund",
		VintageYear:   2025,
		TargetSizeUSD: mustMoney(t, "500000000.00"),
		Status:        "Draft",
	})
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", resp.StatusCode, string(body))
	}

	var errEnvelope apiresponse.ErrorEnvelope
	if err := json.Unmarshal(body, &errEnvelope); err != nil {
		t.Fatalf("decode error envelope: %v", err)
	}
	if errEnvelope.Error.Code != "validation_error" {
		t.Fatalf("code = %s", errEnvelope.Error.Code)
	}
	if errEnvelope.Error.Fields["status"] == "" {
		t.Fatal("expected status validation field")
	}
}

func TestGetFundNotFound(t *testing.T) {
	resetDB(t)

	resp, _ := requestJSON(t, http.MethodGet, "/funds/"+uuid.NewString(), nil)
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d", resp.StatusCode)
	}
}

func TestCreateInvestorAndDuplicateEmail(t *testing.T) {
	resetDB(t)

	resp, body := requestJSON(t, http.MethodPost, "/investors", request.CreateInvestorRequest{
		Name:         "CalPERS",
		InvestorType: enum.InvestorTypeInstitution.String(),
		Email:        "privateequity@calpers.ca.gov",
	})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", resp.StatusCode, string(body))
	}

	var investor apiresponse.Investor
	if err := json.Unmarshal(body, &investor); err != nil {
		t.Fatalf("decode investor: %v", err)
	}
	if investor.Email.String() != "privateequity@calpers.ca.gov" {
		t.Fatalf("email = %q", investor.Email.String())
	}

	resp, body = requestJSON(t, http.MethodPost, "/investors", request.CreateInvestorRequest{
		Name:         "CalPERS Two",
		InvestorType: enum.InvestorTypeInstitution.String(),
		Email:        "privateequity@calpers.ca.gov",
	})
	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("status = %d, body = %s", resp.StatusCode, string(body))
	}

	var errEnvelope apiresponse.ErrorEnvelope
	if err := json.Unmarshal(body, &errEnvelope); err != nil {
		t.Fatalf("decode error envelope: %v", err)
	}
	if errEnvelope.Error.Code != "conflict" {
		t.Fatalf("code = %s", errEnvelope.Error.Code)
	}
}

func TestInvestmentsFlow(t *testing.T) {
	resetDB(t)

	fund := createFund(t, "Titanbay Growth Fund I")
	investor := createInvestor(t, "Goldman Sachs Asset Management", "investments@example.com")

	investment := createInvestment(t, fund.ID.String(), investor.ID.String(), "75000000.00", "2024-09-22")
	if investment.FundID != fund.ID {
		t.Fatalf("fund_id = %s, want %s", investment.FundID, fund.ID)
	}
	if investment.InvestorID != investor.ID {
		t.Fatalf("investor_id = %s, want %s", investment.InvestorID, investor.ID)
	}
	if investment.AmountUSD.String() != "75000000.00" {
		t.Fatalf("amount = %s", investment.AmountUSD.String())
	}
	if investment.InvestmentDate.String() != "2024-09-22" {
		t.Fatalf("investment_date = %s", investment.InvestmentDate.String())
	}

	resp, body := requestJSON(t, http.MethodGet, "/funds/"+fund.ID.String()+"/investments", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, body = %s", resp.StatusCode, string(body))
	}

	var investments []apiresponse.Investment
	if err := json.Unmarshal(body, &investments); err != nil {
		t.Fatalf("decode investments: %v", err)
	}
	if len(investments) != 1 {
		t.Fatalf("len = %d", len(investments))
	}
	if investments[0].ID != investment.ID {
		t.Fatalf("id = %s, want %s", investments[0].ID, investment.ID)
	}

	resp, body = requestJSON(t, http.MethodPost, "/funds/"+fund.ID.String()+"/investments", request.CreateInvestmentRequest{
		InvestorID:     uuid.NewString(),
		AmountUSD:      mustMoney(t, "1000000.00"),
		InvestmentDate: mustDate(t, "2024-09-23"),
	})
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d, body = %s", resp.StatusCode, string(body))
	}

	resp, body = requestJSON(t, http.MethodGet, "/funds/"+uuid.NewString()+"/investments", nil)
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d, body = %s", resp.StatusCode, string(body))
	}
}

func mustDate(t *testing.T, raw string) string {
	t.Helper()
	date, err := vo.ParseDate(raw)
	if err != nil {
		t.Fatalf("parse date: %v", err)
	}
	return date.String()
}

func mustMoney(t *testing.T, raw string) vo.Money {
	t.Helper()
	value, err := decimal.NewFromString(raw)
	if err != nil {
		t.Fatalf("decimal parse: %v", err)
	}
	return vo.NewMoney(value)
}

func TestListEndpointsReturnRawArrays(t *testing.T) {
	resetDB(t)

	fund := createFund(t, "Titanbay Growth Fund III")
	investor := createInvestor(t, "Institution One", "one@example.com")
	createInvestment(t, fund.ID.String(), investor.ID.String(), "1000000.00", "2024-01-01")

	resp, body := requestJSON(t, http.MethodGet, "/funds", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, body = %s", resp.StatusCode, string(body))
	}
	var funds []apiresponse.Fund
	if err := json.Unmarshal(body, &funds); err != nil {
		t.Fatalf("decode funds: %v", err)
	}
	if len(funds) != 1 {
		t.Fatalf("funds len = %d", len(funds))
	}

	resp, body = requestJSON(t, http.MethodGet, "/investors", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, body = %s", resp.StatusCode, string(body))
	}
	var investors []apiresponse.Investor
	if err := json.Unmarshal(body, &investors); err != nil {
		t.Fatalf("decode investors: %v", err)
	}
	if len(investors) != 1 {
		t.Fatalf("investors len = %d", len(investors))
	}

	resp, body = requestJSON(t, http.MethodGet, "/funds/"+fund.ID.String()+"/investments", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, body = %s", resp.StatusCode, string(body))
	}
	var investments []apiresponse.Investment
	if err := json.Unmarshal(body, &investments); err != nil {
		t.Fatalf("decode investments: %v", err)
	}
	if len(investments) != 1 {
		t.Fatalf("investments len = %d", len(investments))
	}
}
