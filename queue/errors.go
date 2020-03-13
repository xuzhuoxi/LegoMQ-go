package queue

import "errors"

var (
	ErrQueueClosed = errors.New("ContextQueue is closed! ")
	ErrQueueFull   = errors.New("Cache is Full! ")
	ErrQueueEmpty  = errors.New("Cache is Empty! ")

	inputCtxArrayNilError  = errors.New("Input Context Array is Nil! ")
	outputCtxArrayNilError = errors.New("Output Context Array is Nil! ")

	ErrSize = errors.New("Max size should be larger than 0. ")
)
