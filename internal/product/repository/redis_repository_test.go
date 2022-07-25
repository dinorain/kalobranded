package repository

import (
	"context"
	"log"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/dinorain/checkoutaja/internal/models"
)

func SetupRedis() *productRedisRepo {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatal(err)
	}
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	productRedisRepository := NewProductRedisRepo(client, nil)
	return productRedisRepository
}

func TestProductRedisRepo_SetProductCtx(t *testing.T) {
	t.Parallel()

	redisRepo := SetupRedis()

	t.Run("SetProductCtx", func(t *testing.T) {
		product := &models.Product{
			ProductID: uuid.New(),
		}

		err := redisRepo.SetProductCtx(context.Background(), redisRepo.createKey(product.ProductID.String()), 10, product)
		require.NoError(t, err)
	})
}

func TestProductRedisRepo_GetByIdCtx(t *testing.T) {
	t.Parallel()

	redisRepo := SetupRedis()

	t.Run("GetByIdCtx", func(t *testing.T) {
		product := &models.Product{
			ProductID: uuid.New(),
		}

		err := redisRepo.SetProductCtx(context.Background(), redisRepo.createKey(product.ProductID.String()), 10, product)
		require.NoError(t, err)

		product, err = redisRepo.GetByIdCtx(context.Background(), redisRepo.createKey(product.ProductID.String()))
		require.NoError(t, err)
		require.NotNil(t, product)
	})
}

func TestProductRedisRepo_DeleteProductCtx(t *testing.T) {
	t.Parallel()

	redisRepo := SetupRedis()

	t.Run("DeleteProductCtx", func(t *testing.T) {
		product := &models.Product{
			ProductID: uuid.New(),
		}

		err := redisRepo.DeleteProductCtx(context.Background(), redisRepo.createKey(product.ProductID.String()))
		require.NoError(t, err)
	})
}
