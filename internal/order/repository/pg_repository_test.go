package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/dinorain/checkoutaja/internal/models"
	"github.com/dinorain/checkoutaja/pkg/utils"
)

func TestOrderRepository_Create(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	orderPGRepository := NewOrderPGRepository(sqlxDB)

	columns := []string{"order_id", "user_id", "seller_id", "item", "quantity", "total_price", "status", "delivery_source_address", "delivery_destination_address", "created_at", "updated_at"}
	orderUUID := uuid.New()
	userUUID := uuid.New()
	sellerUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID:  orderUUID,
		UserID:   userUUID,
		SellerID: sellerUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			SellerID:    sellerUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	rows := sqlmock.NewRows(columns).AddRow(
		orderUUID,
		mockOrder.UserID,
		mockOrder.SellerID,
		mockOrder.Item,
		mockOrder.Quantity,
		mockOrder.TotalPrice,
		mockOrder.Status,
		mockOrder.DeliverySourceAddress,
		mockOrder.DeliveryDestinationAddress,
		time.Now(),
		time.Now(),
	)

	mock.ExpectQuery(createOrderQuery).WithArgs(
		mockOrder.UserID,
		mockOrder.SellerID,
		mockOrder.Item,
		mockOrder.Quantity,
		mockOrder.TotalPrice,
		mockOrder.Status,
		mockOrder.DeliverySourceAddress,
		mockOrder.DeliveryDestinationAddress,
	).WillReturnRows(rows)

	createdOrder, err := orderPGRepository.Create(context.Background(), mockOrder)
	require.NoError(t, err)
	require.NotNil(t, createdOrder)
}

func TestOrderRepository_FindAll(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	orderPGRepository := NewOrderPGRepository(sqlxDB)

	columns := []string{"order_id", "user_id", "seller_id", "item", "quantity", "total_price", "status", "delivery_source_address", "delivery_destination_address", "created_at", "updated_at"}
	orderUUID := uuid.New()
	userUUID := uuid.New()
	sellerUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID:  orderUUID,
		UserID:   userUUID,
		SellerID: sellerUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			SellerID:    sellerUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	rows := sqlmock.NewRows(columns).AddRow(
		orderUUID,
		mockOrder.UserID,
		mockOrder.SellerID,
		mockOrder.Item,
		mockOrder.Quantity,
		mockOrder.TotalPrice,
		mockOrder.Status,
		mockOrder.DeliverySourceAddress,
		mockOrder.DeliveryDestinationAddress,
		time.Now(),
		time.Now(),
	)

	size := 10
	mock.ExpectQuery(findAllQuery).WithArgs(size, 0).WillReturnRows(rows)
	foundOrders, err := orderPGRepository.FindAll(context.Background(), utils.NewPaginationQuery(size, 1))
	require.NoError(t, err)
	require.NotNil(t, foundOrders)
	require.Equal(t, len(foundOrders), 1)

	mock.ExpectQuery(findAllQuery).WithArgs(size, 10).WillReturnRows(rows)
	foundOrders, err = orderPGRepository.FindAll(context.Background(), utils.NewPaginationQuery(size, 2))
	require.NoError(t, err)
	require.Nil(t, foundOrders)
}

