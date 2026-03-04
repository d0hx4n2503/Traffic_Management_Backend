package govagency

import (
	"context"

	"github.com/adohong4/driving-license/internal/models"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/google/uuid"
)

type Repository interface {
	CreateGovAgency(ctx context.Context, gov *models.GovAgency) (*models.GovAgency, error)
	UpdateGovAgency(ctx context.Context, gov *models.GovAgency) (*models.GovAgency, error)
	DeleteGovAgency(ctx context.Context, gov *models.GovAgency) (*models.GovAgency, error)
	RevokeGovAgency(ctx context.Context, gov *models.GovAgency) (*models.GovAgency, error)
	GetGovAgency(ctx context.Context, pq *utils.PaginationQuery) (*models.GovAgencyList, error)
	GetGovAgencyByID(ctx context.Context, Id uuid.UUID) (*models.GovAgency, error)
	SearchByName(ctx context.Context, name string, query *utils.PaginationQuery) (*models.GovAgencyList, error)
	FindAgencyByUserAddress(ctx context.Context, g *models.GovAgency) (*models.GovAgency, error)
}
