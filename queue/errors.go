package queue

import "errors"

var (
	cacheFullError  = errors.New("Cache is Full! ")
	cacheEmptyError = errors.New("Cache is Empty! ")

	inputCtxArrayNilError  = errors.New("Input Context Array is Nil! ")
	outputCtxArrayNilError = errors.New("Output Context Array is Nil! ")

	sizeError       = errors.New("Max size should be larger than 0. ")
	cacheCloseError = errors.New("Cache is closed! ")
)
