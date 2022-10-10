package xcfg

type Opts struct {
	rawList [][]byte
	files   []string
}

type OptFunc func(o *Opts)

func bindOpts(opt *Opts, opts ...OptFunc) {
	for _, f := range opts {
		f(opt)
	}
}

func WithConfigFile(name string) OptFunc {
	return func(o *Opts) {
		o.files = append(o.files, name)
	}
}

func WithRawData(raw []byte) OptFunc {
	return func(o *Opts) {
		o.rawList = append(o.rawList, raw)
	}
}
