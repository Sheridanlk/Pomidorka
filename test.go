package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func GigaTimer(ctx context.Context, DurOfBreak int, DurOfWork int, DurOfPomidorka int) {


	mainTimer := time.Timer{}
	secondTicker := time.NewTicker(1 * time.Second)

	numerOfIntervals := 10
	intervals := make([]time.Duration, numerOfIntervals)

	i := 0
	numberOfSecons := time.Duration(0)
	for {
		select {
		case <-mainTimer.C:
			
			mainTimer.Reset(intervals[i])
			i++
			numberOfSecons = 0
		case <-secondTicker.C:
			passed := numberOfSecons*time.Second
			fmt.Printf("%s / %s\n", passed, intervals[i])
			numberOfSecons++
		case <-ctx.Done():
			mainTimer.Stop()
			secondTicker.Stop()
			return
		}

	}

}

func findTimeDuration(hrs, min, sec int) int {
	return hrs*3600 + min*60 + sec

}

func Output(TimeDuration int) {
	t := time.Tick(time.Second)
	for TimeDuration != 0 {
		fmt.Printf("%d : %d : %d \n", TimeDuration/3600, (TimeDuration-(TimeDuration/3600)*3600-TimeDuration%60)/60, TimeDuration%60)
		TimeDuration -= 1
		<-t
	}
}

func main() {
	var wrkmin, wrksec, wrkhrs int
	fmt.Println("WORK TIME:")
	fmt.Print("Type hours: ")
	fmt.Fscan(os.Stdin, &wrkhrs)
	fmt.Print("Type min: ")
	fmt.Fscan(os.Stdin, &wrkmin)
	fmt.Print("Type seconds: ")
	fmt.Fscan(os.Stdin, &wrksec)

	var brkmin, brksec, brkhrs int
	fmt.Println("BREAK TIME:")
	fmt.Print("Type hours: ")
	fmt.Fscan(os.Stdin, &brkhrs)
	fmt.Print("Type min: ")
	fmt.Fscan(os.Stdin, &brkmin)
	fmt.Print("Type seconds: ")
	fmt.Fscan(os.Stdin, &brksec)

	WorkTimeDuration := findTimeDuration(wrkhrs, wrkmin, wrksec)
	BreakTimeDuration := findTimeDuration(brkhrs, brkmin, brksec)

	const DurOfPomidorka int = 10
	ctx := context.Background()
//	contectWithTimeout, cancelFunc := context.WithTimeout(ctx, 5 * time.Second)
//	defer cancelFunc()

	// signal.Notify()
	signalCtx, cancelFunc := signal.NotifyContext(ctx, os.Interrupt)
	defer cancelFunc()

	GigaTimer(signalCtx, BreakTimeDuration, WorkTimeDuration, int(DurOfPomidorka))
}
