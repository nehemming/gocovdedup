// Package main is the main entry point for the program.
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"

	"golang.org/x/tools/cover"
)

var errHelp = errors.New(`usage: gocovdedup [<file1> <file2> ... <fileN>|-]
files must be in go cover format or if '-' is supplied then read from stdin`)

func processArgs(args []string, stdIn io.Reader) ([]*cover.Profile, error) {
	var profiles []*cover.Profile
	switch len(args) {
	case 1:
		return nil, errHelp
	default:
		var files []string
		readStdin := false
		for _, arg := range args[1:] {
			if arg == "-" {
				readStdin = true
			} else {
				files = append(files, arg)
			}
		}

		if readStdin {
			stdInProfiles, err := cover.ParseProfilesFromReader(stdIn)
			if err != nil {
				return nil, err
			}
			profiles = append(profiles, stdInProfiles...)
		}

		fileProfiles, err := loadProfilesForFiles(files)
		if err != nil {
			return nil, err
		}

		profiles = append(profiles, fileProfiles...)
	}

	return profiles, nil
}

func loadProfilesForFiles(files []string) ([]*cover.Profile, error) {
	profiles := []*cover.Profile{}
	for _, file := range files {
		profile, err := cover.ParseProfiles(file)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile...)
	}
	return profiles, nil
}

type orderedBlocks []cover.ProfileBlock

func (b orderedBlocks) Len() int      { return len(b) }
func (b orderedBlocks) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b orderedBlocks) Less(i, j int) bool {
	if b[i].StartLine < b[j].StartLine {
		return true
	}

	if b[i].StartLine == b[j].StartLine {
		if b[i].StartCol < b[j].StartCol {
			return true
		}
		if b[i].StartCol == b[j].StartCol {
			if b[i].EndLine < b[j].EndLine {
				return true
			}
			if b[i].EndLine == b[j].EndLine {
				if b[i].EndCol < b[j].EndCol {
					return true
				}
			}
		}
	}

	return false
}

func maxEndLine(b1, b2 *cover.ProfileBlock) (int, int) {
	if b1.EndLine > b2.EndLine {
		return b1.EndLine, b1.EndCol
	}
	if b2.EndLine > b1.EndLine {
		return b2.EndLine, b2.EndCol
	}
	if b1.EndCol > b2.EndCol {
		return b1.EndLine, b1.EndCol
	}
	return b2.EndLine, b2.EndCol
}

func combine(profiles []*cover.Profile) map[string]*cover.Profile {
	fileMap := make(map[string]*cover.Profile)

	// combine all blocks by file name
	for _, profile := range profiles {
		if p, found := fileMap[profile.FileName]; found {
			p.Blocks = append(p.Blocks, profile.Blocks...)
		} else {
			fileMap[profile.FileName] = profile
		}
	}

	return fileMap
}

type byFileName []*cover.Profile

func (p byFileName) Len() int           { return len(p) }
func (p byFileName) Less(i, j int) bool { return p[i].FileName < p[j].FileName }
func (p byFileName) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func overlaps(b1, b2 *cover.ProfileBlock) bool {
	// return true if b2 overlaps b1
	if b1.EndLine < b2.StartLine {
		return false
	}
	if b1.EndLine > b2.StartLine {
		return true
	}
	return b1.EndCol >= b2.StartCol
}

func deDuplicate(profiles []*cover.Profile) []*cover.Profile {
	combined := combine(profiles)

	result := make([]*cover.Profile, 0, len(combined))

	// dedup blocks
	for _, profile := range combined {
		sort.Sort(orderedBlocks(profile.Blocks))

		// dedup blocks
		deduped := make([]cover.ProfileBlock, 0, len(profile.Blocks))
		var current *cover.ProfileBlock
		for _, blockIterator := range profile.Blocks {
			block := blockIterator
			if current == nil {
				current = &block
				continue
			}
			if overlaps(current, &block) {
				// merge blocks as overlapping
				current.EndLine, current.EndCol = maxEndLine(current, &block)
				current.NumStmt = max(current.NumStmt, block.NumStmt)
				current.Count = max(current.Count, block.Count)
			} else {
				deduped = append(deduped, *current)
				current = &block
			}
		}
		if current != nil {
			deduped = append(deduped, *current)
		}

		profile.Blocks = deduped
		result = append(result, profile)
	}

	sort.Sort(byFileName(result))
	return result
}

func printProfile(profile *cover.Profile, w io.Writer) {
	// name.go:line.column,line.column numberOfStatements count
	name := profile.FileName
	for _, block := range profile.Blocks {
		fmt.Fprintf(w, "%s:%d.%d,%d.%d %d %d\n", name, block.StartLine, block.StartCol, block.EndLine, block.EndCol, block.NumStmt, block.Count)
	}
}

func printProfiles(profiles []*cover.Profile, w io.Writer) {
	if len(profiles) > 0 {
		fmt.Fprintf(w, "mode: %s\n", profiles[0].Mode)
	}

	for _, profile := range profiles {
		printProfile(profile, w)
	}
}

func checkError(err error, w io.Writer, exit func(code int)) {
	if err != nil {
		if errors.Is(err, errHelp) {
			fmt.Fprintln(w, err)
			exit(99)
			return
		}

		fmt.Fprintln(w, err)
		exit(1)
	}
}

func main() {
	profiles, err := processArgs(os.Args, os.Stdin)
	checkError(err, os.Stderr, os.Exit)
	printProfiles(deDuplicate(profiles), os.Stdout)
}
