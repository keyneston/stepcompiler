package stepcompiler

// This idea sadly doesn't work while maintaining the builder pattern.
//
type ChainableState interface {
	State

	NextChainable(State)
}

func ChainStates(list []ChainableState) State {
	if len(list) == 0 {
		panic("can't handle an empty list")
	}

	for i := 1; i < (len(list)); i++ {
		list[i-1].NextChainable(list[i])
	}

	return list[0]
}
