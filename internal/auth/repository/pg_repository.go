package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/adohong4/driving-license/internal/auth"
	"github.com/adohong4/driving-license/internal/models"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Auth Repository
type authRepo struct {
	db *sqlx.DB
}

// Auth new constructor
func NewAuthRepository(db *sqlx.DB) auth.Repository {
	return &authRepo{db: db}
}

func (r *authRepo) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	u := &models.User{}
	if err := r.db.QueryRowxContext(ctx, createUserQuery,
		user.Id, user.UserAddress, user.IdentityNo,
		user.FullName, user.DateOfBirth, user.Gender, user.Nationality, user.PlaceOfOrigin, user.PlaceOfResidence,
		user.Active, user.Role,
		user.Version, user.CreatorId, user.ModifierId, user.CreatedAt, user.UpdatedAt,
	).StructScan(u); err != nil {
		return nil, errors.Wrap(err, "authRepo.CreateUser.StructScan")
	}
	return u, nil
}

func (r *authRepo) Update(ctx context.Context, user *models.User) (*models.User, error) {
	u := &models.User{}
	if err := r.db.QueryRowxContext(ctx, updateUserQuery,
		user.UserAddress, user.IdentityNo,
		user.FullName, user.DateOfBirth, user.Gender, user.Nationality, user.PlaceOfOrigin, user.PlaceOfResidence,
		user.Active, user.Role,
		user.CreatorId, user.ModifierId, user.Id, user.Version,
	).StructScan(u); err != nil {
		return nil, errors.Wrap(err, "authRepo.Update.QueryRowxContext")
	}
	return u, nil
}

func (r *authRepo) Delete(ctx context.Context, id uuid.UUID, modifierId uuid.UUID, version int) error {
	result, err := r.db.ExecContext(ctx, deleteUserQuery, modifierId, time.Now(), id, version)
	if err != nil {
		return errors.Wrap(err, "authRepo.Delete.ExecContext")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "authRepo.Delete.RowsAffected")
	}
	if rowsAffected == 0 {
		return errors.Wrap(sql.ErrNoRows, "authRepo.Delete.rowsAffected")
	}
	return nil
}

func (r *authRepo) GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	if err := r.db.QueryRowxContext(ctx, getUserQuery, id).StructScan(user); err != nil {
		return nil, errors.Wrap(err, "authRepo.GetByID.QueryRowxContext")
	}
	return user, nil
}

func (r *authRepo) FindByIdentityNO(ctx context.Context, identity string, query *utils.PaginationQuery) (*models.UsersList, error) {
	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTotalCount, identity); err != nil {
		return nil, errors.Wrap(err, "authRepo.FindByIdentityNO.GetContext.totalCount")
	}

	if totalCount == 0 {
		return &models.UsersList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPage(totalCount, query.GetSize()),
			Page:       query.GetPage(),
			Size:       query.GetSize(),
			HasMore:    utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
			Users:      make([]*models.User, 0),
		}, nil
	}

	var users []*models.User
	if err := r.db.SelectContext(ctx, &users, findUsers, identity, query.GetOffset(), query.GetLimit()); err != nil {
		return nil, errors.Wrap(err, "authRepo.FindByIdentityNO.SelectContext")
	}

	return &models.UsersList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPage(totalCount, query.GetSize()),
		Page:       query.GetPage(),
		Size:       query.GetSize(),
		HasMore:    utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
		Users:      users,
	}, nil
}

func (r *authRepo) GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*models.UsersList, error) {
	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTotal); err != nil {
		return nil, errors.Wrap(err, "authRepo.GetUsers.GetContext.totalCount")
	}

	if totalCount == 0 {
		return &models.UsersList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPage(totalCount, pq.GetSize()),
			Page:       pq.GetPage(),
			Size:       pq.GetSize(),
			HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			Users:      make([]*models.User, 0),
		}, nil
	}

	var users []*models.User
	if err := r.db.SelectContext(ctx, &users, getUsers, pq.GetOrderBy(), pq.GetOffset(), pq.GetLimit()); err != nil {
		return nil, errors.Wrap(err, "authRepo.GetUsers.SelectContext")
	}

	return &models.UsersList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPage(totalCount, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		Users:      users,
	}, nil
}

func (r *authRepo) FindByIdentity(ctx context.Context, user *models.User) (*models.User, error) {
	foundUser := &models.User{}
	err := r.db.QueryRowxContext(ctx, findUserByIdentity, user.IdentityNo).StructScan(foundUser)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "authRepo.FindByIdentity.QueryRowxContext")
	}
	return foundUser, nil
}

func (r *authRepo) FindByUserAddress(ctx context.Context, user *models.User) (*models.User, error) {
	foundUser := &models.User{}
	err := r.db.QueryRowxContext(ctx, findUserByUserAddress, user.UserAddress).StructScan(foundUser)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "authRepo.FindByUserAddress.QueryRowxContext")
	}
	return foundUser, nil
}

func (r *authRepo) GetUserIdentityAndNameByAddress(ctx context.Context, userAddress string) (identityNo, fullName string, err error) {
	err = r.db.QueryRowContext(ctx, getUserIdentityAndNameByAddress, userAddress).Scan(&identityNo, &fullName)
	if err == sql.ErrNoRows {
		return "", "", sql.ErrNoRows
	}
	if err != nil {
		return "", "", errors.Wrap(err, "authRepo.GetUserIdentityAndNameByAddress")
	}
	return identityNo, fullName, nil
}

func (r *authRepo) IsUserAddressLinked(ctx context.Context, identityNo string) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, checkUserAddressLinked, identityNo)
	if err != nil {
		return false, errors.Wrap(err, "authRepo.IsUserAddressLinked")
	}
	return exists, nil
}

func (r *authRepo) LinkWalletAddress(ctx context.Context, identityNo, walletAddress string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "authRepo.LinkWalletAddress.BeginTxx")
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	result, err := tx.ExecContext(ctx, linkWalletAddressQuery, walletAddress, identityNo)
	if err != nil {
		return errors.Wrap(err, "authRepo.LinkWalletAddress.updateUser")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "authRepo.LinkWalletAddress.RowsAffected")
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	if _, err := tx.ExecContext(ctx, syncDriverLicenseWalletByIdentityQuery, walletAddress, identityNo); err != nil {
		return errors.Wrap(err, "authRepo.LinkWalletAddress.syncDriverLicenseWallet")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "authRepo.LinkWalletAddress.Commit")
	}
	committed = true
	return nil
}

func (r *authRepo) UnlinkWalletAddress(ctx context.Context, identityNo string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "authRepo.UnlinkWalletAddress.BeginTxx")
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	result, err := tx.ExecContext(ctx, unlinkWalletAddressQuery, identityNo)
	if err != nil {
		return errors.Wrap(err, "authRepo.UnlinkWalletAddress.clearUserWallet")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "authRepo.UnlinkWalletAddress.RowsAffected")
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	if _, err := tx.ExecContext(ctx, clearDriverLicenseWalletByIdentityQuery, identityNo); err != nil {
		return errors.Wrap(err, "authRepo.UnlinkWalletAddress.clearDriverLicenseWallet")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "authRepo.UnlinkWalletAddress.Commit")
	}
	committed = true
	return nil
}
