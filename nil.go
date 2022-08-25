package btrgo

func FirstNotNil[T any](args ...*T) *T {
	for _, arg := range args {
		if arg != nil {
			return arg
		}
	}
	return nil
}
