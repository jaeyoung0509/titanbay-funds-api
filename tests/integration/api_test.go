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

func TestUtilityRoutes(t *testing.T) {
	t.Run("returns health status", func(t *testing.T) {
		resetDB(t)

		resp, body := requestJSON(t, http.MethodGet, "/", nil)

		expectStatus(t, resp, body, http.StatusOK)
		health := decodeJSON[apiresponse.Health](t, body)
		if health.Status != "ok" {
			t.Fatalf("status = %q", health.Status)
		}
	})

	t.Run("serves swagger UI", func(t *testing.T) {
		resetDB(t)

		resp, body := requestJSON(t, http.MethodGet, "/swagger", nil)

		expectStatus(t, resp, body, http.StatusOK)
	})

	t.Run("returns not found for unknown routes", func(t *testing.T) {
		resetDB(t)

		resp, body := requestJSON(t, http.MethodGet, "/does-not-exist", nil)

		expectStatus(t, resp, body, http.StatusNotFound)
		expectError(t, body, "not_found", "")
	})
}

func TestFundsAPI(t *testing.T) {
	t.Run("creates a fund", func(t *testing.T) {
		resetDB(t)
		req := validCreateFundRequest("Titanbay Growth Fund II")

		resp, body := requestJSON(t, http.MethodPost, "/funds", req)

		expectStatus(t, resp, body, http.StatusCreated)
		fund := decodeJSON[apiresponse.Fund](t, body)
		requireFundMatchesCreateRequest(t, fund, req)
	})

	t.Run("updates a fund", func(t *testing.T) {
		resetDB(t)
		fund := createFund(t, "Titanbay Growth Fund I")
		req := validUpdateFundRequest(fund.ID.String(), "Titanbay Growth Fund I")

		resp, body := requestJSON(t, http.MethodPut, "/funds", req)

		expectStatus(t, resp, body, http.StatusOK)
		updated := decodeJSON[apiresponse.Fund](t, body)
		requireFundMatchesUpdateRequest(t, updated, req)
	})

	t.Run("gets a fund by id", func(t *testing.T) {
		resetDB(t)
		fund := createFund(t, "Titanbay Growth Fund I")

		resp, body := requestJSON(t, http.MethodGet, "/funds/"+fund.ID.String(), nil)

		expectStatus(t, resp, body, http.StatusOK)
		got := decodeJSON[apiresponse.Fund](t, body)
		if got.ID != fund.ID {
			t.Fatalf("id = %s, want %s", got.ID, fund.ID)
		}
	})

	t.Run("lists funds as a raw array", func(t *testing.T) {
		resetDB(t)
		fund := createFund(t, "Titanbay Growth Fund III")

		resp, body := requestJSON(t, http.MethodGet, "/funds", nil)

		expectStatus(t, resp, body, http.StatusOK)
		funds := decodeJSON[[]apiresponse.Fund](t, body)
		requireContainsFund(t, funds, fund.ID)
	})

	t.Run("returns not found for missing fund", func(t *testing.T) {
		resetDB(t)

		resp, body := requestJSON(t, http.MethodGet, "/funds/"+uuid.NewString(), nil)

		expectStatus(t, resp, body, http.StatusNotFound)
		expectError(t, body, "not_found", "")
	})

	t.Run("rejects invalid status on create", func(t *testing.T) {
		resetDB(t)
		req := validCreateFundRequest("Broken Fund")
		req.Status = "Draft"

		resp, body := requestJSON(t, http.MethodPost, "/funds", req)

		expectStatus(t, resp, body, http.StatusBadRequest)
		expectError(t, body, "validation_error", "status")
	})

	t.Run("rejects invalid status on update", func(t *testing.T) {
		resetDB(t)
		fund := createFund(t, "Status Validation Fund")
		req := validUpdateFundRequest(fund.ID.String(), "Status Validation Fund")
		req.Status = "Draft"

		resp, body := requestJSON(t, http.MethodPut, "/funds", req)

		expectStatus(t, resp, body, http.StatusBadRequest)
		expectError(t, body, "validation_error", "status")
	})
}

