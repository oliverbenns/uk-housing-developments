package location

import (
	"context"
	"encoding/json"
)

// Cache not strictly needed but in development this speeds it up
// as well as prevents racking up Google Map API cost.

func (c *Client) getFromCache(ctx context.Context, key string) (Geometry, error) {
	val, err := c.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return Geometry{}, err
	}

	geometry := Geometry{}
	err = json.Unmarshal([]byte(val), &geometry)
	if err != nil {
		return Geometry{}, err
	}

	return geometry, nil
}

func (c *Client) saveToCache(ctx context.Context, key string, geometry Geometry) error {
	value, err := json.Marshal(geometry)
	if err != nil {
		return err
	}

	return c.RedisClient.Set(ctx, key, string(value), 0).Err()
}
