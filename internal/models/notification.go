package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	Id         uuid.UUID `json:"id" db:"id" validate:"required"`
	Code       string    `json:"code" db:"code"`
	Title      string    `json:"title" db:"title"`
	Content    string    `json:"content" db:"content"`
	Type       string    `json:"type" db:"type"`
	Target     string    `json:"target" db:"target"`           // Đối tượng nhận (type: all/personal/group)
	TargetUser string    `json:"target_user" db:"target_user"` // CCCD
	Status     string    `json:"status" db:"status"`           // if user -> status = unread, if all --> status = success

	CreatorId  uuid.UUID  `json:"creator_id" db:"creator_id"`
	ModifierID *uuid.UUID `json:"modifier_id" db:"modifier_id"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
	Active     bool       `json:"active" db:"active"`
}

func (n *Notification) PrepareCreate() error {
	n.Title = strings.TrimSpace(n.Title)
	n.Content = strings.TrimSpace(n.Content)
	n.Target = strings.TrimSpace(n.Target)
	n.Code = strings.TrimSpace(n.Code)

	n.Id = uuid.New()
	if n.Code == "" {
		n.Code = "NOTI-" + strings.ToUpper(uuid.NewString()[:8])
	}
	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()
	n.Active = true
	return nil
}

func (n *Notification) PrepareUpdate() error {
	n.Title = strings.TrimSpace(n.Title)
	n.Content = strings.TrimSpace(n.Content)
	n.Target = strings.TrimSpace(n.Target)

	n.UpdatedAt = time.Now()
	return nil
}

type NotificationList struct {
	TotalCount   int             `json:"total_count"`
	TotalPages   int             `json:"total_pages"`
	Page         int             `json:"page"`
	Size         int             `json:"size"`
	HasMore      bool            `json:"has_more"`
	Notification []*Notification `json:"notifications"`
}
