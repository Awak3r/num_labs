package main

import (
	"fmt"
	"io"
	"os"

	"lab3/demo"
)

func main() {
	f, err := os.Create("results.txt")
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	defer f.Close()

	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	os.Stdout = w
	done := make(chan struct{})
	go func() {
		defer close(done)
		io.Copy(io.MultiWriter(orig, f), r)
	}()

	demo.Interpolation()
	demo.Spline()
	demo.LSM()
	demo.Derivative()
	demo.Integrate()

	w.Close()
	<-done
}
