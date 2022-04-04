package db

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/shopspring/decimal"
)

var mycache *cache.Cache

func InitRedis() {

	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": ":6379",
		},
	})

	mycache = cache.New(&cache.Options{
		Redis:      ring,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

}
func SetCacheBalance(walletId int32, balance decimal.Decimal) {
	if err := mycache.Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   string(walletId),
		Value: balance,
	}); err != nil {
		//log.Fatalln(err)
		//panic(err)
	}
}

func GetCacheBalance(walletId int32) (balance decimal.Decimal, err error) {

	cacheError := mycache.Get(context.TODO(), string(walletId), &balance)
	if cacheError != nil {
		log.Infoln(cacheError)
		return balance, cacheError
	}
	return balance, nil
}
