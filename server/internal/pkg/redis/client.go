package redis

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"

	"bbs-go/internal/pkg/config"
)

var (
	globalRedisClient *redis.Client
	ErrCacheNotFound  = errors.New("cache not found")
	ErrCacheMiss      = errors.New("cache miss")
	ErrCircuitOpen    = errors.New("circuit breaker is open")
	ErrLockTimeout    = errors.New("lock acquisition timeout")
	ErrSerializationFailed = errors.New("data serialization failed")
	ErrConnectionFailed    = errors.New("redis connection failed")
)

// Cache key 前缀
const (
	KeyPrefix           = "bbs:"
	UserCacheKey        = KeyPrefix + "user:%d"
	UserTokenCacheKey   = KeyPrefix + "token:%s"
	UserScoreRankKey    = KeyPrefix + "rank:score"
	UserCheckInRankKey  = KeyPrefix + "rank:checkin:%s"
	TopicCacheKey       = KeyPrefix + "topic:%d"
	ArticleCacheKey     = KeyPrefix + "article:%d"
	CommentCacheKey     = KeyPrefix + "comment:%d"
	SysConfigCacheKey   = KeyPrefix + "config:%s"
	TagCacheKey         = KeyPrefix + "tag:%d"
	TagAllCacheKey      = KeyPrefix + "tags:all"
	ForbiddenWordsKey   = KeyPrefix + "forbidden:words"
	ArticleTagCacheKey  = KeyPrefix + "article_tag:%d"

	// 缓存雪崩防护
	BloomFilterKey      = KeyPrefix + "bloom_filter"
	CacheBreakdownKey   = KeyPrefix + "breakdown:%s"
	
	// 分布式锁
	LockKeyPrefix       = KeyPrefix + "lock:"
)

// 缓存过期时间
const (
	UserCacheExpire       = 30 * time.Minute
	TokenCacheExpire      = 60 * time.Minute
	ConfigCacheExpire     = 24 * time.Hour
	TopicCacheExpire      = 15 * time.Minute
	RankCacheExpire       = 10 * time.Minute
	TagCacheExpire        = 60 * time.Minute
	ForbiddenWordsExpire  = 24 * time.Hour
	ArticleTagCacheExpire = 30 * time.Minute

	// 防缓存雪崩随机过期时间范围
	RandomExpireRange     = 5 * time.Minute
	// 防缓存击穿锁过期时间
	BreakdownLockExpire   = 30 * time.Second
)

// InitRedis 初始化Redis客户端
func InitRedis(cfg *config.RedisConfig) error {
	if !cfg.Enabled {
		slog.Info("Redis is disabled")
		return nil
	}

	globalRedisClient = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  time.Duration(cfg.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := globalRedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("redis connection failed: %w", err)
	}

	slog.Info("Redis client initialized successfully")
	return nil
}

// GetRedisClient 获取Redis客户端
func GetRedisClient() *redis.Client {
	return globalRedisClient
}

// IsRedisEnabled 检查Redis是否启用
func IsRedisEnabled() bool {
	return globalRedisClient != nil
}

// Cache 缓存接口
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, keys ...string) (int64, error)
	HGet(ctx context.Context, key, field string) (string, error)
	HSet(ctx context.Context, key string, values ...interface{}) error
	HDel(ctx context.Context, key string, fields ...string) error
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	ZAdd(ctx context.Context, key string, members ...redis.Z) error
	ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) ([]string, error)
	ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	Pipeline() redis.Pipeliner
	
	// JSON 操作
	GetJSON(ctx context.Context, key string, dest interface{}) error
	SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	
	// 防缓存击穿
	GetWithBreakdown(ctx context.Context, key string, expire time.Duration, fn func() (interface{}, error)) (interface{}, error)
	
	// 分布式锁 (安全版本)
	LockSafe(ctx context.Context, key string, expire time.Duration) (*lockInfo, error)
	SafeUnlock(ctx context.Context, lock *lockInfo) error
	
	// 分布式锁 (兼容版本)
	Lock(ctx context.Context, key string, expire time.Duration) (bool, error)
	Unlock(ctx context.Context, key string) error
}

// CircuitBreaker 熔断器状态
type CircuitBreaker struct {
	failureCount     int64        // 失败次数
	successCount     int64        // 成功次数  
	totalCount       int64        // 总请求次数
	lastFailureTime  time.Time    // 最后失败时间
	lastSuccessTime  time.Time    // 最后成功时间
	mutex            sync.RWMutex // 读写锁
	state            int32        // 0=关闭, 1=半开, 2=开启
	failureThreshold int64        // 失败阈值
	timeout          time.Duration // 超时时间
}

