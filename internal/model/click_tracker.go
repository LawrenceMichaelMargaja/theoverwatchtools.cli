package model

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
	"github.com/volatiletech/null/v8"
	"time"
)

type CreateClickTracker struct {
	Name   string `json:"name" validate:"required"`
	UserId int    `json:"user_id"`
}

type ClickTracker struct {
	Id                int       `json:"id" boil:"id"`
	Name              string    `json:"name" boil:"name"`
	UrlName           string    `json:"url_name" boil:"url_name"`
	RedirectUrl       int       `json:"redirect_url" boil:"redirect_url"`
	Clicks            int       `json:"clicks" boil:"clicks"`
	UniqueClicks      int       `json:"unique_clicks" boil:"unique_clicks"`
	LastImpressionAt  time.Time `json:"last_impression_at" boil:"last_impression_at"`
	ClickTrackerSetId int       `json:"click_tracker_set_id" boil:"click_tracker_set_id"`
	CreatedBy         string    `json:"created_by" boil:"created_by"`
	LastUpdatedBy     string    `json:"last_updated_by" boil:"last_updated_by"`
	CreatedAt         time.Time `json:"created_at" boil:"created_at"`
	LastUpdatedAt     null.Time `json:"last_updated_at" boil:"last_updated_at"`
	IsActive          bool      `json:"is_active" boil:"is_active"`
}

func (c *ClickTracker) Validate() error {
	if err := validationutils.Validate(c); err != nil {
		return fmt.Errorf("validate: %v", err)
	}
	return nil
}

type PaginatedClickTrackers struct {
	ClickTrackers []ClickTracker `json:"click_trackers"`
	Pagination    *Pagination    `json:"pagination"`
}

type ClickTrackerFilters struct {
	ClickTrackerNameIn     []string  `query:"click_tracker_name_in" json:"click_tracker_name_in"`
	ClickTrackerIsActive   null.Bool `query:"is_active" json:"is_active"`
	CreatedBy              null.Int  `query:"created_by" json:"created_by"`
	LastUpdatedBy          null.Int  `query:"last_updated_by" json:"last_updated_by"`
	IdsIn                  []int     `query:"ids_in" json:"ids_in"`
	PaginationQueryFilters `swaggerignore:"true"`
}

func (c *ClickTrackerFilters) Validate() error {
	if err := c.ValidatePagination(); err != nil {
		return fmt.Errorf("pagination: %v", err)
	}
	if err := validationutils.Validate(c); err != nil {
		return fmt.Errorf("click tracker filters: %v", err)
	}

	return nil
}
