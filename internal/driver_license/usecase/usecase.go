package usecase

import (
	"context"
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/adohong4/driving-license/config"
	driverlicense "github.com/adohong4/driving-license/internal/driver_license"
	"github.com/adohong4/driving-license/internal/models"
	"github.com/adohong4/driving-license/pkg/httpErrors"
	"github.com/adohong4/driving-license/pkg/logger"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type DriverLicenseUC struct {
	cfg               *config.Config
	DriverLicenseRepo driverlicense.Repository
	logger            logger.Logger
}

func NewDriverLicenseUseCase(cfg *config.Config, DriverLicenseRepo driverlicense.Repository, logger logger.Logger) driverlicense.UseCase {
	return &DriverLicenseUC{cfg: cfg, DriverLicenseRepo: DriverLicenseRepo, logger: logger}
}

func (u *DriverLicenseUC) CreateDriverLicense(ctx context.Context, dl *models.DrivingLicense) (*models.DrivingLicense, error) {
	existsLicenseNO, err := u.DriverLicenseRepo.FindLicenseNO(ctx, dl)
	if existsLicenseNO != nil {
		return nil, httpErrors.NewRestErrorWithMessage(http.StatusBadRequest, httpErrors.ErrLicenseAlreadyExists, nil)
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "DriverLicenseUC.CreateDriverLicense.FindLicenseNO")
	}

	if err = dl.PrepareCreate(); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "DriverLicenseUC.CreateDriverLicense.PrepareCreate"))
	}

	user, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		return nil, httpErrors.NewUnauthorizedError(errors.WithMessage(err, "DriverLicenseUC.Create.GetUserFromCtx"))
	}

	dl.CreatorId = user.Id

	if err = utils.ValidateStruct(ctx, dl); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.WithMessage(err, "DriverLicenseUC.Create.ValidateStruct"))
	}

	n, err := u.DriverLicenseRepo.CreateDriverLicense(ctx, dl)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (u *DriverLicenseUC) UpdateDriverLicense(ctx context.Context, dl *models.DrivingLicense) (*models.DrivingLicense, error) {
	user, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		return nil, httpErrors.NewUnauthorizedError(errors.WithMessage(err, "DriverLicenseUC.UpdateDriverLicense.GetUserFromCtx"))
	}

	dl.ModifierId = &user.Id

	if err = dl.PrepareUpdate(); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "DriverLicenseUC.UpdateDriverLicense.PrepareCreate"))
	}

	if err := utils.ValidateStruct(ctx, dl); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.WithMessage(err, "DriverLicenseUC.UpdateDriverLicense.ValidateStruct"))
	}

	updatedLicense, err := u.DriverLicenseRepo.UpdateDriverLicense(ctx, dl)
	if err != nil {
		return nil, err
	}

	return updatedLicense, nil
}

func (u *DriverLicenseUC) ConfirmBlockchainStorage(ctx context.Context, dl *models.DrivingLicense) (*models.DrivingLicense, error) {
	user, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		return nil, httpErrors.NewUnauthorizedError(errors.WithMessage(err, "DriverLicenseUC.ConfirmBlockchainStorage.GetUserFromCtx"))
	}

	dl.ModifierId = &user.Id
	dl.OnBlockchain = true
	dl.UpdatedAt = time.Now()

	if dl.BlockchainTxHash == "" {
		return nil, httpErrors.NewBadRequestError(errors.New("BlockchainTxHash is required"))
	}
	if !strings.HasPrefix(dl.BlockchainTxHash, "0x") || len(dl.BlockchainTxHash) != 66 {
		return nil, httpErrors.NewBadRequestError(errors.New("BlockchainTxHash is invalid"))
	}

	updatedLicense, err := u.DriverLicenseRepo.ConfirmBlockchainStorage(ctx, dl)
	if err != nil {
		return nil, err
	}

	return updatedLicense, nil
}

func (u *DriverLicenseUC) AddWalletAddress(ctx context.Context, dl *models.DrivingLicense) (*models.DrivingLicense, error) {
	user, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		return nil, httpErrors.NewUnauthorizedError(errors.WithMessage(err, "DriverLicenseUC.AddWalletAddress.GetUserFromCtx"))
	}

	dl.ModifierId = &user.Id
	dl.UpdatedAt = time.Now()

	if dl.WalletAddress == "" {
		return nil, httpErrors.NewBadRequestError(errors.New("WalletAddress is required"))
	}

	updatedLicense, err := u.DriverLicenseRepo.UpdateWalletAddress(ctx, dl)
	if err != nil {
		return nil, err
	}

	return updatedLicense, nil
}

