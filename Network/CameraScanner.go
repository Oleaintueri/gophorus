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
}

func Ulimit() int64 {
	out, err := exec.Command("ulimit", "-n").Output()
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

func ScanPort(ip string, port int, timeout time.Duration) (isOpen bool, err error) {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			return ScanPort(ip, port, timeout)
		}
	}
	err = conn.Close()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (ps *Semaphore) RunHelper(camera []Camera) (openCameras []string) {
	wg := sync.WaitGroup{}
	defer wg.Wait()
	outputChannel := make(chan Camera)
	for i := range camera {
		wg.Add(1)
		ps.lock.Acquire(context.TODO(), 1)
		tmpCamera := camera[i]
		go func(tmpCamera Camera) {
			defer ps.lock.Release(1)
			defer wg.Done()
			isOpen, err := ScanPort(tmpCamera.ip, tmpCamera.port, tmpCamera.timeout)
			if err == nil {
				if isOpen == true {
					outputChannel <- tmpCamera
				}
			}
		}(tmpCamera)
	}
	for {
		val, ok := <-outputChannel
		if ok == false {
			break
		}
		openCameras = append(openCameras, val.ip)
	}
	return
}

func Run(ipRange string, port int) []string {
	cameras := parseIpRange(ipRange, port)
	s := &Semaphore{lock: semaphore.NewWeighted(Ulimit()),}
	return s.RunHelper(cameras)
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

	baseIP := strings.Join(startIPArr[:len(startIPArr)-1], "")
	for i := startInt; i < endInt; i++ {
		cameraScanner = append(cameraScanner, Camera{
			ip:      baseIP + string(i),
			port:    port,
			timeout: 3000 * time.Millisecond,
		})
	}
	return
}
