package util

import (
	"os"
	"time"

	"github.com/nightlyone/lockfile"
)

const (
	CatlightLockFile = "/tmp/.catlight.lock"
)

func LockCatlight() error {
	lock, err := lockfile.New(CatlightLockFile)
	if err != nil {
		return err
	}

	for {
		err := lock.TryLock()
		if err == lockfile.ErrBusy {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		if err != nil {
			return err
		}

		return nil
	}
}

func UnlockCatlight() error {
	lock, err := lockfile.New(CatlightLockFile)
	if err != nil {
		return err
	}

	return lock.Unlock()
}

func CleanupCatlightLock() error {
	// TODO: Warn when removed
	return os.RemoveAll(CatlightLockFile)
}
