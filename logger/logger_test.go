package logger_test

import (
	"context"
	"echoserver/logger"
	"sync"
	"testing"
	"time"
)

// 여러개의 고루틴에서 초당 2번의 로그 쓰기 작업을 할때,
// 초마다 로그 파일이 제대로 생성 & 쓰기 작업이 되는지 테스트
func TestLoggerSecondOption(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	Logger, err := logger.NewLogger(logger.Second, "./logs", "server")
	if err != nil {
		t.Error(err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(ctx context.Context) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}

				Logger.DebugF("DebugF : %s %d", "Test", 1)
				Logger.Debug("Debug : Test 2")
				Logger.InfoF("InfoF : %s %d", "Test", 3)
				Logger.Info("Info : Test 4")

				time.Sleep(time.Second / 2)
			}
		}(ctx)
	}
	wg.Wait()
}
