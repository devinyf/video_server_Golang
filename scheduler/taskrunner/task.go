package taskrunner

import (
	"errors"
	"log"
	"myproject/video_server/scheduler/dbops"
	"os"
	"sync"
	
)

func deleteVideo(vid string) error {
	// 删除服务器中的视频数据
	err := os.Remove(VIDEO_PATH + vid)
	// os.IsNotExist(err) 说明路径不存在--已经删除过了
	if err != nil && !os.IsNotExist(err) { // !os.IsNotExist(err) 说明没有被删除掉
		log.Printf("Deleting video error: %v", err)
		return err
	}

	return nil
}

func VideoClearDispatcher(dc dataChan) error {
	res, err := dbops.ReadVideoDeletionRecord(3)
	if err != nil {
		log.Printf("Video clear dispatcher error: %v", err)
		return err
	}

	if len(res) == 0 {
		return errors.New("All tasks finished")
	}

	for _, id := range res {
		dc <- id
	}

	return nil
}

func VideoClearExecutor(dc dataChan) error {
	errMap := sync.Map{}
	var err error
forloop:
	for {
		select {
		// 接收到分发任务: 删除视频
		case vid := <-dc:
			go func(id interface{}) {
				// 删除视频---- id.(string): 把interface{} 转为string
				if err := deleteVideo(id.(string)); err != nil {
					// 把error带出去做进一步处理
					errMap.Store(id, err)
					return
				}
				// 删除数据库中的视频名
				if err := dbops.DelVideoDeletionRecord(id.(string)); err != nil {
					errMap.Store(id, err)
					return
				}
			}(vid)
		default:
			break forloop

		}
	}
	// 遍历errMap
	errMap.Range(func(k, v interface{}) bool {
		err = v.(error)
		if err != nil {
			return false
		}
		return true
	})
	return err
}
