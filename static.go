// +build static

package main

import staticFiles "github.com/silas/keyfu/static"

func init() {
	static = staticFiles.Data
}
