package segmenttree

func Sum(x Addable, y Addable) Addable {
	return x.Add(y)
}

func InverseSum(x Addable, y Addable) Addable {
	return x.Subtract(y)
}

func Count(x Addable, y Addable) Addable {
	return x.Add(y)
}

func InverseCount(x Addable, y Addable) Addable {
	return x.Subtract(y)
}

func Average(x Addable, y Addable) Addable {
	return x.Add(y)
}

func InverseAverage(x Addable, y Addable) Addable {
	return x.Subtract(y)
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

func Identity(v Addable) Addable {
	return v
}

type Comparable interface {
	Compare(x Comparable) int
}

type Addable interface {
	Add(x Addable) Addable
	Inverse() Addable
	Subtract(x Addable) Addable
}

type AverageTuple struct {
	Sum   int
	Count int
}

func (x AverageTuple) Add(y Addable) Addable {
	return AverageTuple{
		Sum:   x.Sum + y.(AverageTuple).Sum,
		Count: x.Count + y.(AverageTuple).Count,
	}
}

func (x AverageTuple) Subtract(y Addable) Addable {
	return AverageTuple{
		Sum:   x.Sum - y.(AverageTuple).Sum,
		Count: x.Count - y.(AverageTuple).Count,
	}
}
func (x AverageTuple) Inverse() Addable {
	return AverageTuple{
		Sum:   -x.Sum,
		Count: -1,
	}
}

type Float float32

func (x Float) Add(y Addable) Addable {
	return x + y.(Float)
}

func (x Float) Inverse() Addable {
	return -x
}

func (x Float) Subtract(y Addable) Addable {
	return x - y.(Float)
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
