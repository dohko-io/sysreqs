package pkgmnt

import "github.com/emirpasic/gods/lists/arraylist"

//Platform represents
type Platform struct {
	Name      string
	Runtime   string
	Buildtime string
}

// Package comprises a library or application that implements a set of related
// commands or features.
type Package struct {
	Name         string
	Version      string
	Architecture string
	Platforms    []Platform
	Description  string
	Dependencies []Package
}

// Packages is a list of package
type Packages struct {
	arraylist.List
}

// Repository provides access to a package repository
type Repository interface {
	Store(p Package) error
	StoreAll(ps Packages) error
	GetWithName(name string) (Package, error)
	All() (Packages, error)
}
