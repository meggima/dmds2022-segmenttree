package segmenttree

func createAndSortValueTimeTuples(aggregate Aggregate, values []ValueIntervalTuple) []ValueTimeTuple {
	result := make([]ValueTimeTuple, 0, 2*len(values))

	for _, value := range values {
		positiveTuple := ValueTimeTuple{
			value: value.value,
			time:  value.interval.start,
		}

		negativeTuple := ValueTimeTuple{
			value: aggregate.inverseOperation(aggregate.neutralElement, value.value),
			time:  value.interval.end,
		}

		result = insertInOrder(aggregate, positiveTuple, result)
		result = insertInOrder(aggregate, negativeTuple, result)
	}

	return result
}

func insertInOrder(aggregate Aggregate, toInsert ValueTimeTuple, values []ValueTimeTuple) []ValueTimeTuple {
	for i := 0; i < len(values); i++ {
		if toInsert.time > values[i].time {
			continue
		} else if toInsert.time == values[i].time {
			additionResult := aggregate.operation(toInsert.value, values[i].value)

			if additionResult == aggregate.neutralElement {
				// Remove element
				return append(values[:i], values[i+1:]...)
			} else {
				// Replace existing element with combined element
				values[i] = ValueTimeTuple{value: additionResult, time: toInsert.time}
				return values
			}
		} else {
			// Insert element
			values = append(values[:i+1], values[i:]...)
			values[i] = toInsert
			return values
		}
	}

	return append(values, toInsert)
}
