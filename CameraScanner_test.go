package GoNetworkCameraScanner

import (
	"fmt"
	"testing"
)

func TestRunt(t *testing.T) {
	got := Run("216.239.38.1-216.239.38.120", 80)
	fmt.Println("IP Addresses open: ")
	for i := range got {
		fmt.Println(got[i])
	}
}

func BenchmarkRun(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Run("185.60.219.1-41.185.60.219.250", 80)
	}
}
