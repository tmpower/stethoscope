/*
 TODO
	1. max timeout is 1 minute
	2. calculate time passed processing when sleeping
	3. Store host in redis
	4. dynamic update of hosts. update hosts array without exiting loop
	5. while storing hosts in the DB just add a field named 'bool-new'. Then select rows whom 'new' field is true and add them to redis every minute.
		And finally make them false once we added them to redis. (look "database/sql" package)
	6.
*/

package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

func heartbeat(host string) {
	resp, err := http.Get(host)

	if err == nil && resp.StatusCode == 200 {
		fmt.Println(host, ": OK!")
	} else {
		fmt.Println(host, ": DOWN!!!")
	}
}

func main() {

	var hosts []string

	// 1. Read websites to be checked from a file
	hostsFile, err := os.Open("hosts.txt")

	if err != nil {
		fmt.Println(err)
	}
	defer hostsFile.Close()

	scanner := bufio.NewScanner(hostsFile)
	for scanner.Scan() {
		hosts = append(hosts, scanner.Text())
	}

	// 2. Loop through all websites and check each with goroutine (every 1 minute)
	for {
		fmt.Println("\n<=====================================================================================>\n")
		fmt.Println("\nChecking all hosts...\n")

		var wg sync.WaitGroup
		for _, host := range hosts {
			wg.Add(1)
			go func(host string) {
				heartbeat(host)
				wg.Done()
			}(host)
		}
		wg.Wait()
		fmt.Println("\nDone!!!")

		// sleep for 1 minute
		time.Sleep(30 * time.Second)
	}
}
