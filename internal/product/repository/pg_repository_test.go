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

func TestProductRepository_Create(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	productPGRepository := NewProductPGRepository(sqlxDB)

	columns := []string{"product_id", "name", "description", "price", "seller_id", "created_at", "updated_at"}
	productUUID := uuid.New()
	sellerUUID := uuid.New()
	mockProduct := &models.Product{
		ProductID:   productUUID,
		Name:        "Name",
		Description: "Description",
		Price:       10000.00,
		SellerID:    sellerUUID,
	}

	rows := sqlmock.NewRows(columns).AddRow(
		productUUID,
		mockProduct.Name,
		mockProduct.Description,
		mockProduct.Price,
		mockProduct.SellerID,
		time.Now(),
		time.Now(),
	)

	mock.ExpectQuery(createProductQuery).WithArgs(
		mockProduct.Name,
		mockProduct.Description,
		mockProduct.Price,
		mockProduct.SellerID,
	).WillReturnRows(rows)

	createdProduct, err := productPGRepository.Create(context.Background(), mockProduct)
	require.NoError(t, err)
	require.NotNil(t, createdProduct)
}

func TestProductRepository_FindAll(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	productPGRepository := NewProductPGRepository(sqlxDB)

	columns := []string{"product_id", "name", "description", "price", "seller_id", "created_at", "updated_at"}
	productUUID := uuid.New()
	sellerUUID := uuid.New()
	mockProduct := &models.Product{
		ProductID:   productUUID,
		Name:        "Name",
		Description: "Description",
		Price:       10000.00,
		SellerID:    sellerUUID,
	}

	rows := sqlmock.NewRows(columns).AddRow(
		productUUID,
		mockProduct.Name,
		mockProduct.Description,
		mockProduct.Price,
		mockProduct.SellerID,
		time.Now(),
		time.Now(),
	)

	size := 10
	mock.ExpectQuery(findAllQuery).WithArgs(size, 0).WillReturnRows(rows)
	foundProducts, err := productPGRepository.FindAll(context.Background(), utils.NewPaginationQuery(size, 1))
	require.NoError(t, err)
	require.NotNil(t, foundProducts)
	require.Equal(t, len(foundProducts), 1)

	mock.ExpectQuery(findAllQuery).WithArgs(size, 10).WillReturnRows(rows)
	foundProducts, err = productPGRepository.FindAll(context.Background(), utils.NewPaginationQuery(size, 2))
	require.NoError(t, err)
	require.Nil(t, foundProducts)
}

func TestProductRepository_FindAllBySellerId(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	productPGRepository := NewProductPGRepository(sqlxDB)

	columns := []string{"product_id", "name", "description", "price", "seller_id", "created_at", "updated_at"}
	productUUID := uuid.New()
	sellerUUID := uuid.New()
	mockProduct := &models.Product{
		ProductID:   productUUID,
		Name:        "Name",
		Description: "Description",
		Price:       10000.00,
		SellerID:    sellerUUID,
	}

	rows := sqlmock.NewRows(columns).AddRow(
		productUUID,
		mockProduct.Name,
		mockProduct.Description,
		mockProduct.Price,
		mockProduct.SellerID,
		time.Now(),
		time.Now(),
	)

	otherUUID := productUUID
	otherRows := sqlmock.NewRows(columns).AddRow(
		productUUID,
		mockProduct.Name,
		mockProduct.Description,
		mockProduct.Price,
		otherUUID,
		time.Now(),
		time.Now(),
	)

	size := 10
	mock.ExpectQuery(findAllBySellerIdQuery).WithArgs(mockProduct.SellerID, size, 0).WillReturnRows(rows)
	foundProducts, err := productPGRepository.FindAllBySellerId(context.Background(), mockProduct.SellerID, utils.NewPaginationQuery(size, 1))
	require.NoError(t, err)
	require.NotNil(t, foundProducts)
	require.Equal(t, len(foundProducts), 1)

	mock.ExpectQuery(findAllBySellerIdQuery).WithArgs(mockProduct.SellerID, size, 10).WillReturnRows(rows)
	foundProducts, err = productPGRepository.FindAllBySellerId(context.Background(), mockProduct.SellerID, utils.NewPaginationQuery(size, 2))
	require.NoError(t, err)
	require.Nil(t, foundProducts)

	mock.ExpectQuery(findAllBySellerIdQuery).WithArgs(otherUUID, size, 0).WillReturnRows(otherRows)
	foundProducts, err = productPGRepository.FindAllBySellerId(context.Background(), otherUUID, utils.NewPaginationQuery(size, 1))
	require.NoError(t, err)
	require.NotNil(t, foundProducts)
	require.Equal(t, len(foundProducts), 1)

	mock.ExpectQuery(findAllBySellerIdQuery).WithArgs(otherUUID, size, 10).WillReturnRows(otherRows)
	foundProducts, err = productPGRepository.FindAllBySellerId(context.Background(), otherUUID, utils.NewPaginationQuery(size, 2))
	require.NoError(t, err)
	require.Nil(t, foundProducts)
}

