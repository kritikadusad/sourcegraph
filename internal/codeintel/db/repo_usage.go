package db

import (
	"context"
	"sort"

	"github.com/keegancsmith/sqlf"
)

// RepoUsageStatistics pairs a repository identifier with a count of code intelligence events.
type RepoUsageStatistics struct {
	RepositoryID int
	SearchCount  int
	PreciseCount int
}

// RepoUsageStatistics reads recent event log records and returns the number of search-based and precise
// code intelligence activity within the last week grouped by repository. The resulting slice is ordered
// by search then precise event counts.
func (db *dbImpl) RepoUsageStatistics(ctx context.Context) ([]RepoUsageStatistics, error) {
	countsByName, err := scanCountsByName(db.query(ctx, sqlf.Sprintf(`
		SELECT sub.* FROM (
			SELECT
				substring(url from '//[^/]+/(.+)/-/') AS repo_name,
				COUNT(*) FILTER (WHERE name LIKE 'codeintel.search%%%%') AS search_count,
				COUNT(*) FILTER (WHERE name LIKE 'codeintel.lsif%%%%') AS precise_count
			FROM event_logs
			WHERE timestamp >= NOW() - INTERVAL '1 week'
			GROUP BY repo_name
		) sub
		WHERE EXISTS (SELECT * FROM repo r WHERE r.uri = sub.repo_name)
	`)))
	if err != nil {
		return nil, err
	}

	var names []string
	for name := range countsByName {
		names = append(names, name)
	}

	ids, err := db.RepoIDs(ctx, names)
	if err != nil {
		return nil, err
	}

	var stats []RepoUsageStatistics
	for name, count := range countsByName {
		repositoryID, ok := ids[name]
		if !ok {
			continue
		}

		stats = append(stats, RepoUsageStatistics{
			RepositoryID: repositoryID,
			SearchCount:  count.SearchCount,
			PreciseCount: count.PreciseCount,
		})
	}

	sort.Slice(stats, func(i, j int) bool {
		comparisons := [2]int{
			stats[j].SearchCount - stats[i].SearchCount,
			stats[j].PreciseCount - stats[i].PreciseCount,
		}

		for _, cmp := range comparisons {
			if cmp != 0 {
				return cmp < 0
			}
		}

		return false
	})

	return stats, nil
}
