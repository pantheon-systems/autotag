package autotag_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pantheon-systems/autotag"
)

func newRepo(t *testing.T) autotag.GitRepo {
	repo := createTestRepo(t)

	seedTestRepo(t, repo)
	r, err := autotag.NewRepo(repo.Path())
	if err != nil {
		t.Fatal("Error creating repo", err)
	}

	repo.Free() // we don't need to hold this repo open it was just to make the path
	return *r
}

// adds a #major comit to the repo
func newRepoMajor(t *testing.T) autotag.GitRepo {
	repo := createTestRepo(t)

	seedTestRepo(t, repo)
	updateReadme(t, repo, "#major change")

	r, err := autotag.NewRepo(repo.Path())
	if err != nil {
		t.Fatal("Error creating repo", err)
	}

	repo.Free() // we don't need to hold this repo open it was just to make the path
	return *r
}

func TestMajor(t *testing.T) {
	r := newRepo(t)
	defer cleanupTestRepo(t, r.Repo)

	majorTag(t, r.Repo)
	v, err := r.MajorBump()
	if err != nil {
		t.Fatal("MajorBump failed: ", err)
	}

	if v.String() != "2.0.1" {
		fmt.Printf("MajorBump failed expected '2.0.1' got '%s' ", v)
	}

	fmt.Printf("Major is now %s", v)
}

func TestMinor(t *testing.T) {
	r := newRepo(t)
	defer cleanupTestRepo(t, r.Repo)

	majorTag(t, r.Repo)
	v, err := r.MinorBump()
	if err != nil {
		t.Fatal("MinorBump failed: ", err)
	}

	if v.String() != "1.1.0" {
		fmt.Printf("MinorBump failed expected '1.1.0' got '%s' ", v)
	}
}

func TestPatch(t *testing.T) {
	r := newRepo(t)
	defer cleanupTestRepo(t, r.Repo)

	majorTag(t, r.Repo)
	v, err := r.PatchBump()
	if err != nil {
		t.Fatal("PatchBump failed: ", err)
	}

	if v.String() != "1.0.2" {
		t.Fatalf("PatchBump failed expected '1.0.2' got '%s' ", v)
	}
}

func TestAutoTag(t *testing.T) {
	r := newRepo(t)
	defer cleanupTestRepo(t, r.Repo)

	err := r.AutoTag()
	if err != nil {
		t.Fatal("AutoTag failed ", err)
	}

	tags, err := r.Repo.Tags.List()
	expect := []string{"v1.0.1", "v1.0.2"}

	if !reflect.DeepEqual(expect, tags) {
		t.Fatalf("AutoBump expected '%+v' got '%+v'", expect, tags)
	}
}

func TestAutoTagCommits(t *testing.T) {
	r := newRepoMajor(t)
	defer cleanupTestRepo(t, r.Repo)

	err := r.AutoTag()
	if err != nil {
		t.Fatal("AutoTag failed ", err)
	}

	tags, err := r.Repo.Tags.List()
	expect := []string{"v1.0.1", "v2.0.1"}

	if !reflect.DeepEqual(expect, tags) {
		t.Fatalf("AutoBump expected '%+v' got '%+v'", expect, tags)
	}
}
