package pgsql

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/jmoiron/sqlx"
)

func EnsureDB(db *sqlx.DB, assets http.FileSystem) error {
	dbdrv, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	s, err := Source{
		Assets: assets,
	}.Open("")
	if err != nil {
		return err
	}
	migrator, err := migrate.NewWithInstance("vfs", s, "postgres", dbdrv)
	if err != nil {
		return err
	}

	return migrator.Up()
}

type Source struct {
	Assets     http.FileSystem
	opened     http.File
	migrations *source.Migrations
}

func (s Source) Close() error {
	return s.opened.Close()
}

var ErrAppend = errors.New("append to migrations")

func (s Source) Open(url string) (source.Driver, error) {
	ns := s
	var err error
	ns.opened, err = s.Assets.Open("/")
	if err != nil {
		return nil, err
	}
	files, err := ns.opened.Readdir(-1)
	if err != nil {
		return nil, err
	}
	ns.migrations = source.NewMigrations()
	for _, f := range files {
		if !f.IsDir() {
			m, err := source.DefaultParse(f.Name())
			if err != nil {
				continue
			}
			if !ns.migrations.Append(m) {
				return nil, ErrAppend
			}
		}
	}
	return ns, nil
}

func (s Source) First() (version uint, err error) {
	if v, ok := s.migrations.First(); ok {
		return v, nil
	}
	return 0, &os.PathError{Op: "first", Path: "/", Err: os.ErrNotExist}
}

func (s Source) Prev(version uint) (prevVersion uint, err error) {
	if v, ok := s.migrations.Prev(version); ok {
		return v, nil
	}

	return 0, &os.PathError{Op: fmt.Sprintf("prev for version %v", version), Path: "/", Err: os.ErrNotExist}
}

func (s Source) Next(version uint) (nextVersion uint, err error) {
	if v, ok := s.migrations.Next(version); ok {
		return v, nil
	}

	return 0, &os.PathError{Op: fmt.Sprintf("next for version %v", version), Path: "/", Err: os.ErrNotExist}
}

func (s Source) ReadUp(version uint) (r io.ReadCloser, identifier string, err error) {
	if m, ok := s.migrations.Up(version); ok {
		r, err := s.Assets.Open(path.Join("/", m.Raw))
		if err != nil {
			return nil, "", err
		}
		return r, m.Identifier, nil
	}
	return nil, "", &os.PathError{Op: fmt.Sprintf("read version %v", version), Path: "/", Err: os.ErrNotExist}
}

func (s Source) ReadDown(version uint) (r io.ReadCloser, identifier string, err error) {
	if m, ok := s.migrations.Down(version); ok {
		r, err := s.Assets.Open(path.Join("/", m.Raw))
		if err != nil {
			return nil, "", err
		}
		return r, m.Identifier, nil
	}
	return nil, "", &os.PathError{Op: fmt.Sprintf("read version %v", version), Path: "/", Err: os.ErrNotExist}
}
