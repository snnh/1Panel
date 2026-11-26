// Package xpack holds the community implementation of the optional provider
// interfaces used by core. It is a thin pass-through to core/app/auth.
package xpack

import "github.com/1Panel-dev/1Panel/core/utils/xpack/helper"

var AuthProvider = helper.NewIAuthProvider()

// IsMasterProductPro reports whether the panel runs as the commercial master
// product. The community build never does; the symbol is kept so that shared
// handlers can stay identical to upstream.
func IsMasterProductPro() bool {
	return false
}
