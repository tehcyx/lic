package core

//Options defines available options for the command
type Options struct {
	Verbose bool
}

//NewOptions creates options with default values
func NewOptions() *Options {
	return &Options{}
}
