//go:build windows

package lock

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// DefaultTimeout is the default lock acquisition timeout for hook gates.
const DefaultTimeout = 2 * time.Second

type FileLock struct {
	file *os.File
	path string
}

type LockManager struct {
	locks      map[string]*FileLock
	locksMutex sync.RWMutex
}

var (
	globalLockManager *LockManager
	managerOnce       sync.Once
)

func GetLockManager() *LockManager {
	managerOnce.Do(func() {
		globalLockManager = &LockManager{
			locks:      make(map[string]*FileLock),
			locksMutex: sync.RWMutex{},
		}
	})
	return globalLockManager
}

func (lm *LockManager) Acquire(path string) error {
	lm.locksMutex.Lock()
	defer lm.locksMutex.Unlock()

	lockPath := path + ".lock"

	for i := 0; i < 5; i++ {
		file, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0600)
		if err == nil {
			lm.locks[path] = &FileLock{file: file, path: lockPath}
			return nil
		}
		time.Sleep(time.Millisecond * 100)
	}

	return fmt.Errorf("failed to acquire lock after 5 retries: %s", lockPath)
}

func (lm *LockManager) Release(path string) error {
	lm.locksMutex.Lock()
	defer lm.locksMutex.Unlock()

	lock, ok := lm.locks[path]
	if !ok {
		return fmt.Errorf("no lock found for path: %s", path)
	}

	if err := lock.file.Close(); err != nil {
		return err
	}
	os.Remove(lock.path)
	delete(lm.locks, path)
	return nil
}

func (lm *LockManager) AcquireWithTimeout(path string, timeout time.Duration) error {
	type result struct {
		err     error
		success bool
	}
	done := make(chan result, 1)

	go func() {
		err := lm.Acquire(path)
		done <- result{err: err, success: err == nil}
	}()

	select {
	case res := <-done:
		return res.err
	case <-time.After(timeout):
		go func() {
			res := <-done
			if res.success {
				lm.Release(path)
			}
		}()
		return fmt.Errorf("lock acquisition timeout for: %s", path)
	}
}
