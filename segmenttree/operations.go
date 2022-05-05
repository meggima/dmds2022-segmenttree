package segmenttree

func Sum(x Addable, y Addable) Addable {
	return x.Add(y)
}

func Count(x Addable, y Addable) Addable {
	return x.Add(y)
}

func Average(x Addable, y Addable) Addable {
	return x.Add(y)
}

func Min(x Comparable, y Comparable) Comparable {
	res := x.Compare(y)

	if res < 0 {
		return x
	}

	return y
}

func Max(x Comparable, y Comparable) Comparable {
	res := x.Compare(y)
	if res > 0 {
		return x
	}

	return y
}

type Comparable interface {
	Compare(x Comparable) int
}

type Addable interface {
	Add(x Addable) Addable
}

type AverageTuple struct {
	Sum   int
	Count int
}

func (x AverageTuple) Add(y AverageTuple) AverageTuple {
	return AverageTuple{
		Sum:   x.Sum + y.Sum,
		Count: x.Count + y.Count,
	}
}

type Float float32

func (x Float) Add(y Addable) Addable {
	return x + y.(Float)
}

func (x Float) Compare(y Float) int {
	if x > y {
		return 1
	}
	if x < y {
		return -1
	}

	return 0
}
