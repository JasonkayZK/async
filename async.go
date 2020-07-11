package async

import (
	"fmt"
	"reflect"
	"time"
)

func decorateCallbackWithAttempt(decoPtr, fn interface{}, retryInterval int64, attempts int64) (err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			err = fmt.Errorf("decorator err: %v", err1)
		}
	}()

	var decoratedFunc, targetFunc reflect.Value

	decoratedFunc = reflect.ValueOf(decoPtr).Elem()
	targetFunc = reflect.ValueOf(fn)

	v := reflect.MakeFunc(targetFunc.Type(),
		func(in []reflect.Value) (out []reflect.Value) {
			for retry := attempts; retry > 0; retry-- {
				hasErr := false
				// Call callback func
				out = targetFunc.Call(in)

				// Has return val
				if valuesNum := len(out); valuesNum > 0 {
					resultItems := make([]interface{}, valuesNum)
					// Check value

					for k, val := range out {
						resultItems[k] = val.Interface()
						// Has error
						if _, ok := resultItems[k].(error); ok {
							hasErr = true
							break
						}
					}

					// Has err, retry
					if hasErr {
						time.Sleep(time.Duration(retryInterval) * time.Second)
						fmt.Printf("retry %d\n", retry)
						continue
					}
					return
				}
			}
			return out
		})

	decoratedFunc.Set(v)
	return
}

func ExecuteCallbackWithAttempt(fn interface{}, retryInterval int64, attempts int64, params ...interface{}) (err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			err = fmt.Errorf("decorator err: %v", err1)
		}
	}()

	decoPtr := fn
	err = decorateCallbackWithAttempt(&decoPtr, fn, retryInterval, attempts)
	if err != nil {
		return err
	}

	paramNum := len(params)
	callParams := make([]reflect.Value, paramNum)
	if paramNum > 0 {
		for k, v := range params {
			callParams[k] = reflect.ValueOf(v)
		}
	}

	go reflect.ValueOf(decoPtr).Call(callParams)

	return nil
}
