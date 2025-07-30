package main

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type OTP struct {
	Key       string
	CreatedAt time.Time
}

type RetentionMap map[string]OTP

func NewRetentionMap(ctx context.Context, retentionPeriod time.Duration) RetentionMap {
	rm := make(RetentionMap)

	go rm.Retention(ctx, retentionPeriod)

	return rm
}

func (rm RetentionMap) NewOTP() OTP {
	o := OTP{
		Key:       uuid.NewString(),
		CreatedAt: time.Now(),
	}

	rm[o.Key] = o
	return o
}

func (rm RetentionMap) VerifyOTP(opt string) bool {
	if _, ok := rm[opt]; !ok {
		return false // otp not valid
	}

	delete(rm, opt)
	return true
}

func (rm RetentionMap) Retention(ctx context.Context, retentionPeriod time.Duration) {
	thicker := time.NewTicker(400 * time.Millisecond)

	for {
		select {
		case <-thicker.C:
			for _, otp := range rm {
				if otp.CreatedAt.Add(retentionPeriod).Before(time.Now()) {
					delete(rm, otp.Key)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
