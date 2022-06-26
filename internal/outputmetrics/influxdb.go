package outputmetrics

import (
	"time"

	"github.com/rcrowley/go-metrics"
	influxdb "github.com/volvofixthis/go-metrics-influxdb"
)

type InfluxDBOutputMetrics struct {
	Dst   string
	Extra InfluxDBExtra
}

type InfluxDBExtra struct {
	Database    string            `mapstructure:"database"`
	Measurement string            `mapstructure: "measurement"`
	Username    string            `mapstructure:"username"`
	Password    string            `mapstructure:"password"`
	Tags        map[string]string `mapstructure:"tags"`
}

func (o *InfluxDBOutputMetrics) SetResult(name string, value bool) {
	g := metrics.NewGauge()
	metrics.Register(name, g)
	var valueInt int64 = 0
	if value {
		valueInt = 1
	}
	g.Update(valueInt)
}

func (o *InfluxDBOutputMetrics) SetTime(name string, value int64) {
	g := metrics.NewGauge()
	metrics.Register(name+".time", g)
	g.Update(value)
}

func (o *InfluxDBOutputMetrics) Flush() {
	go influxdb.InfluxDBWithTags(metrics.DefaultRegistry,
		25*time.Millisecond,
		o.Dst,
		o.Extra.Database,
		o.Extra.Measurement,
		o.Extra.Username,
		o.Extra.Password,
		o.Extra.Tags,
		true,
	)
	time.Sleep(40 * time.Millisecond)
}
