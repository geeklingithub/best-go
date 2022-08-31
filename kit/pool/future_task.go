package pool

import "context"

type FutureTask struct {
	taskFunc   func() any
	taskResult any
	ctx        context.Context
	cancel     context.CancelFunc
}
