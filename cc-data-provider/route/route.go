package route

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/tizianocitro/climate-change/cc-data-provider/config"
	"github.com/tizianocitro/climate-change/cc-data-provider/controller"
	"github.com/tizianocitro/climate-change/cc-data-provider/repository"
)

func UseRoutes(app *fiber.App, context *config.Context) {
	basePath := app.Group("/cc-data-provider")
	useOrganizations(basePath)
	useEcosystem(basePath, context)
}

func useOrganizations(basePath fiber.Router) {
	organizations := basePath.Group("/organizations")
	organizations.Get("/", func(c *fiber.Ctx) error {
		log.Printf("GET /organizations called")
		return controller.GetOrganizations(c)
	})
	organizations.Get("/no_page", func(c *fiber.Ctx) error {
		log.Printf("GET /organizations/no_page called")
		return controller.GetOrganizationsNoPage(c)
	})
	organizations.Get("/:organizationId", func(c *fiber.Ctx) error {
		log.Printf("GET /organizations/:organizationId called")
		return controller.GetOrganization(c)
	})
	useOrganizationsTemperatures(organizations)
	useOrganizationsDioxide(organizations)
	useOrganizationsSea(organizations)
}

func useOrganizationsTemperatures(organizations fiber.Router) {
	temperatureController := controller.NewTemperatureController()

	temperatures := organizations.Group("/:organizationId/temperatures")
	temperatures.Get("/", func(c *fiber.Ctx) error {
		log.Printf("GET /:organizationId/temperatures called")
		return temperatureController.GetTemperatures(c)
	})
	temperaturesWithId := temperatures.Group("/:temperatureId")
	temperaturesWithId.Get("/", func(c *fiber.Ctx) error {
		log.Printf("GET /:organizationId/temperatures/:temperatureId called")
		return temperatureController.GetTemperature(c)
	})
	temperaturesWithId.Get("/desc", func(c *fiber.Ctx) error {
		log.Printf("GET /:organizationId/temperatures/:temperatureId/desc called")
		return temperatureController.GetTemperatureDescription(c)
	})
	temperaturesWithId.Get("/map", func(c *fiber.Ctx) error {
		log.Printf("GET /:organizationId/temperatures/:temperatureId/map called")
		return temperatureController.GetTemperatureMap(c)
	})
	temperaturesWithId.Get("/chart", func(c *fiber.Ctx) error {
		log.Printf("GET /:organizationId/temperatures/:temperatureId/chart called")
		return temperatureController.GetTemperatureChart(c)
	})
}

func useOrganizationsDioxide(organizations fiber.Router) {
	dioxideController := controller.NewDioxideController()

	dioxide := organizations.Group("/:organizationId/dioxide")
	dioxide.Get("/", func(c *fiber.Ctx) error {
		log.Printf("GET /:organizationId/dioxide called")
		return dioxideController.GetAllDioxide(c)
	})
	dioxideWithId := dioxide.Group("/:dioxideId")
	dioxideWithId.Get("/", func(c *fiber.Ctx) error {
		log.Printf("GET /:organizationId/dioxide/:dioxideId called")
		return dioxideController.GetDioxide(c)
	})
	dioxideWithId.Get("/desc", func(c *fiber.Ctx) error {
		log.Printf("GET /:organizationId/dioxide/:dioxideId/desc called")
		return dioxideController.GetDioxideDescription(c)
	})
	dioxideWithId.Get("/map", func(c *fiber.Ctx) error {
		log.Printf("GET /:organizationId/dioxide/:dioxideId/map called")
		return dioxideController.GetDioxideMap(c)
	})
	dioxideWithId.Get("/chart", func(c *fiber.Ctx) error {
		log.Printf("GET /:organizationId/dioxide/:dioxideId/chart called")
		return dioxideController.GetDioxideChart(c)
	})
}

func useOrganizationsSea(organizations fiber.Router) {
	seaController := controller.NewSeaController()

	sea := organizations.Group("/:organizationId/seas")
	sea.Get("/", func(c *fiber.Ctx) error {
		log.Printf("GET /:organizationId/seas called")
		return seaController.GetSeas(c)
	})
	seaWithId := sea.Group("/:seaId")
	seaWithId.Get("/", func(c *fiber.Ctx) error {
		log.Printf("GET /:organizationId/seas/:seaId called")
		return seaController.GetSea(c)
	})
	seaWithId.Get("/desc", func(c *fiber.Ctx) error {
		log.Printf("GET /:organizationId/seas/:seaId/desc called")
		return seaController.GetSeaDescription(c)
	})
	seaWithId.Get("/map", func(c *fiber.Ctx) error {
		log.Printf("GET /:organizationId/seas/:seaId/map called")
		return seaController.GetSeaMap(c)
	})
}

func useEcosystem(basePath fiber.Router, context *config.Context) {
	issueRepository := context.RepositoriesMap["issues"].(*repository.IssueRepository)
	issueController := controller.NewIssueController(issueRepository)

	ecosystem := basePath.Group("/issues")
	ecosystem.Get("/", func(c *fiber.Ctx) error {
		log.Printf("GET /issues called")
		return issueController.GetIssues(c)
	})
	ecosystem.Get("/:issueId", func(c *fiber.Ctx) error {
		log.Printf("GET /issues/:issueId called")
		return issueController.GetIssue(c)
	})
	ecosystem.Post("/", func(c *fiber.Ctx) error {
		log.Printf("POST /issues called")
		return issueController.SaveIssue(c)
	})
}
