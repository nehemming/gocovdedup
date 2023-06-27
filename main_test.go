package main

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	"golang.org/x/tools/cover"
)

func TestPrintProfile_Format(t *testing.T) {
	// test cases for the printProfile function

	buf := bytes.Buffer{}
	profiles := []*cover.Profile{
		{
			FileName: "file1.go",
			Mode:     "set",
			Blocks: []cover.ProfileBlock{
				{
					StartLine: 1,
					StartCol:  2,
					EndLine:   3,
					EndCol:    4,
					NumStmt:   5,
					Count:     6,
				},
			},
		},
	}

	printProfile(profiles[0], &buf)

	actual := strings.Trim(buf.String(), "\n")

	expected := "file1.go:1.2,3.4 5 6"

	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestPrintProfiles_Format(t *testing.T) {
	buf := bytes.Buffer{}
	profiles := []*cover.Profile{
		{
			FileName: "file1.go",
			Mode:     "count",
			Blocks: []cover.ProfileBlock{
				{
					StartLine: 1,
					StartCol:  2,
					EndLine:   3,
					EndCol:    4,
					NumStmt:   5,
					Count:     6,
				},
				{
					StartLine: 11,
					StartCol:  12,
					EndLine:   13,
					EndCol:    14,
					NumStmt:   15,
					Count:     16,
				},
			},
		},
		{
			FileName: "file2.go",
			Mode:     "set",
			Blocks: []cover.ProfileBlock{
				{
					StartLine: 1,
					StartCol:  2,
					EndLine:   3,
					EndCol:    4,
					NumStmt:   5,
					Count:     6,
				},
				{
					StartLine: 21,
					StartCol:  22,
					EndLine:   23,
					EndCol:    24,
					NumStmt:   25,
					Count:     26,
				},
			},
		},
	}

	expected := []string{
		"mode: count",
		"file1.go:1.2,3.4 5 6",
		"file1.go:11.12,13.14 15 16",
		"file2.go:1.2,3.4 5 6",
		"file2.go:21.22,23.24 25 26",
	}

	var actual []string

	printProfiles(profiles, &buf)

	for reader := bufio.NewReader(&buf); ; {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		actual = append(actual, strings.Trim(line, "\n"))
	}

	if len(expected) != len(actual) {
		t.Errorf("expected %d lines, got %d\n%v", len(actual), len(expected), actual)
	}

	for i, expect := range expected {
		if actual[i] != expect {
			t.Errorf("%d expected %s, got %s", i+1, expect, actual[i])
		}
	}
}

func TestCheckError(t *testing.T) {
	count := 0
	testCases := []struct {
		name   string
		err    error
		output string
		fn     func(i int)
	}{
		{"nil", nil, "", func(i int) {
			t.Error("should not be called")
		}},
		{"help", errHelp, `usage: gocovdedup [<file1> <file2> ... <fileN>|-]
files must be in go cover format or if '-' is supplied then read from stdin`, func(i int) {
			if i != 99 {
				t.Errorf("expected 99, got %d", i)
			}
			count++
		}},
		{"general", errors.New("general"), "general", func(i int) {
			if i != 1 {
				t.Errorf("expected 1, got %d", i)
			}
			count++
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			count = 0
			var buf bytes.Buffer
			checkError(tc.err, &buf, tc.fn)

			actual := strings.Trim(buf.String(), "\n")
			if actual != tc.output {
				t.Errorf("expected %s, got %s", tc.output, actual)
			}
			if tc.err != nil && count != 1 {
				t.Errorf("expected 1, got %d", count)
			}
		})
	}
}

func TestMax(t *testing.T) {
	testCases := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"a", 1, 0, 1},
		{"b", 0, 1, 1},
		{"equal", 1, 1, 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := max(tc.a, tc.b)
			if actual != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, actual)
			}
		})
	}
}

