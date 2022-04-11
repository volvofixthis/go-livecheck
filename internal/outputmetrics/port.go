package outputmetrics

type OutputMetrics interface {
	SetResult(name string, value bool)
	Flush()
}