// RedisCache Redis缓存实现
type RedisCache struct {
	client         *redis.Client
	circuitBreaker *CircuitBreaker
}

// NewRedisCache 创建Redis缓存实例
func NewRedisCache() *RedisCache {
	if globalRedisClient == nil {
		return nil
	}
	return &RedisCache{
		client: globalRedisClient,
		circuitBreaker: &CircuitBreaker{
			failureThreshold: 5,               // 5次失败后开启熔断器
			timeout:          30 * time.Second, // 30秒后尝试恢复
		},
	}
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	var result *redis.StringCmd
	err := r.executeWithCircuitBreaker(func() error {
		result = r.client.Get(ctx, key)
		return result.Err()
	})
	
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrCacheNotFound
		}
		return "", err
	}
	
	return result.Val(), nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	// 添加随机过期时间防止缓存雪崩
	finalExpiration := r.addRandomExpire(expiration)
	
	return r.executeWithCircuitBreaker(func() error {
		return r.client.Set(ctx, key, value, finalExpiration).Err()
	})
}

func (r *RedisCache) Del(ctx context.Context, keys ...string) error {
	return r.executeWithCircuitBreaker(func() error {
		return r.client.Del(ctx, keys...).Err()
	})
}

func (r *RedisCache) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Exists(ctx, keys...).Result()
}

func (r *RedisCache) HGet(ctx context.Context, key, field string) (string, error) {
	result := r.client.HGet(ctx, key, field)
	if errors.Is(result.Err(), redis.Nil) {
		return "", ErrCacheNotFound
	}
	return result.Result()
}

func (r *RedisCache) HSet(ctx context.Context, key string, values ...interface{}) error {
	return r.client.HSet(ctx, key, values...).Err()
}

func (r *RedisCache) HDel(ctx context.Context, key string, fields ...string) error {
	return r.client.HDel(ctx, key, fields...).Err()
}

func (r *RedisCache) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.client.HGetAll(ctx, key).Result()
}

func (r *RedisCache) ZAdd(ctx context.Context, key string, members ...redis.Z) error {
	return r.client.ZAdd(ctx, key, members...).Err()
}

func (r *RedisCache) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) ([]string, error) {
	return r.client.ZRangeByScore(ctx, key, opt).Result()
}

func (r *RedisCache) ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.client.ZRevRange(ctx, key, start, stop).Result()
}

func (r *RedisCache) Pipeline() redis.Pipeliner {
	return r.client.Pipeline()
}

// GetJSON 获取JSON数据
func (r *RedisCache) GetJSON(ctx context.Context, key string, dest interface{}) error {
	val, err := r.Get(ctx, key)
	if err != nil {
		if errors.Is(err, ErrCacheNotFound) {
			slog.Debug("Cache key not found", "key", key)
		} else if errors.Is(err, ErrCircuitOpen) {
			slog.Warn("Circuit breaker open for cache read", "key", key)
		} else {
			slog.Error("Cache read error", "key", key, "error", err)
		}
		return err
	}
	
	// 检查是否是空值缓存
	if r.IsNullValue(val) {
		slog.Debug("Null value cache hit", "key", key)
		return ErrCacheNotFound
	}
	
	if err := json.Unmarshal([]byte(val), dest); err != nil {
		slog.Error("JSON unmarshal failed", "key", key, "value_length", len(val), "error", err)
		return fmt.Errorf("%w: %v", ErrSerializationFailed, err)
	}
	
	slog.Debug("Cache hit", "key", key, "value_length", len(val))
	return nil
}

// SetJSON 设置JSON数据
func (r *RedisCache) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		slog.Error("JSON marshal failed", "key", key, "error", err)
		return fmt.Errorf("%w: %v", ErrSerializationFailed, err)
	}
	
	if err := r.Set(ctx, key, data, expiration); err != nil {
		if errors.Is(err, ErrCircuitOpen) {
			slog.Warn("Circuit breaker open for cache write", "key", key, "data_length", len(data))
		} else {
			slog.Error("Cache write error", "key", key, "data_length", len(data), "expire", expiration, "error", err)
		}
		return err
	}
	
	slog.Debug("Cache set", "key", key, "data_length", len(data), "expire", expiration)
	return nil
}

