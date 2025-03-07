package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/eiannone/keyboard"
)

func IntervalsCalculation(DurOfBreak, DurOfWork, DurOfPomidorka time.Duration) []time.Duration {
	intervals := make([]time.Duration, 0)
	numberOfInterval := 0
	for DurOfWork > 0 {
		if (numberOfInterval+1)%2 == 0 {
			if (numberOfInterval+1)%8 == 0 {
				DurOfWork -= DurOfBreak * 5
				intervals = append(intervals, DurOfBreak*5)

			} else {
				DurOfWork -= DurOfBreak
				intervals = append(intervals, DurOfBreak)

			}
		} else {
			DurOfWork -= DurOfPomidorka
			intervals = append(intervals, DurOfPomidorka)

		}
		numberOfInterval += 1
	}

	return intervals
}

func PomdoroTimer(ctx context.Context, intervals []time.Duration, brake bool) []time.Duration {
	pauseChan := make(chan struct{})
	stopChan := make(chan struct{})
	errChan := make(chan error)

	go KeyListener(pauseChan, stopChan, errChan)

	intervalNumber := 0
	numberOfSecons := time.Duration(0)

	if brake {
		fmt.Println("Break: ")
	} else {
		fmt.Println("Work: ")
	}

	timer := time.NewTimer(intervals[intervalNumber])
	ticker := time.NewTicker(1 * time.Second)

	pause := false

	for {
		select {
		case <-timer.C:
			intervalNumber++
			if intervalNumber >= len(intervals) {
				ticker.Stop()
				return nil
			}
			brake = !brake
			if brake {
				fmt.Println("Break: ")
			} else {
				fmt.Println("Work: ")
			}

			timer.Reset(intervals[intervalNumber])
			numberOfSecons = time.Duration(0)

		case <-ticker.C:
			passed := numberOfSecons * time.Second
			fmt.Printf("%s / %s\n", passed, intervals[intervalNumber])
			numberOfSecons++

		case <-pauseChan:
			if !pause {
				timer.Stop()
				ticker.Stop()
				intervals[intervalNumber] -= numberOfSecons * time.Second
				intervals = intervals[intervalNumber:]
				intervalNumber = 0
				numberOfSecons = time.Duration(0)
				fmt.Println("Pause")
				fmt.Printf("Time remaining:\n %s \n", intervals)

			} else {
				fmt.Println("Play")
				if brake {
					fmt.Println("Break: ")
				} else {
					fmt.Println("Work: ")
				}
				timer.Reset(intervals[intervalNumber])
				ticker.Reset(1 * time.Second)
			}
			pause = !pause

		case <-stopChan:
			timer.Stop()
			ticker.Stop()
			intervals[intervalNumber] -= numberOfSecons * time.Second
			return intervals[intervalNumber:]

		case err := <-errChan:
			fmt.Println("Error: ", err)
			return nil
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

func KeyListener(pauseChan chan<- struct{}, stopChan chan<- struct{}, errChan chan<- error) {
	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			errChan <- err
			return
		}
		if key == keyboard.KeySpace {
			pauseChan <- struct{}{}
		} else if key == keyboard.KeyEsc {
			stopChan <- struct{}{}
		}
	}
}

func main() {
	var WorkTime, BreakTime string
	fmt.Println("Set work time:")
	fmt.Fscan(os.Stdin, &WorkTime)

	fmt.Println("Set break time:")
	fmt.Fscan(os.Stdin, &BreakTime)
	const DurOfPomidorka time.Duration = 10 * time.Second // default 20m-25m, 10s for testing

	initalIntervals := IntervalsCalculation(findTimeDuration(BreakTime), findTimeDuration(WorkTime), DurOfPomidorka)
	fmt.Println(initalIntervals)

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	result := PomdoroTimer(ctx, initalIntervals, false)

	if result != nil {
		fmt.Println("Timer remaining")
		fmt.Println(result)
	} else {
		fmt.Println("Time is up")
	}
	// isBreak = (len(initalIntervals)-len(flag))%2 != 0

	// contectWithTimeout, cancelFunc := context.WithTimeout(ctx, findTimeDuration(WorkTime))
	// defer cancelFunc()

}
