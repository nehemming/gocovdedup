package main

import "golang.org/x/tools/cover"

func newProfileOne() []*cover.Profile {
	return []*cover.Profile{
		{
			FileName: "github.com/repo/gocovdedup/main.go",
			Mode:     "set",
			Blocks: []cover.ProfileBlock{
				{ // 17->19
					StartLine: 17,
					StartCol:  76,
					EndLine:   19,
					EndCol:    22,
					NumStmt:   2,
					Count:     0,
				},
				{ // 20->21
					StartLine: 20,
					StartCol:  9,
					EndLine:   21,
					EndCol:    22,
					NumStmt:   1,
					Count:     0,
				},
				{ // 22->25
					StartLine: 22,
					StartCol:  10,
					EndLine:   25,
					EndCol:    35,
					NumStmt:   3,
					Count:     0,
				},
				{ // 25 ->26
					StartLine: 25,
					StartCol:  35,
					EndLine:   26,
					EndCol:    18,
					NumStmt:   1,
					Count:     0,
				},
			},
		},
	}
}

func newCombined(filename string) []*cover.Profile {
	if filename == "" {
		filename = "github.com/repo/gocovdedup/main.go"
	}
	return []*cover.Profile{
		{
			FileName: filename,
			Mode:     "set",
			Blocks: []cover.ProfileBlock{
				{ // 17->19
					StartLine: 17,
					StartCol:  76,
					EndLine:   19,
					EndCol:    22,
					NumStmt:   2,
					Count:     0,
				},
				{ // 20->21
					StartLine: 20,
					StartCol:  9,
					EndLine:   21,
					EndCol:    22,
					NumStmt:   1,
					Count:     0,
				},
				{ // 22->26
					StartLine: 22,
					StartCol:  10,
					EndLine:   26,
					EndCol:    18,
					NumStmt:   3,
					Count:     0,
				},
			},
		},
	}
}

func newDisjoint() []*cover.Profile {
	return []*cover.Profile{
		{
			FileName: "github.com/repo/gocovdedup/alt.go",
			Mode:     "set",
			Blocks: []cover.ProfileBlock{
				{
					StartLine: 17,
					StartCol:  76,
					EndLine:   19,
					EndCol:    22,
					NumStmt:   2,
					Count:     0,
				},
				{
					StartLine: 20,
					StartCol:  9,
					EndLine:   21,
					EndCol:    22,
					NumStmt:   1,
					Count:     0,
				},
				{
					StartLine: 22,
					StartCol:  10,
					EndLine:   25,
					EndCol:    35,
					NumStmt:   3,
					Count:     0,
				},
				{
					StartLine: 25,
					StartCol:  35,
					EndLine:   26,
					EndCol:    18,
					NumStmt:   1,
					Count:     0,
				},
			},
		},
	}
}