func TestByFileName(t *testing.T) {
	// test byFileName sorting

	testCases := []struct {
		name     string
		order    []*cover.Profile
		expected []*cover.Profile
	}{
		{"empty", []*cover.Profile{}, []*cover.Profile{}},
		{"one", []*cover.Profile{{FileName: "a"}}, []*cover.Profile{{FileName: "a"}}},
		{"two", []*cover.Profile{{FileName: "b"}, {FileName: "a"}}, []*cover.Profile{{FileName: "a"}, {FileName: "b"}}},
		{"three", []*cover.Profile{{FileName: "b"}, {FileName: "c"}, {FileName: "a"}}, []*cover.Profile{{FileName: "a"}, {FileName: "b"}, {FileName: "c"}}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sort.Sort(byFileName(tc.order))
			if !reflect.DeepEqual(tc.order, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, tc.order)
			}
		})
	}
}

func TestCombine(t *testing.T) {
	// test function combine
	testCases := []struct {
		name     string
		order    []*cover.Profile
		expected map[string]*cover.Profile
	}{
		{"empty", []*cover.Profile{}, map[string]*cover.Profile{}},
		{"one", []*cover.Profile{{FileName: "a"}}, map[string]*cover.Profile{"a": {FileName: "a"}}},
		{"two", []*cover.Profile{{FileName: "b"}, {FileName: "a"}}, map[string]*cover.Profile{"a": {FileName: "a"}, "b": {FileName: "b"}}},
		{"three", []*cover.Profile{{FileName: "b"}, {FileName: "c"}, {FileName: "a"}}, map[string]*cover.Profile{"a": {FileName: "a"}, "b": {FileName: "b"}, "c": {FileName: "c"}}},
		{
			"blocks",
			[]*cover.Profile{
				{FileName: "b", Blocks: []cover.ProfileBlock{{StartLine: 1}}},
				{FileName: "b", Blocks: []cover.ProfileBlock{{StartLine: 2}}},
				{FileName: "a"},
			},
			map[string]*cover.Profile{
				"a": {FileName: "a"},
				"b": {FileName: "b", Blocks: []cover.ProfileBlock{{StartLine: 1}, {StartLine: 2}}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := combine(tc.order)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func TestMaxEndLine(t *testing.T) {
	// test function maxEndLine
	testCases := []struct {
		name      string
		a, b      cover.ProfileBlock
		line, col int
	}{
		{
			name: "equal",
			a:    cover.ProfileBlock{EndLine: 1, EndCol: 1},
			b:    cover.ProfileBlock{EndLine: 1, EndCol: 1},
			line: 1,
			col:  1,
		},
		{
			name: "disjoint",
			a:    cover.ProfileBlock{EndLine: 2, EndCol: 25},
			b:    cover.ProfileBlock{EndLine: 1, EndCol: 1},
			line: 2,
			col:  25,
		},
		{
			name: "overlapping",
			a:    cover.ProfileBlock{EndLine: 1, EndCol: 25},
			b:    cover.ProfileBlock{EndLine: 7, EndCol: 1},
			line: 7,
			col:  1,
		},
		{
			name: "same line",
			a:    cover.ProfileBlock{EndLine: 7, EndCol: 25},
			b:    cover.ProfileBlock{EndLine: 7, EndCol: 1},
			line: 7,
			col:  25,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			line, col := maxEndLine(&tc.a, &tc.b)
			if line != tc.line || col != tc.col {
				t.Errorf("expected %d:%d, got %d:%d", tc.line, tc.col, line, col)
			}
		})
	}
}

func TestOrderedBlocks(t *testing.T) {
	// test sorting of orderedBlocks
	testCases := []struct {
		name     string
		blocks   []cover.ProfileBlock
		expected []cover.ProfileBlock
	}{
		{
			name: "basic",
			blocks: []cover.ProfileBlock{
				{StartLine: 17, StartCol: 76, EndLine: 19, EndCol: 22, NumStmt: 2, Count: 0},
			},
			expected: []cover.ProfileBlock{
				{StartLine: 17, StartCol: 76, EndLine: 19, EndCol: 22, NumStmt: 2, Count: 0},
			},
		},
		{
			name: "overlapping",
			blocks: []cover.ProfileBlock{
				{StartLine: 17, StartCol: 76, EndLine: 19, EndCol: 22, NumStmt: 2, Count: 0},
				{StartLine: 16, StartCol: 76, EndLine: 17, EndCol: 22, NumStmt: 2, Count: 0},
			},
			expected: []cover.ProfileBlock{
				{StartLine: 16, StartCol: 76, EndLine: 17, EndCol: 22, NumStmt: 2, Count: 0},
				{StartLine: 17, StartCol: 76, EndLine: 19, EndCol: 22, NumStmt: 2, Count: 0},
			},
		},
		{
			name: "same start",
			blocks: []cover.ProfileBlock{
				{StartLine: 17, StartCol: 76, EndLine: 19, EndCol: 22, NumStmt: 2, Count: 0},
				{StartLine: 17, StartCol: 32, EndLine: 17, EndCol: 22, NumStmt: 2, Count: 0},
			},
			expected: []cover.ProfileBlock{
				{StartLine: 17, StartCol: 32, EndLine: 17, EndCol: 22, NumStmt: 2, Count: 0},
				{StartLine: 17, StartCol: 76, EndLine: 19, EndCol: 22, NumStmt: 2, Count: 0},
			},
		},
		{
			name: "same start col",
			blocks: []cover.ProfileBlock{
				{StartLine: 17, StartCol: 32, EndLine: 19, EndCol: 22, NumStmt: 2, Count: 0},
				{StartLine: 17, StartCol: 32, EndLine: 17, EndCol: 22, NumStmt: 2, Count: 0},
			},
			expected: []cover.ProfileBlock{
				{StartLine: 17, StartCol: 32, EndLine: 17, EndCol: 22, NumStmt: 2, Count: 0},
				{StartLine: 17, StartCol: 32, EndLine: 19, EndCol: 22, NumStmt: 2, Count: 0},
			},
		},
		{
			name: "same end line",
			blocks: []cover.ProfileBlock{
				{StartLine: 17, StartCol: 32, EndLine: 17, EndCol: 23, NumStmt: 2, Count: 0},
				{StartLine: 17, StartCol: 32, EndLine: 17, EndCol: 22, NumStmt: 2, Count: 0},
			},
			expected: []cover.ProfileBlock{
				{StartLine: 17, StartCol: 32, EndLine: 17, EndCol: 22, NumStmt: 2, Count: 0},
				{StartLine: 17, StartCol: 32, EndLine: 17, EndCol: 23, NumStmt: 2, Count: 0},
			},
		},
		{
			name: "same end line and col",
			blocks: []cover.ProfileBlock{
				{StartLine: 17, StartCol: 32, EndLine: 17, EndCol: 23, NumStmt: 2, Count: 0},
				{StartLine: 17, StartCol: 32, EndLine: 17, EndCol: 23, NumStmt: 2, Count: 0},
			},
			expected: []cover.ProfileBlock{
				{StartLine: 17, StartCol: 32, EndLine: 17, EndCol: 23, NumStmt: 2, Count: 0},
				{StartLine: 17, StartCol: 32, EndLine: 17, EndCol: 23, NumStmt: 2, Count: 0},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := orderedBlocks(tc.blocks)
			sort.Sort(actual)
			if len(actual) != len(tc.expected) {
				t.Fatal("len wrong", tc.expected, actual)
			}
			for i, v := range tc.expected {
				if v != actual[i] {
					t.Error("order fail", i, v, actual[i])
				}
			}
		})
	}
}

func sp(s string) *string {
	return &s
}

func ps(sp *string) string {
	if sp == nil {
		return ""
	}
	return *sp
}

func TestLoadProfilesForFiles(t *testing.T) {
	// test loadProfilesForFiles
	testCases := []struct {
		name     string
		files    []string
		expected []*cover.Profile
		err      *string
	}{
		{
			name:  "missing",
			files: []string{"testdata/notfound.out"},
			err:   sp("open testdata/notfound.out: no such file or directory"),
		},
		{
			name:     "one",
			files:    []string{"testdata/cover_1.out"},
			expected: newProfileOne(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := loadProfilesForFiles(tc.files)
			if err != nil {
				if tc.err == nil || err.Error() != *tc.err {
					t.Fatalf("unexpected:\n%s\n%s\n", ps(tc.err), err)
				}
			} else if tc.err != nil {
				t.Fatal("expected", err)
			}
			if len(actual) != len(tc.expected) {
				t.Errorf("len wrong %v %+v\n", tc.expected, actual)
			}

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("wrong %v %+v\n", tc.expected, actual[0])
			}
		})
	}
}

func TestProcessArgsFiles(t *testing.T) {
	args := []string{
		"app",
		"testdata/cover_1.out",
		"testdata/cover_2.out",
	}
	p, err := processArgs(args, nil)
	if err != nil {
		t.Error("unexpected err", err)
	}
	if len(p) != 2 {
		t.Error("profiles len != 2", len(p))
	}
}

func TestProcessArgsFilesNone(t *testing.T) {
	args := []string{
		"app",
	}
	p, err := processArgs(args, nil)
	if !errors.Is(err, errHelp) {
		t.Error("unexpected not errHelp", err)
	}

	if len(p) != 0 {
		t.Error("profiles len != 0", len(p))
	}
}

func TestProcessArgsStdIn(t *testing.T) {
	f, err := os.Open("testdata/cover_2.out")
	if err != nil {
		t.Fatal("open file", err)
	}
	defer f.Close()
	args := []string{
		"app",
		"-",
	}

	p, err := processArgs(args, f)
	if err != nil {
		t.Error("unexpected err", err)
	}
	if len(p) != 1 {
		t.Error("profiles len != 1", len(p))
	}
}

func TestOverlaps(t *testing.T) {
	// test overlaps
	testCases := []struct {
		name     string
		b1       cover.ProfileBlock
		b2       cover.ProfileBlock
		expected bool
	}{
		{
			name:     "none",
			expected: true,
		},
		{
			name:     "same",
			b1:       cover.ProfileBlock{StartLine: 1, StartCol: 1, EndLine: 1, EndCol: 1},
			b2:       cover.ProfileBlock{StartLine: 1, StartCol: 1, EndLine: 1, EndCol: 1},
			expected: true,
		},
		{
			name:     "adjacent",
			b1:       cover.ProfileBlock{StartLine: 1, StartCol: 1, EndLine: 1, EndCol: 1},
			b2:       cover.ProfileBlock{StartLine: 1, StartCol: 2, EndLine: 1, EndCol: 2},
			expected: false,
		},
		{
			name:     "overlap",
			b1:       cover.ProfileBlock{StartLine: 1, StartCol: 1, EndLine: 1, EndCol: 2},
			b2:       cover.ProfileBlock{StartLine: 1, StartCol: 2, EndLine: 1, EndCol: 3},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := overlaps(&tc.b1, &tc.b2)
			if actual != tc.expected {
				t.Error("overlaps fail", tc.expected, actual)
			}
		})
	}
}

func TestDeDuplicate(t *testing.T) {
	disjoint := append(newProfileOne(), newDisjoint()...)
	disjointExpected := append(newCombined("github.com/repo/gocovdedup/alt.go"), newCombined("")...)
	overlapped := append(newProfileOne(), newProfileOne()...)
	// Test the deduplication logic
	testCases := []struct {
		name     string
		profiles []*cover.Profile
		expected []*cover.Profile
	}{
		{
			name:     "empty",
			expected: []*cover.Profile{},
		},
		{
			name:     "combined",
			profiles: newProfileOne(),
			expected: newCombined(""),
		},
		{
			name:     "disjoint",
			profiles: disjoint,
			expected: disjointExpected,
		},
		{
			name:     "overlapped",
			profiles: overlapped,
			expected: newCombined(""),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := deDuplicate(tc.profiles)
			if len(actual) != len(tc.expected) {
				t.Fatal("len wrong", tc.expected, actual)
			}
			for i, v := range tc.expected {
				if !reflect.DeepEqual(v, actual[i]) {
					t.Errorf("order fail %d\n%v\n%v", i, v, actual[i])
				}
			}
		})
	}
}
