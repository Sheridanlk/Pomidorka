package main

import (
	"fmt"
	"os"
	"time"
)

func TimeSegment(DurOfPause int, Iteration int, MaxNumOfIterations int) uint {
	DurOfPomidorka := 10
	buffer := DurOfPause
	if Iteration%4 == 0 && Iteration != 0 {
		buffer = 5 * DurOfPause
	}
	PomidorkaTimer := time.NewTimer(time.Duration(DurOfPomidorka * int(time.Second)))
	PomidorkaTimerAndBreakTimer := time.NewTimer(time.Duration((buffer + DurOfPomidorka) * int(time.Second)))
	fmt.Println("Work")
	Output(DurOfPomidorka)

	<-PomidorkaTimer.C
	fmt.Println("Break")
	Output(buffer)
	<-PomidorkaTimerAndBreakTimer.C
	return TimeSegment(DurOfPause, Iteration+1, MaxNumOfIterations)

}

func findTimeDuration(hrs, min, sec int) int {
	return hrs*3600 + min*60 + sec

}

func Output(TimeDuration int) {
	tick := time.Tick(time.Second)
	for TimeDuration != 0 {
		<-tick
		fmt.Printf("%d : %d : %d \n", TimeDuration/3600, (TimeDuration-(TimeDuration/3600)*3600-TimeDuration%60)/60, TimeDuration%60)
		TimeDuration -= 1
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
	numberOfSector := WorkTimeDuration / (BreakTimeDuration + 10)

	timer1 := time.NewTimer(time.Duration(WorkTimeDuration * int(time.Second)))
	go TimeSegment(BreakTimeDuration, 0, numberOfSector)

	<-timer1.C
	fmt.Println("Time end")
	return

}
