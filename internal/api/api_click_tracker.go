package api

import (
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/utilities/errs"
	"github.com/gofiber/fiber/v2"
	"github.com/volatiletech/null/v8"
	"net/http"
)

func (a *Api) ListClickTrackers(ctx *fiber.Ctx) error {
	filter := model.ClickTrackerFilters{
		ClickTrackerIsActive: null.BoolFrom(true),
	}
	if err := ctx.QueryParser(&filter); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}

	if err := filter.Validate(); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}
	filter.SetPaginationDefaults()

	clickTrackers, err := a.cfg.ClickTrackerService.ListClickTrackers(ctx.Context(), &filter)
	return a.WriteResponse(ctx, http.StatusOK, clickTrackers, err)
}

func (a *Api) AddClickTracker(ctx *fiber.Ctx) error {
	var body model.CreateClickTracker
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}
	clickTracker, err := a.cfg.ClickTrackerService.AddClickTracker(ctx.Context(), &body)
	return a.WriteResponse(ctx, http.StatusCreated, clickTracker, err)
}
