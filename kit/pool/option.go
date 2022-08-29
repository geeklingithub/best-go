package pool

import "time"

// Option 配置项
type Option struct {
	WorkerExpiredDuration time.Duration // 工人过期间隔
	MaxTaskSize           int           // 最大任务数
}
