package megach

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkMegaChannel(b *testing.B) {
	ctx := context.Background()

	ch := NewMegaChannel()
	ch.Run(ctx)

	b.N = 100

	for n := 0; n < b.N; n++ {
		ch.Send <- n
	}

	for n := 0; n < b.N; n++ {
		assert.Equal(b, n, <-ch.Recv)
	}
}
