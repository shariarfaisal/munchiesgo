package api

import (
	"net/http"
	"time"

	db "github.com/Munchies-Engineering/backend/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createOrderItemRequest struct {
	ProductID int64 `json:"productId" binding:"required"`
	Quantity  int32 `json:"quantity" binding:"required"`
}

type createOrderRequest struct {
	PaymentMethod     string                   `json:"paymentMethod" binding:"required,oneof=cash online"`
	RiderNote         string                   `json:"riderNote"`
	DispatchTime      string                   `json:"dispatchTime"`
	DeliveryAddressID int64                    `json:"deliveryAddressId" binding:"required"`
	OrderItems        []createOrderItemRequest `json:"orderItems" binding:"required,min=1"`
	AddressId         int64                    `json:"addressId" binding:"required,min=1"`
}

func (server *Server) createOrder(ctx *gin.Context) {
	var req createOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := []db.PlaceOrderItemParams{}
	for _, item := range req.OrderItems {
		arg = append(arg, db.PlaceOrderItemParams{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	orderItemsProducts, err := server.store.OrderItemProducts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	itemsTotal := 0.0
	totalDiscount := 0.0
	orderItemsById := map[int64]*db.OrderItemProduct{}

	for _, item := range orderItemsProducts {
		orderItemsById[item.ID] = item

		if item.Parent.ID != 0 {
			orderItemsById[item.Parent.ID] = item
		}
	}

	for _, item := range orderItemsById {
		itemsTotal += item.Price
	}

	orderItemsArg := []db.PlaceOrderItemParams{}

	for _, item := range req.OrderItems {
		product := orderItemsById[item.ProductID]
		orderItemsArg = append(orderItemsArg, db.PlaceOrderItemParams{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			BrandID:   product.BrandID,
			Price:     product.Price,
			Discount:  0, // FIXME: Discount
		})
	}

	ordersArg := db.PlaceOrderParams{
		PaymentMethod:     req.PaymentMethod,
		RiderNote:         req.RiderNote,
		DispatchTime:      time.Now(), // FIXME:
		DeliveryAddressID: req.AddressId,
		OrderItems:        orderItemsArg,
		Total:             itemsTotal,
		TotalDiscount:     totalDiscount,
		ServiceCharge:     0,                          // FIXME: ServiceCharge
		Payable:           itemsTotal - totalDiscount, // FIXME: ServiceCharge
	}

	server.store.PlaceOrder(ctx, ordersArg)

	ctx.JSON(http.StatusOK, orderItemsProducts)
}
