package scanner

type ParallelScanner struct {
	workers int
}

func NewParallelScanner(workers int) *ParallelScanner {
	if workers < 1 {
		workers = 1
	}
	return &ParallelScanner{workers: workers}
}

func (p *ParallelScanner) Workers() int {
	return p.workers
}
