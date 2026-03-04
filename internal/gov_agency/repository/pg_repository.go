package repository

import (
	"context"
	"database/sql"

	govagency "github.com/adohong4/driving-license/internal/gov_agency"
	"github.com/adohong4/driving-license/internal/models"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type GovAgencyRepo struct {
	db *sqlx.DB
}

// Goverment Agency new constructor
func NewGovAgencyRepo(db *sqlx.DB) govagency.Repository {
	return &GovAgencyRepo{db: db}
}

func (r *GovAgencyRepo) CreateGovAgency(ctx context.Context, gov *models.GovAgency) (*models.GovAgency, error) {
	g := &models.GovAgency{}
	if err := r.db.QueryRowxContext(ctx, createGovAgencyQuery,
		gov.Id, gov.Name, gov.UserAddress, gov.Address, gov.City, gov.Type, gov.Phone, gov.Email, gov.Status,
		gov.Version, gov.CreatedAt, gov.UpdatedAt, gov.Active,
	).StructScan(g); err != nil {
		return nil, errors.Wrap(err, "GovAgencyRepo.CreateGovAgency.StructScan")
	}
	return g, nil
}

func (r *GovAgencyRepo) UpdateGovAgency(ctx context.Context, gov *models.GovAgency) (*models.GovAgency, error) {
	g := &models.GovAgency{}
	if err := r.db.QueryRowxContext(ctx, updateGovAgencyQuery,
		gov.Name, gov.UserAddress, gov.Address, gov.City, gov.Type, gov.Phone, gov.Email, gov.Status, gov.UpdatedAt, gov.Id,
	).StructScan(g); err != nil {
		return nil, errors.Wrap(err, "GovAgencyRepo.UpdateGovAgency.StructScan")
	}
	return g, nil
}

func (r *GovAgencyRepo) DeleteGovAgency(ctx context.Context, gov *models.GovAgency) (*models.GovAgency, error) {
	g := &models.GovAgency{}
	if err := r.db.QueryRowxContext(ctx, deleteGovAgencyQuery, gov.UpdatedAt, gov.Id).StructScan(g); err != nil {
		return nil, errors.Wrap(err, "GovAgencyRepo.DeleteGovAgency.StructScan")
	}
	return g, nil
}

func (r *GovAgencyRepo) RevokeGovAgency(ctx context.Context, gov *models.GovAgency) (*models.GovAgency, error) {
	g := &models.GovAgency{}
	if err := r.db.QueryRowxContext(ctx, revokeGovAgencyQuery, gov.UpdatedAt, gov.Id).StructScan(g); err != nil {
		return nil, errors.Wrap(err, "GovAgencyRepo.RevokeGovAgency.StructScan")
	}
	return g, nil
}

func (r *GovAgencyRepo) GetGovAgency(ctx context.Context, pq *utils.PaginationQuery) (*models.GovAgencyList, error) {
	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTotalGovAgencyCount); err != nil {
		return nil, errors.Wrap(err, "GovAgencyRepo.GetGovAgency.GetContext.totalCount")
	}

	if totalCount == 0 {
		return &models.GovAgencyList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPage(totalCount, pq.GetSize()),
			Page:       pq.GetPage(),
			Size:       pq.GetSize(),
			HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			GovAgency:  make([]*models.GovAgency, 0),
		}, nil
	}

	var NewGovAgency = make([]*models.GovAgency, 0, pq.GetSize())
	rows, err := r.db.QueryxContext(ctx, getAllGovAgency, pq.GetOffset(), pq.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "GovAgencyRepo.GetGovAgency.NewGovAgency")
	}
	defer rows.Close()

	for rows.Next() {
		n := &models.GovAgency{}
		if err := rows.StructScan(n); err != nil {
			return nil, errors.Wrap(err, "GovAgencyRepo.GetGovAgency.StructScan")
		}
		NewGovAgency = append(NewGovAgency, n)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "NewGovAgency.rows.err")
	}

	return &models.GovAgencyList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPage(totalCount, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		GovAgency:  NewGovAgency,
	}, nil
}

func (r *GovAgencyRepo) GetGovAgencyByID(ctx context.Context, Id uuid.UUID) (*models.GovAgency, error) {
	g := &models.GovAgency{}
	if err := r.db.GetContext(ctx, g, getGovAgencyQuery, Id); err != nil {
		return nil, errors.Wrap(err, "GovAgencyRepo.GetGovAgencyByID.GetContext")
	}
	return g, nil
}

func (r *GovAgencyRepo) SearchByName(ctx context.Context, name string, query *utils.PaginationQuery) (*models.GovAgencyList, error) {
	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, searchGovAgencyByNameCount, name); err != nil {
		return nil, errors.Wrap(err, "GovAgencyRepo.SearchByName.GetContext.totalCount")
	}

	if totalCount == 0 {
		return &models.GovAgencyList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPage(totalCount, query.GetSize()),
			Page:       query.GetPage(),
			Size:       query.GetSize(),
			HasMore:    utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
			GovAgency:  make([]*models.GovAgency, 0),
		}, nil
	}

	var NewGovAgency = make([]*models.GovAgency, 0, query.GetSize())
	rows, err := r.db.QueryxContext(ctx, searchGovAgencyByName, name, query.GetOffset(), query.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "GovAgencyRepo.SearchByName.NewGovAgency")
	}
	defer rows.Close()

	for rows.Next() {
		n := &models.GovAgency{}
		if err := rows.StructScan(n); err != nil {
			return nil, errors.Wrap(err, "GovAgencyRepo.SearchByName.StructScan")
		}
		NewGovAgency = append(NewGovAgency, n)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "NewGovAgency.rows.err")
	}

	return &models.GovAgencyList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPage(totalCount, query.GetSize()),
		Page:       query.GetPage(),
		Size:       query.GetSize(),
		HasMore:    utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
		GovAgency:  NewGovAgency,
	}, nil
}

func (r *GovAgencyRepo) FindAgencyByUserAddress(ctx context.Context, g *models.GovAgency) (*models.GovAgency, error) {
	foundAgency := &models.GovAgency{}
	err := r.db.QueryRowxContext(ctx, findAgencyByUserAddress, g.UserAddress).StructScan(foundAgency)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "agencyRepo.FindAgencyByUserAddress.QueryRowxContext")
	}
	return foundAgency, nil
}
