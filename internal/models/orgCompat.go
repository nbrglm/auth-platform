package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/nbrglm/auth-platform/db"
)

type OrgCompat struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Slug        string         `json:"slug"`
	Description *string        `json:"description,omitempty"`
	AvatarURL   *string        `json:"avatarURL,omitempty"`
	Settings    map[string]any `json:"settings,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   *time.Time     `json:"deletedAt,omitempty"`
}

func NewOrgCompat(o *db.Org) *OrgCompat {
	var deletedAt *time.Time
	if o.DeletedAt.Valid {
		t := o.DeletedAt.Time
		deletedAt = &t
	}

	var description *string
	if o.Description != nil {
		description = o.Description
	}

	var avatarURL *string
	if o.AvatarUrl != nil {
		avatarURL = o.AvatarUrl
	}

	settings := make(map[string]any)
	if o.Settings != nil {
		err := json.Unmarshal(o.Settings, &settings)
		if err != nil {
			settings = make(map[string]any)
		}
	}

	return &OrgCompat{
		ID:          o.ID,
		Name:        o.Name,
		Slug:        o.Slug,
		Description: description,
		AvatarURL:   avatarURL,
		Settings:    settings,
		CreatedAt:   o.CreatedAt.Time,
		UpdatedAt:   o.UpdatedAt.Time,
		DeletedAt:   deletedAt,
	}
}
