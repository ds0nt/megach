package megach

import (
	"context"
	"fmt"
)

const buffSize = 8

type MegaChannel struct {
	Send    chan interface{}
	Recv    chan interface{}
	buffers []chan interface{}
}

func NewMegaChannel() *MegaChannel {
	return &MegaChannel{
		Send: make(chan interface{}),
		Recv: make(chan interface{}),
		buffers: []chan interface{}{
			make(chan interface{}, buffSize),
		},
	}
}

func (m *MegaChannel) Run(ctx context.Context) {
	// recv immediately, or to buffer
	go func() {
		for {
			var x interface{}
			select {
			case <-ctx.Done():
				return
			case x = <-m.Send:
			}
			lenBuffs := len(m.buffers)

			select {
			case m.buffers[lenBuffs-1] <- x:
				fmt.Printf("%d -> %d\n", x, lenBuffs-1)
			default:
				m.buffers = append(m.buffers, make(chan interface{}, buffSize))
				m.buffers[lenBuffs] <- x
				fmt.Printf("%d -> %d\n", x, lenBuffs)
				go func(to, from chan interface{}) {
					for {
						select {
						case <-ctx.Done():
							return
						default:
						}
						to <- <-from
					}
				}(m.buffers[lenBuffs], m.buffers[lenBuffs-1])
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case m.Recv <- <-m.buffers[0]:
			}
		}
	}()
}
