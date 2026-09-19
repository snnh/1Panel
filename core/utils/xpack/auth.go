// Package xpack holds the community implementation of the optional provider
// interfaces used by core. It is a thin pass-through to core/app/auth.
package xpack

import "github.com/1Panel-dev/1Panel/core/utils/xpack/helper"

var AuthProvider = helper.NewIAuthProvider()
