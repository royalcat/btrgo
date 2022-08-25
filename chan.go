package btrgo

// drain the channel
func Drain[T any](c chan T) {
L:
	for {
		select {
		case <-c:
		default:
			break L
		}
	}
}