func TestInvestorsAPI(t *testing.T) {
	t.Run("creates an investor", func(t *testing.T) {
		resetDB(t)
		req := validCreateInvestorRequest("CalPERS", "privateequity@calpers.ca.gov")

		resp, body := requestJSON(t, http.MethodPost, "/investors", req)

		expectStatus(t, resp, body, http.StatusCreated)
		investor := decodeJSON[apiresponse.Investor](t, body)
		requireInvestorMatchesCreateRequest(t, investor, req)
	})

	t.Run("lists investors as a raw array", func(t *testing.T) {
		resetDB(t)
		investor := createInvestor(t, "Institution One", "one@example.com")

		resp, body := requestJSON(t, http.MethodGet, "/investors", nil)

		expectStatus(t, resp, body, http.StatusOK)
		investors := decodeJSON[[]apiresponse.Investor](t, body)
		requireContainsInvestor(t, investors, investor.ID)
	})

	t.Run("rejects duplicate email", func(t *testing.T) {
		resetDB(t)
		req := validCreateInvestorRequest("CalPERS", "privateequity@calpers.ca.gov")
		createInvestor(t, req.Name, req.Email)
		req.Name = "CalPERS Two"

		resp, body := requestJSON(t, http.MethodPost, "/investors", req)

		expectStatus(t, resp, body, http.StatusConflict)
		expectError(t, body, "conflict", "")
	})

	t.Run("rejects invalid investor type", func(t *testing.T) {
		resetDB(t)
		req := validCreateInvestorRequest("Invalid Investor", "invalid@example.com")
		req.InvestorType = "Draft"

		resp, body := requestJSON(t, http.MethodPost, "/investors", req)

		expectStatus(t, resp, body, http.StatusBadRequest)
		expectError(t, body, "validation_error", "investor_type")
	})
}

