//go:build production

package development

func Assert(condition bool, msg string) {
	// No-op in production
}
