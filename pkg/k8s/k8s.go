package k8s

import "time"

const (
	sortByCPU            = "cpu"
	sortByMemory         = "memory"
	metricsCreationDelay = 2 * time.Minute
)
