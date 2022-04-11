package outputmetrics

import (
	"bitbucket.rbc.ru/go/go-livecheck/internal/config"
	"github.com/mitchellh/mapstructure"
)

func NewOutputMetrics(c *config.Config) (OutputMetrics, error) {
	o := InfluxDBOutputMetrics{Dst: c.OutputMetrics.Dst}
	if len(c.OutputMetrics.Extra) > 0 {
		err := mapstructure.Decode(c.OutputMetrics.Extra, &o.Extra)
		if err != nil {
			return nil, err
		}
	}
	return &o, nil
}