func TestInvestmentsAPI(t *testing.T) {
	t.Run("creates an investment", func(t *testing.T) {
		resetDB(t)
		fund := createFund(t, "Titanbay Growth Fund I")
		investor := createInvestor(t, "Goldman Sachs Asset Management", "investments@example.com")
		req := validCreateInvestmentRequest(investor.ID.String(), "75000000.00", "2024-09-22")

		resp, body := requestJSON(t, http.MethodPost, "/funds/"+fund.ID.String()+"/investments", req)

		expectStatus(t, resp, body, http.StatusCreated)
		investment := decodeJSON[apiresponse.Investment](t, body)
		requireInvestmentMatchesCreateRequest(t, investment, fund.ID, req)
	})

	t.Run("lists investments for a fund as a raw array", func(t *testing.T) {
		resetDB(t)
		fund := createFund(t, "Titanbay Growth Fund III")
		investor := createInvestor(t, "Institution One", "one@example.com")
		investment := createInvestment(t, fund.ID.String(), investor.ID.String(), "1000000.00", "2024-01-01")

		resp, body := requestJSON(t, http.MethodGet, "/funds/"+fund.ID.String()+"/investments", nil)

		expectStatus(t, resp, body, http.StatusOK)
		investments := decodeJSON[[]apiresponse.Investment](t, body)
		requireContainsInvestment(t, investments, investment.ID)
	})

	t.Run("returns not found when creating with a missing investor", func(t *testing.T) {
		resetDB(t)
		fund := createFund(t, "Titanbay Growth Fund I")
		req := validCreateInvestmentRequest(uuid.NewString(), "1000000.00", "2024-09-23")

		resp, body := requestJSON(t, http.MethodPost, "/funds/"+fund.ID.String()+"/investments", req)

		expectStatus(t, resp, body, http.StatusNotFound)
		expectError(t, body, "not_found", "")
	})

	t.Run("returns not found when creating for a missing fund", func(t *testing.T) {
		resetDB(t)
		investor := createInvestor(t, "Goldman Sachs Asset Management", "investments@example.com")
		req := validCreateInvestmentRequest(investor.ID.String(), "1000000.00", "2024-09-23")

		resp, body := requestJSON(t, http.MethodPost, "/funds/"+uuid.NewString()+"/investments", req)

		expectStatus(t, resp, body, http.StatusNotFound)
		expectError(t, body, "not_found", "")
	})

	t.Run("returns not found when listing investments for a missing fund", func(t *testing.T) {
		resetDB(t)

		resp, body := requestJSON(t, http.MethodGet, "/funds/"+uuid.NewString()+"/investments", nil)

		expectStatus(t, resp, body, http.StatusNotFound)
		expectError(t, body, "not_found", "")
	})
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

func expectStatus(t *testing.T, resp *http.Response, body []byte, wantStatus int) {
	t.Helper()
	if resp.StatusCode != wantStatus {
		t.Fatalf("status = %d, want %d, body = %s", resp.StatusCode, wantStatus, string(body))
	}
}

func decodeJSON[T any](t *testing.T, body []byte) T {
	t.Helper()

	var out T
	if err := json.Unmarshal(body, &out); err != nil {
		t.Fatalf("decode JSON: %v, body = %s", err, string(body))
	}
	return out
}

func expectError(t *testing.T, body []byte, wantCode string, wantField string) {
	t.Helper()

	errEnvelope := decodeJSON[apiresponse.ErrorEnvelope](t, body)
	if errEnvelope.Error.Code != wantCode {
		t.Fatalf("code = %s, want %s", errEnvelope.Error.Code, wantCode)
	}
	if wantField != "" && errEnvelope.Error.Fields[wantField] == "" {
		t.Fatalf("expected %s validation field", wantField)
	}
}

func createFund(t *testing.T, name string) apiresponse.Fund {
	t.Helper()

	resp, body := requestJSON(t, http.MethodPost, "/funds", validCreateFundRequest(name))
	expectStatus(t, resp, body, http.StatusCreated)
	return decodeJSON[apiresponse.Fund](t, body)
}

func createInvestor(t *testing.T, name string, email string) apiresponse.Investor {
	t.Helper()

	resp, body := requestJSON(t, http.MethodPost, "/investors", validCreateInvestorRequest(name, email))
	expectStatus(t, resp, body, http.StatusCreated)
	return decodeJSON[apiresponse.Investor](t, body)
}

func createInvestment(t *testing.T, fundID, investorID string, amount string, date string) apiresponse.Investment {
	t.Helper()

	resp, body := requestJSON(t, http.MethodPost, "/funds/"+fundID+"/investments", validCreateInvestmentRequest(investorID, amount, date))
	expectStatus(t, resp, body, http.StatusCreated)
	return decodeJSON[apiresponse.Investment](t, body)
}

func validCreateFundRequest(name string) request.CreateFundRequest {
	return request.CreateFundRequest{
		Name:          name,
		VintageYear:   2025,
		TargetSizeUSD: mustMoneyValue("500000000.00"),
		Status:        enum.FundStatusFundraising.String(),
	}
}

func validUpdateFundRequest(id string, name string) request.UpdateFundRequest {
	return request.UpdateFundRequest{
		ID:            id,
		Name:          name,
		VintageYear:   2025,
		TargetSizeUSD: mustMoneyValue("300000000.00"),
		Status:        enum.FundStatusInvesting.String(),
	}
}

func validCreateInvestorRequest(name string, email string) request.CreateInvestorRequest {
	return request.CreateInvestorRequest{
		Name:         name,
		InvestorType: enum.InvestorTypeInstitution.String(),
		Email:        email,
	}
}

func validCreateInvestmentRequest(investorID string, amount string, date string) request.CreateInvestmentRequest {
	return request.CreateInvestmentRequest{
		InvestorID:     investorID,
		AmountUSD:      mustMoneyValue(amount),
		InvestmentDate: date,
	}
}

func requireFundMatchesCreateRequest(t *testing.T, fund apiresponse.Fund, req request.CreateFundRequest) {
	t.Helper()

	if fund.Name != req.Name {
		t.Fatalf("name = %q, want %q", fund.Name, req.Name)
	}
	if fund.VintageYear != req.VintageYear {
		t.Fatalf("vintage_year = %d, want %d", fund.VintageYear, req.VintageYear)
	}
	if fund.TargetSizeUSD.String() != req.TargetSizeUSD.String() {
		t.Fatalf("target_size_usd = %s, want %s", fund.TargetSizeUSD.String(), req.TargetSizeUSD.String())
	}
	if fund.Status.String() != req.Status {
		t.Fatalf("status = %s, want %s", fund.Status, req.Status)
	}
	if fund.CreatedAt.String() == "" {
		t.Fatal("expected created_at")
	}
}

func requireFundMatchesUpdateRequest(t *testing.T, fund apiresponse.Fund, req request.UpdateFundRequest) {
	t.Helper()

	if fund.ID.String() != req.ID {
		t.Fatalf("id = %s, want %s", fund.ID, req.ID)
	}
	requireFundMatchesCreateRequest(t, fund, request.CreateFundRequest{
		Name:          req.Name,
		VintageYear:   req.VintageYear,
		TargetSizeUSD: req.TargetSizeUSD,
		Status:        req.Status,
	})
}

func requireInvestorMatchesCreateRequest(t *testing.T, investor apiresponse.Investor, req request.CreateInvestorRequest) {
	t.Helper()

	if investor.Name != req.Name {
		t.Fatalf("name = %q, want %q", investor.Name, req.Name)
	}
	if investor.InvestorType.String() != req.InvestorType {
		t.Fatalf("investor_type = %s, want %s", investor.InvestorType, req.InvestorType)
	}
	if investor.Email.String() != req.Email {
		t.Fatalf("email = %q, want %q", investor.Email.String(), req.Email)
	}
	if investor.CreatedAt.String() == "" {
		t.Fatal("expected created_at")
	}
}

func requireInvestmentMatchesCreateRequest(t *testing.T, investment apiresponse.Investment, fundID vo.ID, req request.CreateInvestmentRequest) {
	t.Helper()

	if investment.FundID != fundID {
		t.Fatalf("fund_id = %s, want %s", investment.FundID, fundID)
	}
	if investment.InvestorID.String() != req.InvestorID {
		t.Fatalf("investor_id = %s, want %s", investment.InvestorID, req.InvestorID)
	}
	if investment.AmountUSD.String() != req.AmountUSD.String() {
		t.Fatalf("amount_usd = %s, want %s", investment.AmountUSD.String(), req.AmountUSD.String())
	}
	if investment.InvestmentDate.String() != req.InvestmentDate {
		t.Fatalf("investment_date = %s, want %s", investment.InvestmentDate.String(), req.InvestmentDate)
	}
}

func requireContainsFund(t *testing.T, funds []apiresponse.Fund, id vo.ID) {
	t.Helper()
	for _, fund := range funds {
		if fund.ID == id {
			return
		}
	}
	t.Fatalf("fund %s not found in response", id)
}

func requireContainsInvestor(t *testing.T, investors []apiresponse.Investor, id vo.ID) {
	t.Helper()
	for _, investor := range investors {
		if investor.ID == id {
			return
		}
	}
	t.Fatalf("investor %s not found in response", id)
}

func requireContainsInvestment(t *testing.T, investments []apiresponse.Investment, id vo.ID) {
	t.Helper()
	for _, investment := range investments {
		if investment.ID == id {
			return
		}
	}
	t.Fatalf("investment %s not found in response", id)
}

func mustMoneyValue(raw string) vo.Money {
	value, err := decimal.NewFromString(raw)
	if err != nil {
		panic(fmt.Sprintf("parse test money %q: %v", raw, err))
	}
	return vo.NewMoney(value)
}
