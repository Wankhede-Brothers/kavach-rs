//go:build !windows

package lock

import (
	"fmt"
	"os"
	"sync"
	"syscall"
	"time"
)

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

// DefaultTimeout is the default lock acquisition timeout for hook gates.
const DefaultTimeout = 2 * time.Second

func (lm *LockManager) Acquire(path string) error {
	lm.locksMutex.Lock()
	defer lm.locksMutex.Unlock()

	lockPath := path + ".lock"
	file, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return fmt.Errorf("failed to create lock file: %w", err)
	}

	for i := 0; i < 10; i++ {
		err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if err == nil {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}

	if err != nil {
		file.Close()
		return fmt.Errorf("failed to acquire lock after retries: %w", err)
	}

	lm.locks[path] = &FileLock{file: file, path: lockPath}
	return nil
}

func (lm *LockManager) Release(path string) error {
	lm.locksMutex.Lock()
	defer lm.locksMutex.Unlock()

	lock, ok := lm.locks[path]
	if !ok {
		return fmt.Errorf("no lock found for path: %s", path)
	}

	syscall.Flock(int(lock.file.Fd()), syscall.LOCK_UN)
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
