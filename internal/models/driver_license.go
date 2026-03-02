package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type DrivingLicense struct {
	Id               uuid.UUID  `json:"id" db:"id" validate:"required"`
	Name             string     `json:"full_name" db:"full_name"`
	Avatar           string     `json:"avatar" db:"avatar"`
	DOB              string     `json:"dob" db:"dob"`                     // Ngày sinh
	IdentityNo       string     `json:"identity_no" db:"identity_no"`     // Căn cước công dân
	OwnerAddress     string     `json:"owner_address" db:"owner_address"` // Địa chỉ
	OwnerCity        string     `json:"owner_city" db:"owner_city"`
	LicenseNo        string     `json:"license_no" db:"license_no"`               // Số bằng lái
	IssueDate        string     `json:"issue_date" db:"issue_date"`               // Ngày cấp
	ExpiryDate       *string    `json:"expiry_date" db:"expiry_date"`             // Ngày hết hạn (có thời hạn, vô thời hạn)
	Status           string     `json:"status" db:"status"`                       // Trạng thái (pending: chờ đợi, expired: hết hạn, active: hoạt động, pause: tạm dừng (point = 0), revoke: thu hồi)
	LicenseType      string     `json:"license_type" db:"license_type"`           // Loại bằng lái (A1, B1, B2, ...)
	AuthorityId      uuid.UUID  `json:"authority_id" db:"authority_id"`           // Mã nơi cấp
	IssuingAuthority string     `json:"issuing_authority" db:"issuing_authority"` // Nơi cấp
	Nationality      string     `json:"nationality" db:"nationality"`             // Quốc tịch (Việt Nam, Hàn Quốc, ....)
	Point            int        `json:"point" db:"point"`                         // Điểm bằng lái xe (0 < point < 12)
	WalletAddress    string     `json:"wallet_address" db:"wallet_address"`
	OnBlockchain     bool       `json:"on_blockchain" db:"on_blockchain"`         // Trạng thái lưu ở blockchain (lưa/ chưa lưu)
	BlockchainTxHash string     `json:"blockchain_txhash" db:"blockchain_txhash"` // Mã lưu ở blockchain
	Version          int        `json:"version" db:"version"`                     // Phiên bản, tự động tăng
	CreatorId        uuid.UUID  `json:"creator_id" db:"creator_id"`               // ID của người tạo
	ModifierId       *uuid.UUID `json:"modifier_id" db:"modifier_id"`             // ID của người sửa
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`               // Thời gian tạo
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`               // Thời gian cập nhật
	Active           bool       `json:"active" db:"active"`
}

// Prepare the driver license for creation
func (d *DrivingLicense) PrepareCreate() error {
	d.IdentityNo = strings.TrimSpace(d.IdentityNo)
	d.LicenseNo = strings.TrimSpace(d.LicenseNo)
	d.LicenseType = strings.TrimSpace(d.LicenseType)
	d.DOB = strings.TrimSpace(d.DOB)
	d.IssueDate = strings.TrimSpace(d.IssueDate)
	if d.ExpiryDate != nil {
		expiry := strings.TrimSpace(*d.ExpiryDate)
		if expiry == "" {
			d.ExpiryDate = nil
		} else {
			d.ExpiryDate = &expiry
		}
	}

	d.Id = uuid.New()
	d.Point = 12
	d.OnBlockchain = false
	d.BlockchainTxHash = ""
	d.WalletAddress = ""
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()
	d.Active = true
	d.Version = 1
	return nil
}

// Prepare the driver license for updating
func (d *DrivingLicense) PrepareUpdate() error {
	d.IdentityNo = strings.TrimSpace(d.IdentityNo)
	d.LicenseNo = strings.TrimSpace(d.LicenseNo)
	d.LicenseType = strings.TrimSpace(d.LicenseType)
	d.DOB = strings.TrimSpace(d.DOB)
	d.IssueDate = strings.TrimSpace(d.IssueDate)
	if d.ExpiryDate != nil {
		expiry := strings.TrimSpace(*d.ExpiryDate)
		if expiry == "" {
			d.ExpiryDate = nil
		} else {
			d.ExpiryDate = &expiry
		}
	}

	d.UpdatedAt = time.Now()
	return nil
}

// All driver license response
type DrivingLicenseList struct {
	TotalCount     int               `json:"total_count"`
	TotalPages     int               `json:"total_pages"`
	Page           int               `json:"page"`
	Size           int               `json:"size"`
	HasMore        bool              `json:"has_more"`
	DrivingLicense []*DrivingLicense `json:"driver_licenses"`
}

// Status Distribution Item
type StatusDistributionItem struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

// Status Distribution Response
type StatusDistributionResponse struct {
	Distribution []StatusDistributionItem `json:"distribution"`
	Total        int                      `json:"total"`
}

// License Type Distribution Item
type LicenseTypeDistributionItem struct {
	LicenseType string `json:"license_type"`
	Count       int    `json:"count"`
}

// License Type Distribution Response
type LicenseTypeDistributionResponse struct {
	Distribution []LicenseTypeDistributionItem `json:"distribution"`
	Total        int                           `json:"total"`
}

// License Type Detail Distribution
type LicenseTypeDetailDistribution struct {
	LicenseType string                   `json:"license_type"`
	Total       int                      `json:"total"`
	ByStatus    []StatusDistributionItem `json:"by_status"`
}

// License Type Detail Distribution Response
type LicenseTypeDetailDistributionResponse struct {
	Distribution []LicenseTypeDetailDistribution `json:"distribution"`
	GrandTotal   int                             `json:"grand_total"`
}

// City Status Item
type CityStatusItem struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

// City Detail Distribution
type CityDetailDistribution struct {
	OwnerCity string           `json:"owner_city"`
	Total     int              `json:"total"`
	ByStatus  []CityStatusItem `json:"by_status"`
}

// City Detail Distribution Response
type CityDetailDistributionResponse struct {
	Distribution []CityDetailDistribution `json:"distribution"`
	GrandTotal   int                      `json:"grand_total"`
}
