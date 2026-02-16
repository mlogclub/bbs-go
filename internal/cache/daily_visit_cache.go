package cache

import (
	"encoding/binary"
	"sync"
	"time"

	"github.com/bits-and-blooms/bloom/v3"
)

// dailyVisitBloom 用布隆过滤器记录用户当日是否已发送过「登录态访问」事件。
// - 无假阴性：未标记时一定不会误判为已发送，不会漏发事件。
// - 有极低假阳性：可能误判为已发送而少发一次，对「每日登录」任务影响可接受。
// - 内存固定：按日轮换，每天一个固定大小的 filter，不随用户数增长。
const (
	dailyVisitCapacity      = 1000000 // 预期当日访问用户数上限
	dailyVisitFalsePositive = 0.0001  // 目标假阳性率 0.01%
)

var DailyVisitCache = &dailyVisitBloomWrapper{}

type dailyVisitBloomWrapper struct {
	mu     sync.Mutex
	day    string // "20060102"，用于按日轮换
	filter *bloom.BloomFilter
}

func (w *dailyVisitBloomWrapper) filterForToday() *bloom.BloomFilter {
	today := time.Now().In(time.Local).Format("20060102")
	if w.day != today {
		w.filter = bloom.NewWithEstimates(dailyVisitCapacity, dailyVisitFalsePositive)
		w.day = today
	}
	return w.filter
}

func userIdToBytes(userId int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(userId))
	return b
}

// IsSentToday 返回该用户今日是否已标记为已发送（存在即视为已发送）。可能有假阳性。
func (w *dailyVisitBloomWrapper) IsSentToday(userId int64) bool {
	if userId <= 0 {
		return true
	}
	w.mu.Lock()
	f := w.filterForToday()
	ok := f.Test(userIdToBytes(userId))
	w.mu.Unlock()
	return ok
}

// MarkSentToday 标记该用户今日已发送，用于登录接口发事件后与 trySendDailyVisitEvent 去重。
func (w *dailyVisitBloomWrapper) MarkSentToday(userId int64) {
	if userId <= 0 {
		return
	}
	w.mu.Lock()
	w.filterForToday().Add(userIdToBytes(userId))
	w.mu.Unlock()
}

// TryMarkAndReturnIfNew 若今日尚未标记则标记并返回 true（首次），否则返回 false。
// 用于「今日访问」事件去重：返回 true 时调用方发送事件。在 w.mu 下原子完成 Test+Add。
func (w *dailyVisitBloomWrapper) TryMarkAndReturnIfNew(userId int64) bool {
	if userId <= 0 {
		return false
	}
	w.mu.Lock()
	f := w.filterForToday()
	already := f.TestOrAdd(userIdToBytes(userId))
	w.mu.Unlock()
	return !already
}
