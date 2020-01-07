package processor

import (
	"net/http"

	config "../config"
	connection "../connection"
)

type CreateCountFunc func() (int64, error)

type Processor interface {
	ProcessRequest(method, request_url string, params map[string]string, body []byte, w http.ResponseWriter, r *http.Request) error
}

type BaseProcessor struct {
	RedisCli      *connection.RedisAdaptor
	Configure     *config.Configure
	Lru           *connection.LRU
	CountFunction CreateCountFunc
}

func CreateCounter(redis *connection.RedisAdaptor) CreateCountFunc {
	return func() (int64, error) {
		count, err := redis.NewShortUrlCount()
		if err != nil {
			return 0, err
		}
		return count, nil
	}
}
