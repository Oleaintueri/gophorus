/*
Author Alano Terblanche (Benehiko) with the guidance of Kent Gruber's Medium article (https://medium.com/@KentGruber/building-a-high-performance-port-scanner-with-golang-9976181ec39d)
License is the Apache License and can be found in the project as LICENSE.md
Extra information can be found in README.md
 */
package GoNetworkCameraScanner

import (
	"fmt"
	"testing"
)

//Run Test on "Run" function
func TestRunt(t *testing.T) {
	got := Run("216.58.223.1-216.58.223.142", 80)
	fmt.Println("IP Addresses open: ")
	for i := range got {
		fmt.Println(got[i])
	}
}

//Run Benchmark on "Run" function
func BenchmarkRun(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Run("216.58.223.1-216.58.223.142", 80)
	}
}
