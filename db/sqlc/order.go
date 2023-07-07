package db

import (
	"context"
	"time"
)

// Order Status
const (
	OrderStatusPending    = "PENDING"
	OrderStatusAccepted   = "ACCEPTED"
	OrderStatusRejected   = "REJECTED"
	OrderStatusCancelled  = "CANCELLED"
	OrderStatusReady      = "READY"
	OrderStatusDispatched = "PICKED"
	OrderStatusDelivered  = "DELIVERED"
)

// Payment Method
const (
	PaymentMethodCash   = "CASH"
	PaymentMethodOnline = "ONLINE"
)

// Payment Status
const (
	PaymentStatusPending = "PENDING"
	PaymentStatusPaid    = "PAID"
)

type PlaceOrderItemParams struct {
	ProductID int64 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
	BrandID   int64
	Price     float64
	Discount  float64
}

type PlaceOrderParams struct {
	PaymentMethod     string
	RiderNote         string
	DispatchTime      time.Time
	DeliveryAddressID int64
	OrderItems        []PlaceOrderItemParams
	Total             float64
	TotalDiscount     float64
	ServiceCharge     float64
	Payable           float64
}

func (store *SqlStore) PlaceOrder(ctx context.Context, arg PlaceOrderParams) (*Order, error) {

	// create order

	err := store.execTx(ctx, func(q *Queries) error {
		orderArg := CreateOrderParams{
			CustomerID:    1, // FIXME: get customer id from token
			Status:        OrderStatusPending,
			PaymentMethod: arg.PaymentMethod,
			PaymentStatus: PaymentStatusPending,
			RiderNote:     arg.RiderNote,
			DispatchTime:  arg.DispatchTime,
			Total:         arg.Total,
			TotalDiscount: arg.TotalDiscount,
			ServiceCharge: arg.ServiceCharge,
			Payable:       arg.Payable,
		}

		order, err := q.CreateOrder(ctx, orderArg)
		if err != nil {
			return err
		}

		for _, item := range arg.OrderItems {
			orderItemArg := CreateOrderItemParams{
				OrderID:   order.ID,
				ProductID: item.ProductID,
				BrandID:   item.BrandID,
				Price:     item.Price,
				Quantity:  item.Quantity,
				Discount:  item.Discount,
			}

			_, err := q.CreateOrderItem(ctx, orderItemArg)
			if err != nil {
				return err
			}
		}

		groupItemsByBrand := make(map[int64][]PlaceOrderItemParams)
		for _, item := range arg.OrderItems {
			groupItemsByBrand[item.BrandID] = append(groupItemsByBrand[item.BrandID], item)
		}

		for brandID, items := range groupItemsByBrand {
			total := 0.0
			for _, item := range items {
				total += item.Price * float64(item.Quantity)
			}
			brandOrderArg := CreateBrandOrderParams{
				BrandID: brandID,
				OrderID: order.ID,
				Status:  OrderStatusPending,
			}

			_, err := q.CreateBrandOrder(ctx, brandOrderArg)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// // create delivery address
	// deliveryAddressArg := db.CreateDeliveryAddressParams{
	// 	OrderID
	// 	CustomerID
	// 	Address
	// 	GeoPoint
	// 	Apartment
	// 	Area
	// 	Floor
	// 	Phone
	// }

	// // create order items
	// orderItemArg := db.CreateOrderItemParams{
	// 	OrderID
	// 	ProductID
	// 	BrandID
	// 	Price
	// 	Quantity
	// 	Discount
	// }
	// // create brand order
	// brandOrderArg := db.CreateBrandOrderParams{
	// 	BrandID
	// 	OrderID
	// 	Status
	// 	Total
	// 	Discount
	// 	Note
	// }

	return nil, nil

}
