package segmenttree

func MinUint32(a uint32, b uint32) uint32 {
	if a < b {
		return a
	}

	return b
}

func MaxUint32(a uint32, b uint32) uint32 {
	if a > b {
		return a
	}

	return b
}
