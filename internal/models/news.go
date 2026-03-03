package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
)

type News struct {
	Id         uuid.UUID      `json:"id" db:"id" validate:"required"`
	Code       string         `json:"code" db:"code"`
	Image      string         `json:"image" db:"image"`
	Title      string         `json:"title" db:"title"`
	Content    string         `json:"content" db:"content"`
	Category   string         `json:"category" db:"category"`
	Author     string         `json:"author" db:"author"`
	Type       string         `json:"type" db:"type"`
	Tag        types.JSONText `json:"tag" db:"tag"`
	View       int            `json:"view" db:"view"`
	Status     string         `json:"status" db:"status"`
	Version    int            `json:"version" db:"version"`
	CreatorId  uuid.UUID      `json:"creator_id" db:"creator_id"`
	ModifierID *uuid.UUID     `json:"modifier_id" db:"modifier_id"`
	CreatedAt  time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at" db:"updated_at"`
	Active     bool           `json:"active" db:"active"`
}

func (n *News) PrepareCreate() error {
	n.Title = strings.TrimSpace(n.Title)
	n.Content = strings.TrimSpace(n.Content)
	n.Code = strings.TrimSpace(n.Code)

	n.Id = uuid.New()
	if n.Code == "" {
		n.Code = "NEWS-" + strings.ToUpper(uuid.NewString()[:8])
	}
	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()
	n.Active = true
	n.Version = 1
	return nil
}

func (n *News) PrepareUpdate() error {
	n.Title = strings.TrimSpace(n.Title)
	n.Content = strings.TrimSpace(n.Content)

	n.UpdatedAt = time.Now()
	return nil
}

func (n *News) CheckView() error {
	n.View += 1
	return nil
}

type NewsList struct {
	TotalCount int     `json:"total_count"`
	TotalPages int     `json:"total_pages"`
	Page       int     `json:"page"`
	Size       int     `json:"size"`
	HasMore    bool    `json:"has_more"`
	News       []*News `json:"news"`
}
