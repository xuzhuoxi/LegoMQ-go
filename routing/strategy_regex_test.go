package routing

import (
	"fmt"
	"testing"
)

type element string

func (e element) Id() string {
	return string(e)
}

func (e element) SetId(Id string) {
	return
}

func (e element) Formats() []string {
	return []string{string(e)}
}

func (e element) SetFormat(formats []string) {
	return
}

const (
	ele0 element = "192.168.3.2"
	ele1 element = "*.168.3.2"
	ele2 element = "192.*.3.2"
	ele3 element = "192.168.*.2"
	ele4 element = "192.168.3.*"
	ele5 element = "*.*.*.*"
	ele6 element = "*"
)

var (
	routingKey0 = "192.168.3.2"
	routingKey1 = "192.168.3.1"
)

func TestRegexStrategy_Route(t *testing.T) {
	s, err := RegexRouting.NewRoutingStrategy()
	if nil != err {
		t.Fatal(err)
	}
	cfg := s.Config()
	cfg.AppendRoutingTargets([]IRoutingElement{ele0, ele1, ele2, ele3, ele4, ele5, ele6})
	targets, err := s.Route(routingKey0)
	if nil != err {
		t.Fatal(err)
	}
	fmt.Println(targets)
	targets, err = s.Route(routingKey1)
	if nil != err {
		t.Fatal(err)
	}
	fmt.Println(targets)
}
