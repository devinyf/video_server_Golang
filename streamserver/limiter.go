/*
流控机制， 使用 channel 来确保并发访问量不会超过限制
流控机制通过 中间件 middleware 加入到 http中
*/
package main

import (
	"log"
)

type ConnLimiter struct {
	concurrentConn int
	bucket         chan int
}

// 仿构造函啊(初始化)
func NewConnLimiter(cc int) *ConnLimiter {
	return &ConnLimiter{
		concurrentConn: cc,
		bucket:         make(chan int, cc),
	}
}

// 返回true 表示拿到一个 token
func (cl *ConnLimiter) GetConn() bool {
	if len(cl.bucket) >= cl.concurrentConn {
		log.Printf("Reached the rate limitation. ")
		return false
	}
	// 如果 bucket 没满，随意写进去一个值
	cl.bucket <- 1
	return true
}

func (cl *ConnLimiter) ReleaseConn() {
	c := <-cl.bucket
	log.Printf("New connction coming: %d", c)
}

