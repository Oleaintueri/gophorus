package Network

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

func Ulimit() int64 {
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

func ScanPort(ip string, port int, timeout time.Duration, wg *sync.WaitGroup) (isOpen bool, err error) {
	defer wg.Done()
	defer println("Done with Scan..." + ip)
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			println("too many files open..")
			time.Sleep(timeout)
			return ScanPort(ip, port, timeout, wg)
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

func (semaphore *Semaphore) RunHelper(camera []Camera) (openCameras []string) {
	wg := sync.WaitGroup{}
	defer wg.Wait()
	outputChannel := make(chan Camera)

	for i := range camera {
		wg.Add(1)
		err := semaphore.lock.Acquire(context.Background(), 1)
		if err != nil {
			panic(err)
		}
		tmpCamera := camera[i]
		go func(tmpCamera Camera) {
			defer semaphore.lock.Release(1)
			fmt.Println("Testing port " + strconv.Itoa(tmpCamera.port) + " with IP " + tmpCamera.ip)
			isOpen, err := ScanPort(tmpCamera.ip, tmpCamera.port, tmpCamera.timeout, &wg)
			if err == nil {
				if isOpen == true {
					tmpCamera.isOpen = true
				} else {
					fmt.Println("Port of IP " + tmpCamera.ip + " closed")
				}
			} else {
				println(err.Error())
			}
		}(tmpCamera)
	}
	count := 0
	println("len: " + strconv.Itoa(len(outputChannel)))
	for i := range camera {
		val := camera[i]
		openCameras = append(openCameras, val.ip)
		fmt.Println("Port of IP "+val.ip+" open\n", "Calls: "+strconv.Itoa(count))
		count++
	}
	println("Ended ScanHelper")
	//close(outputChannel)
	return openCameras
}

func Run(ipRange string, port int) []string {
	cameras := parseIpRange(ipRange, port)
	s := &Semaphore{lock: semaphore.NewWeighted(Ulimit()),}
	opens := s.RunHelper(cameras)
	fmt.Println("Length of opens: " + strconv.Itoa(len(opens)))
	return opens
}

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
