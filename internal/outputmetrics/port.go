package outputmetrics

type OutputMetrics interface {
	SetResult(name string, value bool)
	SetTime(name string, value int64)
	Flush()
}
