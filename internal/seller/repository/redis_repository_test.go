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

func SetupRedis() *sellerRedisRepo {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatal(err)
	}
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	sellerRedisRepository := NewSellerRedisRepo(client, nil)
	return sellerRedisRepository
}

func TestSellerRedisRepo_SetSellerCtx(t *testing.T) {
	t.Parallel()

	redisRepo := SetupRedis()

	t.Run("SetSellerCtx", func(t *testing.T) {
		seller := &models.Seller{
			SellerID: uuid.New(),
		}

		err := redisRepo.SetSellerCtx(context.Background(), redisRepo.createKey(seller.SellerID.String()), 10, seller)
		require.NoError(t, err)
	})
}

func TestSellerRedisRepo_GetByIdCtx(t *testing.T) {
	t.Parallel()

	redisRepo := SetupRedis()

	t.Run("GetByIdCtx", func(t *testing.T) {
		seller := &models.Seller{
			SellerID: uuid.New(),
		}

		err := redisRepo.SetSellerCtx(context.Background(), redisRepo.createKey(seller.SellerID.String()), 10, seller)
		require.NoError(t, err)

		seller, err = redisRepo.GetByIdCtx(context.Background(), redisRepo.createKey(seller.SellerID.String()))
		require.NoError(t, err)
		require.NotNil(t, seller)
	})
}

func TestSellerRedisRepo_DeleteSellerCtx(t *testing.T) {
	t.Parallel()

	redisRepo := SetupRedis()

	t.Run("DeleteSellerCtx", func(t *testing.T) {
		seller := &models.Seller{
			SellerID: uuid.New(),
		}

		err := redisRepo.DeleteSellerCtx(context.Background(), redisRepo.createKey(seller.SellerID.String()))
		require.NoError(t, err)
	})
}
