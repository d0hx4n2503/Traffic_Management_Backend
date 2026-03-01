package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type VehicleRegistration struct {
	ID                uuid.UUID  `json:"id" db:"id" validate:"required"`
	OwnerID           *uuid.UUID `json:"owner_id" db:"owner_id"`
	UserAddress       *string    `json:"user_address,omitempty" db:"user_address"`
	Brand             string     `json:"brand" db:"brand"`
	TypeVehicle       string     `json:"type_vehicle" db:"type_vehicle"`
	VehiclePlateNo    string     `json:"vehicle_no" db:"vehicle_no"`
	ColorPlate        string     `json:"color_plate" db:"color_plate"`
	ChassisNo         string     `json:"chassis_no" db:"chassis_no"`
	EngineNo          string     `json:"engine_no" db:"engine_no"`
	ColorVehicle      string     `json:"color_vehicle" db:"color_vehicle"`
	OwnerName         string     `json:"owner_name" db:"owner_name"`
	Seats             *int       `json:"seats,omitempty" db:"seats"`
	IssueDate         string     `json:"issue_date" db:"issue_date"`
	Issuer            string     `json:"issuer" db:"issuer"`
	RegistrationCode  *string    `json:"registration_code" db:"registration_code"`
	RegistrationDate  *string    `json:"registration_date" db:"registration_date"`
	ExpiryDate        *string    `json:"expiry_date" db:"expiry_date"`
	RegistrationPlace *string    `json:"registration_place" db:"registration_place"`
	OnBlockchain      bool       `json:"on_blockchain" db:"on_blockchain"`
	BlockchainTxHash  string     `json:"blockchain_txhash" db:"blockchain_txhash"`
	Status            string     `json:"status" db:"status"`
	Version           int        `json:"version" db:"version"`
	CreatorId         uuid.UUID  `json:"creator_id" db:"creator_id"`
	ModifierId        *uuid.UUID `json:"modifier_id" db:"modifier_id"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
	Active            bool       `json:"active" db:"active"`
}

func (v *VehicleRegistration) PrepareCreate() error {
	v.VehiclePlateNo = strings.TrimSpace(v.VehiclePlateNo)
	v.ChassisNo = strings.TrimSpace(v.ChassisNo)
	v.EngineNo = strings.TrimSpace(v.EngineNo)
	v.Brand = strings.TrimSpace(v.Brand)
	v.ColorPlate = strings.TrimSpace(v.ColorPlate)
	v.ColorVehicle = strings.TrimSpace(v.ColorVehicle)
	v.OwnerName = strings.TrimSpace(v.OwnerName)
	v.Issuer = strings.TrimSpace(v.Issuer)
	v.Status = strings.TrimSpace(v.Status)

	v.OnBlockchain = false
	v.BlockchainTxHash = ""
	v.ID = uuid.New()
	v.CreatedAt = time.Now()
	v.UpdatedAt = time.Now()
	v.Active = true
	v.Version = 1
	return nil
}

func (v *VehicleRegistration) PrepareUpdate() error {
	v.VehiclePlateNo = strings.TrimSpace(v.VehiclePlateNo)
	v.ChassisNo = strings.TrimSpace(v.ChassisNo)
	v.EngineNo = strings.TrimSpace(v.EngineNo)
	v.Brand = strings.TrimSpace(v.Brand)
	v.ColorPlate = strings.TrimSpace(v.ColorPlate)
	v.ColorVehicle = strings.TrimSpace(v.ColorVehicle)
	v.OwnerName = strings.TrimSpace(v.OwnerName)
	v.Issuer = strings.TrimSpace(v.Issuer)
	v.Status = strings.TrimSpace(v.Status)

	v.UpdatedAt = time.Now()
	v.Version++
	return nil
}

type VehicleRegistrationList struct {
	TotalCount      int                    `json:"total_count"`
	TotalPages      int                    `json:"total_pages"`
	Page            int                    `json:"page"`
	Size            int                    `json:"size"`
	HasMore         bool                   `json:"has_more"`
	VehicleDocument []*VehicleRegistration `json:"vehicle_registration"`
}

type ConfirmBlockchainRequest struct {
	BlockchainTxHash string `json:"blockchain_txhash" validate:"required"`
	OnBlockchain     bool   `json:"on_blockchain"`
}

type CountItem struct {
	Key   string `json:"key"`
	Count int    `json:"count"`
}

type VehicleTypeCounts []*CountItem

type StatusCounts []*CountItem

type BrandCounts []*CountItem
