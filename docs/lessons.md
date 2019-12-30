# Lessons learned

This has been a great learning experience with Go and programming in general.


I was getting this weird bug that was two parts. 

In spawning `wiz executor` with `exec.Command`, for some reason, the `log.Println` statements in Wiz Executor wouldn't work, probably because the same default `log` package or default logger or some shit was being used by the main program.

This then made me blind to the main issue which was that the executor was dying immediately. Once I figured out the loggin issue I saw the error:  `listen tcp :8080: bind: address already in use`.

## Logging
**Lesson:** When in doubt, use `fmt.Println` instead of the standard `log` package. This may also apply to certain 3rd party loggers.

Then, there was the question of why the address was in use. I checked before and after running the program (with `netstat -tulpn`) and nothing else was using it.

But this was a quicker fix. I remembered adding a helper function to test if an address is open for listening. I didn't remember to call `listener.Close()` though, so the listener must have still been there until the main program died.

A simple `defer l.Close()` did the trick.

## Error checking and Defer

**Lesson:** Take care of your business and close everything your done with. For files and listeners, as soon as you open, you also `defer smth.Close()`!

**Lesson:** Correction: check for error on the thing you open, and then call `defer thing.Close()` Otherwise you can get a nil pointer exception when `thing = nil, err = smth`

E.g. do this:
```go
l, err := net.Listen("tcp", ":"+strconv.FormatUint(uint64(port), 10))
// check error before defer, b/c on err, l will be nil
if err != nil {
    return false
}
defer l.Close()
```
not this:

```go

l, err := net.Listen("tcp", ":"+strconv.FormatUint(uint64(port), 10))
defer l.Close()
if err != nil {
    return false
    // now l.Close() will be called when l is nil
}
```

## Embedded structs and interfaces
Beware the embedded struct, for it may also be an interface.

It may seem quicker to write

```go
type Run struct {
	RunID         string
	Configuration Configuration

	State
}
```

than 


```go
type Run struct {
	RunID         string
	Configuration Configuration

	State State
}
```

but remember:
> Embedding structs will automatically promote the child struct's functions,

and since our State type is an enum that overrides the default `Marshal` and `Unmarshal` functions, it overwrites it for the parent type as well, resulting in some weird serialization bugs


## JSON/Serialization

Next bug:
```bash
2019/12/30 14:54:49 json: unsupported type: map[interface {}]interface {}
```

This was a source of a few bugs for me. And since it was usually happening during the very beginning or end of the program, at the (de)serialization steps, it was sometimes hard to catch (e.g. was happening inside a defer stmt or see weird logging above)

The fix is to switch to a JSON library that can handle more edge cases like this, namely
```go
import (
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary
```

The code's API is exactly the same, but bug is fixed!

TODO: look into whether the idiomatic/raw json-iter API is better than the regular `Marshall`/`Unmarshall`