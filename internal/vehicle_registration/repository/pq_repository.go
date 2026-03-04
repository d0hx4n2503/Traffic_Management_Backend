package repository

import (
	"context"
	"database/sql"

	"github.com/adohong4/driving-license/internal/models"
	vehiclelicense "github.com/adohong4/driving-license/internal/vehicle_registration"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Vehicle Document Repository
type vehicleDocRepo struct {
	db *sqlx.DB
}

// Vehicle document new constructor
func NewVehicleDocRepository(db *sqlx.DB) vehiclelicense.Repository {
	return &vehicleDocRepo{db: db}
}

func (r *vehicleDocRepo) CreateVehicleDoc(ctx context.Context, veDoc *models.VehicleRegistration) (*models.VehicleRegistration, error) {
	v := &models.VehicleRegistration{}
	if err := r.db.QueryRowxContext(ctx, createLicenseQuery,
		veDoc.ID, veDoc.OwnerID, veDoc.Brand, veDoc.TypeVehicle, veDoc.VehiclePlateNo, veDoc.ColorPlate, veDoc.ChassisNo, veDoc.EngineNo, veDoc.ColorVehicle,
		veDoc.OwnerName, veDoc.Seats, veDoc.IssueDate, veDoc.Issuer, veDoc.RegistrationCode, veDoc.RegistrationDate, veDoc.ExpiryDate, veDoc.RegistrationPlace, veDoc.OnBlockchain, veDoc.BlockchainTxHash,
		veDoc.Status, veDoc.Version, veDoc.CreatorId, veDoc.ModifierId, veDoc.CreatedAt, veDoc.UpdatedAt,
	).StructScan(v); err != nil {
		return nil, errors.Wrap(err, "vehicleDocRepo.CreateVehicleDoc.StructScan")
	}
	return v, nil
}

func (r *vehicleDocRepo) UpdateVehicleDoc(ctx context.Context, veDoc *models.VehicleRegistration) (*models.VehicleRegistration, error) {
	v := &models.VehicleRegistration{}
	if err := r.db.QueryRowxContext(ctx, updateLicenseQuery,
		veDoc.OwnerID, veDoc.Brand, veDoc.TypeVehicle, veDoc.VehiclePlateNo, veDoc.ColorPlate, veDoc.ChassisNo, veDoc.EngineNo,
		veDoc.ColorVehicle, veDoc.OwnerName, veDoc.Seats, veDoc.IssueDate, veDoc.Issuer, veDoc.RegistrationCode, veDoc.RegistrationDate, veDoc.ExpiryDate, veDoc.RegistrationPlace,
		veDoc.Status, veDoc.ModifierId, veDoc.ID,
	).StructScan(v); err != nil {
		return nil, errors.Wrap(err, "vehicleDocRepo.UpdateVehicleDoc.StructScan")
	}
	return v, nil
}

func (r *vehicleDocRepo) ConfirmBlockchainStorage(ctx context.Context, v *models.VehicleRegistration) (*models.VehicleRegistration, error) {
	d := &models.VehicleRegistration{}
	if err := r.db.QueryRowxContext(ctx, updateBlockchainConfirmationQuery,
		v.BlockchainTxHash, v.OnBlockchain, v.ModifierId, v.UpdatedAt, v.ID,
	).StructScan(d); err != nil {
		return nil, errors.Wrap(err, "VehicleDocRepo.ConfirmBlockchainStorage.StructScan")
	}
	return d, nil
}

func (r *vehicleDocRepo) DeleteVehicleDoc(ctx context.Context, veDoc *models.VehicleRegistration) (*models.VehicleRegistration, error) {
	v := &models.VehicleRegistration{}
	if err := r.db.QueryRowxContext(ctx, deleteLicenseQuery,
		veDoc.ModifierId, veDoc.UpdatedAt, veDoc.ID,
	).StructScan(v); err != nil {
		return nil, errors.Wrap(err, "VehicleDocRepo.DeleteVehicle.StructScan")
	}
	return v, nil
}

func (r *vehicleDocRepo) GetVehicleDocs(ctx context.Context, pq *utils.PaginationQuery) (*models.VehicleRegistrationList, error) {
	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTotalCount); err != nil {
		return nil, errors.Wrap(err, "VehicleDocRepo.GetContext.totalCount")
	}

	if totalCount == 0 {
		return &models.VehicleRegistrationList{
			TotalCount:      totalCount,
			TotalPages:      utils.GetTotalPage(totalCount, pq.GetSize()),
			Page:            pq.GetPage(),
			Size:            pq.GetSize(),
			HasMore:         utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			VehicleDocument: make([]*models.VehicleRegistration, 0),
		}, nil
	}

	var NewVehicleDocs = make([]*models.VehicleRegistration, 0, pq.GetSize())
	rows, err := r.db.QueryxContext(ctx, getVehicleDocuments, pq.GetOffset(), pq.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "VehicleDocRepo.GetVehicleDocs.QueryRowxContext")
	}
	defer rows.Close()
	for rows.Next() {
		n := &models.VehicleRegistration{}
		if err = rows.StructScan(n); err != nil {
			return nil, errors.Wrap(err, "VehicleDocRepo.GetNews.StructScan")
		}
		NewVehicleDocs = append(NewVehicleDocs, n)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "VehicleDocRepo.GetNews.rows.Err")
	}
	return &models.VehicleRegistrationList{
		TotalCount:      totalCount,
		TotalPages:      utils.GetTotalPage(totalCount, pq.GetSize()),
		Page:            pq.GetPage(),
		Size:            pq.GetSize(),
		HasMore:         utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		VehicleDocument: NewVehicleDocs,
	}, nil
}

