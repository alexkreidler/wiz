# Processor internals

The library has two types, Processor and ChunkProcessor

Processor is the basic type that can be configured and describe itself.

ChunkProcessor can only handle one chunk. The idea is that a new ChunkProcessor will be created for each
chunk and disposed of on completion.

The State function returns a channel with the chunk state, allowing the processor to update the manager
or controlling entity easily, without having to have knowledge of the data structure used to track the state of each chunk (in this case, a Sync.Map or a map with a mutex). We decided this would make the most sense in the library.

For example, then a third-party user could do something like this:

```go
package processors

import (
	"fmt"

	"github.com/alexkreidler/wiz/processors/processor"
)

func main() {
	var p processor.ChunkProcessor
	p := NewSpecialProcessor()

	fmt.Println(p.Metadata())
	
	data := map[string]interface{}{
		"test":"test",
	}
	
	go p.Run(data)
	for elem := range p.State() {
		fmt.Println("current state:", elem)
	}
}
```

You could even do something like the Go pipelines pattern

## Builtin Processors

TODO:

- [x] basic noop
- [ ] custom http downloader with Accept Ranges support
- [ ] aria2 downloader - expose as much config as possible
- [ ] ftp/sftp
- [ ] git
- [ ] (un)archive
- [ ] basic file mangement (move, copy, etc)