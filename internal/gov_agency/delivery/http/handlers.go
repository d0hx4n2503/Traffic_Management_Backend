package http

import (
	"net/http"

	"github.com/adohong4/driving-license/config"
	govagency "github.com/adohong4/driving-license/internal/gov_agency"
	"github.com/adohong4/driving-license/internal/models"
	"github.com/adohong4/driving-license/pkg/httpErrors"
	"github.com/adohong4/driving-license/pkg/logger"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Goverment Agency handlers
type GovAgencyHandlers struct {
	cfg         *config.Config
	GovAgencyUC govagency.UseCase
	logger      logger.Logger
}

func NewGovAgencyHandlers(cfg *config.Config, GovAgencyUC govagency.UseCase, logger logger.Logger) govagency.Handlers {
	return &GovAgencyHandlers{cfg: cfg, GovAgencyUC: GovAgencyUC, logger: logger}
}

// Create godoc
// @Summary Create  Agency
// @Description Create GovAgency handler
// @Tags Goverment Agency
// @Accept json
// @Produce json
// @Success 201 {object} models.GovAgency
// @Router /agency/create [post]
func (h GovAgencyHandlers) CreateGovAgency() echo.HandlerFunc {
	return func(c echo.Context) error {
		n := &models.GovAgency{}
		if err := c.Bind(n); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		ctx := c.Request().Context()
		CreatedGovAgency, err := h.GovAgencyUC.CreateGovAgency(ctx, n)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusCreated, CreatedGovAgency)
	}
}

// Update godoc
// @Summary Update Goverment Agency
// @Description Update GovAgency handler
// @Tags Goverment Agency
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} models.GovAgency
// @Router /agency/{id} [put]
func (h GovAgencyHandlers) UpdateGovAgency() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		GovAgencyUUID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		n := &models.GovAgency{}
		if err = c.Bind(n); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		n.Id = GovAgencyUUID
		UpdatedGovAgency, err := h.GovAgencyUC.UpdateGovAgency(ctx, n)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, UpdatedGovAgency)
	}
}

// Delete godoc
// @Summary Delete Goverment Agency
// @Description Delete by id GovAgency handler
// @Tags Goverment Agency
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} models.GovAgency
// @Router /agency/{id} [Delete]
func (h GovAgencyHandlers) DeleteGovAgency() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		GovAgencyUUID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		n := &models.GovAgency{}
		if err = c.Bind(n); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		n.Id = GovAgencyUUID
		DeletedGovAgency, err := h.GovAgencyUC.DeleteGovAgency(ctx, n)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, DeletedGovAgency)
	}
}

// Revoke godoc
// @Summary Revoke Goverment Agency
// @Description Revoke GovAgency by id
// @Tags Goverment Agency
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} models.GovAgency
// @Router /agency/{id}/revoke [put]
func (h GovAgencyHandlers) RevokeGovAgency() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		GovAgencyUUID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		n := &models.GovAgency{Id: GovAgencyUUID}
		revokedGovAgency, err := h.GovAgencyUC.RevokeGovAgency(ctx, n)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, revokedGovAgency)
	}
}

// GetByID godoc
// @Summary Get by Goverment Agency ID
// @Description Get by Goverment Agency handler
// @Tags Goverment Agency
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} models.GovAgency
// @Router /agency/{id} [get]
func (h GovAgencyHandlers) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		govAgencyUUID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		getGovAgencyID, err := h.GovAgencyUC.GetGovAgencyByID(ctx, govAgencyUUID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusCreated, getGovAgencyID)
	}
}

// GetAllGovermentAgency godoc
// @Summary Get all Goverment Agency
// @Description Get all Goverment Agency with pagination
// @Tags Goverment Agency
// @Accept json
// @Produce json
// @Param page query int false "page number" Format(page)
// @Param size query int false "number of elements per page" Format(size)
// @Param orderBy query int false "filter name" Format(orderBy)
// @Success 200 {object} models.GovAgency
// @Router /agency/getAll [get]
func (h GovAgencyHandlers) GetAllGovAgency() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		newList, err := h.GovAgencyUC.GetGovAgency(ctx, pq)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, newList)
	}
}

// SearchByName godoc
// @Summary Search by Name
// @Description Search Goverment Agency by Name
// @Tags Goverment Agency
// @Accept json
// @Produce json
// @Param page query int false "page number" Format(page)
// @Param size query int false "number of elements per page" Format(size)
// @Param orderBy query int false "filter name" Format(orderBy)
// @Success 200 {object} models.GovAgency
// @Router /agency/search [get]
func (h GovAgencyHandlers) SearchByName() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		newList, err := h.GovAgencyUC.SearchByName(ctx, c.QueryParam("name"), pq)

		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, newList)
	}
}

func (h *GovAgencyHandlers) ConnectWallet() echo.HandlerFunc {
	type ConnectWallet struct {
		UserAddress string `json:"user_address" db:"user_address"`
	}
	return func(c echo.Context) error {
		connectWallet := &ConnectWallet{}
		if err := utils.ReadRequest(c, connectWallet); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		ctx := c.Request().Context()
		userWithToken, err := h.GovAgencyUC.ConnectWallet(ctx, &models.GovAgency{
			UserAddress: connectWallet.UserAddress,
		})
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, userWithToken)
	}
}
