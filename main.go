package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func Retry(shell string, args string, timeout time.Duration, attempt int, waittime time.Duration) {
	for index, _ := range make([]int, attempt) {
		fmt.Printf("try %d time\n", index+1)
		err := func() error {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			return Execute(ctx, shell, args)
		}()
		if err != nil {
			fmt.Printf("task of %d time failed with %v\n", index+1, err)
			Wait(waittime)
			continue
		}
		fmt.Println("task success!!")
		break
	}
}

func Execute(ctx context.Context, shell string, args string) error {
	cmd := exec.Command(shell, args)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	errChan := make(chan error)

	go func() {
		err := cmd.Run()
		if err != nil {
			errChan <- err
		}
		close(errChan)
	}()

	select {
	case <-ctx.Done():
		cmd.Process.Kill()
		return ctx.Err()
	case err, ok := <-errChan:
		if ok {
			return err
		}
		return nil
	}
}

func Wait(duration time.Duration) {
	fmt.Printf("wait %v\n", duration)
	time.Sleep(duration)
}

func main() {
	shell := "sleep"
	args := "40"
	timeout := time.Duration(5 * time.Second)
	attempt := 5
	waittime := time.Duration(5 * time.Second)
	Retry(shell, args, timeout, attempt, waittime)
}