func TestProductRepository_FindById(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	productPGRepository := NewProductPGRepository(sqlxDB)

	columns := []string{"product_id", "name", "description", "price", "seller_id", "created_at", "updated_at"}
	productUUID := uuid.New()
	sellerUUID := uuid.New()
	mockProduct := &models.Product{
		ProductID:   productUUID,
		Name:        "Name",
		Description: "Description",
		Price:       10000.00,
		SellerID:    sellerUUID,
	}

	rows := sqlmock.NewRows(columns).AddRow(
		productUUID,
		mockProduct.Name,
		mockProduct.Description,
		mockProduct.Price,
		mockProduct.SellerID,
		time.Now(),
		time.Now(),
	)

	mock.ExpectQuery(findByIdQuery).WithArgs(mockProduct.ProductID).WillReturnRows(rows)

	foundProduct, err := productPGRepository.FindById(context.Background(), mockProduct.ProductID)
	require.NoError(t, err)
	require.NotNil(t, foundProduct)
	require.Equal(t, foundProduct.ProductID, mockProduct.ProductID)
}

func TestProductRepository_UpdateById(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	productPGRepository := NewProductPGRepository(sqlxDB)

	columns := []string{"product_id", "name", "description", "price", "seller_id", "created_at", "updated_at"}
	productUUID := uuid.New()
	sellerUUID := uuid.New()
	mockProduct := &models.Product{
		ProductID:   productUUID,
		Name:        "Name",
		Description: "Description",
		Price:       10000.00,
		SellerID:    sellerUUID,
	}

	_ = sqlmock.NewRows(columns).AddRow(
		productUUID,
		mockProduct.Name,
		mockProduct.Description,
		mockProduct.Price,
		mockProduct.SellerID,
		time.Now(),
		time.Now(),
	)

	mockProduct.Name = "NameChanged"
	mock.ExpectExec(updateByIdQuery).WithArgs(
		mockProduct.ProductID,
		mockProduct.Name,
		mockProduct.Description,
		mockProduct.Price,
		mockProduct.SellerID,
	).WillReturnResult(sqlmock.NewResult(0, 1))

	updatedProduct, err := productPGRepository.UpdateById(context.Background(), mockProduct)
	require.NoError(t, err)
	require.NotNil(t, mockProduct)
	require.Equal(t, updatedProduct.Name, mockProduct.Name)
	require.Equal(t, updatedProduct.ProductID, mockProduct.ProductID)
}

func TestProductRepository_DeleteById(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	productPGRepository := NewProductPGRepository(sqlxDB)

	columns := []string{"product_id", "name", "description", "price", "seller_id", "created_at", "updated_at"}
	productUUID := uuid.New()
	sellerUUID := uuid.New()
	mockProduct := &models.Product{
		ProductID:   productUUID,
		Name:        "Name",
		Description: "Description",
		Price:       10000.00,
		SellerID:    sellerUUID,
	}

	_ = sqlmock.NewRows(columns).AddRow(
		productUUID,
		mockProduct.Name,
		mockProduct.Description,
		mockProduct.Price,
		mockProduct.SellerID,
		time.Now(),
		time.Now(),
	)

	mock.ExpectExec(deleteByIdQuery).WithArgs(mockProduct.ProductID).WillReturnResult(sqlmock.NewResult(0, 1))

	err = productPGRepository.DeleteById(context.Background(), mockProduct.ProductID)
	require.NoError(t, err)
	require.NotNil(t, mockProduct)
}
