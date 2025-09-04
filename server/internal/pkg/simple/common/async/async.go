package async

import (
	"context"
	"time"
)

type Future[T any] struct {
	await func(context.Context) (T, error)
}

func (f Future[T]) Await() (T, error) {
	return f.await(context.Background())
}

func (f Future[T]) AwaitTimeout(timeout time.Duration) (T, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()
	return f.await(ctx)
}

func Exec[T any](f func() (T, error)) Future[T] {
	var (
		result T
		err    error
	)

	c := make(chan struct{})
	go func() {
		defer close(c)
		result, err = f()
	}()
	return Future[T]{
		await: func(ctx context.Context) (T, error) {
			select {
			case <-ctx.Done():
				return result, ctx.Err()
			case <-c:
				return result, err
			}
		},
	}
}

func (f Future[T]) AwaitNoError() T {
	v, _ := f.await(context.Background())
	return v
}

func ExecNoErr[T any](f func() T) Future[T] {
	return Exec(func() (T, error) {
		return f(), nil
	})
}
