package tair

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestTairPipeline(t *testing.T) {
	opt := &redis.Options{
		Addr: "localhost:6379",
	}
	client := NewTairClient(opt)
	ctx := context.Background()
	pipe := client.TairPipeline()
	pipe.Set(ctx, "key", "value", 0)
	pipe.Get(ctx, "key")
	cmds, err := pipe.Exec(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "value", cmds[1].(*redis.StringCmd).Val())
}

func TestTairPipelined(t *testing.T) {
	opt := &redis.Options{
		Addr: "localhost:6379",
	}
	client := NewTairClient(opt)
	ctx := context.Background()

	cmds, err := client.TairPipelined(ctx, func(p redis.Pipeliner) error {
		p.Set(ctx, "key", "value", 0)
		p.Get(ctx, "key")
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "value", cmds[1].(*redis.StringCmd).Val())
}
