package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/tizianocitro/climate-change/cc-data-provider/model"
)

func GetOrganizations(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"totalCount": 4,
		"pageCount":  1,
		"hasMore":    false,
		"items":      organizations,
	})
}

func GetOrganizationsNoPage(c *fiber.Ctx) error {
	return c.JSON(organizations)
}

func GetOrganization(c *fiber.Ctx) error {
	id := c.Params("organizationId")
	index, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(model.Organization{})
	}
	organization := organizations[index]
	return c.JSON(organization)

}

var organizations = []model.Organization{
	{
		ID:          "0",
		Name:        "Ecosystem",
		Description: "This organization is for the whole ecosystem",
	},
	{
		ID:          "1",
		Name:        "IMF Climate Data",
		Description: "The organization at https://climatedata.imf.org/pages/climatechange-data",
	},
}
