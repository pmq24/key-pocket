package cfg

import "github.com/spf13/pflag"

type SetGlobalOpts struct {
	Dir *pflag.Flag
	Profile *pflag.Flag
}

func SetGlobal(opts SetGlobalOpts) {
	dir = opts.Dir
	profile = opts.Profile
}

type Global struct {
	Dir string
	Profile string
}

func ReadGlobal() Global {
	return Global{
		Dir: dir.Value.String(),
		Profile: profile.Value.String(),
	}
}

var (
	dir *pflag.Flag
	profile *pflag.Flag
)
