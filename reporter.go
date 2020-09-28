package main

import (
	"fmt"
	"strings"
)

const (
	deletedColor = "\033[1;31m%s\033[0m\n"
	addedColor   = "\033[1;36m%s\033[0m\n"

	defaultDistance = 10
)

type Option interface {
	optionType() string
}

type HideMatchedLines struct {
	Distance int
}

func (hide HideMatchedLines) optionType() string {
	return "hide-matches"
}

type HideLineNumbers struct {
}

func (hide HideLineNumbers) optionType() string {
	return "hide-line-numbers"
}


func PrintDiff(diff Diff, options ...Option) {
	var hideDistance int
	var hideLineNumbers bool

	for _, option := range options {
		switch option.optionType() {
		case HideMatchedLines{}.optionType():
			if option.(HideMatchedLines).Distance > 0 {
				hideDistance = option.(HideMatchedLines).Distance
			} else {
				hideDistance = defaultDistance
			}

		case HideLineNumbers{}.optionType():
			hideLineNumbers = true
		}
	}

	for _, value := range diff.Blocks {
		if value.IsEqual() {
			if hideDistance > 0 && len(value.(EqualBlock).Equals) > 2 * hideDistance{
				for i := 0; i < hideDistance; i++ {
					if hideLineNumbers {
						fmt.Printf(" %s\n", value.(EqualBlock).Equals[i].Value)
					} else {
						fmt.Printf(" [%4d] [%4d] : %s\n", value.(EqualBlock).Equals[i].LeftLine + 1, value.(EqualBlock).Equals[i].RightLine + 1, value.(EqualBlock).Equals[i].Value)
					}
				}
				fmt.Println(" " + strings.Repeat(".", 120))

				for i := len(value.(EqualBlock).Equals) - hideDistance; i < len(value.(EqualBlock).Equals); i++ {
					if hideLineNumbers {
						fmt.Printf(" %s\n", value.(EqualBlock).Equals[i].Value)
					} else {
						fmt.Printf(" [%4d] [%4d] : %s\n", value.(EqualBlock).Equals[i].LeftLine + 1, value.(EqualBlock).Equals[i].RightLine + 1, value.(EqualBlock).Equals[i].Value)
					}
				}
			} else {
				for _, equal := range value.(EqualBlock).Equals {
					if hideLineNumbers {
						fmt.Printf(" %s\n", equal.Value)
					} else {
						fmt.Printf(" [%4d] [%4d] : %s\n", equal.LeftLine + 1, equal.RightLine + 1, equal.Value)
					}
				}
			}
		} else {
			//fmt.Printf("Diff counts left: %d, right: %d \n", len(value.(DiffBlock).LeftChanges), len(value.(DiffBlock).RightChanges))
			for _, change := range value.(DiffBlock).LeftChanges {
				if hideLineNumbers {
					fmt.Printf(deletedColor, fmt.Sprintf(" %s", change.Value))
				} else {
					fmt.Printf(deletedColor, fmt.Sprintf(" [%4d] [    ] : %s", change.Line + 1, change.Value))
				}
			}
			for _, change := range value.(DiffBlock).RightChanges {
				if hideLineNumbers {
					fmt.Printf(addedColor, fmt.Sprintf(" %s", change.Value))
				} else {
					fmt.Printf(addedColor, fmt.Sprintf(" [    ] [%4d] : %s", change.Line + 1, change.Value))
				}
			}
		}
	}
}

