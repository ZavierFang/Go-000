package Week06

import (
	"sync"
	"time"
)

//滑动窗口，默认粒度为1s
type RollingWindow struct {
	Buckets map[int64]*numberBucket
	Mutex   *sync.RWMutex
	//滑动窗口长度
	Size int64
}

type numberBucket struct {
	Value float64
}

func NewRollingWindow(size int64) *RollingWindow {
	r := &RollingWindow{
		Buckets: make(map[int64]*numberBucket),
		Mutex:   &sync.RWMutex{},
		Size:    size,
	}

	return r
}

func (rw *RollingWindow) Incr(i float64) {
	if i == 0 {
		return
	}

	rw.Mutex.Lock()
	defer rw.Mutex.Unlock()

	bucket := rw.getCurrentBucket()
	bucket.Value += i
	rw.removeBuckets()
}

func (rw *RollingWindow) getCurrentBucket() *numberBucket {
	now := time.Now().Unix()
	var bucket *numberBucket
	var ok bool
	if bucket, ok = rw.Buckets[now]; !ok {
		bucket = &numberBucket{}
		rw.Buckets[now] = bucket
	}

	return bucket
}

func (rw *RollingWindow) removeBuckets() {
	expiredTime := time.Now().Unix() - rw.Size
	for timestamp := range rw.Buckets {
		if timestamp < expiredTime {
			delete(rw.Buckets, timestamp)
		}
	}
}

func (rw *RollingWindow) Sum(now time.Time) float64 {
	sum := float64(0)
	rw.Mutex.RLock()
	defer rw.Mutex.RUnlock()

	for k, v := range rw.Buckets {
		if k >= now.Unix() {
			sum += v.Value
		}
	}

	return sum
}

func (rw *RollingWindow) Avg(now time.Time) float64 {
	return rw.Sum(now) / float64(rw.Size)
}

func (rw *RollingWindow) Max(now time.Time) float64 {
	var max float64

	rw.Mutex.RLock()
	defer rw.Mutex.RUnlock()

	for i, bucket := range rw.Buckets {
		if i >= now.Unix()-rw.Size {
			if bucket.Value > max {
				max = bucket.Value
			}
		}
	}

	return max
}
