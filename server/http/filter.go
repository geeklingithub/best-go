package best_http

import (
	"fmt"
	"time"
)

type Filter func(reqBody any, c *NewContext)

type FilterHandle func(next Filter) Filter

func MetricFilterHandle(next Filter) Filter {
	f := func(reqBody any, c *NewContext) {
		// 执行前的时间
		startTime := time.Now().UnixNano()
		next(reqBody, c)
		// 执行后的时间
		endTime := time.Now().UnixNano()
		fmt.Printf("run time: %d \n", endTime-startTime)
	}
	return f
}

func (filter Filter) AddFilter(fs ...FilterHandle) Filter {
	for _, f := range fs {
		filter = f(filter)
	}

	return filter
}
