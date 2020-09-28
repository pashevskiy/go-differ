package main

import (
	"math"
	"strings"
)

type Diff struct {
	Blocks []Block
}

// Block is common interface fo EqualBlock and DiffBlock
type Block interface {
	// Returns true on EqualBlock, false on DiffBlock
	IsEqual() bool
}

type EqualBlock struct {
	Equals []Equal
}

func (block EqualBlock) IsEqual() bool {
	return true
}

type DiffBlock struct {
	LeftChanges  []Change
	RightChanges []Change
}

func (block DiffBlock) IsEqual() bool {
	return false
}

type Equal struct {
	LeftLine  int
	RightLine int
	Value     string
}

type Change struct {
	Line  int
	Value string
}

func NewDiffLines(left, right string) Diff {
	var rawBlocks []Block

	leftLines := strings.Split(strings.TrimSuffix(left, "\n"), "\n")

	rightLines := strings.Split(strings.TrimSuffix(right, "\n"), "\n")

	for leftIndex, leftLine := range leftLines {
		for rightIndex, rightLine := range rightLines {
			if leftLine == rightLine {
				matchIsValid := checkAndCleanLessFitItems(&rawBlocks, leftIndex, rightIndex)
				if matchIsValid {
					addMatchToBlocks(&rawBlocks, leftIndex, rightIndex, leftLine)
					break
				}
			}
		}
	}

	// Remove empty equal blocks
	var blocks []Block

	for index, block := range rawBlocks {
		if block.IsEqual() {
			if len(block.(EqualBlock).Equals) > 0 {
				blocks = append(blocks, block)
			}
		} else {
			blocks = append(blocks, block)
		}

		rawBlocks[index] = nil
	}

	//Filling in diffs
	for index, block := range blocks {
		if !block.IsEqual() {
			var startLeft, startRight, endLeft, endRight int

			if index == 0 {
				startLeft = 0
				startRight = 0
			} else {
				startLeft = blocks[index-1].(EqualBlock).Equals[len(blocks[index-1].(EqualBlock).Equals)-1].LeftLine + 1
				startRight = blocks[index-1].(EqualBlock).Equals[len(blocks[index-1].(EqualBlock).Equals)-1].RightLine + 1
			}

			if index == len(blocks)-1 {
				endLeft = len(leftLines) - 1
				endRight = len(rightLines) - 1
			} else {
				endLeft = blocks[index+1].(EqualBlock).Equals[0].LeftLine
				endRight = blocks[index+1].(EqualBlock).Equals[0].RightLine
			}

			for i := startLeft; i < endLeft; i++ {
				blocks[index] = DiffBlock{LeftChanges: append(blocks[index].(DiffBlock).LeftChanges, Change{
					Line:  i,
					Value: leftLines[i],
				}),
					RightChanges: []Change{}}
			}

			for i := startRight; i < endRight; i++ {
				blocks[index] = DiffBlock{LeftChanges: blocks[index].(DiffBlock).LeftChanges,
					RightChanges: append(blocks[index].(DiffBlock).RightChanges, Change{
						Line:  i,
						Value: rightLines[i],
					})}
			}
		}
	}

	return Diff{Blocks: blocks}
}

// Just removes an equal item from the slice with saving order
func removeItemFromEquals(slice []Equal, index int) []Equal {
	copy(slice[index:], slice[index+1:])
	slice[len(slice)-1] = Equal{}
	slice = slice[:len(slice)-1]

	return slice
}

// checkAndCleanLessFitItems function checks all possible collisions (cases where matched line is already used or there are matches where right line > mathedLine)
// then if those ones were found - function decides which match is more fit (existing or current), removes existing ones ore returns false to block adding current match into the blocks
func checkAndCleanLessFitItems(blocks *[]Block, curentLine, matchedLine int) bool {
	for i, block := range *blocks {
		if block.IsEqual() {
			for j, equal := range block.(EqualBlock).Equals {
				if equal.RightLine >= matchedLine {
					if math.Abs(float64(equal.LeftLine-equal.RightLine)) > math.Abs(float64(curentLine-matchedLine)) {
						(*blocks)[i] = EqualBlock{removeItemFromEquals(block.(EqualBlock).Equals, j)}
					} else {
						return false
					}
				}
			}
		}
	}

	return true
}

// addMatchToBlocks adds new match to blocks, if previous match lines are missing - adds an empty diff block before creating new equal block
func addMatchToBlocks(blocks *[]Block, left, right int, value string) {
	if len(*blocks) > 0 && (*blocks)[len(*blocks)-1].IsEqual() {
		if len(*blocks) > 0 && (left-(*blocks)[len(*blocks)-1].(EqualBlock).Equals[len((*blocks)[len(*blocks)-1].(EqualBlock).Equals)-1].LeftLine > 1 ||
			right-(*blocks)[len(*blocks)-1].(EqualBlock).Equals[len((*blocks)[len(*blocks)-1].(EqualBlock).Equals)-1].RightLine > 1) {
			*blocks = append(*blocks, DiffBlock{}, EqualBlock{[]Equal{{
				LeftLine:  left,
				RightLine: right,
				Value:     value,
			}}})
		} else {
			(*blocks)[len(*blocks)-1] = EqualBlock{
				Equals: append((*blocks)[len(*blocks)-1].(EqualBlock).Equals, Equal{
					LeftLine:  left,
					RightLine: right,
					Value:     value,
				}),
			}
		}
	} else {
		*blocks = append(*blocks, EqualBlock{[]Equal{{
			LeftLine:  left,
			RightLine: right,
			Value:     value,
		}}})
	}
}
