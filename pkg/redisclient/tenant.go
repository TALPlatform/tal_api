package redisclient

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/TALPlatform/tal_api/db"
)

func (r *RedisClient) DeleteAllTenants(ctx context.Context) error {
	pattern := "tenant--*"
	var cursor uint64
	var keys []string
	for {
		// Scan for keys matching the pattern
		var err error
		keys, cursor, err = r.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}
		// If there are keys to delete, delete them
		if len(keys) > 0 {
			_, err := r.client.Del(ctx, keys...).Result()
			if err != nil {
				return err
			}
		}
		// If the cursor is 0, we've finished scanning
		if cursor == 0 {
			break
		}
	}

	return nil
}
func (r *RedisClient) TenantDelete(ctx context.Context, tenantId int32) error {
	key := fmt.Sprintf("tenant--%d", tenantId)
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}
func (r *RedisClient) TenantDeleteBulk(ctx context.Context, tenantIds []int32) error {
	for _, tenantID := range tenantIds {
		err := r.TenantDelete(ctx, tenantID)
		if err != nil {
			return err
		}
	}
	return nil
}
func (r *RedisClient) TenantCreate(ctx context.Context, tenantId int32, tenant *db.TenantFindRow) (*db.TenantFindRow, error) {
	key := fmt.Sprintf("tenant--%d", tenantId)
	if tenant == nil {
		r.client.Del(ctx, key)
		return nil, nil
	}
	response, err := json.Marshal(tenant)
	if err != nil {
		return nil, err
	}
	err = r.client.Set(ctx, key, response, 0).Err()
	if err != nil {
		return nil, err
	}
	return tenant, nil
}

func (r *RedisClient) TenantFind(ctx context.Context, tenantId int32) (*db.TenantFindRow, error) {
	if r.isDisabled {
		return nil, nil
	}
	var parsedStruct db.TenantFindRow
	key := fmt.Sprintf("tenant--%d", tenantId)
	jsonBytes, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(jsonBytes, &parsedStruct); err != nil {
		return nil, err
	}

	return &parsedStruct, nil
}
