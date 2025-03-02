package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/eiannone/keyboard"
)

func IntervalsCalculation(DurOfBreak, DurOfWork, DurOfPomidorka int) []time.Duration {
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

	return intervals
}

func GigaTimer(ctx context.Context, intervals []time.Duration, isBreak bool, infChan chan<- []time.Duration) {
	if isBreak {
		fmt.Println("Break: ")
	} else {
		fmt.Println("Work: ")
	}
	isBreak = !isBreak
	interval_number := 0
	Timer := time.NewTimer(intervals[interval_number])
	Ticker := time.NewTicker(1 * time.Second)

	numberOfSecons := time.Duration(1)
	for {
		select {
		case <-Timer.C:
			interval_number++
			if interval_number >= len(intervals) {
				fmt.Println("Time is end")
				Ticker.Stop()
				infChan <- nil
				return
			}
			if isBreak {
				fmt.Println("Break: ")
			} else {
				fmt.Println("Work: ")
			}
			isBreak = !isBreak
			Timer.Reset(intervals[interval_number])
			numberOfSecons = time.Duration(1)

		case <-Ticker.C:
			passed := numberOfSecons * time.Second
			fmt.Printf("%s / %s\n", passed, intervals[interval_number])
			numberOfSecons++
		case <-ctx.Done():
			Timer.Stop()
			Ticker.Stop()
			intervals[interval_number] -= numberOfSecons * time.Second
			infChan <- intervals[interval_number:]
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

func KeyListener(cChan chan<- bool, errChan chan<- error) {
	Pause := false
	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			errChan <- err
			return
		}
		if key == keyboard.KeySpace {
			Pause = !Pause
			cChan <- Pause
		}
	}
}

func main() {
	var WorkTime, BreakTime string
	fmt.Println("Set work time:")
	fmt.Fscan(os.Stdin, &WorkTime)

	fmt.Println("Set break time:")
	fmt.Fscan(os.Stdin, &BreakTime)
	const DurOfPomidorka time.Duration = 10 * time.Second

	initalIntervals := IntervalsCalculation(int(findTimeDuration(BreakTime)), int(findTimeDuration(WorkTime)), int(DurOfPomidorka))
	fmt.Println(initalIntervals)

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	chChan := make(chan bool, 1)
	errChan := make(chan error, 1)
	infChan := make(chan []time.Duration, 1)
	var isBreak bool
	Pause := false
	flag := initalIntervals
	go KeyListener(chChan, errChan)

	for flag != nil {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		isBreak = (len(initalIntervals)-len(flag))%2 != 0

		if !Pause {
			go GigaTimer(ctx, flag, isBreak, infChan)
		}

		select {
		case Pause = <-chChan:
			if Pause {
				fmt.Println("Pause")
				cancel()
				flag = <-infChan
				fmt.Println(flag)
				if flag == nil {
					return
				}
			} else {
				fmt.Println("Play")
				continue
			}
		case err := <-errChan:
			fmt.Println(err)
			return
		}
	}

	// contectWithTimeout, cancelFunc := context.WithTimeout(ctx, findTimeDuration(WorkTime))
	// defer cancelFunc()

}
