package main

import (
	"fmt"
	"os"
	"time"
)

func GigaTimer(DurOfBreak, DurOfWork time.Duration, DurOfPomidorka int) {
	WorkTimer := time.NewTimer(time.Duration(DurOfWork) * time.Second)
	SegmentTimer := time.NewTimer(time.Duration(DurOfPomidorka) * time.Second)

	fmt.Println("Work")
	go Output(DurOfPomidorka)

	NumberOfSector := 1
	IsWorkTimer := true
	for {
		select {
		case <-SegmentTimer.C:
			if IsWorkTimer {
				fmt.Println("Break")
				if NumberOfSector%4 == 0 {
					SegmentTimer.Reset(time.Duration(DurOfBreak*5) * time.Second)
					go Output(DurOfBreak * 5)
				} else {
					SegmentTimer.Reset(time.Duration(DurOfBreak) * time.Second)
					go Output(DurOfBreak)
				}
				NumberOfSector += 1
				IsWorkTimer = false
			} else {
				fmt.Println("Work")
				SegmentTimer.Reset(time.Duration(DurOfPomidorka) * time.Second)
				go Output(DurOfPomidorka)
				IsWorkTimer = true
			}

		case <-WorkTimer.C:
			fmt.Println("Time is out")
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

func Output(TimeDuration int) {
	t := time.Tick(time.Second)
	for TimeDuration != 0 {
		fmt.Printf("%d : %d : %d \n", TimeDuration/3600, (TimeDuration-(TimeDuration/3600)*3600-TimeDuration%60)/60, TimeDuration%60)
		TimeDuration -= 1
		<-t
	}
}

func main() {
	var WorkTime, BreakTime string
	fmt.Println("Set work time:")
	fmt.Fscan(os.Stdin, &WorkTime)

	fmt.Println("Set break time:")
	fmt.Fscan(os.Stdin, &BreakTime)

	const DurOfPomidorka int = 10
	GigaTimer(int(findTimeDuration(WorkTime)), int(findTimeDuration(BreakTime)), int(DurOfPomidorka))

}
