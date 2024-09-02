package api

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/model"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . categoryService
type categoryService interface {
	ListCategories(ctx context.Context, filters *model.CategoryFilters) (*model.PaginatedCategories, error)
	CreateCategory(ctx context.Context, category *model.CreateCategory) (*model.Category, error)
	UpdateCategory(ctx context.Context, category *model.UpdateCategory) (*model.Category, error)
	DeleteCategory(ctx context.Context, params *model.DeleteCategory) error
	RestoreCategory(ctx context.Context, params *model.RestoreCategory) error
}

//counterfeiter:generate . organizationService
type organizationService interface {
	ListOrganizations(ctx context.Context, filters *model.OrganizationFilters) (*model.PaginatedOrganizations, error)
	DeleteOrganization(ctx context.Context, params *model.DeleteOrganization) error
	UpdateOrganization(ctx context.Context, organization *model.UpdateOrganization) (*model.Organization, error)
	CreateOrganization(ctx context.Context, organization *model.CreateOrganization) (*model.Organization, error)
	RestoreOrganization(ctx context.Context, params *model.RestoreOrganization) error
}

//counterfeiter:generate . capturePageService
type capturePageService interface {
	ListCapturePages(ctx context.Context, filters *model.CapturePageFilters) (*model.PaginatedCapturePages, error)
	DeleteCapturePage(ctx context.Context, params *model.DeleteCapturePage) error
	UpdateCapturePage(ctx context.Context, capturePage *model.UpdateCapturePage) (*model.CapturePage, error)
	AddCapturePage(ctx context.Context, capturePage *model.CreateCapturePage) (*model.CapturePage, error)
	RestoreCapturePage(ctx context.Context, params *model.RestoreCapturePage) error
}
