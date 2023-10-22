# ulimit

[![godoc](https://godoc.org/github.com/calmdocs/ulimit?status.svg)](https://godoc.org/github.com/calmdocs/ulimit)

ulimit.SetMax() allows users to set the maximum possible filesystem ulimit.

```
oldLimit, err := ulimit.Get()
if err != nil {
	log.Fatal(err)
}
err = ulimit.SetMax()
if err != nil {
	log.Fatal(err)
}
newLimit, err := ulimit.Get()
if err != nil {
	log.Fatal(err)
}
fmt.Println(oldLimit, newLimit)
// -> 256 24576
```

## The problem

In theory, when we run the following go code, we should be able to set the ulimit (rLimit.Cur) to any number up to rLimit.Max: 

```
var rLimit syscall.Rlimit
err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
...
```

Unfortunately, if we try to set rLimit.Cur as rLimit.Max, the system may return an error because the actual limit is lower than rLimit.Max.  

## The solution

ulimit.SetMax() uses brute force to find and then set the maximum possible filesystem ulimit.
