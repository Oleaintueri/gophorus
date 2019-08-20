package GoNetworkCameraScanner

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Semaphore struct {
	lock *semaphore.Weighted
}

type Camera struct {
	ip      string
	timeout time.Duration
	port    int
	isOpen  bool
}

/*
ulimit is to limit the amount of concurrent processes
*/
func ulimit() int64 {
	out, err := exec.Command("/bin/sh", "-c", "ulimit -n").Output()

	if err != nil {
		panic(err)
	}
	s := strings.TrimSpace(string(out))
	i, err := strconv.ParseInt(s, 10, 64)

	if err != nil {
		panic(err)
	}
	return i
}

/*
Scan the port of the ip
*/
func scanPort(ip string, port int, timeout time.Duration, wg *sync.WaitGroup) (isOpen bool, err error) {
	// wg will call Done at the end of the function's execution using defer
	defer wg.Done()
	defer println("Done with Scan..." + ip)

	// Check the port. If the connection throws an error, the port is closed.
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			println("too many files open..")
			time.Sleep(timeout)
			return scanPort(ip, port, timeout, wg)
		}
		return false, err

	}
	if conn != nil {
		err = conn.Close()
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

/*
Run helper executes using a Semaphore. It also is the parent of all the goroutines
*/
func (semaphore *Semaphore) runHelper(camera []Camera) (openCameras []string) {
	wg := sync.WaitGroup{}
	// Wait for all goroutines to finish

	// Loop through all the cameras and execute the scanPort function
	for i := range camera {
		// Add to WaitGroup
		wg.Add(1)
		// Lock the Semaphore
		err := semaphore.lock.Acquire(context.Background(), 1)
		if err != nil {
			panic(err)
		}
		//tmpCamera := &camera[i]
		// Execute an anonymous goroutine
		go func(tc int) {
			// Once anonymous function is done executing the semaphore will release
			defer semaphore.lock.Release(1)
			//fmt.Println("Testing port " + strconv.Itoa(tmpCamera.port) + " with IP " + tmpCamera.ip)
			isOpen, err := scanPort(camera[tc].ip, camera[tc].port, camera[tc].timeout, &wg)
			fmt.Println("Testing port " + strconv.Itoa(camera[tc].port) + " with IP " + camera[tc].ip + " STATUS: " + strconv.FormatBool(isOpen))
			if err == nil {
				if isOpen {
					camera[tc].isOpen = true
				}
			} else {
				println(err.Error())
			}
		}(i)
	}
	wg.Wait()
	count := 0

	for i := range camera {
		val := camera[i]
		if val.isOpen {
			openCameras = append(openCameras, val.ip)
			fmt.Println("Port of IP "+val.ip+" open\n", "Calls: "+strconv.Itoa(count))
		} else {
			fmt.Println("Port of IP "+val.ip+" closed\n", "Calls: "+strconv.Itoa(count))
		}

		count++
	}
	println("Ended ScanHelper")
	//close(outputChannel)
	return openCameras
}

/*
Start point
*/
func Run(ipRange string, port int) []string {
	cameras := parseIpRange(ipRange, port)
	s := &Semaphore{lock: semaphore.NewWeighted(ulimit()),}
	opens := s.runHelper(cameras)
	fmt.Println("Length of opens: " + strconv.Itoa(len(opens)))
	return opens
}

/*
Parse IP Range string eg. "192.168.1.100-192.168.1.200"
*/
func parseIpRange(ipRange string, port int) (cameraScanner []Camera) {
	// ipRange := "41.188.226.1-41.188.226.250"
	splitIPArr := strings.Split(ipRange, "-")

	start, end := splitIPArr[0], splitIPArr[1]
	startIPArr := strings.Split(start, ".")
	endIPArr := strings.Split(end, ".")

	startInt, err := strconv.ParseInt(startIPArr[len(startIPArr)-1:][0], 10, 32)
	if err != nil {
		panic(err)
	}
	endInt, err := strconv.ParseInt(endIPArr[len(endIPArr)-1:][0], 10, 32)
	if err != nil {
		panic(err)
	}

	baseIP := strings.Join(startIPArr[:len(startIPArr)-1], ".")
	for i := startInt; i <= endInt; i++ {
		cameraScanner = append(cameraScanner, Camera{
			ip:      baseIP + "." + strconv.FormatInt(int64(i), 10),
			port:    port,
			timeout: 3000 * time.Millisecond,
			isOpen:  false,
		})
	}
	return
}
