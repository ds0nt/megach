package megach

import (
	"context"
)

const buffSize = 10

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

func (m *MegaChannel) Stats() (_len, _cap int) {
	for i := 0; i < len(m.buffers); i++ {
		_len += len(m.buffers[i])
		_cap += cap(m.buffers[i])
	}
	return
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
			n := len(m.buffers)

			select {
			case m.buffers[n-1] <- x:
				// fmt.Printf("appending [%d] <- %d\n", n-1, x)
			default:
				m.buffers = append(m.buffers, make(chan interface{}, buffSize<<uint(n)))
				// fmt.Printf("new buffer [%d]\n", n)
				m.buffers[n] <- x
				// fmt.Printf("appending [%d] <- %d\n", n, x)
				go func(to, from int) {
					for {
						select {
						case <-ctx.Done():
							return
						case y := <-m.buffers[from]:
							// fmt.Printf("moving %d := <- [%d]\n", x, from)

							select {
							case <-ctx.Done():
								return
							case m.buffers[to] <- y:
								// fmt.Printf("moving [%d] <- %d\n", to, x)
							}
						}
					}
				}(n-1, n)
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case x := <-m.buffers[0]:
				// fmt.Printf("staging %d\n", x)
				select {
				case <-ctx.Done():
				case m.Recv <- x:
					// fmt.Printf("sending %d\n", x)
				}
			}
		}
	}()
}
