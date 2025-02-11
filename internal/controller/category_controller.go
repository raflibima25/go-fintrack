package controller

import (
	"go-fintrack/internal/payload/request"
	"go-fintrack/internal/payload/response"
	"go-fintrack/internal/service"
	"go-fintrack/internal/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	CategoryService *service.CategoryService
}

// GetAllCategoriesHandler godoc
// @Summary 	Get all categories
// @Description Get all categories for logged in user
// @Tags 		categories
// @Accept 		json
// @Produce 	json
// @Security 	BearerAuth
// @Success 	200 {object} response.SuccessResponse
// @Failure 	400 {object} response.SuccessResponse
// @Failure 	401 {object} response.SuccessResponse
// @Router 		/category [get]
func (c *CategoryController) GetAllCategoriesHandler(ctx *gin.Context) {
	userID, err := utility.GetUserIDFromContext(ctx)
	if err != nil {
		utility.ErrorResponse(ctx, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	categories, err := c.CategoryService.GetCategories(userID)
	if err != nil {
		utility.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		ResponseStatus:  true,
		ResponseMessage: "Get categories successful",
		Data: response.CategoryListResponse{
			Categories: categories,
		},
	})
}

// GetCategoryIdHandler godoc
// @Summary 	Get category by ID
// @Description Get category by ID for logged in user
// @Tags 		categories
// @Accept 		json
// @Produce 	json
// @Security 	BearerAuth
// @Param 		id path int true "Category ID"
// @Success 	200 {object} response.SuccessResponse
// @Failure 	400 {object} response.SuccessResponse
// @Failure 	401 {object} response.SuccessResponse
// @Router 		/category/{id} [get]
func (c *CategoryController) GetCategoryIdHandler(ctx *gin.Context) {
	userID, err := utility.GetUserIDFromContext(ctx)
	if err != nil {
		utility.ErrorResponse(ctx, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utility.ErrorResponse(ctx, http.StatusBadRequest, "Invalid category ID", nil)
		return
	}

	category, err := c.CategoryService.GetCategoryByID(uint(id), userID)
	if err != nil {
		utility.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		ResponseStatus:  true,
		ResponseMessage: "Get category by id success",
		Data:            category,
	})
}

// CreateCategoryHandler godoc
// @Summary 	Create category
// @Description Create category for logged in user
// @Tags 		categories
// @Accept 		json
// @Produce 	json
// @Security 	BearerAuth
// @Param 		request body request.CategoryRequest true "Category data"
// @Success 	201 {object} response.SuccessResponse
// @Failure 	400 {object} response.SuccessResponse
// @Failure 	401 {object} response.SuccessResponse
// @Router 		/category [post]
func (c *CategoryController) CreateCategoryHandler(ctx *gin.Context) {
	userID, err := utility.GetUserIDFromContext(ctx)
	if err != nil {
		utility.ErrorResponse(ctx, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	var req request.CategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utility.ErrorResponse(ctx, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	category, err := c.CategoryService.CreateCategory(&req, userID)
	if err != nil {
		utility.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	ctx.JSON(http.StatusCreated, response.SuccessResponse{
		ResponseStatus:  true,
		ResponseMessage: "Category created!",
		Data:            category,
	})
}

// UpdateCategoryHandler godoc
// @Summary 	Update category
// @Description Update category for logged in user
// @Tags 		categories
// @Accept 		json
// @Produce 	json
// @Security 	BearerAuth
// @Param 		id path int true "Category ID"
// @Param 		request body request.UpdateCategoryRequest true "Category data"
// @Success 	200 {object} response.SuccessResponse
// @Failure 	400 {object} response.SuccessResponse
// @Failure 	401 {object} response.SuccessResponse
// @Router 		/category/{id} [put]
func (c *CategoryController) UpdateCategoryHandler(ctx *gin.Context) {
	userID, err := utility.GetUserIDFromContext(ctx)
	if err != nil {
		utility.ErrorResponse(ctx, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utility.ErrorResponse(ctx, http.StatusBadRequest, "Invalid category ID", nil)
		return
	}

	var req request.UpdateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utility.ErrorResponse(ctx, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	category, err := c.CategoryService.UpdateCategory(uint(id), userID, &req)
	if err != nil {
		utility.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		ResponseStatus:  true,
		ResponseMessage: "Category updated",
		Data:            category,
	})
}

// DeleteCategoryHandler godoc
// @Summary 	Delete category
// @Description Delete category for logged in user
// @Tags 		categories
// @Accept 		json
// @Produce 	json
// @Security 	BearerAuth
// @Param 		id path int true "Category ID"
// @Success 	200 {object} response.SuccessResponse
// @Failure 	400 {object} response.SuccessResponse
// @Failure 	401 {object} response.SuccessResponse
// @Router 		/category/{id} [delete]
func (c *CategoryController) DeleteCategoryHandler(ctx *gin.Context) {
	userID, err := utility.GetUserIDFromContext(ctx)
	if err != nil {
		utility.ErrorResponse(ctx, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utility.ErrorResponse(ctx, http.StatusBadRequest, "Invalid category ID", nil)
		return
	}

	if err := c.CategoryService.DeleteCategory(uint(id), userID); err != nil {
		utility.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		ResponseStatus:  true,
		ResponseMessage: "Category deleted",
		Data:            nil,
	})
}
