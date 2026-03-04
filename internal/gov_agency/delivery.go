package govagency

import "github.com/labstack/echo/v4"

type Handlers interface {
	CreateGovAgency() echo.HandlerFunc
	UpdateGovAgency() echo.HandlerFunc
	DeleteGovAgency() echo.HandlerFunc
	RevokeGovAgency() echo.HandlerFunc
	GetByID() echo.HandlerFunc
	GetAllGovAgency() echo.HandlerFunc
	SearchByName() echo.HandlerFunc
	ConnectWallet() echo.HandlerFunc
}
