package main

import (
	"fmt"
	"time"
)

// Retry выполняет операцию с повторными попытками и экспоненциальным отступлением
func Retry(attempts int, sleep time.Duration, function func() error) error {
	for i := 0; ; i++ {
		err := function()
		if err == nil {
			return nil
		}

		if i >= (attempts - 1) {
			return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
		}

		time.Sleep(sleep)
		sleep = sleep * 2 // Экспоненциальное увеличение времени ожидания
	}
}
