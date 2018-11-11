package taskrunner

import (
	"errors"
	"log"
	"testing"
	"time"
)

func TestRunner(t *testing.T) {
	d := func(dc dataChan) error {
		for i := 0; i < 30; i++ {
			dc <- i
			log.Println("Dispatcher sented: ", i)
		}
		return nil
	}

	e := func(dc dataChan) error {
	forloop: // break标签：必须写在 要跳出循环的前面
		for {
			select {
			case dl := <-dc:
				log.Printf("Executor received: %v", dl)
			default:
				// break 标签： 直接跳出for循环, 不加标签只跳出 select
				break forloop
			}
		}
		return errors.New("Executor")
	}

	runner := NewRunner(30, false, d, e)
	go runner.startAll() //startAll 是死循环  需要goroutine
	time.Sleep(3 * time.Second)
}
