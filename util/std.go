package util

import (
	"io"
	"log"
)

// Closer closes c. If that fails, it will log the error.
// The intended usage is for convinient defer calls only!
// It gives only little knowledge about where the error is,
// but it's slightly better than a bare defer xyz.Close()
func Closer(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Printf("Error on close `%v`: %v", c, err)
	}
}
