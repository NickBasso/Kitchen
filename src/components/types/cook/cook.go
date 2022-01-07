package cook

import "sync"

type Cook struct {
	ID           int
	Rank         int
	Proficiency  int
	Name         string
	CatchPhrase  string
	WorkingCount int
	Mutex 			 sync.Mutex
}
