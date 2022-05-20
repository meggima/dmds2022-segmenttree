package segmenttree

type Aggregate struct {
	operation        func(Addable, Addable) Addable
	inverseOperation func(Addable, Addable) Addable
	additionElement  func(Addable) Addable
	neutralElement   Addable
}
