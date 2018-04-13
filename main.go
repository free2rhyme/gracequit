package main

import (
	"sync"
)

func main() {

	handler := GetHandler()

	if err := handler.Start(); nil != err {
		return
	}

	waitGroup := &sync.WaitGroup{}

	waitGroup.Add(1)
	go func() {
		for i := 0; i < 10; i++ {
			handler.OnEvent1(i)
			handler.OnEvent2(i * 10)
			handler.OnEvent3(i * 100)

			if i%2 == 4 {
				handler.Stop()
			}
		}
		waitGroup.Done()
	}()

	waitGroup.Add(1)
	go func() {
		for i := 0; i < 10; i++ {
			handler.OnEvent1(i)
			handler.OnEvent2(i * 10)
			handler.OnEvent3(i * 100)

			if i%2 == 2 {
				handler.Stop()
			}
		}
		waitGroup.Done()
	}()

	waitGroup.Add(1)
	go func() {
		for i := 0; i < 10; i++ {
			handler.OnEvent1(i)
			handler.OnEvent2(i * 10)
			handler.OnEvent3(i * 100)

			if i%2 == 0 {
				handler.Stop()
			}
		}
		waitGroup.Done()
	}()

	for i := 0; i < 5; i++ {
		handler.OnEvent1(i)
		handler.OnEvent2(i * 10)
		handler.OnEvent3(i * 100)
	}

	handler.Stop()

	waitGroup.Wait()
}
