package main

import (
	"bytes"
	"fmt"
	"io"
	"syscall/js"
)

type envReader struct {
	buffer bytes.Buffer
}

func (r *envReader) Read(dest []byte) (n int, err error) {
	return r.buffer.Read(dest)
}

func (r *envReader) Close() error {
	return nil
}

type envWriter struct {
	buffer bytes.Buffer
	name   string
}

func (w *envWriter) Write(dest []byte) (n int, err error) {
	return w.buffer.Write(dest)
}

func (w *envWriter) Close() error {
	js.Global().Call("saveFile", w.name, w.buffer.String())
	return nil
}

func envRead(name string) (io.ReadCloser, error) {
	r := &envReader{}

	fileRead := make(chan string)
	fileError := make(chan error)

	defer close(fileRead)
	defer close(fileError)

	go js.Global().Call("openFile").Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) > 0 && args[0].Type() == js.TypeString {
			fileRead <- args[0].String()
		} else {
			fileError <- fmt.Errorf("canceled by user")
		}
		return nil
	}))

	select {
	case content := <-fileRead:
		_, err := r.buffer.WriteString(content)
		return r, err
	case err := <-fileError:
		return nil, err
	}
}

func envWrite(name string) (io.WriteCloser, error) {
	w := &envWriter{
		name: name,
	}
	return w, nil
}