// GetWithBreakdown 防缓存击穿获取数据 - 改进版本
func (r *RedisCache) GetWithBreakdown(ctx context.Context, key string, expire time.Duration, fn func() (interface{}, error)) (interface{}, error) {
	// 先尝试从缓存获取
	var result interface{}
	err := r.GetJSON(ctx, key, &result)
	if err == nil {
		return result, nil
	}
	
	if !errors.Is(err, ErrCacheNotFound) {
		// 缓存读取错误，回退到数据库
		slog.Warn("Cache read error, fallback to database", "key", key, "error", err)
		return fn()
	}
	
	// 缓存不存在，使用分布式锁防击穿
	lockKey := fmt.Sprintf("%s%s", CacheBreakdownKey, key)
	lock, err := r.LockSafe(ctx, lockKey, BreakdownLockExpire)
	if err != nil {
		slog.Warn("Failed to acquire lock, fallback to database", "key", key, "error", err)
		return fn()
	}
	
	if lock == nil {
		// 未获取到锁，使用指数退避重试缓存
		maxRetries := 3
		baseDelay := 10 * time.Millisecond
		
		for i := 0; i < maxRetries; i++ {
			delay := time.Duration(int64(baseDelay) * (1 << i)) // 指数退避: 10ms, 20ms, 40ms
			select {
			case <-time.After(delay):
				err := r.GetJSON(ctx, key, &result)
				if err == nil {
					return result, nil
				}
				// 如果到了最后一次重试，直接访问数据库
				if i == maxRetries-1 {
					return fn()
				}
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
		// 重试失败，回退到数据库
		return fn()
	}
	
	defer func() {
		if unlockErr := r.SafeUnlock(ctx, lock); unlockErr != nil {
			slog.Warn("Failed to unlock", "key", lockKey, "error", unlockErr)
		}
	}()
	
	// 再次检查缓存（双重检查）
	err = r.GetJSON(ctx, key, &result)
	if err == nil {
		return result, nil
	}
	
	// 执行数据加载函数
	data, err := fn()
	if err != nil {
		return nil, err
	}
	
	// 缓存数据（异步执行避免阻塞）
	if data != nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if setErr := r.SetJSON(ctx, key, data, expire); setErr != nil {
				slog.Warn("Failed to set cache", "key", key, "error", setErr)
			}
		}()
	} else {
		// 设置短期空值缓存防止缓存穿透
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			r.SetNullValue(ctx, key, 5*time.Minute)
		}()
	}
	
	return data, nil
}

// 锁信息结构
type lockInfo struct {
	key   string
	value string
}

// LockSafe 安全分布式锁 - 返回锁信息用于安全释放
func (r *RedisCache) LockSafe(ctx context.Context, key string, expire time.Duration) (*lockInfo, error) {
	lockKey := fmt.Sprintf("%s%s", LockKeyPrefix, key)
	value := r.generateLockValue()
	
	result := r.client.SetNX(ctx, lockKey, value, expire)
	locked, err := result.Result()
	if err != nil {
		return nil, err
	}
	
	if !locked {
		return nil, nil // 未获取到锁
	}
	
	return &lockInfo{
		key:   lockKey,
		value: value,
	}, nil
}

// Lock 兼容版本分布式锁 - 返回bool用于兼容性
func (r *RedisCache) Lock(ctx context.Context, key string, expire time.Duration) (bool, error) {
	lockKey := fmt.Sprintf("%s%s", LockKeyPrefix, key)
	value := r.generateLockValue()
	
	result := r.client.SetNX(ctx, lockKey, value, expire)
	return result.Result()
}

// Unlock 兼容版本释放锁 - 简单删除锁
func (r *RedisCache) Unlock(ctx context.Context, key string) error {
	lockKey := fmt.Sprintf("%s%s", LockKeyPrefix, key)
	return r.client.Del(ctx, lockKey).Err()
}

// SafeUnlock 安全释放分布式锁 - 只释放自己持有的锁
func (r *RedisCache) SafeUnlock(ctx context.Context, lock *lockInfo) error {
	if lock == nil {
		return nil
	}
	
	// Lua脚本确保只释放自己的锁
	luaScript := `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end
	`
	
	result := r.client.Eval(ctx, luaScript, []string{lock.key}, lock.value)
	return result.Err()
}


// addRandomExpire 添加随机过期时间防止缓存雪崩
func (r *RedisCache) addRandomExpire(expiration time.Duration) time.Duration {
	if expiration <= 0 {
		return expiration
	}
	
	// 添加0到RandomExpireRange的随机时间
	randomSecs := r.randomInt(int(RandomExpireRange.Seconds()))
	return expiration + time.Duration(randomSecs)*time.Second
}

