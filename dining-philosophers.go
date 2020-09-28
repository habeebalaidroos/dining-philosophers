package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var mutex = &sync.Mutex{}
	var wg sync.WaitGroup

	var forks [5]bool
	// forks is assigned "false" for each of the array element. "false" indicates that the fork is avaiable
	//Setting the initial status of all philosophers to "Initial"
	philosophers := [5]string{"Initial", "Initial", "Initial", "Initial", "Initial"}

	// Getting all the philosophers started
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go activatePhilosopher(i, mutex, &wg, &philosophers, &forks)
	}

	// Variable to check if all philosophers have finished thinking/eating cycle
	// Once they have all reached the status of "Finsihed", the program terminates
	allFinished := false
	for allFinished == false {
		// start by assuming all philosophers are finished with their thinking/eating cycle
		allFinished = true
		//Make a copy of the philosophers array so as to not have it corrupted by activatePhilosopher function
		myPhilosophers := philosophers
		for i := 0; i < 5; i++ {
			if myPhilosophers[i] != "Finished" {
				// Even one philosopher not having a status of "Finished" indicates all philosophers are not done with their thinking/eating cycles
				allFinished = false
			}
			// print each philosophers status
			fmt.Print(i, ":", myPhilosophers[i], "\t\t")
		}
		fmt.Println(" ")
		time.Sleep(1500 * time.Millisecond)
	}

	wg.Wait()

}

func activatePhilosopher(i int, mutex *sync.Mutex, wg *sync.WaitGroup, philosophers *[5]string, forks *[5]bool) {
	for j := 0; j < 3; j++ {

		// Start by putting the philosopher in "Thinking" mode
		philosophers[i] = "Thinking"
		time.Sleep(6 * time.Second)

		philosophers[i] = "Hungry"
		gotFork := false
		for gotFork == false {
			// Attemp mutex lock to see if forks available. This ensure only one philosopher is able to check the status of forks at any given time
			mutex.Lock()
			if forks[i] == false && forks[(i+4)%5] == false {
				gotFork = true
				// Set the left (i) and right ((i+4)%5) values to true so that others philosophers cannot use the same forks
				forks[i] = true
				forks[(i+4)%5] = true
			} else {
				// If you want to display the forks that a philosopher is waiting on, the code below helps
				left := ""
				right := ""
				if forks[i] == true {
					right = "Right busy"
				}
				if forks[(i+4)%5] == true {
					left = "Left busy"
				}
				//If you want to print which fork the philosopher is waiting on, uncomment the line below
				//philosophers[i] = "Hungry-" + left + " " + right
				//adding the line below to avoid "declared but not used" error. If you uncomment the line above, you dont need to add this line
				_, _ = left, right
			}
			mutex.Unlock()
			time.Sleep(2 * time.Second)
		}
		philosophers[i] = "Eating"
		time.Sleep(4 * time.Second)

		mutex.Lock()
		forks[i] = false
		forks[(i+4)%5] = false
		mutex.Unlock()

		philosophers[i] = "Done Eating"
		time.Sleep(2 * time.Second)
	}
	philosophers[i] = "Finished"
	wg.Done()
}
