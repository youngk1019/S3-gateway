package full_path

import (
	"fmt"
	"testing"
)

func TestSplitFullPath(t *testing.T) {
	ret := SplitFullPath("aa")
	fmt.Printf("%#v\n", ret)

	ret = SplitFullPath("")
	fmt.Printf("%#v\n", ret)

	ret = SplitFullPath("aa/bb/")
	fmt.Printf("%#v\n", ret)

	ret = SplitFullPath("aa/bb/cc")
	fmt.Printf("%#v\n", ret)
}
