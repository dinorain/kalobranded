package repository

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/dinorain/kalobranded/internal/models"
	"github.com/dinorain/kalobranded/pkg/utils"
)

func TestOrderRepository_Create(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	orderPGRepository := NewOrderPGRepository(sqlxDB)

	columns := []string{"order_id", "user_id", "brand_id", "item", "quantity", "total_price", "status", "delivery_source_address", "delivery_destination_address", "created_at", "updated_at"}
	orderUUID := uuid.New()
	userUUID := uuid.New()
	brandUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID: orderUUID,
		UserID:  userUUID,
		BrandID: brandUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			BrandID:     brandUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	valueJson, _ := json.Marshal(mockOrder.Item)

	rows := sqlmock.NewRows(columns).AddRow(
		orderUUID,
		mockOrder.UserID,
		mockOrder.BrandID,
		valueJson,
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
		mockOrder.BrandID,
		valueJson,
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

	columns := []string{"order_id", "user_id", "brand_id", "item", "quantity", "total_price", "status", "delivery_source_address", "delivery_destination_address", "created_at", "updated_at"}
	orderUUID := uuid.New()
	userUUID := uuid.New()
	brandUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID: orderUUID,
		UserID:  userUUID,
		BrandID: brandUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			BrandID:     brandUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	valueJson, _ := json.Marshal(mockOrder.Item)

	rows := sqlmock.NewRows(columns).AddRow(
		orderUUID,
		mockOrder.UserID,
		mockOrder.BrandID,
		valueJson,
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

func TestOrderRepository_FindAllByBrandId(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	orderPGRepository := NewOrderPGRepository(sqlxDB)

	columns := []string{"order_id", "user_id", "brand_id", "item", "quantity", "total_price", "status", "delivery_source_address", "delivery_destination_address", "created_at", "updated_at"}
	orderUUID := uuid.New()
	userUUID := uuid.New()
	brandUUID := uuid.New()
	otherBrandUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID: orderUUID,
		UserID:  userUUID,
		BrandID: otherBrandUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			BrandID:     otherBrandUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	mockOtherOrder := &models.Order{
		OrderID: orderUUID,
		UserID:  userUUID,
		BrandID: brandUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			BrandID:     brandUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	valueJson, _ := json.Marshal(mockOtherOrder.Item)

	otherRows := sqlmock.NewRows(columns).AddRow(
		orderUUID,
		mockOtherOrder.UserID,
		mockOtherOrder.BrandID,
		valueJson,
		mockOtherOrder.Quantity,
		mockOtherOrder.TotalPrice,
		mockOtherOrder.Status,
		mockOtherOrder.DeliverySourceAddress,
		mockOtherOrder.DeliveryDestinationAddress,
		time.Now(),
		time.Now(),
	)

	valueJson, _ = json.Marshal(mockOrder.Item)

	rows := sqlmock.NewRows(columns).AddRow(
		orderUUID,
		mockOrder.UserID,
		mockOrder.BrandID,
		valueJson,
		mockOrder.Quantity,
		mockOrder.TotalPrice,
		mockOrder.Status,
		mockOrder.DeliverySourceAddress,
		mockOrder.DeliveryDestinationAddress,
		time.Now(),
		time.Now(),
	)

	size := 10
	mock.ExpectQuery(findAllByBrandIdQuery).WithArgs(mockOrder.BrandID, size, 0).WillReturnRows(rows)
	foundOrders, err := orderPGRepository.FindAllByBrandId(context.Background(), mockOrder.BrandID, utils.NewPaginationQuery(size, 1))
	require.NoError(t, err)
	require.NotNil(t, foundOrders)
	require.Equal(t, len(foundOrders), 1)

	mock.ExpectQuery(findAllByBrandIdQuery).WithArgs(mockOtherOrder.BrandID, size, 0).WillReturnRows(otherRows)
	foundOrders, err = orderPGRepository.FindAllByBrandId(context.Background(), mockOtherOrder.BrandID, utils.NewPaginationQuery(size, 1))
	require.NoError(t, err)
	require.NotNil(t, foundOrders)
	require.Equal(t, len(foundOrders), 1)

	mock.ExpectQuery(findAllByBrandIdQuery).WithArgs(mockOtherOrder.BrandID, size, 10).WillReturnRows(otherRows)
	foundOrders, err = orderPGRepository.FindAllByBrandId(context.Background(), mockOtherOrder.BrandID, utils.NewPaginationQuery(size, 2))
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

	columns := []string{"order_id", "user_id", "brand_id", "item", "quantity", "total_price", "status", "delivery_source_address", "delivery_destination_address", "created_at", "updated_at"}
	orderUUID := uuid.New()
	userUUID := uuid.New()
	brandUUID := uuid.New()
	otherUserUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID: orderUUID,
		UserID:  otherUserUUID,
		BrandID: brandUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			BrandID:     brandUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	mockOtherOrder := &models.Order{
		OrderID: orderUUID,
		UserID:  userUUID,
		BrandID: brandUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			BrandID:     brandUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	valueJson, _ := json.Marshal(mockOtherOrder.Item)

	otherRows := sqlmock.NewRows(columns).AddRow(
		orderUUID,
		mockOtherOrder.UserID,
		mockOtherOrder.BrandID,
		valueJson,
		mockOtherOrder.Quantity,
		mockOtherOrder.TotalPrice,
		mockOtherOrder.Status,
		mockOtherOrder.DeliverySourceAddress,
		mockOtherOrder.DeliveryDestinationAddress,
		time.Now(),
		time.Now(),
	)

	valueJson, _ = json.Marshal(mockOrder.Item)

	rows := sqlmock.NewRows(columns).AddRow(
		orderUUID,
		mockOrder.UserID,
		mockOrder.BrandID,
		valueJson,
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

func TestOrderRepository_FindAllByUserIdBrandId(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	orderPGRepository := NewOrderPGRepository(sqlxDB)

	columns := []string{"order_id", "user_id", "brand_id", "item", "quantity", "total_price", "status", "delivery_source_address", "delivery_destination_address", "created_at", "updated_at"}
	orderUUID := uuid.New()
	userUUID := uuid.New()
	brandUUID := uuid.New()
	otherUserUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID: orderUUID,
		UserID:  otherUserUUID,
		BrandID: brandUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			BrandID:     brandUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	mockOtherOrder := &models.Order{
		OrderID: orderUUID,
		UserID:  userUUID,
		BrandID: brandUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			BrandID:     brandUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	valueJson, _ := json.Marshal(mockOtherOrder.Item)

	otherRows := sqlmock.NewRows(columns).AddRow(
		orderUUID,
		mockOtherOrder.UserID,
		mockOtherOrder.BrandID,
		valueJson,
		mockOtherOrder.Quantity,
		mockOtherOrder.TotalPrice,
		mockOtherOrder.Status,
		mockOtherOrder.DeliverySourceAddress,
		mockOtherOrder.DeliveryDestinationAddress,
		time.Now(),
		time.Now(),
	)

	valueJson, _ = json.Marshal(mockOrder.Item)

	rows := sqlmock.NewRows(columns).AddRow(
		orderUUID,
		mockOrder.UserID,
		mockOrder.BrandID,
		valueJson,
		mockOrder.Quantity,
		mockOrder.TotalPrice,
		mockOrder.Status,
		mockOrder.DeliverySourceAddress,
		mockOrder.DeliveryDestinationAddress,
		time.Now(),
		time.Now(),
	)

	size := 10
	mock.ExpectQuery(findAllByUserIdBrandIDQuery).WithArgs(mockOrder.UserID, mockOrder.BrandID, size, 0).WillReturnRows(rows)
	foundOrders, err := orderPGRepository.FindAllByUserIdBrandId(context.Background(), mockOrder.UserID, mockOrder.BrandID, utils.NewPaginationQuery(size, 1))
	require.NoError(t, err)
	require.NotNil(t, foundOrders)
	require.Equal(t, len(foundOrders), 1)

	mock.ExpectQuery(findAllByUserIdBrandIDQuery).WithArgs(mockOtherOrder.UserID, mockOtherOrder.BrandID, size, 0).WillReturnRows(otherRows)
	foundOrders, err = orderPGRepository.FindAllByUserIdBrandId(context.Background(), mockOtherOrder.UserID, mockOtherOrder.BrandID, utils.NewPaginationQuery(size, 1))
	require.NoError(t, err)
	require.NotNil(t, foundOrders)
	require.Equal(t, len(foundOrders), 1)

	mock.ExpectQuery(findAllByUserIdBrandIDQuery).WithArgs(mockOtherOrder.UserID, mockOtherOrder.BrandID, size, 10).WillReturnRows(otherRows)
	foundOrders, err = orderPGRepository.FindAllByUserIdBrandId(context.Background(), mockOtherOrder.UserID, mockOtherOrder.BrandID, utils.NewPaginationQuery(size, 2))
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

	columns := []string{"order_id", "user_id", "brand_id", "item", "quantity", "total_price", "status", "delivery_source_address", "delivery_destination_address", "created_at", "updated_at"}
	orderUUID := uuid.New()
	userUUID := uuid.New()
	brandUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID: orderUUID,
		UserID:  userUUID,
		BrandID: brandUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			BrandID:     brandUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	valueJson, _ := json.Marshal(mockOrder.Item)

	rows := sqlmock.NewRows(columns).AddRow(
		orderUUID,
		mockOrder.UserID,
		mockOrder.BrandID,
		valueJson,
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

	columns := []string{"order_id", "user_id", "brand_id", "item", "quantity", "total_price", "status", "delivery_source_address", "delivery_destination_address", "created_at", "updated_at"}
	orderUUID := uuid.New()
	userUUID := uuid.New()
	brandUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID: orderUUID,
		UserID:  userUUID,
		BrandID: brandUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			BrandID:     brandUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	valueJson, _ := json.Marshal(mockOrder.Item)

	_ = sqlmock.NewRows(columns).AddRow(
		orderUUID,
		mockOrder.UserID,
		mockOrder.BrandID,
		valueJson,
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
		mockOrder.BrandID,
		valueJson,
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

	columns := []string{"order_id", "user_id", "brand_id", "item", "quantity", "total_price", "status", "delivery_source_address", "delivery_destination_address", "created_at", "updated_at"}
	orderUUID := uuid.New()
	userUUID := uuid.New()
	brandUUID := uuid.New()
	productUUID := uuid.New()
	mockOrder := &models.Order{
		OrderID: orderUUID,
		UserID:  userUUID,
		BrandID: brandUUID,
		Item: models.OrderItem{
			ProductID:   productUUID,
			Name:        "Name",
			Description: "Description",
			Price:       10000.00,
			BrandID:     brandUUID,
		},
		Quantity:                   1,
		TotalPrice:                 10000.0,
		Status:                     models.OrderStatusPending,
		DeliverySourceAddress:      "DeliverySourceAddress",
		DeliveryDestinationAddress: "DeliveryDestinationAddress",
	}

	valueJson, _ := json.Marshal(mockOrder.Item)

	_ = sqlmock.NewRows(columns).AddRow(
		orderUUID,
		mockOrder.UserID,
		mockOrder.BrandID,
		valueJson,
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
