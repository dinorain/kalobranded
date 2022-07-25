package repository

import (
	"context"
	"log"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/dinorain/kalobranded/internal/models"
)

func SetupRedis() *brandRedisRepo {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatal(err)
	}
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	brandRedisRepository := NewBrandRedisRepo(client, nil)
	return brandRedisRepository
}

func TestBrandRedisRepo_SetBrandCtx(t *testing.T) {
	t.Parallel()

	redisRepo := SetupRedis()

	t.Run("SetBrandCtx", func(t *testing.T) {
		brand := &models.Brand{
			BrandID: uuid.New(),
		}

		err := redisRepo.SetBrandCtx(context.Background(), redisRepo.createKey(brand.BrandID.String()), 10, brand)
		require.NoError(t, err)
	})
}

func TestBrandRedisRepo_GetByIdCtx(t *testing.T) {
	t.Parallel()

	redisRepo := SetupRedis()

	t.Run("GetByIdCtx", func(t *testing.T) {
		brand := &models.Brand{
			BrandID: uuid.New(),
		}

		err := redisRepo.SetBrandCtx(context.Background(), redisRepo.createKey(brand.BrandID.String()), 10, brand)
		require.NoError(t, err)

		brand, err = redisRepo.GetByIdCtx(context.Background(), redisRepo.createKey(brand.BrandID.String()))
		require.NoError(t, err)
		require.NotNil(t, brand)
	})
}

func TestBrandRedisRepo_DeleteBrandCtx(t *testing.T) {
	t.Parallel()

	redisRepo := SetupRedis()

	t.Run("DeleteBrandCtx", func(t *testing.T) {
		brand := &models.Brand{
			BrandID: uuid.New(),
		}

		err := redisRepo.DeleteBrandCtx(context.Background(), redisRepo.createKey(brand.BrandID.String()))
		require.NoError(t, err)
	})
}
