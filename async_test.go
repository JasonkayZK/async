package async

import (
	"fmt"
	"testing"
	"time"
)

func TestAsyncCallbackWithAttempt(t *testing.T) {
	_ = ExecuteCallbackWithAttempt(errorCallbackTest, 1, 3, "good!")
	_ = ExecuteCallbackWithAttempt(errorCallbackTest, 1, 3, "error")
	_ = ExecuteCallbackWithAttempt(errorCallbackTest, 1, 3, "panic")

	time.Sleep(time.Second * 5)
}

func errorCallbackTest(str string) (res string, err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			err = fmt.Errorf("callback err: %v", err1)
		}
	}()

	if str == "error" {
		return "", fmt.Errorf("mock err: %s", str)
	}
	if str == "panic" {
		panic(fmt.Errorf("panic: %s", str))
	}

	fmt.Printf("Test msg: %s\n", str)
	return fmt.Sprintf("Test msg: %s", str), nil
}
