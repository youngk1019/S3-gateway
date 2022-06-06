package util

type ArrayFlags []string

func (arr *ArrayFlags) String() string {
	return "array flags"
}

func (arr *ArrayFlags) Set(value string) error {
	*arr = append(*arr, value)
	return nil
}
