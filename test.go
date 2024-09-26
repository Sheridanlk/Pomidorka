package main

import (
	"fmt"
	"os"
	"time"
)

func GigaTimer(DurOfBreak int, DurOfWork int, DurOfPomidorka int) {
	WorkTimer := time.NewTimer(time.Duration(DurOfWork) * time.Second)
	PomidorkaTimer := time.NewTimer(10 * time.Second)

	fmt.Println("Work")
	go Output(DurOfPomidorka)

	BreakTimer := time.NewTimer(time.Duration(DurOfBreak))
	BreakTimer.Stop()

	NumberOfSector := 1
	for {
		select {
		case <-PomidorkaTimer.C:
			buffer := DurOfBreak
			if NumberOfSector%4 == 0 {
				buffer = DurOfBreak * 5
			}
			fmt.Println("Break")
			BreakTimer.Reset(time.Duration(buffer) * time.Second)
			go Output(buffer)
		case <-BreakTimer.C:
			fmt.Println("Work")
			PomidorkaTimer.Reset(time.Duration(DurOfPomidorka) * time.Second)
			go Output(DurOfPomidorka)
		case <-WorkTimer.C:
			fmt.Println("Time is out")
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
	GigaTimer(BreakTimeDuration, WorkTimeDuration, int(DurOfPomidorka))

}
