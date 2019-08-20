package GoNetworkCameraScanner

import (
	"fmt"
	"testing"
)

func TestRunt(t *testing.T) {
	got := Run("216.58.223.1-216.58.223.142", 80)
	fmt.Println("IP Addresses open: ")
	for i := range got {
		fmt.Println(got[i])
	}
}

func BenchmarkRun(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Run("216.58.223.1-216.58.223.142", 80)
	}
}