func TestOrderRepository_FindAllBySellerId(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	orderPGRepository := NewOrderPGRepository(sqlxDB)

	columns := []string{"order_id", "user_id", "seller_id", "item", "quantity", "total_price", "status", "delivery_source_address", "delivery_destination_address", "created_at", "updated_at"}
	orderUUID := uuid.New()
	userUUID := uuid.New()
	sellerUUID := uuid.New()
	otherSellerUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID:  orderUUID,
		UserID:   userUUID,
		SellerID: otherSellerUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			SellerID:    otherSellerUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	mockOtherOrder := &models.Order{
		OrderID:  orderUUID,
		UserID:   userUUID,
		SellerID: sellerUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			SellerID:    sellerUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	otherRows := sqlmock.NewRows(columns).AddRow(
		orderUUID,
		mockOtherOrder.UserID,
		mockOtherOrder.SellerID,
		mockOtherOrder.Item,
		mockOtherOrder.Quantity,
		mockOtherOrder.TotalPrice,
		mockOtherOrder.Status,
		mockOtherOrder.DeliverySourceAddress,
		mockOtherOrder.DeliveryDestinationAddress,
		time.Now(),
		time.Now(),
	)

	rows := sqlmock.NewRows(columns).AddRow(
		orderUUID,
		mockOrder.UserID,
		mockOrder.SellerID,
		mockOrder.Item,
		mockOrder.Quantity,
		mockOrder.TotalPrice,
		mockOrder.Status,
		mockOrder.DeliverySourceAddress,
		mockOrder.DeliveryDestinationAddress,
		time.Now(),
		time.Now(),
	)

	size := 10
	mock.ExpectQuery(findAllBySellerIdQuery).WithArgs(mockOrder.SellerID, size, 0).WillReturnRows(rows)
	foundOrders, err := orderPGRepository.FindAllBySellerId(context.Background(), mockOrder.SellerID, utils.NewPaginationQuery(size, 1))
	require.NoError(t, err)
	require.NotNil(t, foundOrders)
	require.Equal(t, len(foundOrders), 1)

	mock.ExpectQuery(findAllBySellerIdQuery).WithArgs(mockOtherOrder.SellerID, size, 0).WillReturnRows(otherRows)
	foundOrders, err = orderPGRepository.FindAllBySellerId(context.Background(), mockOtherOrder.SellerID, utils.NewPaginationQuery(size, 1))
	require.NoError(t, err)
	require.NotNil(t, foundOrders)
	require.Equal(t, len(foundOrders), 1)

	mock.ExpectQuery(findAllBySellerIdQuery).WithArgs(mockOtherOrder.SellerID, size, 10).WillReturnRows(otherRows)
	foundOrders, err = orderPGRepository.FindAllBySellerId(context.Background(), mockOtherOrder.SellerID, utils.NewPaginationQuery(size, 2))
	require.NoError(t, err)
	require.Nil(t, foundOrders)
}

func TestOrderRepository_FindAllByUserId(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	orderPGRepository := NewOrderPGRepository(sqlxDB)

	columns := []string{"order_id", "user_id", "seller_id", "item", "quantity", "total_price", "status", "delivery_source_address", "delivery_destination_address", "created_at", "updated_at"}
	orderUUID := uuid.New()
	userUUID := uuid.New()
	sellerUUID := uuid.New()
	otherUserUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID:  orderUUID,
		UserID:   otherUserUUID,
		SellerID: sellerUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			SellerID:    sellerUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	mockOtherOrder := &models.Order{
		OrderID:  orderUUID,
		UserID:   userUUID,
		SellerID: sellerUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			SellerID:    sellerUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	otherRows := sqlmock.NewRows(columns).AddRow(
		orderUUID,
		mockOtherOrder.UserID,
		mockOtherOrder.SellerID,
		mockOtherOrder.Item,
		mockOtherOrder.Quantity,
		mockOtherOrder.TotalPrice,
		mockOtherOrder.Status,
		mockOtherOrder.DeliverySourceAddress,
		mockOtherOrder.DeliveryDestinationAddress,
		time.Now(),
		time.Now(),
	)

	rows := sqlmock.NewRows(columns).AddRow(
		orderUUID,
		mockOrder.UserID,
		mockOrder.SellerID,
		mockOrder.Item,
		mockOrder.Quantity,
		mockOrder.TotalPrice,
		mockOrder.Status,
		mockOrder.DeliverySourceAddress,
		mockOrder.DeliveryDestinationAddress,
		time.Now(),
		time.Now(),
	)

	size := 10
	mock.ExpectQuery(findByUserIdQuery).WithArgs(mockOrder.UserID, size, 0).WillReturnRows(rows)
	foundOrders, err := orderPGRepository.FindAllByUserId(context.Background(), mockOrder.UserID, utils.NewPaginationQuery(size, 1))
	require.NoError(t, err)
	require.NotNil(t, foundOrders)
	require.Equal(t, len(foundOrders), 1)

	mock.ExpectQuery(findByUserIdQuery).WithArgs(mockOtherOrder.UserID, size, 0).WillReturnRows(otherRows)
	foundOrders, err = orderPGRepository.FindAllByUserId(context.Background(), mockOtherOrder.UserID, utils.NewPaginationQuery(size, 1))
	require.NoError(t, err)
	require.NotNil(t, foundOrders)
	require.Equal(t, len(foundOrders), 1)

	mock.ExpectQuery(findByUserIdQuery).WithArgs(mockOtherOrder.UserID, size, 10).WillReturnRows(otherRows)
	foundOrders, err = orderPGRepository.FindAllByUserId(context.Background(), mockOtherOrder.UserID, utils.NewPaginationQuery(size, 2))
	require.NoError(t, err)
	require.Nil(t, foundOrders)
}

func TestOrderRepository_FindAllByUserIdSellerId(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	orderPGRepository := NewOrderPGRepository(sqlxDB)

	columns := []string{"order_id", "user_id", "seller_id", "item", "quantity", "total_price", "status", "delivery_source_address", "delivery_destination_address", "created_at", "updated_at"}
	orderUUID := uuid.New()
	userUUID := uuid.New()
	sellerUUID := uuid.New()
	otherUserUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID:  orderUUID,
		UserID:   otherUserUUID,
		SellerID: sellerUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			SellerID:    sellerUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	mockOtherOrder := &models.Order{
		OrderID:  orderUUID,
		UserID:   userUUID,
		SellerID: sellerUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			SellerID:    sellerUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	otherRows := sqlmock.NewRows(columns).AddRow(
		orderUUID,
		mockOtherOrder.UserID,
		mockOtherOrder.SellerID,
		mockOtherOrder.Item,
		mockOtherOrder.Quantity,
		mockOtherOrder.TotalPrice,
		mockOtherOrder.Status,
		mockOtherOrder.DeliverySourceAddress,
		mockOtherOrder.DeliveryDestinationAddress,
		time.Now(),
		time.Now(),
	)

	rows := sqlmock.NewRows(columns).AddRow(
		orderUUID,
		mockOrder.UserID,
		mockOrder.SellerID,
		mockOrder.Item,
		mockOrder.Quantity,
		mockOrder.TotalPrice,
		mockOrder.Status,
		mockOrder.DeliverySourceAddress,
		mockOrder.DeliveryDestinationAddress,
		time.Now(),
		time.Now(),
	)

	size := 10
	mock.ExpectQuery(findAllByUserIdSellerIDQuery).WithArgs(mockOrder.UserID, mockOrder.SellerID, size, 0).WillReturnRows(rows)
	foundOrders, err := orderPGRepository.FindAllByUserIdSellerId(context.Background(), mockOrder.UserID, mockOrder.SellerID, utils.NewPaginationQuery(size, 1))
	require.NoError(t, err)
	require.NotNil(t, foundOrders)
	require.Equal(t, len(foundOrders), 1)

	mock.ExpectQuery(findAllByUserIdSellerIDQuery).WithArgs(mockOtherOrder.UserID, mockOtherOrder.SellerID, size, 0).WillReturnRows(otherRows)
	foundOrders, err = orderPGRepository.FindAllByUserIdSellerId(context.Background(), mockOtherOrder.UserID, mockOtherOrder.SellerID, utils.NewPaginationQuery(size, 1))
	require.NoError(t, err)
	require.NotNil(t, foundOrders)
	require.Equal(t, len(foundOrders), 1)

	mock.ExpectQuery(findAllByUserIdSellerIDQuery).WithArgs(mockOtherOrder.UserID, mockOtherOrder.SellerID, size, 10).WillReturnRows(otherRows)
	foundOrders, err = orderPGRepository.FindAllByUserIdSellerId(context.Background(), mockOtherOrder.UserID, mockOtherOrder.SellerID, utils.NewPaginationQuery(size, 2))
	require.NoError(t, err)
	require.Nil(t, foundOrders)
}

func TestOrderRepository_FindById(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	orderPGRepository := NewOrderPGRepository(sqlxDB)

	columns := []string{"order_id", "user_id", "seller_id", "item", "quantity", "total_price", "status", "delivery_source_address", "delivery_destination_address", "created_at", "updated_at"}
	orderUUID := uuid.New()
	userUUID := uuid.New()
	sellerUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID:  orderUUID,
		UserID:   userUUID,
		SellerID: sellerUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			SellerID:    sellerUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	rows := sqlmock.NewRows(columns).AddRow(
		orderUUID,
		mockOrder.UserID,
		mockOrder.SellerID,
		mockOrder.Item,
		mockOrder.Quantity,
		mockOrder.TotalPrice,
		mockOrder.Status,
		mockOrder.DeliverySourceAddress,
		mockOrder.DeliveryDestinationAddress,
		time.Now(),
		time.Now(),
	)

	mock.ExpectQuery(findByIdQuery).WithArgs(mockOrder.OrderID).WillReturnRows(rows)

	foundOrder, err := orderPGRepository.FindById(context.Background(), mockOrder.OrderID)
	require.NoError(t, err)
	require.NotNil(t, foundOrder)
	require.Equal(t, foundOrder.OrderID, mockOrder.OrderID)
}

func TestOrderRepository_UpdateById(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	orderPGRepository := NewOrderPGRepository(sqlxDB)

	columns := []string{"order_id", "user_id", "seller_id", "item", "quantity", "total_price", "status", "delivery_source_address", "delivery_destination_address", "created_at", "updated_at"}
	orderUUID := uuid.New()
	userUUID := uuid.New()
	sellerUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID:  orderUUID,
		UserID:   userUUID,
		SellerID: sellerUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			SellerID:    sellerUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	_ = sqlmock.NewRows(columns).AddRow(
		orderUUID,
		mockOrder.UserID,
		mockOrder.SellerID,
		mockOrder.Item,
		mockOrder.Quantity,
		mockOrder.TotalPrice,
		mockOrder.Status,
		mockOrder.DeliverySourceAddress,
		mockOrder.DeliveryDestinationAddress,
		time.Now(),
		time.Now(),
	)

	mockOrder.Status = models.OrderStatusAccepted
	mock.ExpectExec(updateByIdQuery).WithArgs(
		mockOrder.OrderID,
		mockOrder.UserID,
		mockOrder.SellerID,
		mockOrder.Item,
		mockOrder.Quantity,
		mockOrder.TotalPrice,
		mockOrder.Status,
		mockOrder.DeliverySourceAddress,
		mockOrder.DeliveryDestinationAddress,
	).WillReturnResult(sqlmock.NewResult(0, 1))

	updatedOrder, err := orderPGRepository.UpdateById(context.Background(), mockOrder)
	require.NoError(t, err)
	require.NotNil(t, mockOrder)
	require.Equal(t, updatedOrder.Status, mockOrder.Status)
	require.Equal(t, updatedOrder.OrderID, mockOrder.OrderID)
}

func TestOrderRepository_DeleteById(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	orderPGRepository := NewOrderPGRepository(sqlxDB)

	columns := []string{"order_id", "user_id", "seller_id", "item", "quantity", "total_price", "status", "delivery_source_address", "delivery_destination_address", "created_at", "updated_at"}
	orderUUID := uuid.New()
	userUUID := uuid.New()
	sellerUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID:  orderUUID,
		UserID:   userUUID,
		SellerID: sellerUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			SellerID:    sellerUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	_ = sqlmock.NewRows(columns).AddRow(
		orderUUID,
		mockOrder.UserID,
		mockOrder.SellerID,
		mockOrder.Item,
		mockOrder.Quantity,
		mockOrder.TotalPrice,
		mockOrder.Status,
		mockOrder.DeliverySourceAddress,
		mockOrder.DeliveryDestinationAddress,
		time.Now(),
		time.Now(),
	)

	mock.ExpectExec(deleteByIdQuery).WithArgs(mockOrder.OrderID).WillReturnResult(sqlmock.NewResult(0, 1))

	err = orderPGRepository.DeleteById(context.Background(), mockOrder.OrderID)
	require.NoError(t, err)
	require.NotNil(t, mockOrder)
}