func (r *vehicleDocRepo) GetVehicleByID(ctx context.Context, vehicleID uuid.UUID) (*models.VehicleRegistration, error) {
	v := &models.VehicleRegistration{}
	if err := r.db.GetContext(ctx, v, getLicenseQuery, vehicleID); err != nil {
		return nil, errors.Wrap(err, "VehicleDocRepo.GetVehicleByID.GetContext")
	}
	return v, nil
}

func (r *vehicleDocRepo) SearchByVehiclePlateNO(ctx context.Context, vePlaNO string, query *utils.PaginationQuery) (*models.VehicleRegistrationList, error) {
	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, findByVehiclePlateNOCount, vePlaNO); err != nil {
		return nil, errors.Wrap(err, "VehicleDocRepo.FindByVehiclePlateNOCount.GetContext")
	}

	if totalCount == 0 {
		return &models.VehicleRegistrationList{
			TotalCount:      totalCount,
			TotalPages:      utils.GetTotalPage(totalCount, query.GetSize()),
			Page:            query.GetPage(),
			Size:            query.GetSize(),
			HasMore:         utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
			VehicleDocument: make([]*models.VehicleRegistration, 0),
		}, nil
	}

	var NewVehicleDocs = make([]*models.VehicleRegistration, 0, query.GetSize())
	rows, err := r.db.QueryxContext(ctx, searchByVehiclePlateNO, vePlaNO, query.GetOffset(), query.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "NewVehicleDocs.FindByVehiclePlateNOCount.QueryxContext")
	}
	defer rows.Close()

	for rows.Next() {
		n := &models.VehicleRegistration{}
		if err = rows.StructScan(n); err != nil {
			return nil, errors.Wrap(err, "NewVehicleDocs.FindByVehiclePlateNOCount.StructScan")
		}
		NewVehicleDocs = append(NewVehicleDocs, n)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "NewVehicleDocs.FindByVehiclePlateNOCount.rows.err")
	}

	return &models.VehicleRegistrationList{
		TotalCount:      totalCount,
		TotalPages:      utils.GetTotalPage(totalCount, query.GetSize()),
		Page:            query.GetPage(),
		Size:            query.GetSize(),
		HasMore:         utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
		VehicleDocument: NewVehicleDocs,
	}, nil
}

func (r *vehicleDocRepo) FindVehiclePlateNO(ctx context.Context, veDoc *models.VehicleRegistration) (*models.VehicleRegistration, error) {
	foundVehicleReq := &models.VehicleRegistration{}
	err := r.db.QueryRowxContext(ctx, findVehiclePlateNO, veDoc.VehiclePlateNo).StructScan(foundVehicleReq)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "vehicleDocRepo.findVehiclePlateNO.QueryRowxContext")
	}
	return foundVehicleReq, nil
}

func (r *vehicleDocRepo) GetCountByType(ctx context.Context) ([]*models.CountItem, error) {
	var items []*models.CountItem
	rows, err := r.db.QueryxContext(ctx, getCountByType)
	if err != nil {
		return nil, errors.Wrap(err, "vehicleDocRepo.GetCountByType.QueryxContext")
	}
	defer rows.Close()

	for rows.Next() {
		var typeVehicle string
		var count int
		if err := rows.Scan(&typeVehicle, &count); err != nil {
			return nil, errors.Wrap(err, "vehicleDocRepo.GetCountByType.Scan")
		}
		items = append(items, &models.CountItem{
			Key:   typeVehicle,
			Count: count,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "vehicleDocRepo.GetCountByType.rows.Err")
	}

	return items, nil
}

func (r *vehicleDocRepo) GetTopBrands(ctx context.Context) ([]*models.CountItem, error) {
	var items []*models.CountItem
	rows, err := r.db.QueryxContext(ctx, getTopBrands)
	if err != nil {
		return nil, errors.Wrap(err, "vehicleDocRepo.GetTopBrands.QueryxContext")
	}
	defer rows.Close()

	var topSum int
	for rows.Next() {
		var brand string
		var count int
		if err := rows.Scan(&brand, &count); err != nil {
			return nil, errors.Wrap(err, "vehicleDocRepo.GetTopBrands.Scan")
		}
		items = append(items, &models.CountItem{Key: brand, Count: count})
		topSum += count
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "vehicleDocRepo.GetTopBrands.rows.Err")
	}

	// Calculate others
	var total int
	if err := r.db.GetContext(ctx, &total, getTotalActiveVehicles); err != nil {
		return nil, errors.Wrap(err, "vehicleDocRepo.GetTopBrands.GetTotalActiveVehicles")
	}
	others := total - topSum
	items = append(items, &models.CountItem{Key: "khác", Count: others})

	return items, nil
}

