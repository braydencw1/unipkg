package unipkg

type Manager interface {
	Install(pkg string, opts *Options) error
	Refresh(opts *Options) error
	Remove(pkg string, opts *Options) error
	Update(opts *Options) error
}

type Options struct {
	UseSudo bool
	DryRun  bool
	Logger  func(string)
}
