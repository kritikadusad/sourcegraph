package db

import (
	"context"

	"github.com/keegancsmith/sqlf"
)

// RepoName returns the name for the repo with the given identifier.
func (db *dbImpl) RepoName(ctx context.Context, repositoryID int) (string, error) {
	name, exists, err := scanFirstString(db.query(
		ctx,
		sqlf.Sprintf(`SELECT name FROM repo WHERE id = %s`, repositoryID),
	))
	if err != nil {
		return "", err
	}
	if !exists {
		return "", ErrUnknownRepository
	}
	return name, nil
}

// RepoIDs returns the identifiers for the repos with the given names.
func (db *dbImpl) RepoIDs(ctx context.Context, names []string) (map[string]int, error) {
	if len(names) == 0 {
		return nil, nil
	}

	var qs []*sqlf.Query
	for _, name := range names {
		qs = append(qs, sqlf.Sprintf("%s", name))
	}

	return scanRepoIDs(db.query(
		ctx,
		sqlf.Sprintf(`SELECT name, id FROM repo WHERE name IN (%s)`, sqlf.Join(qs, ",")),
	))
}
