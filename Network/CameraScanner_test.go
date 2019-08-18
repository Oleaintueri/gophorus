package Network

import (
	"fmt"
	"testing"
)

func TestRunt(t *testing.T) {
	got := Run("41.188.226.1-41.188.226.250", 1935)
	fmt.Println(got)
}

func BenchmarkRun(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Run("41.188.226.1-41.188.226.250", 1935)
	}
}