func (u *DriverLicenseUC) DeleteDriverLicense(ctx context.Context, dl *models.DrivingLicense) (*models.DrivingLicense, error) {
	user, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		return nil, httpErrors.NewUnauthorizedError(errors.WithMessage(err, "DriverLicenseUC.DeleteDriverLicense.GetUserFromCtx"))
	}

	dl.ModifierId = &user.Id

	if err = dl.PrepareUpdate(); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "DriverLicenseUC.DeleteDriverLicense.PrepareCreate"))
	}

	if err := utils.ValidateStruct(ctx, dl); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.WithMessage(err, "DriverLicenseUC.DeleteDriverLicense.ValidateStruct"))
	}

	updatedLicense, err := u.DriverLicenseRepo.DeleteDriverLicense(ctx, dl)
	if err != nil {
		return nil, err
	}

	return updatedLicense, nil
}

func (u *DriverLicenseUC) GetDriverLicense(ctx context.Context, pq *utils.PaginationQuery) (*models.DrivingLicenseList, error) {
	return u.DriverLicenseRepo.GetDriverLicense(ctx, pq)
}

func (u *DriverLicenseUC) GetDriverLicenseById(ctx context.Context, Id uuid.UUID) (*models.DrivingLicense, error) {
	n, err := u.DriverLicenseRepo.GetDriverLicenseById(ctx, Id)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (u *DriverLicenseUC) GetDriverLicenseByWalletAddress(ctx context.Context, address string) (*models.DrivingLicense, error) {
	n, err := u.DriverLicenseRepo.GetDriverLicenseByWalletAddress(ctx, address)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (u *DriverLicenseUC) GetDriverLicenseByLicenseNO(ctx context.Context, address string) (*models.DrivingLicense, error) {
	n, err := u.DriverLicenseRepo.GetDriverLicenseByLicenseNO(ctx, address)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (u *DriverLicenseUC) SearchByLicenseNo(ctx context.Context, lno string, query *utils.PaginationQuery) (*models.DrivingLicenseList, error) {
	return u.DriverLicenseRepo.SearchByLicenseNo(ctx, lno, query)
}

func (u *DriverLicenseUC) GetStatusDistribution(ctx context.Context) (*models.StatusDistributionResponse, error) {
	return u.DriverLicenseRepo.GetStatusDistribution(ctx)
}

func (u *DriverLicenseUC) GetLicenseTypeDistribution(ctx context.Context) (*models.LicenseTypeDistributionResponse, error) {
	return u.DriverLicenseRepo.GetLicenseTypeDistribution(ctx)
}

func (u *DriverLicenseUC) GetLicenseTypeStatusDistribution(ctx context.Context) (*models.LicenseTypeDetailDistributionResponse, error) {
	return u.DriverLicenseRepo.GetLicenseTypeStatusDistribution(ctx)
}

func (u *DriverLicenseUC) GetCityStatusDistribution(ctx context.Context) (*models.CityDetailDistributionResponse, error) {
	return u.DriverLicenseRepo.GetCityStatusDistribution(ctx)
}

func (u *DriverLicenseUC) GetMyDrivingLicenses(ctx context.Context, identityNo string, pq *utils.PaginationQuery) (*models.DrivingLicenseList, error) {
	return u.DriverLicenseRepo.GetDrivingLicensesByIdentityNo(ctx, identityNo, pq)
}

func (u *DriverLicenseUC) GetMyDrivingLicenseById(ctx context.Context, identityNo string, id uuid.UUID) (*models.DrivingLicense, error) {
	dl, err := u.DriverLicenseRepo.GetDriverLicenseById(ctx, id)
	if err != nil {
		return nil, err
	}
	if dl == nil || dl.IdentityNo != identityNo {
		return nil, nil
	}
	return dl, nil
}

func (u *DriverLicenseUC) GetMyDrivingLicenseByLicenseNo(ctx context.Context, identityNo, licenseNo string) (*models.DrivingLicense, error) {
	dl, err := u.DriverLicenseRepo.GetDriverLicenseByLicenseNO(ctx, licenseNo)
	if err != nil {
		return nil, err
	}
	if dl == nil || dl.IdentityNo != identityNo {
		return nil, nil
	}
	return dl, nil
}
