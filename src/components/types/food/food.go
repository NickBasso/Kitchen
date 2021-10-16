package food

import "kitchen/src/components/types/apparatus"

type Food struct {
	Id              int16
	Name            string
	PreparationTime int
	Complexity      byte
	Apparatus       apparatus.Apparatus
}
