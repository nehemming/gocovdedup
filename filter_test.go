package main

import "testing"

func TestFilterNoFile(t *testing.T) {
	files := []string{"testdata/cover_1.out"}
	profiles, err := loadProfilesForFiles(files)
	if err != nil {
		t.Fatal("fatal profile read", err)
	}

	if len(profiles) != 1 {
		t.Fatal("profile len not 1", len(profiles))
	}

	ret, err := filterProfiles(profiles, "nofile")
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if len(ret) != len(profiles) {
		t.Fatal("len wrong", len(profiles), len(ret))
	}
}

func TestFilterIncludeAll(t *testing.T) {
	files := []string{"testdata/cover_1.out"}
	profiles, err := loadProfilesForFiles(files)
	if err != nil {
		t.Fatal("fatal profile read", err)
	}

	if len(profiles) != 1 {
		t.Fatal("profile len not 1", len(profiles))
	}

	ret, err := filterProfiles(profiles, "testdata/filter")
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if len(ret) != len(profiles) {
		t.Fatal("len wrong", len(profiles), len(ret))
	}
}

func TestFilterIncludeNoAlt(t *testing.T) {
	files := []string{"testdata/cover_multi.out"}
	profiles, err := loadProfilesForFiles(files)
	if err != nil {
		t.Fatal("fatal profile read", err)
	}

	if len(profiles) != 3 {
		t.Fatal("profile len not 3", len(profiles))
	}

	ret, err := filterProfiles(profiles, "testdata/filter")
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if len(ret) != 1 {
		t.Fatal("len wrong alt still in?", len(profiles), len(ret))
	}
}
