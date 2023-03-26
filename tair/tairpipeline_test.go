package tair_test

import (
	"context"
	"testing"

	"github.com/alibaba/tair-go/tair"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PipelineTestSuite struct {
	suite.Suite
	tairClient        *tair.TairClient
	tairClusterClient *tair.TairClusterClient
}

func (suite *PipelineTestSuite) SetupTest() {
	suite.tairClient = tair.NewTairClient(redisOptions())
	assert.Equal(suite.T(), "OK", suite.tairClient.FlushDB(ctx).Val())
	suite.tairClusterClient = tair.NewTairClusterClient(&tair.TairClusterOptions{ClusterOptions: redisClusterOptions()})
	err := suite.tairClusterClient.ForEachMaster(ctx, func(ctx context.Context, master *redis.Client) error {
		return master.FlushDB(ctx).Err()
	})
	assert.NoError(suite.T(), err)
}

func (suite *PipelineTestSuite) TestTairPipeline(t *testing.T) {
	pipe := suite.tairClient.TairPipeline()
	pipe.Set(ctx, "key", "value", 0)
	pipe.Get(ctx, "key")
	cmds, err := pipe.Exec(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "value", cmds[1].(*redis.StringCmd).Val())
}

func (suite *PipelineTestSuite) TestTairPipelined(t *testing.T) {
	cmds, err := suite.tairClient.TairPipelined(ctx, func(p redis.Pipeliner) error {
		p.Set(ctx, "key", "value", 0)
		p.Get(ctx, "key")
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "value", cmds[1].(*redis.StringCmd).Val())
}

func (suite *PipelineTestSuite) TestTairClusterPipeline(t *testing.T) {
	pipe := suite.tairClusterClient.TairPipeline()
	pipe.Set(ctx, "key", "value", 0)
	pipe.Get(ctx, "key")
	cmds, err := pipe.Exec(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "value", cmds[1].(*redis.StringCmd).Val())
}

func (suite *PipelineTestSuite) TestTairClusterPipelined(t *testing.T) {
	cmds, err := suite.tairClusterClient.TairPipelined(ctx, func(p redis.Pipeliner) error {
		p.Set(ctx, "key", "value", 0)
		p.Get(ctx, "key")
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "value", cmds[1].(*redis.StringCmd).Val())
}

func TestTairPipelineTestSuite(t *testing.T) {
	suite.Run(t, new(PipelineTestSuite))
}
