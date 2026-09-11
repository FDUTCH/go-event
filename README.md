go-event

A tiny library for event-driven development

## Packages: event, storage

### Event:

Contains simple thread-safe event bus

Example of using event bus:
```go
package main

import "fmt"
import "github.com/FDUTCH/go-event/event"

type MyEvent struct {
  Message string
}

func main() {
  bus := event.NewBus()

  listener := func(ev MyEvent) {
    fmt.Println("received:", ev.Message)
  }

  bus.Subscribe(listener)
  bus.Publish(Event{"Hi!"})
  bus.Unsubscribe(listener)
}
```

### Storage:

Contains simple typed storage

Example of using typed storage:

```go
package main

import "fmt"
import "github.com/FDUTCH/go-event/storage"

type MyData struct {
  Message string
}

func main() {
  st := storage.NewStorage()

  st.Set(MyData{"Hi!")

  data := st.Get[MyData]()

  fmt.Println(data.Message)
}
```
