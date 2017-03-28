package cumulative_converter

import (
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
	"github.com/influxdata/telegraf/plugins/processors"
)

type CumulativeConverter struct {
	cache map[uint64]*state
}

func NewCumulativeConverter() telegraf.Processor {
	cc := &CumulativeConverter{}
	cc.cache = make(map[uint64]*state)
	return cc
}

type state struct {
	fields	map[string]float64
	name string
	tags	map[string]string
}

var sampleConfig = `
`

func (p *CumulativeConverter) SampleConfig() string {
	return sampleConfig
}

func (p *CumulativeConverter) Description() string {
	return "Convert cumulative field values to interval field values."
}

func (p *CumulativeConverter) Apply(in ...telegraf.Metric) []telegraf.Metric {
	out := make([]telegraf.Metric, 0, len(in))
	for _, m := range in {
		id := m.HashID()
		s := p.cache[id];
		if s == nil {
			// hit a new metric, save its value and discard it
			s = &state{
				name:   m.Name(),
				tags:   m.Tags(),
				fields: make(map[string]float64),
			}
			for k, v := range m.Fields() {
				if fv, ok := convert(v); ok {
					s.fields[k] = fv
				}
			}
			p.cache[id] = s
			continue
		}
		metric_reset := false
		fields := make(map[string]interface{})
		for k, v := range m.Fields() {
			if fv, ok := convert(v); ok {
				pv := s.fields[k]
				if fv < pv {
					// stat reset, discard this
					if !metric_reset {
						metric_reset = true
					}
				} else {
					fields[k] = fv - pv
				}
				s.fields[k] = fv
			}
		}
		if metric_reset {
			continue
		}
		// TODO handle error _
		newmetric, _ := metric.New(m.Name(), m.Tags(), fields, m.Time())
		out = append(out, newmetric)
	}
	return out
}

func convert(in interface{}) (float64, bool) {
	switch v := in.(type) {
	case float64:
		return v, true
	case int64:
		return float64(v), true
	default:
		return 0, false
	}
}

func init() {
	processors.Add("cumulative_converter", func() telegraf.Processor {
		return NewCumulativeConverter()
	})
}
