package taskrunner

const (
	READY_TO_DISPATCH = "d"
	READY_TO_EXECUTE = "e"
	CLOSE = "c"

	VIDEO_PATH = "./videos/"
)

type controlChan chan string

// 需要下发的数据
type dataChan chan interface{}

// 分发和执行
type fn func(dc dataChan) error
