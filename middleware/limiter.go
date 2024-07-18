package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tiant-developer/go-tiant/redis"
	"github.com/tiant-developer/go-tiant/zlog"
	"strconv"
	"time"
)

var limitRedis *redis.Redis

func LimitMiddleware(inputRedis *redis.Redis) gin.HandlerFunc {
	limitRedis = inputRedis
	InitLimiter()
	return func(ctx *gin.Context) {
		for !ApiLimiter.Allow(ctx) {
			zlog.Infof(ctx, "接口限流中，稍等")
			time.Sleep(100 * time.Millisecond)
		}
		ctx.Next()
	}
}

type Limiter struct {
	limit  int
	expire time.Duration
}

var ApiLimiter *Limiter

func InitLimiter() {
	ApiLimiter = NewLimiter(100, 60*time.Second)
}

func NewLimiter(limit int, expire time.Duration) *Limiter {
	return &Limiter{
		limit:  limit,
		expire: expire,
	}
}

func (limit *Limiter) Allow(ctx *gin.Context) bool {
	clientIP := ctx.ClientIP() //zSet的key
	urlPath := ctx.Request.URL.Path
	redisKey := clientIP + urlPath
	now := time.Now().UnixNano()
	// 移除时间戳小于(now - interval)的所有成员  (回收令牌)
	max := fmt.Sprintf("%d", now-int64(limit.expire))
	limitRedis.ZRemRangeByScore(redisKey, "0", max)
	// 获取当前集合长度，即令牌数
	size, _ := limitRedis.ZCard(redisKey)
	if size < int64(limit.limit) {
		// 如果令牌数小于最大请求数，则添加一个新的时间戳成为成员  (相当于往令牌桶中放令牌)
		limitRedis.ZAdd(redisKey, map[string]float64{strconv.FormatInt(now, 10): float64(now)})
		return true
	}
	// 否则拒绝该请求  (当达到limit放不了令牌时就阻塞)
	return false
}
