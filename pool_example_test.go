package dj_test

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-dj/dj"
)

func ExampleNewWorkerPool() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a pool with 2 workers which converts strings to ints.
	pool := dj.NewWorkerPool(ctx, 2, func(ctx context.Context, in string) (int, error) {
		return strconv.Atoi(in)
	})

	// Submit 3 jobs to the pool.
	job1 := pool.Submit(ctx, "1")
	job2 := pool.Submit(ctx, "foo")
	job3 := pool.Submit(ctx, "3")

	// Each job should be processed.
	fmt.Println(job1.R())
	fmt.Println(job2.R())
	fmt.Println(job3.R())

	// Submit a list of jobs to the pool.
	fmt.Println(pool.Process(ctx, "4", "5", "6"))

	// If any job fails, the entire list of jobs will fail.
	fmt.Println(pool.Process(ctx, "7", "foo", "9"))

	// Output:
	// 1 <nil>
	// 0 strconv.Atoi: parsing "foo": invalid syntax
	// 3 <nil>
	// [4 5 6] <nil>
	// [] strconv.Atoi: parsing "foo": invalid syntax
}
