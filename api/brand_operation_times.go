package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	db "github.com/Munchies-Engineering/backend/db/sqlc"
	"github.com/Munchies-Engineering/backend/token"
	"github.com/gin-gonic/gin"
)

type addOperationTimeRequest struct {
	DayOfWeek string `json:"dayOfWeek" binding:"required,oneof=monday tuesday wednesday thursday friday saturday sunday"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

type operationTimeResponseType struct {
	ID        int64  `json:"id"`
	BrandID   int64  `json:"brandId"`
	DayOfWeek string `json:"dayOfWeek"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

func operationTimeResponse(opt db.OperationTime) *operationTimeResponseType {
	day, ok := daysByIndex[opt.DayOfWeek]
	if !ok {
		panic(fmt.Sprintf("invalid day of week %d", opt.DayOfWeek))
	}

	return &operationTimeResponseType{
		ID:        opt.ID,
		BrandID:   opt.BrandID,
		DayOfWeek: day,
		StartTime: opt.StartTime.Format("15:04:05"),
		EndTime:   opt.EndTime.Format("15:04:05"),
	}
}

func (server *Server) addOperationTime(ctx *gin.Context) {
	var req addOperationTimeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	brandId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)

	day, ok := days[req.DayOfWeek]
	if !ok {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid day of week")))
		return
	}

	brand, err := server.store.GetBrand(ctx, brandId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if brand.VendorID != authPayload.VendorID {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("forbidden")))
		return
	}

	startTime, err := time.Parse("2006-01-02T15:04:05Z", fmt.Sprintf("2006-01-02T%sZ", req.StartTime))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid start time")))
		return
	}

	endTime, err := time.Parse("2006-01-02T15:04:05Z", fmt.Sprintf("2006-01-02T%sZ", req.EndTime))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid end time")))
		return
	}

	arg := db.CreateOperationTimeParams{
		BrandID:   brand.ID,
		DayOfWeek: day,
		StartTime: startTime,
		EndTime:   endTime,
	}

	operationTime, err := server.store.CreateOperationTime(ctx, arg)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, operationTimeResponse(operationTime))
}

type getOperationTimesRequest struct {
	PageSize int32 `form:"pageSize" binding:"required"`
	PageID   int32 `form:"pageId" binding:"required,min=1"`
}

func (server *Server) getOperationTimes(ctx *gin.Context) {
	var req getOperationTimesRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	brandId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListOperationTimesByBrandIdParams{
		BrandID: brandId,
		Limit:   req.PageSize,
		Offset:  (req.PageID - 1) * req.PageSize,
	}

	times, err := server.store.ListOperationTimesByBrandId(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var response []*operationTimeResponseType
	for _, time := range times {
		response = append(response, operationTimeResponse(time))
	}

	ctx.JSON(http.StatusOK, response)
}

type updateOperationTimeRequest struct {
	DayOfWeek string `json:"dayOfWeek" binding:"required,oneof=monday tuesday wednesday thursday friday saturday sunday"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

func (server *Server) updateOperationTime(ctx *gin.Context) {
	var req updateOperationTimeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("Invalid id")))
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)

	operationTime, err := server.store.GetOperationTime(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	brand, err := server.store.GetBrand(ctx, operationTime.BrandID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if brand.VendorID != authPayload.VendorID {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("forbidden")))
		return
	}

	day, ok := days[req.DayOfWeek]
	if !ok {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid day of week")))
		return
	}

	arg := db.UpdateOperationTimeParams{
		ID:        id,
		DayOfWeek: day,
		StartTime: operationTime.StartTime,
		EndTime:   operationTime.EndTime,
		BrandID:   operationTime.BrandID,
	}

	var startTime time.Time
	if req.StartTime != "" {
		t, err := time.Parse("2006-01-02T15:04:05Z", fmt.Sprintf("2006-01-02T%sZ", req.StartTime))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid start time")))
			return
		}

		startTime = t
	}

	var endTime time.Time
	if req.EndTime != "" {
		t, err := time.Parse("2006-01-02T15:04:05Z", fmt.Sprintf("2006-01-02T%sZ", req.EndTime))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid end time")))
			return
		}

		endTime = t
	}

	if !startTime.IsZero() {
		arg.StartTime = startTime
	}

	if !endTime.IsZero() {
		arg.EndTime = endTime
	}

	operationTime, err = server.store.UpdateOperationTime(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, operationTimeResponse(operationTime))
}

func (server *Server) deleteOperationTime(ctx *gin.Context) {
	var id int64
	var err error

	if id, err = strconv.ParseInt(ctx.Param("id"), 10, 64); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)

	operationTime, err := server.store.GetOperationTime(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	brand, err := server.store.GetBrand(ctx, operationTime.BrandID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if brand.VendorID != authPayload.VendorID {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("forbidden")))
		return
	}

	if err := server.store.DeleteOperationTime(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, operationTimeResponse(operationTime))
}
