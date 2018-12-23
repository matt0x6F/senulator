package senulator

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/cisco/senml"
)

func TestNew(t *testing.T) {
	type args struct {
		request Request
	}
	tests := []struct {
		name    string
		args    args
		want    Request
		wantErr bool
	}{
		// Test Table
		{
			"error-on-negative-duration",
			args{
				Request{
					Name:  "test",
					Start: time.Now().Unix(),
					End:   time.Now().Add(time.Hour * 24).Unix(),
				},
			},
			Request{
				Name:     "test",
				Start:    time.Now().Unix(),
				End:      time.Now().Add(time.Hour * 24).Unix(),
				Duration: time.Now().Add(time.Hour*24).Unix() - time.Now().Unix(),
			},
			false,
		},
		{
			"error-on-end-before-start",
			args{
				Request{
					End:   time.Now().Unix(),
					Start: time.Now().Add(time.Hour * 24).Unix(),
				},
			},
			Request{},
			true,
		},
		{
			"calculate-interval",
			args{
				Request{
					End:   time.Now().Add(time.Hour * 24).Unix(),
					Start: time.Now().Unix(),
					Units: []Unit{
						Unit{
							Interval: 900,
						},
					},
				},
			},
			Request{
				End:         time.Now().Add(time.Hour * 24).Unix(),
				Start:       time.Now().Unix(),
				Duration:    time.Now().Add(time.Hour*24).Unix() - time.Now().Unix(),
				RecordCount: 96,
				Units: []Unit{
					Unit{
						Interval: 900,
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_Generate(t *testing.T) {
	type fields struct {
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
	tests := []struct {
		name    string
		fields  fields
		want    []senml.SenMLRecord
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Request{
				Name:        tt.fields.Name,
				Start:       tt.fields.Start,
				End:         tt.fields.End,
				Duration:    tt.fields.Duration,
				Version:     tt.fields.Version,
				Debug:       tt.fields.Debug,
				RecordCount: tt.fields.RecordCount,
				Records:     tt.fields.Records,
				Units:       tt.fields.Units,
			}
			got, err := r.Generate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Request.Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Request.Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_generateUnit(t *testing.T) {
	type fields struct {
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
	type args struct {
		unit   *Unit
		source rand.Source
		reader *rand.Rand
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []senml.SenMLRecord
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Request{
				Name:        tt.fields.Name,
				Start:       tt.fields.Start,
				End:         tt.fields.End,
				Duration:    tt.fields.Duration,
				Version:     tt.fields.Version,
				Debug:       tt.fields.Debug,
				RecordCount: tt.fields.RecordCount,
				Records:     tt.fields.Records,
				Units:       tt.fields.Units,
			}
			got, err := r.generateUnit(tt.args.unit, tt.args.source, tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("Request.generateUnit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Request.generateUnit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_createRecord(t *testing.T) {
	testValue := 1.0
	type fields struct {
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
	type args struct {
		value float64
		unit  string
		time  int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    senml.SenMLRecord
		wantErr bool
	}{
		{
			"create-record",
			fields{
				Name:     "test",
				Start:    time.Now().Unix(),
				End:      time.Now().Add(time.Hour * 24).Unix(),
				Duration: time.Now().Add(time.Hour*24).Unix() - time.Now().Unix(),
				Version:  1,
			},
			args{
				value: testValue,
				unit:  "T",
				time:  time.Now().Add(time.Minute * 15).Unix(),
			},
			senml.SenMLRecord{
				BaseName:    "test",
				BaseTime:    float64(time.Now().Unix()),
				BaseVersion: 1,
				Time:        float64(time.Now().Add(time.Minute*15).Unix() - time.Now().Unix()),
				Value:       &testValue,
				Unit:        "T",
			},
			false,
		},
		{
			"error-on-record-time-before-start-time",
			fields{
				Name:     "test",
				Start:    time.Now().Add(time.Minute * 15).Unix(),
				End:      time.Now().Add(time.Hour * 24).Unix(),
				Duration: time.Now().Add(time.Hour*24).Unix() - time.Now().Add(time.Minute*15).Unix(),
				Version:  1,
			},
			args{
				value: testValue,
				unit:  "T",
				time:  time.Now().Unix(),
			},
			senml.SenMLRecord{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Request{
				Name:        tt.fields.Name,
				Start:       tt.fields.Start,
				End:         tt.fields.End,
				Duration:    tt.fields.Duration,
				Version:     tt.fields.Version,
				Debug:       tt.fields.Debug,
				RecordCount: tt.fields.RecordCount,
				Records:     tt.fields.Records,
				Units:       tt.fields.Units,
			}
			got, err := r.createRecord(tt.args.value, tt.args.unit, tt.args.time)
			if (err != nil) != tt.wantErr {
				t.Errorf("Request.createRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Request.createRecord() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Examples

func Example() {
	// Create a new request.
	// At least Start, End, and 1 unit are required or New() will return an error
	request := Request{
		Name:  "test:meter:1",
		Start: time.Now().Unix(),
		End:   time.Now().Add(time.Hour * 24).Unix(),
		Units: []Unit{
			Unit{
				Name:        "Volume",
				Symbol:      "L",
				Probability: []float64{.7, .2, .1},
				Categories: map[int][]float64{
					0: []float64{0.0, 0.0},
					1: []float64{0.1, 19},
					2: []float64{19.1, 56.7812},
				},
				Reading:  0,
				Interval: 900,
			},
		},
	}

	meter, err := New(request)
	if err != nil {
		fmt.Println("Error occured!")
	}

	records, err := meter.Generate()
	if err != nil {
		fmt.Println("Error generating records!")
	}

	pack := senml.SenML{
		Records: records,
	}
	options := senml.OutputOptions{
		PrettyPrint: true,
	}
	jsonData, err := senml.Encode(pack, senml.JSON, options)
	if err != nil {
		fmt.Println("Error: Unable to encode JSON")
	}

	fmt.Printf("data: %s\n", jsonData)
}
