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

func TestSellerRepository_Create(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	sellerPGRepository := NewSellerPGRepository(sqlxDB)

	columns := []string{"seller_id", "first_name", "last_name", "email", "password", "avatar", "pickup_address", "created_at", "updated_at"}
	sellerUUID := uuid.New()
	mockSeller := &models.Seller{
		SellerID:      sellerUUID,
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Avatar:        nil,
		Password:      "123456",
		PickupAddress: "PickupAddress",
	}

	rows := sqlmock.NewRows(columns).AddRow(
		sellerUUID,
		mockSeller.FirstName,
		mockSeller.LastName,
		mockSeller.Email,
		mockSeller.Password,
		mockSeller.Avatar,
		mockSeller.PickupAddress,
		time.Now(),
		time.Now(),
	)

	mock.ExpectQuery(createSellerQuery).WithArgs(
		mockSeller.FirstName,
		mockSeller.LastName,
		mockSeller.Email,
		mockSeller.Password,
		mockSeller.Avatar,
		mockSeller.PickupAddress,
	).WillReturnRows(rows)

	createdSeller, err := sellerPGRepository.Create(context.Background(), mockSeller)
	require.NoError(t, err)
	require.NotNil(t, createdSeller)
}

func TestSellerRepository_FindByEmail(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	sellerPGRepository := NewSellerPGRepository(sqlxDB)

	columns := []string{"seller_id", "first_name", "last_name", "email", "password", "avatar", "pickup_address", "created_at", "updated_at"}
	sellerUUID := uuid.New()
	mockSeller := &models.Seller{
		SellerID:      sellerUUID,
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Avatar:        nil,
		Password:      "123456",
		PickupAddress: "PickupAddress",
	}

	rows := sqlmock.NewRows(columns).AddRow(
		sellerUUID,
		mockSeller.FirstName,
		mockSeller.LastName,
		mockSeller.Email,
		mockSeller.Password,
		mockSeller.Avatar,
		mockSeller.PickupAddress,
		time.Now(),
		time.Now(),
	)

	mock.ExpectQuery(findByEmailQuery).WithArgs(mockSeller.Email).WillReturnRows(rows)

	foundSeller, err := sellerPGRepository.FindByEmail(context.Background(), mockSeller.Email)
	require.NoError(t, err)
	require.NotNil(t, foundSeller)
	require.Equal(t, foundSeller.Email, mockSeller.Email)
}

func TestSellerRepository_FindAll(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	sellerPGRepository := NewSellerPGRepository(sqlxDB)

	columns := []string{"seller_id", "first_name", "last_name", "email", "password", "avatar", "pickup_address", "created_at", "updated_at"}
	sellerUUID := uuid.New()
	mockSeller := &models.Seller{
		SellerID:      sellerUUID,
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Avatar:        nil,
		Password:      "123456",
		PickupAddress: "PickupAddress",
	}

	rows := sqlmock.NewRows(columns).AddRow(
		sellerUUID,
		mockSeller.FirstName,
		mockSeller.LastName,
		mockSeller.Email,
		mockSeller.Password,
		mockSeller.Avatar,
		mockSeller.PickupAddress,
		time.Now(),
		time.Now(),
	)

	size := 10
	mock.ExpectQuery(findAllQuery).WithArgs(size, 0).WillReturnRows(rows)
	foundSellers, err := sellerPGRepository.FindAll(context.Background(), utils.NewPaginationQuery(size, 1))
	require.NoError(t, err)
	require.NotNil(t, foundSellers)
	require.Equal(t, len(foundSellers), 1)

	mock.ExpectQuery(findAllQuery).WithArgs(size, 10).WillReturnRows(rows)
	foundSellers, err = sellerPGRepository.FindAll(context.Background(), utils.NewPaginationQuery(size, 2))
	require.NoError(t, err)
	require.Nil(t, foundSellers)
}

func TestSellerRepository_FindById(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	sellerPGRepository := NewSellerPGRepository(sqlxDB)

	columns := []string{"seller_id", "first_name", "last_name", "email", "password", "avatar", "pickup_address", "created_at", "updated_at"}
	sellerUUID := uuid.New()
	mockSeller := &models.Seller{
		SellerID:      sellerUUID,
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Avatar:        nil,
		Password:      "123456",
		PickupAddress: "PickupAddress",
	}

	rows := sqlmock.NewRows(columns).AddRow(
		sellerUUID,
		mockSeller.FirstName,
		mockSeller.LastName,
		mockSeller.Email,
		mockSeller.Password,
		mockSeller.Avatar,
		mockSeller.PickupAddress,
		time.Now(),
		time.Now(),
	)

	mock.ExpectQuery(findByIdQuery).WithArgs(mockSeller.SellerID).WillReturnRows(rows)

	foundSeller, err := sellerPGRepository.FindById(context.Background(), mockSeller.SellerID)
	require.NoError(t, err)
	require.NotNil(t, foundSeller)
	require.Equal(t, foundSeller.SellerID, mockSeller.SellerID)
}

func TestSellerRepository_UpdateById(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	sellerPGRepository := NewSellerPGRepository(sqlxDB)

	columns := []string{"seller_id", "first_name", "last_name", "email", "password", "avatar", "pickup_address", "created_at", "updated_at"}
	sellerUUID := uuid.New()
	mockSeller := &models.Seller{
		SellerID:      sellerUUID,
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Avatar:        nil,
		Password:      "123456",
		PickupAddress: "PickupAddress",
	}

	_ = sqlmock.NewRows(columns).AddRow(
		sellerUUID,
		mockSeller.FirstName,
		mockSeller.LastName,
		mockSeller.Email,
		mockSeller.Password,
		mockSeller.Avatar,
		mockSeller.PickupAddress,
		time.Now(),
		time.Now(),
	)

	mockSeller.FirstName = "FirstNameChanged"
	mock.ExpectExec(updateByIdQuery).WithArgs(
		mockSeller.SellerID,
		mockSeller.FirstName,
		mockSeller.LastName,
		mockSeller.Email,
		mockSeller.Password,
		mockSeller.Avatar,
		mockSeller.PickupAddress,
	).WillReturnResult(sqlmock.NewResult(0, 1))

	updatedSeller, err := sellerPGRepository.UpdateById(context.Background(), mockSeller)
	require.NoError(t, err)
	require.NotNil(t, mockSeller)
	require.Equal(t, updatedSeller.FirstName, mockSeller.FirstName)
	require.Equal(t, updatedSeller.SellerID, mockSeller.SellerID)
}

func TestSellerRepository_DeleteById(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	sellerPGRepository := NewSellerPGRepository(sqlxDB)

	columns := []string{"seller_id", "first_name", "last_name", "email", "password", "avatar", "pickup_address", "created_at", "updated_at"}
	sellerUUID := uuid.New()
	mockSeller := &models.Seller{
		SellerID:      sellerUUID,
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Avatar:        nil,
		Password:      "123456",
		PickupAddress: "PickupAddress",
	}

	_ = sqlmock.NewRows(columns).AddRow(
		sellerUUID,
		mockSeller.FirstName,
		mockSeller.LastName,
		mockSeller.Email,
		mockSeller.Password,
		mockSeller.Avatar,
		mockSeller.PickupAddress,
		time.Now(),
		time.Now(),
	)

	mock.ExpectExec(deleteByIdQuery).WithArgs(mockSeller.SellerID).WillReturnResult(sqlmock.NewResult(0, 1))

	err = sellerPGRepository.DeleteById(context.Background(), mockSeller.SellerID)
	require.NoError(t, err)
	require.NotNil(t, mockSeller)
}
