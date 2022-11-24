package id

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/huskyrobotdog/toolbox-go/cache"
	"github.com/huskyrobotdog/toolbox-go/inner"
)

func StartClusterAndGetNodeID(cacheKeyPrefix string) int64 {
	hostname, err := os.Hostname()
	if err != nil {
		inner.Fatal(err.Error())
	}
	id := int64(-1)
	for i := 0; i < 1024; i++ {
		ok, err := cache.Instance.SetNX(context.Background(), fmt.Sprintf("%s%d", cacheKeyPrefix, i), hostname, time.Minute).Result()
		if err != nil {
			inner.Fatal(err.Error())
		}
		if ok {
			id = int64(i)
			inner.Debug(fmt.Sprintf("find snowflake node id : %d", id))
			go startCluster(fmt.Sprintf("%s%d", cacheKeyPrefix, id))
			break
		}
	}
	if id == -1 {
		inner.Fatal("all ids are occupied")
	}
	return id
}

func startCluster(cacheKey string) {
	time.Sleep(13 * time.Second)
	for {
		inner.Debug(fmt.Sprintf("keeplive snowflake id : `%v`", cacheKey))
		ok, err := cache.Instance.Expire(context.Background(), cacheKey, time.Minute).Result()
		if err != nil {
			inner.Fatal(err.Error())
		}
		if !ok {
			inner.Fatal(fmt.Sprintf("`%s` is not exists", cacheKey))
		}
		time.Sleep(13 * time.Second)
	}
}
