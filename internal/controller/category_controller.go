package controller

import (
	"go-manajemen-keuangan/internal/payload/request"
	"go-manajemen-keuangan/internal/payload/response"
	"go-manajemen-keuangan/internal/service"
	"go-manajemen-keuangan/internal/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	CategoryService *service.CategoryService
}

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

	ctx.JSON(http.StatusOK, response.ApiResponse{
		ResponseStatus:  true,
		ResponseMessage: "Get categories successful",
		Data: response.CategoryListResponse{
			Categories: categories,
		},
	})
}

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

	ctx.JSON(http.StatusOK, response.ApiResponse{
		ResponseStatus:  true,
		ResponseMessage: "Get category by id success",
		Data:            category,
	})
}

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

	category, err := c.CategoryService.CreateCategory(req.Name, userID)
	if err != nil {
		utility.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	ctx.JSON(http.StatusCreated, response.ApiResponse{
		ResponseStatus:  true,
		ResponseMessage: "Category created!",
		Data:            category,
	})
}

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

	category, err := c.CategoryService.UpdateCategory(uint(id), userID, req.Name)
	if err != nil {
		utility.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	ctx.JSON(http.StatusOK, response.ApiResponse{
		ResponseStatus:  true,
		ResponseMessage: "Category updated",
		Data:            category,
	})
}

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

	ctx.JSON(http.StatusOK, response.ApiResponse{
		ResponseStatus:  true,
		ResponseMessage: "Category deleted",
		Data:            nil,
	})
}
