package Network

import (
	"fmt"
	"testing"
)

func TestRunt(t *testing.T) {
	got := Run("185.60.219.1-41.185.60.219.4", 80)
	fmt.Println(got)
}

func BenchmarkRun(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Run("185.60.219.1-41.185.60.219.250", 80)
	}
}
