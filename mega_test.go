package megach

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMegaChannel(t *testing.T) {
	ctx := context.Background()

	ch := NewMegaChannel()
	ch.Run(ctx)

	n := 100000

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for i := 0; i < n; i++ {
			assert.Equal(t, i, <-ch.Recv)
		}
		wg.Done()
	}()

	for i := 0; i < n; i++ {
		ch.Send <- i
	}

	wg.Wait()
}

func TestBitshift(t *testing.T) {
	fmt.Println(8 << 0)
}

func BenchmarkMegaChannelSend(b *testing.B) {
	ctx := context.Background()
	ch := NewMegaChannel()
	ch.Run(ctx)

	for i := 0; i < b.N; i++ {
		ch.Send <- i
	}
}

func BenchmarkMegaChannelSendRecv(b *testing.B) {
	ctx := context.Background()
	ch := NewMegaChannel()
	ch.Run(ctx)

	for i := 0; i < b.N; i++ {
		ch.Send <- i
	}
	for i := 0; i < b.N; i++ {
		assert.Equal(b, i, <-ch.Recv)
	}
}
