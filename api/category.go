package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	db "github.com/Munchies-Engineering/backend/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

/*
	CreateCategory
	UpdateCategory
	DeleteCategory
	GetCategory
	ListCategories
	SearchCategories
*/

type createCategoryRequest struct {
	Name  string `json:"name" binding:"required"`
	Image string `json:"image" binding:"required"`
}

func (server *Server) createCategory(ctx *gin.Context) {
	var req createCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateCategoryParams{
		Name:  req.Name,
		Image: req.Image,
	}

	// TODO: add admin authentication

	category, err := server.store.CreateCategory(ctx, arg)
	if err != nil {
		if err := err.(*pq.Error); err != nil {
			switch err.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusConflict, errorResponse(errors.New("category already exists")))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, category)
}

type updateCategoryRequest struct {
	Name  string `json:"name" binding:"required"`
	Image string `json:"image" binding:"required"`
}

func (server *Server) updateCategory(ctx *gin.Context) {
	var req updateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("Invalid category ID")))
		return
	}

	arg := db.UpdateCategoryParams{
		ID:    id,
		Name:  req.Name,
		Image: req.Image,
	}

	// TODO: add admin authentication

	category, err := server.store.UpdateCategory(ctx, arg)
	if err != nil {
		fmt.Println(err)
		if err := err.(*pq.Error); err != nil {
			switch err.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusConflict, errorResponse(errors.New("category already exists")))
				return
			}
		} else if err := sql.ErrNoRows; err != nil {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("category not found")))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, category)
}

func (server *Server) getCategory(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("Invalid category ID")))
		return
	}

	category, err := server.store.GetCategory(ctx, id)
	if err != nil {
		if err := sql.ErrNoRows; err != nil {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("category not found")))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, category)
}

type listCategoriesRequest struct {
	PageSize int32 `form:"pageSize" binding:"required,min=1,max=100"`
	PageID   int32 `form:"pageId" binding:"required,min=1"`
}

func (server *Server) listCategories(ctx *gin.Context) {
	var req listCategoriesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListCategoriesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	categories, err := server.store.ListCategories(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, categories)
}

func (server *Server) deleteCategory(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("Invalid category ID")))
		return
	}

	category, err := server.store.GetCategory(ctx, id)
	if err != nil {
		if err := sql.ErrNoRows; err != nil {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("category not found")))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.DeleteCategory(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, category)
}
