package app

import (
	"io"
	"path/filepath"

	"github.com/gofiber/contrib/v3/swaggerui"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/rs/zerolog"

	"github.com/jaeyoung0509/titanbay-funds-api/internal/handler"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/middleware"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/port"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/usecase"
)

type Dependencies struct {
	Repo            port.Repository
	Logger          *zerolog.Logger
	SwaggerFilePath string
}

func New(deps Dependencies) *fiber.App {
	baseLogger := normalizedLogger(deps.Logger)

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.NewErrorHandler(baseLogger),
	})
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(middleware.NewRequestLogger(baseLogger))

	health := handler.NewHealthHandler()
	swaggerFilePath := deps.SwaggerFilePath
	if swaggerFilePath == "" {
		swaggerFilePath = filepath.Join("docs", "swagger", "openapi.yaml")
	}
	swagger := swaggerui.New(swaggerui.Config{
		BasePath: "/",
		Path:     "swagger",
		FilePath: swaggerFilePath,
	})
	app.Get("/", health.Check)
	app.Get("/health", health.Check)
	app.Get("/swagger", swagger)
	app.Get("/docs/swagger/openapi.yaml", swagger)

	funds := handler.NewFundHandler(usecase.NewFundService(deps.Repo))
	investors := handler.NewInvestorHandler(usecase.NewInvestorService(deps.Repo))
	investments := handler.NewInvestmentHandler(usecase.NewInvestmentService(deps.Repo))

	app.Get("/funds", funds.ListFunds)
	app.Post("/funds", funds.CreateFund)
	app.Put("/funds", funds.UpdateFund)
	app.Get("/funds/:id", funds.GetFund)

	app.Get("/investors", investors.ListInvestors)
	app.Post("/investors", investors.CreateInvestor)

	app.Get("/funds/:fund_id/investments", investments.ListInvestments)
	app.Post("/funds/:fund_id/investments", investments.CreateInvestment)

	return app
}

func normalizedLogger(base *zerolog.Logger) zerolog.Logger {
	if base != nil {
		return *base
	}

	return zerolog.New(io.Discard).With().Timestamp().Logger()
}
