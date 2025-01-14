//go:build !production

package development

import "fmt"

func Assert(condition bool, msg string) {
	if !condition {
		panic(fmt.Sprintf("DEV: %s", msg))
	}
}
