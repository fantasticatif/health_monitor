package data

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HP_STATUS string

const (
	HP_STATUS_OK   HP_STATUS = "ok"
	HP_STATUS_DOWN HP_STATUS = "down"
	HP_STATUS_LATE HP_STATUS = "late"
	HP_STATUS_NONE HP_STATUS = "none"
)

type HitPoint struct {
	gorm.Model
	UUID           string `gorm:"type:varchar(100);primaryKey;<-:create"`
	Name           string `gorm:"not null"` // add composite constraint
	Description    string `gorm:"type:text"`
	ProjectID      uint
	Project        Project
	AccountID      uint `gorm:"not null;"`
	Account        Account
	Status         HP_STATUS              `gorm:"type:varchar(20)"`
	Schedule       map[string]interface{} `gorm:"type:json;serializer:json"`
	LastHitEventID uint
	LastHitEventAt *time.Time
	// LastHitEvent   *HitEvent `gorm:"foreignKey:LastHitEventID"`
}

func (h *HitPoint) BeforeCreate(db *gorm.DB) error {
	h.UUID = uuid.New().String()
	return nil
}

func (h *HitPoint) Create(db *gorm.DB) error {
	h.BeforeCreate(db)
	tx := db.Create(h)
	return tx.Error
}

func (h *HitPoint) GetSchedule() (*HitpointSchedule, error) {
	if h.Schedule == nil {
		return nil, nil
	}
	b, err := json.Marshal(h.Schedule)
	if err != nil {
		return nil, err
	}
	var s HitpointSchedule
	err = json.Unmarshal(b, &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (h *HitPoint) GetStatus() (HP_STATUS, error) {
	s, err := h.GetSchedule()
	if err != nil {
		return HP_STATUS_NONE, err
	}
	if s.IntervalInMinute != nil {
		timeFrom := h.LastHitEventAt
		if timeFrom != nil {
			timeFrom = &h.CreatedAt
		}
		duration := time.Since(*timeFrom)
		if duration.Minutes() < float64(*s.IntervalInMinute) {
			return HP_STATUS_OK, nil
		} else if duration.Minutes() < (float64(*s.IntervalInMinute) + float64(s.GracePeriodInMinute)) {
			return HP_STATUS_LATE, nil
		} else {
			return HP_STATUS_DOWN, nil
		}
	}
	return HP_STATUS_NONE, nil
}
