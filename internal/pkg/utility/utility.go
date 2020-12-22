package utility

import (
	"os/exec"
	"strconv"
	"strings"
)

/*
ulimit is to limit the amount of concurrent processes
*/
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