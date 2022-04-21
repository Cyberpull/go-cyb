package dbo

func Driver(opts *Options) DRIVER {
	if opts.Driver == "" {
		opts.Driver = DRIVER_PGSQL
	}

	return opts.Driver
}
