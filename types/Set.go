package types

type Set[T comparable] map[T]bool

func NewSet[T comparable](T) Set[T] {
	return make(Set[T])
}

//  Add elements

func (set Set[T]) Add1(values ...T) {
	for _, value := range values {
		set[value] = true
	}
}
func (set Set[T]) Add(v []T) {
	for _, value := range v {
		set[value] = true
	}
}

// Delete elements

func (set Set[T]) Delete(values ...T) {
	for _, value := range values {
		delete(set, value)
	}
}

// length

func (set Set[T]) Len() int {
	return len(set)
}

func (set Set[T]) Has(value T) bool {
	_, ok := set[value]
	return ok
}

func (set Set[T]) Iterate(it func(T)) {
	for v := range set {
		it(v)
	}
}

/*

 */
func (set Set[T]) Values() []T {
	values := make([]T, set.Len())
	set.Iterate(func(value T) {
		values = append(values, value)
	})
	return values
}

func (set Set[T]) Clone() Set[T] {
	set1 := make(Set[T])
	set1.Add1(set.Values()...)
	return set1
}

/*
	union of 2 sets
*/

func (set Set[T]) Union(other Set[T]) Set[T] {
	set1 := set.Clone()
	set1.Add1(other.Values()...)
	return set1
}

func (set Set[T]) Intersection(other Set[T]) Set[T] {
	set1 := make(Set[T])
	set.Iterate(func(value T) {
		if other.Has(value) {
			set1.Add1(value)
		}
	})
	return set1
}
