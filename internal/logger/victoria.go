package logger

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type victoriaLogsWriter struct {
	endpoint string
	queue    chan []byte
	errs     chan error
}

func NewVictoriaLogsWriter(endpoint string) *victoriaLogsWriter {
	w := &victoriaLogsWriter{
		endpoint: endpoint,
		queue:    make(chan []byte, 100),
		errs:     make(chan error),
	}
	go w.send()
	return w
}

func (w victoriaLogsWriter) send() {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	for p := range w.queue {
		resp, err := client.Post(w.endpoint, "application/json", bytes.NewBuffer(p))
		if err != nil {
			w.errs <- fmt.Errorf("failed to do request: %w", err)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			w.errs <- fmt.Errorf("received non-OK response: %s", resp.Status)
		}
	}
}

func (w victoriaLogsWriter) Write(p []byte) (n int, err error) {
	select {
	case err := <-w.errs:
		return 0, err
	default:
	}

	select {
	case w.queue <- p:
		return len(p), nil
	case <-time.After(100 * time.Millisecond):
		return 0, errors.New("log queue is full, dropping log")
	}
}

func (w victoriaLogsWriter) Close() error {
	close(w.queue)
	close(w.errs)
	return nil
}
