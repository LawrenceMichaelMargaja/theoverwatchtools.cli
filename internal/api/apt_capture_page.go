package api

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/utilities/errs"
	"github.com/gofiber/fiber/v2"
	"github.com/volatiletech/null/v8"
	"net/http"
	"strconv"
)

func (a *Api) ListCapturePages(ctx *fiber.Ctx) error {
	filter := model.CapturePageFilters{
		CapturePageIsActive: null.BoolFrom(true),
	}
	if err := ctx.QueryParser(&filter); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}

	if err := filter.Validate(); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}
	filter.SetPaginationDefaults()

	capturePages, err := a.cfg.CapturePageService.ListCapturePages(ctx.Context(), &filter)
	return a.WriteResponse(ctx, http.StatusOK, capturePages, err)
}

func (a *Api) CreateCapturePage(ctx *fiber.Ctx) error {
	var body model.CreateCapturePage
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}
	capturePage, err := a.cfg.CapturePageService.CreateCapturePage(ctx.Context(), &body)
	return a.WriteResponse(ctx, http.StatusCreated, capturePage, err)
}

func (a *Api) DeleteCapturePage(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	capturePageId, err := strconv.Atoi(id)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}

	deleteParams := &model.DeleteCapturePage{
		ID: capturePageId,
	}

	err = a.cfg.CapturePageService.DeleteCapturePage(ctx.Context(), deleteParams)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(errs.ToArr(err))
	}

	return ctx.SendStatus(http.StatusNoContent)
}

func (a *Api) UpdateCapturePage(ctx *fiber.Ctx) error {
	var body model.UpdateCapturePage
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}
	capturePage, err := a.cfg.CapturePageService.UpdateCapturePage(ctx.Context(), &body)
	return a.WriteResponse(ctx, http.StatusOK, capturePage, err)
}

func (a *Api) RestoreCapturePage(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	capturePageId, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid capture page ID: %v", err)
	}

	restoreParams := &model.RestoreCapturePage{ID: capturePageId}

	err = a.cfg.CapturePageService.RestoreCapturePage(ctx.Context(), restoreParams)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(errs.ToArr(err))
	}

	return ctx.SendStatus(http.StatusNoContent)
}
