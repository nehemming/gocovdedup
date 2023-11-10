package main

import (
	"fmt"
	"os"

	"github.com/denormal/go-gitignore"
	"golang.org/x/tools/cover"
)

// filterProfiles is used to filter profiles using a gitignore style exclusion file
// if the file is not found the input list ius returned unaltered.
func filterProfiles(profiles []*cover.Profile, ignoreFile string) ([]*cover.Profile, error) {
	if _, err := os.Stat(ignoreFile); err != nil || len(profiles) == 0 {
		return profiles, nil // no include file.
	}

	// ignore, err := createIgnore(ignoreFile)
	ignore, err := gitignore.NewFromFile(ignoreFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read ignore file:%s", err)
	}

	filtered := make([]*cover.Profile, 0, len(profiles))
	for _, p := range profiles {
		m := ignore.Relative(p.FileName, false)
		if m == nil {
			filtered = append(filtered, p)
		}
	}

	return filtered, nil
}
