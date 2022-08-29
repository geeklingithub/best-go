package pool

import "time"

// DefaultCleanIntervalTime 协程池默认配置
const (
	DefaultCleanIntervalTime = time.Second
)

// Options 配置项
type Options struct {
	WorkerExpiredDuration time.Duration // 工人过期间隔
	MaxTaskSize           int           // 最大任务数
}

type OptionFunc func(opts *Options)

func initOptions(options ...OptionFunc) (*Options, error) {
	opts := &Options{
		WorkerExpiredDuration: DefaultCleanIntervalTime,
		MaxTaskSize:           100000,
	}
	for _, option := range options {
		option(opts)
	}

	if expiry := opts.WorkerExpiredDuration; expiry < 0 {
		return nil, ErrInvalidWorkerExpiredDuration
	}

	return opts, nil
}
