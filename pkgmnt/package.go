package pkgmnt

import "github.com/emirpasic/gods/lists/arraylist"

//Platform represents an Operating System family. Examples of
// platforms include Red Hat, Suse, Debian, among others.
type Platform struct {
	Name         string
	Runtime      string
	Buildtime    string
	Repositories Repositories
}

// Repository represents an address where a package
// is available on.
type Repository struct {
	Name        string
	Description string
	BaseURL     string
	Enabled     bool
	GPGCheck    bool
	GPGKey      string
}

// Repositories is a list of Repository
type Repositories struct {
	arraylist.List
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

// RepositoryStore provides access to a package store.
type RepositoryStore interface {
	Store(p Package) error
	StoreAll(ps Packages) error
	GetWithName(name string) (Package, error)
	All() (Packages, error)
}
