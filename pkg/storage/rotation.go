package main

import (
	"os"
	"sync"
	"time"
)

// rotationpolicy defines the interface for log rotation polcies
type RotationPolicy interface {
	ShouldRotate(logFilePath string) bool
}

// rotates logs when they reach a specified size
type SizeBasedRotationPolicy struct {
	MaxSizeBytes int64
}

func NewSizeBasedRotationPolicy(maxSizeBytes int64) *SizeBasedRotationPolicy {
	return &SizeBasedRotationPolicy{MaxSizeBytes: maxSizeBytes}
}

func (p *SizeBasedRotationPolicy) ShouldRotate(logFilePath string) bool {
	info, err := os.Stat(logFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return info.Size() >= p.MaxSizeBytes
}

type TimeBasedRotationPolicy struct {
	IntervalSeconds  int64
	LastRotationTime time.Time
	Mu               sync.Mutex
}

func NewTimeBasedRotationPolicy(intervalSeconds int64) *TimeBasedRotationPolicy {
	return &TimeBasedRotationPolicy{
		IntervalSeconds:  intervalSeconds,
		LastRotationTime: time.Now(),
	}
}

func (p *TimeBasedRotationPolicy) ShouldRotate(logFilePath string) bool {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	currentTIme := time.Now()
	elapsed := currentTIme.Unix() - p.LastRotationTime.Unix()
	if elapsed >= p.IntervalSeconds {
		p.LastRotationTime = currentTIme
		return true
	}

	return false
}
