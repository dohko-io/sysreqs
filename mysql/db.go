package mysql

import (
	"database/sql"
	"encoding/json"

	"github.com/dohko-io/sysreqs/pkgmnt"
)

const (
	sqlInsertPackage     = "INSERT INTO package (name, version, architecture, description, platforms, dependencies) VALUES (?, ?, ?, ?, ?, ?)"
	sqlSelectAllPackages = "SELECT name, version, architecture, description, platforms, dependencies FROM package"
)

type packageRepository struct {
	db *sql.DB
}

func (r *packageRepository) Store(pkg pkgmnt.Package) error {
	stmt, err := r.db.Prepare(sqlInsertPackage)
	defer stmt.Close()

	if err != nil {
		return err
	}
	return store(stmt, pkg)
}

func (r *packageRepository) StoreAll(ps pkgmnt.Packages) error {

	stmt, err := r.db.Prepare(sqlInsertPackage)
	defer stmt.Close()

	if err != nil {
		return err
	}

	iter := ps.Iterator()

	for iter.Next() {
		err := store(stmt, iter.Value().(pkgmnt.Package))

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *packageRepository) Update(p pkgmnt.Package) error {

	return nil
}

func (r *packageRepository) GetWithName(name string) (pkgmnt.Package, error) {
	var pkg pkgmnt.Package

	stmt, err := r.db.Prepare(sqlSelectAllPackages + " WHERE lower(name) = lower(?)")
	defer stmt.Close()

	if err != nil {
		return pkg, err
	}

	rows, err := stmt.Query(name)

	if err != nil {
		return pkg, err
	}

	ps, err := rowsMapper(rows)

	if err != nil {
		return pkg, err
	}

	if !ps.Empty() {
		p, _ := ps.Get(0)
		pkg = p.(pkgmnt.Package)
	}

	return pkg, nil
}

func (r *packageRepository) All() (pkgmnt.Packages, error) {
	var packages pkgmnt.Packages

	stmt, err := r.db.Prepare(sqlSelectAllPackages + " ORDER BY name, version")
	defer stmt.Close()

	if err != nil {
		return packages, err
	}

	rows, err := stmt.Query()

	if err != nil {
		return packages, err
	}

	packages, err = rowsMapper(rows)
	return packages, err
}

// NewPackageRepository creates and returns an instance of package's
// RepositoryStore.
func NewPackageRepository(ds string) (pkgmnt.RepositoryStore, error) {
	db, err := sql.Open("mysql", ds)

	if err != nil {
		return nil, err
	}

	r := &packageRepository{db: db}
	err = db.Ping()

	return r, err
}

func store(stmt *sql.Stmt, pkg pkgmnt.Package) error {

	platform, err := json.Marshal(pkg.Platforms)

	if err != nil {
		return err
	}

	dependencies, err := json.Marshal(pkg.Dependencies)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(pkg.Name, pkg.Version, pkg.Architecture, pkg.Description, platform, dependencies)

	return err
}

func rowsMapper(rows *sql.Rows) (pkgmnt.Packages, error) {
	var packages pkgmnt.Packages

	// name, version, architecture, description, platforms, dependencies

	for rows.Next() {
		var pkg pkgmnt.Package
		//http://go-database-sql.org/nulls.html
		var platforms sql.NullString
		var dependencies sql.NullString

		err := rows.Scan(&pkg.Name, &pkg.Version, &pkg.Architecture, &pkg.Description, &platforms, &dependencies)

		if err != nil {
			return packages, err
		}

		if platforms.Valid {
			var pl []pkgmnt.Platform
			err = json.Unmarshal([]byte(platforms.String), &pl)

			if err != nil {
				return packages, err
			}

			pkg.Platforms = pl
		}

		if dependencies.Valid {
			var deps []pkgmnt.Package

			if len(dependencies.String) > 2 {
				err = json.Unmarshal([]byte(dependencies.String), deps)

				if err != nil {
					return packages, err
				}

				pkg.Dependencies = deps
			}
		}

		packages.Add(pkg)
	}

	return packages, nil
}
