package apparatus

type ApparatusState struct {
	Ovens  OvenState
	Stoves StoveState
}

type OvenState struct {
	TotalCount int
	FreeCount  int
}

type StoveState struct {
	TotalCount int
	FreeCount  int
}