func (r *vehicleDocRepo) GetRegistrationStatusStats(ctx context.Context) (*models.StatusCounts, error) {
	var valid, expired, pending int

	err := r.db.QueryRowxContext(ctx, getRegistrationStatusStats).Scan(
		&valid,
		&expired,
		&pending,
	)
	if err != nil {
		return nil, errors.Wrap(err, "vehicleDocRepo.GetRegistrationStatusStats.QueryRowxContext")
	}

	items := []*models.CountItem{
		{Key: "hợp lệ", Count: valid},
		{Key: "hết hạn", Count: expired},
		{Key: "chờ đăng kiểm", Count: pending},
	}

	return (*models.StatusCounts)(&items), nil
}

func (r *vehicleDocRepo) GetVehiclesByOwnerID(ctx context.Context, ownerID uuid.UUID, pq *utils.PaginationQuery) (*models.VehicleRegistrationList, error) {
	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTotalCountByOwnerID, ownerID); err != nil {
		return nil, errors.Wrap(err, "vehicleDocRepo.GetVehiclesByOwnerID.totalCount")
	}

	if totalCount == 0 {
		return &models.VehicleRegistrationList{
			TotalCount:      totalCount,
			TotalPages:      utils.GetTotalPage(totalCount, pq.GetSize()),
			Page:            pq.GetPage(),
			Size:            pq.GetSize(),
			HasMore:         utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			VehicleDocument: make([]*models.VehicleRegistration, 0),
		}, nil
	}

	var vehicles []*models.VehicleRegistration
	rows, err := r.db.QueryxContext(ctx, getVehiclesByOwnerID, ownerID, pq.GetOffset(), pq.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "vehicleDocRepo.GetVehiclesByOwnerID.QueryxContext")
	}
	defer rows.Close()

	for rows.Next() {
		v := &models.VehicleRegistration{}
		if err := rows.StructScan(v); err != nil {
			return nil, errors.Wrap(err, "vehicleDocRepo.GetVehiclesByOwnerID.StructScan")
		}
		vehicles = append(vehicles, v)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "vehicleDocRepo.GetVehiclesByOwnerID.rows.Err")
	}

	return &models.VehicleRegistrationList{
		TotalCount:      totalCount,
		TotalPages:      utils.GetTotalPage(totalCount, pq.GetSize()),
		Page:            pq.GetPage(),
		Size:            pq.GetSize(),
		HasMore:         utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		VehicleDocument: vehicles,
	}, nil
}

func (r *vehicleDocRepo) GetVehicleByIDAndOwnerID(ctx context.Context, vehicleID, ownerID uuid.UUID) (*models.VehicleRegistration, error) {
	v := &models.VehicleRegistration{}
	err := r.db.GetContext(ctx, v, getVehicleByIDAndOwner, vehicleID, ownerID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "vehicleDocRepo.GetVehicleByIDAndOwnerID.GetContext")
	}
	return v, nil
}

func (r *vehicleDocRepo) GetInspections(ctx context.Context, pq *utils.PaginationQuery) (*models.VehicleRegistrationList, error) {
	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getInspectionsCount); err != nil {
		return nil, errors.Wrap(err, "vehicleDocRepo.GetInspections.totalCount")
	}

	if totalCount == 0 {
		return &models.VehicleRegistrationList{
			TotalCount:      totalCount,
			TotalPages:      utils.GetTotalPage(totalCount, pq.GetSize()),
			Page:            pq.GetPage(),
			Size:            pq.GetSize(),
			HasMore:         utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			VehicleDocument: make([]*models.VehicleRegistration, 0),
		}, nil
	}

	var inspections = make([]*models.VehicleRegistration, 0, pq.GetSize())
	rows, err := r.db.QueryxContext(ctx, getInspections, pq.GetOffset(), pq.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "vehicleDocRepo.GetInspections.QueryxContext")
	}
	defer rows.Close()
	for rows.Next() {
		n := &models.VehicleRegistration{}
		if err = rows.StructScan(n); err != nil {
			return nil, errors.Wrap(err, "vehicleDocRepo.GetInspections.StructScan")
		}
		inspections = append(inspections, n)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "vehicleDocRepo.GetInspections.rows.Err")
	}
	return &models.VehicleRegistrationList{
		TotalCount:      totalCount,
		TotalPages:      utils.GetTotalPage(totalCount, pq.GetSize()),
		Page:            pq.GetPage(),
		Size:            pq.GetSize(),
		HasMore:         utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		VehicleDocument: inspections,
	}, nil
}

func (r *vehicleDocRepo) GetByRegistrationCode(ctx context.Context, code string) (*models.VehicleRegistration, error) {
	v := &models.VehicleRegistration{}
	err := r.db.QueryRowxContext(ctx, getByRegistrationCode, code).StructScan(v)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "vehicleDocRepo.GetByRegistrationCode.QueryRowxContext")
	}
	return v, nil
}
