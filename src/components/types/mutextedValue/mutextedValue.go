package mutextedValue

import "sync"

type MutextedValue struct {
	Value 	int
	Mutex 	sync.Mutex
}
