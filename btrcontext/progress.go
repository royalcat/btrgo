package btrcontext

import (
	"context"

	"github.com/royalcat/btrgo/btrchannels"
)

type UpdateContext[T any] interface {
	context.Context
	Update(T)
}

type updateContext[T any] struct {
	context.Context
	updateBuf *btrchannels.RingChannel[T] // OPTIMIZATION create single element chan
}

var _ UpdateContext[any] = (*updateContext[any])(nil)

func (ctx *updateContext[T]) Update(update T) {
	ctx.updateBuf.In() <- update
}

func (ctx *updateContext[T]) Close(update T) {
	ctx.updateBuf.Close()
}

func WithUpdates[T any](parent context.Context) (ctx context.Context, updates <-chan T) {
	updateBuf := btrchannels.NewRingChannel[T](1)
	ctx = &updateContext[T]{
		Context:   parent,
		updateBuf: updateBuf,
	}
	return ctx, updateBuf.Out()
}
