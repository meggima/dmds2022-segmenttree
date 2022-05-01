package segmenttree

type Interval struct {
	start uint32
	end   uint32
}

func NewInterval(start uint32, end uint32) Interval {
	return Interval{
		start: start,
		end:   end,
	}
}
