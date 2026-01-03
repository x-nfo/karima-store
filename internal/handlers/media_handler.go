package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/karima-store/internal/services"
)

type MediaHandler struct {
	mediaService services.MediaService
}

func NewMediaHandler(mediaService services.MediaService) *MediaHandler {
	return &MediaHandler{
		mediaService: mediaService,
	}
}

// UploadMedia uploads an image file
// @Summary Upload media
// @Description Upload an image file and create a media record
// @Tags media
// @Accept multipart/form-data
// @Produce json
// @Param product_id formData int true "Product ID"
// @Param file formData file true "Image file"
// @Param position formData int false "Position"
// @Param is_primary formData bool false "Is primary"
// @Success 200 {object} map[string]interface{} "Upload response"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/media/upload [post]
func (h *MediaHandler) UploadMedia(c *fiber.Ctx) error {
	// Parse multipart form
	productIDStr := c.FormValue("product_id")
	if productIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "product_id is required",
			"code":    400,
		})
	}

	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid product_id",
			"code":    400,
		})
	}

	// Parse position (optional)
	positionStr := c.FormValue("position")
	position := 0
	if positionStr != "" {
		position, err = strconv.Atoi(positionStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid position",
				"code":    400,
			})
		}
	}

	// Parse is_primary (optional)
	isPrimaryStr := c.FormValue("is_primary")
	isPrimary := false
	if isPrimaryStr == "true" {
		isPrimary = true
	}

	// Get file
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "No file provided",
			"code":    400,
		})
	}

	// Validate file
	if err := h.mediaService.ValidateImageFile(file); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    400,
		})
	}

	// Upload file
	response, err := h.mediaService.UploadImage(file, uint(productID), position, isPrimary)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    500,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}

// DeleteMedia deletes a media record and its file
// @Summary Delete media
// @Description Delete a media record and its associated file
// @Tags media
// @Accept json
// @Produce json
// @Param id path int true "Media ID"
// @Success 200 {object} map[string]interface{} "Success response"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Media not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/media/:id [delete]
func (h *MediaHandler) DeleteMedia(c *fiber.Ctx) error {
	mediaIDStr := c.Params("id")
	if mediaIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "media_id is required",
			"code":    400,
		})
	}

	mediaID, err := strconv.ParseUint(mediaIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid media_id",
			"code":    400,
		})
	}

	if err := h.mediaService.DeleteMedia(uint(mediaID)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    404,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Media deleted successfully",
	})
}

// GetMediaByProduct retrieves all media for a product
// @Summary Get media by product
// @Description Retrieve all media files for a specific product
// @Tags media
// @Accept json
// @Produce json
// @Param product_id path int true "Product ID"
// @Success 200 {object} map[string]interface{} "Media list"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/products/:product_id/media [get]
func (h *MediaHandler) GetMediaByProduct(c *fiber.Ctx) error {
	productIDStr := c.Params("product_id")
	if productIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "product_id is required",
			"code":    400,
		})
	}

	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid product_id",
			"code":    400,
		})
	}

	mediaList, err := h.mediaService.GetMediaByProduct(uint(productID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    500,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   mediaList,
	})
}

// SetPrimaryMedia sets a media item as primary for a product
// @Summary Set primary media
// @Description Set a media item as the primary image for a product
// @Tags media
// @Accept json
// @Produce json
// @Param product_id path int true "Product ID"
// @Param media_id path int true "Media ID"
// @Success 200 {object} map[string]interface{} "Success response"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Media or product not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/products/:product_id/media/:media_id/primary [put]
func (h *MediaHandler) SetPrimaryMedia(c *fiber.Ctx) error {
	productIDStr := c.Params("product_id")
	if productIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "product_id is required",
			"code":    400,
		})
	}

	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid product_id",
			"code":    400,
		})
	}

	mediaIDStr := c.Params("media_id")
	if mediaIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "media_id is required",
			"code":    400,
		})
	}

	mediaID, err := strconv.ParseUint(mediaIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid media_id",
			"code":    400,
		})
	}

	if err := h.mediaService.SetPrimaryMedia(uint(mediaID), uint(productID)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    404,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Primary media set successfully",
	})
}
