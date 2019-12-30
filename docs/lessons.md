# Lessons learned

This has been a great learning experience with Go and programming in general.

## 1. 
I was getting this weird bug that was two parts. 

In spawning `wiz executor` with `exec.Command`, for some reason, the `log.Println` statements in Wiz Executor wouldn't work, probably because the same default `log` package or default logger or some shit was being used by the main program.

This then made me blind to the main issue which was that the executor was dying immediately. Once I figured out the loggin issue I saw the error:  `listen tcp :8080: bind: address already in use`.

**Lesson:** When in doubt, use `fmt.Println` instead of the standard `log` package. This may also apply to certain 3rd party loggers.

Then, there was the question of why the address was in use. I checked before and after running the program (with `netstat -tulpn`) and nothing else was using it.

But this was a quicker fix. I remembered adding a helper function to test if an address is open for listening. I didn't remember to call `listener.Close()` though, so the listener must have still been there until the main program died.

A simple `defer l.Close()` did the trick.

**Lesson:** Take care of your business and close everything your done with. For files and listeners, as soon as you open, you also `defer smth.Close()`!