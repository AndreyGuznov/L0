package postgre

import "time"

func AttempingForConn(fn func() error, attempt int, delay time.Duration) (err error) {
	for attempt > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attempt--
			continue
		}
		return nil
	}
	return
}
