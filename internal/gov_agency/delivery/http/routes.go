package http

import (
	govagency "github.com/adohong4/driving-license/internal/gov_agency"
	"github.com/labstack/echo/v4"
)

func MapGovAgencyRoutes(GovAgencyGroup *echo.Group, h govagency.Handlers) {
	GovAgencyGroup.POST("/create", h.CreateGovAgency())
	GovAgencyGroup.PUT("/:id", h.UpdateGovAgency())
	GovAgencyGroup.PUT("/:id/revoke", h.RevokeGovAgency())
	GovAgencyGroup.DELETE("/:id", h.DeleteGovAgency())
	GovAgencyGroup.GET("/:id", h.GetByID())
	GovAgencyGroup.GET("/getAll", h.GetAllGovAgency())
	GovAgencyGroup.GET("/search", h.SearchByName())
	GovAgencyGroup.POST("/connect-wallet", h.ConnectWallet())
}