// randomInt 生成随机整数
func (r *RedisCache) randomInt(max int) int {
	if max <= 0 {
		return 0
	}
	
	bytes := make([]byte, 4)
	rand.Read(bytes)
	
	// 简单的随机数生成
	num := int(bytes[0])<<24 | int(bytes[1])<<16 | int(bytes[2])<<8 | int(bytes[3])
	if num < 0 {
		num = -num
	}
	
	return num % max
}

// generateLockValue 生成锁值
func (r *RedisCache) generateLockValue() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// 空值缓存防穿透
const NullValue = "NULL"

// SetNullValue 设置空值缓存
func (r *RedisCache) SetNullValue(ctx context.Context, key string, expiration time.Duration) error {
	// 空值缓存时间较短，防止长期占用内存
	nullExpiration := time.Duration(math.Min(float64(expiration), float64(5*time.Minute)))
	return r.Set(ctx, key, NullValue, nullExpiration)
}

// IsNullValue 检查是否是空值
func (r *RedisCache) IsNullValue(value string) bool {
	return value == NullValue
}

// 熔断器方法
func (cb *CircuitBreaker) canExecute() bool {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	
	state := atomic.LoadInt32(&cb.state)
	
	switch state {
	case 0: // 关闭状态
		return true
	case 1: // 半开状态
		return time.Since(cb.lastFailureTime) >= cb.timeout
	case 2: // 开启状态
		return time.Since(cb.lastFailureTime) >= cb.timeout
	default:
		return false
	}
}

func (cb *CircuitBreaker) recordSuccess() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	
	atomic.AddInt64(&cb.successCount, 1)
	atomic.AddInt64(&cb.totalCount, 1)
	atomic.StoreInt64(&cb.failureCount, 0)
	cb.lastSuccessTime = time.Now()
	
	state := atomic.LoadInt32(&cb.state)
	if state != 0 {
		atomic.StoreInt32(&cb.state, 0) // 返回关闭状态
		slog.Info("Circuit breaker closed", "success_count", atomic.LoadInt64(&cb.successCount))
	}
}

func (cb *CircuitBreaker) recordFailure() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	
	cb.lastFailureTime = time.Now()
	failureCount := atomic.AddInt64(&cb.failureCount, 1)
	totalCount := atomic.AddInt64(&cb.totalCount, 1)
	
	if failureCount >= cb.failureThreshold {
		prevState := atomic.SwapInt32(&cb.state, 2) // 设置为开启状态
		if prevState != 2 {
			failureRate := float64(failureCount) / float64(totalCount) * 100
			slog.Warn("Circuit breaker opened", 
				"failure_count", failureCount, 
				"total_count", totalCount,
				"failure_rate", fmt.Sprintf("%.2f%%", failureRate),
				"threshold", cb.failureThreshold)
		}
	}
}

// executeWithCircuitBreaker 使用熔断器执行操作
func (r *RedisCache) executeWithCircuitBreaker(operation func() error) error {
	if r.circuitBreaker == nil {
		return operation()
	}
	
	if !r.circuitBreaker.canExecute() {
		return ErrCircuitOpen
	}
	
	err := operation()
	if err != nil {
		r.circuitBreaker.recordFailure()
		return err
	}
	
	r.circuitBreaker.recordSuccess()
	return nil
}

// GetCircuitBreakerStats 获取熔断器统计信息
func (r *RedisCache) GetCircuitBreakerStats() map[string]interface{} {
	if r.circuitBreaker == nil {
		return nil
	}
	
	cb := r.circuitBreaker
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	
	state := atomic.LoadInt32(&cb.state)
	failureCount := atomic.LoadInt64(&cb.failureCount)
	successCount := atomic.LoadInt64(&cb.successCount)
	totalCount := atomic.LoadInt64(&cb.totalCount)
	
	var stateStr string
	switch state {
	case 0:
		stateStr = "CLOSED"
	case 1:
		stateStr = "HALF_OPEN"
	case 2:
		stateStr = "OPEN"
	default:
		stateStr = "UNKNOWN"
	}
	
	var failureRate float64
	if totalCount > 0 {
		failureRate = float64(failureCount) / float64(totalCount) * 100
	}
	
	return map[string]interface{}{
		"state":              stateStr,
		"failure_count":      failureCount,
		"success_count":      successCount,
		"total_count":        totalCount,
		"failure_rate":       fmt.Sprintf("%.2f%%", failureRate),
		"failure_threshold":  cb.failureThreshold,
		"last_failure_time":  cb.lastFailureTime.Format("2006-01-02 15:04:05"),
		"last_success_time":  cb.lastSuccessTime.Format("2006-01-02 15:04:05"),
		"timeout_duration":   cb.timeout.String(),
	}
}