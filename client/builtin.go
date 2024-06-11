package client

import (
	"fmt"
	"sync"
)

func newSpecularPackage() (*Package, error) {
	pk := NewPackage(
		ModulePathFromTrustedValues("spec", "proto"),
	)
	if _, err := pk.NewType(
		"err",
		TypeBuilder(func() Struct {
			return NewError()
		}),
	); err != nil {
		return nil, err
	}
	if _, err := pk.NewType(
		"hb",
		TypeBuilder(func() Struct {
			return NewHeartbeat()
		}),
	); err != nil {
		return nil, err
	}
	return pk, nil
}

var builtinPackageOnce = sync.OnceValue(func() *Package {
	pk, err := newSpecularPackage()
	if err != nil {
		panic(fmt.Errorf("failed to initialize built-in package, %w", err))
	}
	return pk
})

// BuiltinPackage returns the built-in package
func BuiltinPackage() *Package {
	return builtinPackageOnce()
}
