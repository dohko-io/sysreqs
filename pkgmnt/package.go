package pkgmnt

import "github.com/emirpasic/gods/lists/arraylist"

//Platform represents an Operating System family. Examples of
// platforms include Red Hat, Suse, Debian, among others.
type Platform struct {
	Name         string       `json:"name"`
	Runtime      string       `json:"runtime"`
	Buildtime    string       `json:"buildtime"`
	Repositories []Repository `json:"repositories"`
}

// Repository represents an address where a package
// is available on.
type Repository struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	BaseURL     string `json:"baseurl"`
	Enabled     bool   `json:"enabled"`
	GPGCheck    bool   `json:"gpgcheck"`
	GPGKey      string `json:"gpgkey"`
}

// Repositories is a list of Repository
type Repositories struct {
	arraylist.List
}

// Package comprises a library or application that implements a set of related
// commands or features.
type Package struct {
	Name         string     `json:"name"`
	Version      string     `json:"version"`
	Architecture string     `json:"architecture"`
	Platforms    []Platform `json:"platforms"`
	Description  string     `json:"description"`
	Dependencies []Package  `json:"dependencies"`
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
