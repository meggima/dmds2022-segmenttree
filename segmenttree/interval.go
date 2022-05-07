package segmenttree

var EmptyInterval Interval = NewInterval(0, 0)

type Interval struct {
	start uint32
	end   uint32
}

func NewInterval(start uint32, end uint32) Interval {
	if start > end {
		panic("Interval start must be before end")
	}

	return Interval{
		start: start,
		end:   end,
	}
}

func (interval Interval) IntersectionWith(otherInterval Interval) Interval {
	if interval.end < otherInterval.start ||
		otherInterval.end < interval.start {
		return EmptyInterval
	}

	start := MaxUint32(interval.start, otherInterval.start)
	end := MinUint32(interval.end, otherInterval.end)

	return NewInterval(start, end)
}

func (interval Interval) IsSubsetOf(otherInterval Interval) bool {
	return interval.start >= otherInterval.start &&
		interval.end <= otherInterval.end
}
