package main

import (
	"context"
	"fmt"
	"os"
	"time"
)

func GigaTimer(ctx context.Context, DurOfBreak int, DurOfWork int, DurOfPomidorka int) {
	intervals := make([]time.Duration, 0)
	numerOfIntervals := 0
	for DurOfWork > 0 {
		if (numerOfIntervals+1)%2 == 0 {
			if (numerOfIntervals+1)%8 == 0 {
				DurOfWork -= DurOfBreak * 5
				intervals = append(intervals, time.Duration(DurOfBreak*5))

			} else {
				DurOfWork -= DurOfBreak
				intervals = append(intervals, time.Duration(DurOfBreak))

			}
		} else {
			DurOfWork -= DurOfPomidorka
			intervals = append(intervals, time.Duration(DurOfPomidorka))

		}
		numerOfIntervals += 1
	}

	mainTimer := time.NewTimer(time.Duration(0) * time.Second)
	secondTicker := time.NewTicker(1 * time.Second)
	i := 0
	numberOfSecons := time.Duration(0)
	for {
		select {
		case <-mainTimer.C:
			if (i+1)%2 == 0 {
				fmt.Println("Break: ")
			} else {
				fmt.Println("Work: ")
			}
			mainTimer.Reset(intervals[i])
			numberOfSecons = time.Duration(0)
			i++
		case <-secondTicker.C:
			passed := numberOfSecons * time.Second
			fmt.Printf("%s / %s\n", passed, intervals[i-1])
			numberOfSecons++
		case <-ctx.Done():
			mainTimer.Stop()
			secondTicker.Stop()
			return
		}

	}

}

func findTimeDuration(StringTime string) time.Duration {
	time, err := time.ParseDuration(StringTime)
	if err != nil {
		fmt.Println("Error", err)
		return 0
	}
	return time

}

func main() {
	var WorkTime, BreakTime string
	fmt.Println("Set work time:")
	fmt.Fscan(os.Stdin, &WorkTime)

	fmt.Println("Set break time:")
	fmt.Fscan(os.Stdin, &BreakTime)
	const DurOfPomidorka time.Duration = 10 * time.Second
	ctx := context.Background()

	contectWithTimeout, cancelFunc := context.WithTimeout(ctx, findTimeDuration(WorkTime))
	defer cancelFunc()

	// signalCtx, cancelFunc := signal.NotifyContext(ctx, os.Interrupt)
	// defer cancelFunc()

	GigaTimer(contectWithTimeout, int(findTimeDuration(BreakTime)), int(findTimeDuration(WorkTime)), int(DurOfPomidorka))

}
