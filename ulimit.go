package ulimit

import (
	"fmt"
	"syscall"
)

// SetMax sets the maximum possible filesystem ulimit (rLimit.Cur).
func SetMax() error {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return fmt.Errorf("error getting Rlimit: %s", err.Error())
	}

	// Try to set the max
	rLimit.Cur = rLimit.Max
	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err == nil {
		return nil
	}

	// If that fails, set the maximum possible using brute force
	i := rLimit.Max
	lowerBound := uint64(0)
	for {
		rLimit.Cur = i
		err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
		if err == nil {
			if lowerBound == 0 {
				lowerBound = uint64(i)
				i = i * 10
			} else {
				break
			}
		}
		if lowerBound == 0 {
			i = i / 10
		} else {
			i--
		}
		if i <= 1 {
			return fmt.Errorf("error setting Rlimit: %s", err.Error())
		}
	}
	return set(rLimit, uint64(i))
}

// Get gets the current ulimit.
func Get() (i uint64, err error) {
	var rLimit syscall.Rlimit
	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return 0, fmt.Errorf("error getting Rlimit: %s", err.Error())
	}
	return uint64(rLimit.Cur), nil
}

// Set sets the ulimit.
func Set(i uint64) error {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return fmt.Errorf("error getting Rlimit: %s", err.Error())
	}
	return set(rLimit, i)
}

func set(rLimit syscall.Rlimit, i uint64) error {
	if uint64(rLimit.Cur) < i {
		rLimit.Cur = i
		err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
		if err != nil {
			return fmt.Errorf("error setting Rlimit: %s", err.Error())
		}
	}
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return fmt.Errorf("error getting Rlimit: %s", err.Error())
	}
	if uint64(rLimit.Cur) < i {
		return fmt.Errorf("rLimit.Cur update failure - expected at least: %d result: %d", i, rLimit.Cur)
	}
	return nil
}
