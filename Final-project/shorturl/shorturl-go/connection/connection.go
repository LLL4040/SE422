package redis

// import (
// 	"github.com/go-redis/redis"
// )

// var client = redis.NewClient(&redis.Options{
// 	Addr:     "192.168.99.100:6379",
// 	Password: "redis-shorturl",
// 	DB:       0,
// })

// func ReadRedis(key string) (string, bool) {
// 	val, err := client.Get(key).Result()
// 	if err != nil {
// 		return val, false
// 	}
// 	return val, true
// }

// func SetRedis(key string, value string) bool {
// 	err := client.Set(key, value, 0).Err()
// 	if err != nil {
// 		return false
// 	}
// 	return true
// }
