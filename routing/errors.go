package routing

import "errors"

var (
	ErrRougingRegister = errors.New("RoutingMode Unregister! ")

	ErrRoutingTargetEmpty = errors.New("RoutingStrategy TargetIds is empty! ")

	ErrRoutingFail = errors.New("RoutingStrategy route failed! ")
)
