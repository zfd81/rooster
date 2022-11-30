package stream

const (
	defaultWorkers = 16
	minWorkers     = 1
)

type rxOptions struct {
	unlimitedWorkers bool
	workers          int
}

// UnlimitedWorkers lets the caller use as many workers as the tasks.
func UnlimitedWorkers() Option {
	return func(opts *rxOptions) {
		opts.unlimitedWorkers = true
	}
}

// WithWorkers lets the caller customize the concurrent workers.
func WithWorkers(workers int) Option {
	return func(opts *rxOptions) {
		if workers < minWorkers {
			opts.workers = minWorkers
		} else {
			opts.workers = workers
		}
	}
}

// buildOptions returns a rxOptions with given customizations.
func buildOptions(opts ...Option) *rxOptions {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}

	return options
}

// newOptions returns a default rxOptions.
func newOptions() *rxOptions {
	return &rxOptions{
		workers: defaultWorkers,
	}
}
