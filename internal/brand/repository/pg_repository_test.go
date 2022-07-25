package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/dinorain/kalobranded/internal/models"
	"github.com/dinorain/kalobranded/pkg/utils"
)

func TestBrandRepository_Create(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	brandPGRepository := NewBrandPGRepository(sqlxDB)

	columns := []string{"brand_id", "brand_name", "logo", "pickup_address", "created_at", "updated_at"}
	brandUUID := uuid.New()
	mockBrand := &models.Brand{
		BrandID:       brandUUID,
		BrandName:     "BrandName",
		Logo:          nil,
		PickupAddress: "PickupAddress",
	}

	rows := sqlmock.NewRows(columns).AddRow(
		brandUUID,
		mockBrand.BrandName,
		mockBrand.Logo,
		mockBrand.PickupAddress,
		time.Now(),
		time.Now(),
	)

	mock.ExpectQuery(createBrandQuery).WithArgs(
		mockBrand.BrandName,
		mockBrand.Logo,
		mockBrand.PickupAddress,
	).WillReturnRows(rows)

	createdBrand, err := brandPGRepository.Create(context.Background(), mockBrand)
	require.NoError(t, err)
	require.NotNil(t, createdBrand)
}

func TestBrandRepository_FindAll(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	brandPGRepository := NewBrandPGRepository(sqlxDB)

	columns := []string{"brand_id", "brand_name", "logo", "pickup_address", "created_at", "updated_at"}
	brandUUID := uuid.New()
	mockBrand := &models.Brand{
		BrandID:       brandUUID,
		BrandName:     "BrandName",
		Logo:          nil,
		PickupAddress: "PickupAddress",
	}

	rows := sqlmock.NewRows(columns).AddRow(
		brandUUID,
		mockBrand.BrandName,
		mockBrand.Logo,
		mockBrand.PickupAddress,
		time.Now(),
		time.Now(),
	)

	size := 10
	mock.ExpectQuery(findAllQuery).WithArgs(size, 0).WillReturnRows(rows)
	foundBrands, err := brandPGRepository.FindAll(context.Background(), utils.NewPaginationQuery(size, 1))
	require.NoError(t, err)
	require.NotNil(t, foundBrands)
	require.Equal(t, len(foundBrands), 1)

	mock.ExpectQuery(findAllQuery).WithArgs(size, 10).WillReturnRows(rows)
	foundBrands, err = brandPGRepository.FindAll(context.Background(), utils.NewPaginationQuery(size, 2))
	require.NoError(t, err)
	require.Nil(t, foundBrands)
}

func TestBrandRepository_FindById(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	brandPGRepository := NewBrandPGRepository(sqlxDB)

	columns := []string{"brand_id", "brand_name", "logo", "pickup_address", "created_at", "updated_at"}
	brandUUID := uuid.New()
	mockBrand := &models.Brand{
		BrandID:       brandUUID,
		BrandName:     "BrandName",
		Logo:          nil,
		PickupAddress: "PickupAddress",
	}

	rows := sqlmock.NewRows(columns).AddRow(
		brandUUID,
		mockBrand.BrandName,
		mockBrand.Logo,
		mockBrand.PickupAddress,
		time.Now(),
		time.Now(),
	)

	mock.ExpectQuery(findByIdQuery).WithArgs(mockBrand.BrandID).WillReturnRows(rows)

	foundBrand, err := brandPGRepository.FindById(context.Background(), mockBrand.BrandID)
	require.NoError(t, err)
	require.NotNil(t, foundBrand)
	require.Equal(t, foundBrand.BrandID, mockBrand.BrandID)
}

func TestBrandRepository_UpdateById(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	brandPGRepository := NewBrandPGRepository(sqlxDB)

	columns := []string{"brand_id", "brand_name", "logo", "pickup_address", "created_at", "updated_at"}
	brandUUID := uuid.New()
	mockBrand := &models.Brand{
		BrandID:       brandUUID,
		BrandName:     "BrandName",
		Logo:          nil,
		PickupAddress: "PickupAddress",
	}

	_ = sqlmock.NewRows(columns).AddRow(
		brandUUID,
		mockBrand.BrandName,
		mockBrand.Logo,
		mockBrand.PickupAddress,
		time.Now(),
		time.Now(),
	)

	mockBrand.BrandName = "BrandNameChanged"
	mock.ExpectExec(updateByIdQuery).WithArgs(
		mockBrand.BrandID,
		mockBrand.BrandName,
		mockBrand.Logo,
		mockBrand.PickupAddress,
	).WillReturnResult(sqlmock.NewResult(0, 1))

	updatedBrand, err := brandPGRepository.UpdateById(context.Background(), mockBrand)
	require.NoError(t, err)
	require.NotNil(t, mockBrand)
	require.Equal(t, updatedBrand.BrandName, mockBrand.BrandName)
	require.Equal(t, updatedBrand.BrandID, mockBrand.BrandID)
}

func TestBrandRepository_DeleteById(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	brandPGRepository := NewBrandPGRepository(sqlxDB)

	columns := []string{"brand_id", "brand_name", "logo", "pickup_address", "created_at", "updated_at"}
	brandUUID := uuid.New()
	mockBrand := &models.Brand{
		BrandID:       brandUUID,
		BrandName:     "BrandName",
		Logo:          nil,
		PickupAddress: "PickupAddress",
	}

	_ = sqlmock.NewRows(columns).AddRow(
		brandUUID,
		mockBrand.BrandName,
		mockBrand.Logo,
		mockBrand.PickupAddress,
		time.Now(),
		time.Now(),
	)

	mock.ExpectExec(deleteByIdQuery).WithArgs(mockBrand.BrandID).WillReturnResult(sqlmock.NewResult(0, 1))

	err = brandPGRepository.DeleteById(context.Background(), mockBrand.BrandID)
	require.NoError(t, err)
	require.NotNil(t, mockBrand)
}
