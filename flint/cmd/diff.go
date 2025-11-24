package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/SuperTapood/Flint/core/generated/common"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/spf13/cobra"
)

func printColored(color string, format string, a ...any) {
	if !noColor {
		fmt.Printf(color)
		defer fmt.Printf(colorReset)
	}

	fmt.Fprintf(os.Stdout, format, a...)
}

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "display the difference between the given stack and the existing one (if exists)",
	Long:  `display the difference between the given stack and the existing one (if exists)`,
	Run: func(cmd *cobra.Command, args []string) {
		stack, conn, stackName := StackConnFromApp()
		revision := conn.GetCurrentRevision(stackName)
		fmt.Println("generating changeset for stack '" + stackName + "' (" + strconv.Itoa(revision) + " -> " + strconv.Itoa(revision+1) + "):")
		_, objMap := stack.GetActual().Synth(stackName)

		added, removed, changed := conn.Diff(objMap, stack.GetActual().GetMetadata(), stackName)
		if len(added) == 0 && len(removed) == 0 && len(changed) == 0 {
			fmt.Println("empty changeset nothing to do")
			return
		}
		for _, add := range added {
			printColored(colorGreen, "[+] %s\n", add)
		}
		for _, rem := range removed {
			printColored(colorRed, "[-] %s\n", rem)
		}
		prettyChangeDiff(conn.GetActual(), stack.GetActual().GetMetadata(), changed)
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
	diffCmd.Flags().SortFlags = false

	diffCmd.Flags().StringVarP(&app, "app", "a", "", "the app to diff ")
	diffCmd.MarkFlagRequired("app")
	diffCmd.Flags().StringVarP(&dir, "dir", "d", ".", "the directory to run the app at")
	diffCmd.Flags().BoolVarP(&noColor, "no-color", "c", false, "turn off diff coloring")
}

func prettyChangeDiff(conn common.ConnectionType, stackMetadata map[string]any, changeset []map[string]map[string]any) {
	for _, change := range changeset {
		newObj := change["new"]
		name := conn.PrettyName(newObj, stackMetadata)
		newBytes, err := json.MarshalIndent(newObj, " ", "\t")
		if err != nil {
			panic(err)
		}

		oldObj := change["old"]
		oldBytes, err := json.MarshalIndent(oldObj, " ", "\t")
		if err != nil {
			panic(err)
		}

		diff := difflib.UnifiedDiff{
			A:       difflib.SplitLines(string(oldBytes)),
			B:       difflib.SplitLines(string(newBytes)),
			Context: len(difflib.SplitLines(string(newBytes))),
		}

		printColored(colorYellow, "[~] %s", name)

		PrintCDKDiff(diff)
	}
}

// PrintCDKDiff formats and prints a difflib.UnifiedDiff in a style
// similar to 'cdk diff', complete with ANSI color codes.
// It writes the formatted output to the provided io.Writer.
func PrintCDKDiff(diff difflib.UnifiedDiff) {
	// Create a matcher to compare the two string slices
	m := difflib.NewMatcher(diff.A, diff.B)

	// Use a default context if not provided (or negative)
	context := diff.Context
	if context < 0 {
		context = 3 // A common default context
	}

	// Get the operations grouped by hunks
	groups := m.GetGroupedOpCodes(context)

	// If there are no diffs, we're done after the headers
	if len(groups) == 0 {
		return
	}

	// Iterate over each group (hunk) of changes
	for _, group := range groups {

		// Iterate over each operation within the hunk
		for _, code := range group {
			switch code.Tag {
			case 'e': // 'equal' - context line
				// Lines from A[code.I1:code.I2] are the same in B
				for _, line := range diff.A[code.I1:code.I2] {
					// fmt.Fprintf(w, "%s  %s%s", unchagedColor, line, colorReset) // Indent context lines
					printColored(unchagedColor, " %s", line)
				}
			case 'd': // 'delete' - line removed from 'A'
				for _, line := range diff.A[code.I1:code.I2] {
					//fmt.Fprintf(w, "%s[-] %s%s", colorRed, line, colorReset)
					printColored(colorRed, "[-] %s", line)
				}
			case 'i': // 'insert' - line added to 'B'
				for _, line := range diff.B[code.J1:code.J2] {
					//fmt.Fprintf(w, "%s[+] %s%s", colorGreen, line, colorReset)
					printColored(colorGreen, "[+] %s", line)
				}
			case 'r': // 'replace' - lines from 'A' replaced by lines in 'B'
				// Show as a deletion followed by an insertion
				for _, line := range diff.A[code.I1:code.I2] {
					//fmt.Fprintf(w, "%s[-] %s%s", colorRed, line, colorReset)
					printColored(colorRed, "[-] %s", line)
				}
				for _, line := range diff.B[code.J1:code.J2] {
					//fmt.Fprintf(w, "%s[+] %s%s", colorGreen, line, colorReset)
					printColored(colorGreen, "[+] %s", line)
				}
			}
		}
	}
}
