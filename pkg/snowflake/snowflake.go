// Package snowflake 提供雪花 ID 生成器及其配置工厂。
package snowflake

import (
	"errors"
	"fmt"
	randv2 "math/rand/v2"
	"sync"
	"time"
)

var (
	// InvalidSystemClock 表示本地时钟回退，当前无法继续安全生成新 ID。
	InvalidSystemClock = errors.New("invalid system clock")
)

// NewFactory 创建一个按给定位宽和时间粒度生成雪花 ID 的工厂。
func NewFactory(bitLenWorkerID, bitLenSequence, gapMs uint, startTime time.Time) *Factory {
	return &Factory{
		bitLenWorkerID:  bitLenWorkerID,
		bitLenSequence:  bitLenSequence,
		bitLenTimestamp: 63 - bitLenSequence - bitLenWorkerID,
		startTime:       startTime,
		unit:            time.Duration(gapMs) * time.Millisecond,
	}
}

// Factory 保存雪花 ID 的位宽分配、起始时间和时间粒度配置。
type Factory struct {
	bitLenWorkerID, bitLenTimestamp, bitLenSequence uint
	startTime                                       time.Time
	unit                                            time.Duration
}

// MaskSequence 将序列号截断到当前工厂允许的最大位宽范围内。
func (f *Factory) MaskSequence(sequence uint32) uint32 {
	return sequence & f.MaxSequence()
}

// FlakeTimestamp 将时间转换为当前工厂时间粒度下的时间戳。
func (f *Factory) FlakeTimestamp(t time.Time) uint64 {
	return uint64(t.UnixNano() / int64(f.unit))
}

// CurrentElapsedTime 返回当前时间相对起始时间的已过时间片数。
func (f *Factory) CurrentElapsedTime() uint64 {
	return f.FlakeTimestamp(time.Now()) - f.FlakeTimestamp(f.startTime)
}

// SleepTime 计算当序列耗尽后，为进入下一个时间片需要等待多久。
func (f *Factory) SleepTime(overtime time.Duration) time.Duration {
	return overtime*f.unit - time.Duration(time.Now().UnixNano())%f.unit*time.Nanosecond
}

// BuildID 按 worker、时间片和序列号拼装出一个完整的雪花 ID。
func (f *Factory) BuildID(workerID uint32, elapsedTime uint64, sequence uint32) (uint64, error) {
	if elapsedTime >= 1<<f.bitLenTimestamp {
		return 0, errors.New("over the time limit")
	}
	return elapsedTime<<(f.bitLenSequence+f.bitLenWorkerID) | uint64(sequence)<<f.bitLenWorkerID | uint64(workerID), nil
}

// MaxSequence 返回当前工厂允许的最大序列号。
func (f *Factory) MaxSequence() uint32 {
	return 1<<f.bitLenSequence - 1
}

// MaxWorkerID 返回当前工厂允许的最大 worker ID。
func (f *Factory) MaxWorkerID() uint32 {
	return 1<<f.bitLenWorkerID - 1
}

// MaxTime 返回在当前位宽配置下可表示的最大时间点。
func (f *Factory) MaxTime() time.Time {
	maxTime := uint64(1<<f.bitLenTimestamp - 1)
	return time.Unix(int64(time.Duration(f.FlakeTimestamp(f.startTime)+maxTime)*f.unit/time.Second), 0)
}

// NewSnowflake 基于工厂配置创建一个绑定指定 worker ID 的生成器。
func (f *Factory) NewSnowflake(workerID uint32) (*Snowflake, error) {
	maxWorkerID := f.MaxWorkerID()
	if workerID > maxWorkerID {
		return nil, fmt.Errorf("worker id can't be large than %d", maxWorkerID)
	}
	return &Snowflake{f: f, workerID: workerID, syncMutex: &sync.Mutex{}}, nil
}

// NewSnowflake 使用默认雪花配置创建一个生成器。
func NewSnowflake(workerID uint32) (*Snowflake, error) {
	startTime, _ := time.Parse(time.RFC3339, "2010-11-04T01:42:54.657Z")
	return NewFactory(10, 12, 1, startTime).NewSnowflake(workerID)
}

// Snowflake 是线程安全的雪花 ID 生成器。
type Snowflake struct {
	f           *Factory
	workerID    uint32
	elapsedTime uint64
	sequence    uint32
	syncMutex   *sync.Mutex
}

// WorkerID 返回当前生成器绑定的 worker ID。
func (sf *Snowflake) WorkerID() uint32 {
	return sf.workerID
}

// ID 生成下一个雪花 ID；当检测到系统时钟回退时返回错误。
func (sf *Snowflake) ID() (uint64, error) {
	sf.syncMutex.Lock()
	defer sf.syncMutex.Unlock()

	currentElapsedTime := sf.f.CurrentElapsedTime()

	if sf.elapsedTime < currentElapsedTime {
		sf.elapsedTime = currentElapsedTime
		sf.sequence = generateRandomSequence(9)

		return sf.f.BuildID(sf.workerID, sf.elapsedTime, sf.sequence)
	}

	if sf.elapsedTime > currentElapsedTime {
		currentElapsedTime = sf.f.CurrentElapsedTime()
		if sf.elapsedTime > currentElapsedTime {
			return 0, InvalidSystemClock
		}
	}

	// ==

	sf.sequence = sf.f.MaskSequence(sf.sequence + 1)
	if sf.sequence == 0 {
		sf.elapsedTime = sf.elapsedTime + 1
		time.Sleep(sf.f.SleepTime(time.Duration(sf.elapsedTime - currentElapsedTime)))
	}

	return sf.f.BuildID(sf.workerID, sf.elapsedTime, sf.sequence)
}

func generateRandomSequence(n int32) uint32 {
	return uint32(randv2.New(&source{}).Int32N(n))
}

type source struct {
}

func (s source) Uint64() uint64 {
	return uint64(time.Now().UnixNano())
}
