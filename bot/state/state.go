package state

import (
	"sync"
)

type singleton struct {
	states map[string]StateType
}

var instance *singleton
var once sync.Once

// getInstance returns the servers state instance
func getInstance() *singleton {
	once.Do(func() {
		instance = &singleton{states: map[string]StateType{}}
	})

	return instance
}

// TODO: must be thread safe
func HasGuildID(id string) bool {
	_, constains := getInstance().states[id]

	return constains
}
func AddGuildID(id string) {
	getInstance().states[id] = Normal
}
func SetState(id string, state StateType) {
	if !HasGuildID(id) {
		AddGuildID(id)
	}

	getInstance().states[id] = state
}
func GetState(id string) StateType {
	if !HasGuildID(id) {
		return Normal
	}

	return getInstance().states[id]
}
