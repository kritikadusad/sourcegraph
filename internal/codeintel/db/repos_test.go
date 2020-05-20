package db

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sourcegraph/sourcegraph/internal/db/dbconn"
	"github.com/sourcegraph/sourcegraph/internal/db/dbtesting"
)

func TestRepoName(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	dbtesting.SetupGlobalTestDB(t)
	db := testDB()

	if _, err := dbconn.Global.Exec(`INSERT INTO repo (id, name) VALUES (50, 'github.com/foo/bar')`); err != nil {
		t.Fatalf("unexpected error inserting repo: %s", err)
	}

	name, err := db.RepoName(context.Background(), 50)
	if err != nil {
		t.Fatalf("unexpected error getting repo name: %s", err)
	}
	if name != "github.com/foo/bar" {
		t.Errorf("unexpected repo name. want=%s have=%s", "github.com/foo/bar", name)
	}
}

func TestRepoIDs(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	dbtesting.SetupGlobalTestDB(t)
	db := testDB()

	for id, name := range map[int]string{1: "foo", 2: "baz", 3: "bonk"} {
		if _, err := dbconn.Global.Exec(`INSERT INTO repo (id, name) VALUES ($1, $2)`, id, name); err != nil {
			t.Fatalf("unexpected error inserting repo: %s", err)
		}
	}

	ids, err := db.RepoIDs(context.Background(), []string{"foo", "bar", "baz"})
	if err != nil {
		t.Fatalf("unexpected error getting repo ids: %s", err)
	}

	expectedIDs := map[string]int{
		"foo": 1,
		"baz": 2,
	}
	if diff := cmp.Diff(expectedIDs, ids); diff != "" {
		t.Errorf("unexpected ids (-want +got):\n%s", diff)
	}
}
