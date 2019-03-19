# megach

A self-expanding channel written Go that will never block on send.

```go
import "context"
  
ctx := context.Background()
ch := NewMegaChannel()
ch.Run(ctx)

for i := 0; i < 100; i++ {
	ch.Send <- i
}

for i := 0; i < 100; i++ {
	fmt.Println(<-ch.Recv)
}  
```
