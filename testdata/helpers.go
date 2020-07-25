package testdata

import (
	"bytes"
	"io"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

func ReadStdOut(f func()) string {
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	stderr := os.Stderr
	logrus.SetOutput(w)
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
		logrus.SetOutput(os.Stderr)
	}()
	os.Stdout = w
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		io.Copy(&buf, r)
		out <- buf.String()
	}()
	wg.Wait()
	f()
	w.Close()
	return <-out
}
