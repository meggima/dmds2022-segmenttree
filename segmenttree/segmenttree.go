package segmenttree

type SegmentTree interface {
	GetAtInstant(instant uint32) float32
	GetWithinInterval(interval Interval) []ValueIntervalTuple
	Insert(value ValueIntervalTuple)
	Delete(value ValueIntervalTuple)
}
