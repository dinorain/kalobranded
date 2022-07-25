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

func SetupRedis() *orderRedisRepo {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatal(err)
	}
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	orderRedisRepository := NewOrderRedisRepo(client, nil)
	return orderRedisRepository
}

func TestOrderRedisRepo_SetOrderCtx(t *testing.T) {
	t.Parallel()

	redisRepo := SetupRedis()

	t.Run("SetOrderCtx", func(t *testing.T) {
		order := &models.Order{
			OrderID: uuid.New(),
		}

		err := redisRepo.SetOrderCtx(context.Background(), redisRepo.createKey(order.OrderID.String()), 10, order)
		require.NoError(t, err)
	})
}

func TestOrderRedisRepo_GetByIdCtx(t *testing.T) {
	t.Parallel()

	redisRepo := SetupRedis()

	t.Run("GetByIdCtx", func(t *testing.T) {
		order := &models.Order{
			OrderID: uuid.New(),
		}

		err := redisRepo.SetOrderCtx(context.Background(), redisRepo.createKey(order.OrderID.String()), 10, order)
		require.NoError(t, err)

		order, err = redisRepo.GetByIdCtx(context.Background(), redisRepo.createKey(order.OrderID.String()))
		require.NoError(t, err)
		require.NotNil(t, order)
	})
}

func TestOrderRedisRepo_DeleteOrderCtx(t *testing.T) {
	t.Parallel()

	redisRepo := SetupRedis()

	t.Run("DeleteOrderCtx", func(t *testing.T) {
		order := &models.Order{
			OrderID: uuid.New(),
		}

		err := redisRepo.DeleteOrderCtx(context.Background(), redisRepo.createKey(order.OrderID.String()))
		require.NoError(t, err)
	})
}
