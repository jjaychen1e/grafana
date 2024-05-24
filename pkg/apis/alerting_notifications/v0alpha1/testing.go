package v0alpha1

import (
	"fmt"
	"math/rand"
	"slices"
	"strings"
	"time"

	"github.com/grafana/grafana/pkg/util"
)

type IntervalMutator func(spec *Interval)

type IntervalGenerator struct {
	mutators []IntervalMutator
}

func (t IntervalGenerator) With(mutators ...IntervalMutator) IntervalGenerator {
	return IntervalGenerator{
		mutators: append(t.mutators, mutators...),
	}
}

func (t IntervalGenerator) generateDaysOfMonth() string {
	isRange := rand.Int()%2 == 0
	if !isRange {
		return fmt.Sprintf("%d", rand.Intn(30)+1)
	}
	from := rand.Intn(15) + 1
	to := rand.Intn(31-from) + from + 1
	return fmt.Sprintf("%d:%d", from, to)
}

func (t IntervalGenerator) generateTimeRange() TimeRange {
	from := rand.Int63n(1440 / 2)
	to := rand.Int63n(1440-from) + from + 1
	return TimeRange{
		StartTime: time.Unix(from*60, 0).UTC().Format("15:04"),
		EndTime:   time.Unix(to*60, 0).UTC().Format("15:04"),
	}
}

func (t IntervalGenerator) generateWeekday() string {
	day := rand.Intn(7)
	return strings.ToLower(time.Weekday(day).String())
}

func (t IntervalGenerator) generateYear() string {
	from := 1970 + rand.Intn(100)
	if rand.Int()%3 == 0 {
		to := 1970 + from + rand.Intn(10) + 1
		return fmt.Sprintf("%d:%d", from, to)
	}
	return fmt.Sprintf("%d", from)
}

func (t IntervalGenerator) generateLocation() *string {
	if rand.Int()%3 == 0 {
		return nil
	}
	return util.Pointer("UTC")
}

func (t IntervalGenerator) generateMonth() string {
	return fmt.Sprintf("%d", rand.Intn(12)+1)
}

func (t IntervalGenerator) GenerateMany(count int) []Interval {
	result := make([]Interval, 0, count)
	for i := 0; i < count; i++ {
		result = append(result, t.Generate())
	}
	return result
}

func (t IntervalGenerator) Generate() Interval {
	i := Interval{
		DaysOfMonth: generateMany(rand.Intn(6), true, t.generateDaysOfMonth),
		Location:    t.generateLocation(),
		Months:      generateMany(rand.Intn(3), true, t.generateMonth),
		Times:       generateMany(rand.Intn(6), true, t.generateTimeRange),
		Weekdays:    generateMany(rand.Intn(3), true, t.generateWeekday),
		Years:       generateMany(rand.Intn(3), true, t.generateYear),
	}
	for _, mutator := range t.mutators {
		mutator(&i)
	}
	return i
}

func generateMany[T comparable](repeatTimes int, unique bool, f func() T) []T {
	qty := repeatTimes + 1
	result := make([]T, 0, qty)
	for i := 0; i < qty; i++ {
		r := f()
		if unique && slices.Contains(result, r) {
			continue
		}
		result = append(result, f())
	}
	return result
}

func CopyWith(in Interval, mutators ...IntervalMutator) Interval {
	r := *in.DeepCopy()
	for _, mut := range mutators {
		mut(&r)
	}
	return r
}
