// Package senulator generates realistic, random data, and outputs it in the form of SenML.
/*
Special Thanks

A special thanks is owed to Ashley Troggio who, without her help, this package would not have been possible.
She taught me distributions and gave me the idea for categorizing generated data. The rest is just an abstraction.
*/
package senulator

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/cisco/senml"
	discrete "github.com/dgryski/go-discreterand"
)

// Request is the base struct that contains parameter for generating data that are applicable across all units.
type Request struct {
	Name        string
	Start       int64
	End         int64
	Duration    int64
	Version     int
	Debug       bool
	RecordCount int64
	Records     []senml.SenMLRecord
	Units       []Unit
}

// Unit represents an individual measureable unit with configurable fields to control things such as Ceiling and Floor.
type Unit struct {
	Name        string
	Symbol      string
	UseFloor    bool
	UseCeiling  bool
	Floor       float64
	Ceiling     float64
	Probability []float64
	Categories  map[int][]float64
	Reading     float64
	Interval    int64
}

// New creates a new senulator instance.
/*
New also checks for some basic errors, calculates Duration, and RecordCount.

RecordCount is part of a performance enhancing effort to presize the slices for larger data sizes (think a year or more.)
*/
func New(request Request) (Request, error) {
	// error if time is backwards. Although SenML tolerates this, we're a bit opinionated.
	if request.End-request.Start < 0 {
		return Request{}, errors.New("end time before start time")
	}

	// it is possible, currently, for people to modify the struct retroactively
	if len(request.Units) < 1 {
		return request, errors.New("no units provided")
	}

	request.Duration = request.End - request.Start

	// This was added to increase performance, it will need to be benchmarked
	for _, unit := range request.Units {
		request.RecordCount = request.RecordCount + (request.Duration / unit.Interval)
	}
	return request, nil
}

// Generate kicks off the generation of data given a request.
func (r *Request) Generate() ([]senml.SenMLRecord, error) {
	seed := rand.NewSource(time.Now().UnixNano())
	reader := rand.New(seed)
	r.Records = make([]senml.SenMLRecord, 0, r.RecordCount)

	// Generate volume data
	for _, unit := range r.Units {
		records, err := r.generateUnit(&unit, seed, reader)
		if err != nil {
			return []senml.SenMLRecord{}, errors.New("failed to generate unit record set")
		}
		r.Records = append(r.Records, records...)
	}

	return r.Records, nil
}

func (r *Request) generateUnit(unit *Unit, source rand.Source, reader *rand.Rand) ([]senml.SenMLRecord, error) {
	// Create an alias table from the probabilities
	table := discrete.NewAlias(unit.Probability, source)
	// Calculate iterations to the nearest integer
	iterations := int(r.Duration / unit.Interval)

	records := make([]senml.SenMLRecord, 0, iterations)

	for i := 0; i < iterations; i++ {
		// Determine the next category
		category := table.Next()
		// Fetch the ranges of the category
		categoryMap := unit.Categories[category]
		// Explicit which is which, for clarity
		lower := categoryMap[0]
		upper := categoryMap[1]
		// Generate the potential reading
		reading := reader.Float64() * (upper - lower)
		// Determine what to do with the reading with the Ceiling and Floor rules
		if unit.UseCeiling && unit.Reading+reading > unit.Ceiling {
			reading = -reading
		}
		if unit.UseFloor && unit.Reading-reading < unit.Floor {
			reading = -reading
		}
		// Calculate the record time
		time := r.Start + int64(i)*unit.Interval
		// Generate the record
		record, err := r.createRecord(reading, unit.Symbol, time)
		if err != nil {
			return []senml.SenMLRecord{}, errors.New("failed to create record")
		}
		records = append(records, record)
		// Keep track of the totals so we can do fancy stuff afterwards
		unit.Reading = unit.Reading + reading
	}

	if r.Debug {
		fmt.Printf("Final %s reading: %f\n", unit.Name, unit.Reading)
	}

	return records, nil
}

func (r *Request) createRecord(value float64, unit string, time int64) (senml.SenMLRecord, error) {
	if time < r.Start {
		return senml.SenMLRecord{}, errors.New("record time before start time")
	}
	// Bootstrap the record for any Base info
	record := senml.SenMLRecord{
		BaseName:    r.Name,
		BaseTime:    float64(r.Start),
		BaseVersion: r.Version,
		Time:        float64(time - r.Start),
		Value:       &value,
		Unit:        unit,
	}

	return record, nil
}
