package cmd

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/SuperTapood/Flint/core/generated/general"
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

		if len(changed) == 0 {
			return
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

func prettyChangeDiff(conn general.ConnectionType, stackMetadata map[string]any, changeset []map[string]map[string]any) {
	for _, change := range changeset {
		newObj := change["new"]
		name := conn.PrettyName(newObj, stackMetadata)

		printColored(colorYellow, "[~] %s\n", name)

		NewDiff(change["old"], change["new"])
	}
}

func toString(v any) string {
	return fmt.Sprintf("%v", v)
}

type difference struct {
	Old string
	New string
}

func NewDiff(old, new map[string]any) {
	// Estimate capacity based on input size
	estimatedSize := len(old) + len(new)
	oldFlat := make(map[string]any, estimatedSize)
	newFlat := make(map[string]any, estimatedSize)

	walk(old, nil, oldFlat)
	walk(new, nil, newFlat)

	// Single pass to find all differences
	dontMatch := make(map[string]difference, estimatedSize/2)

	// Check all keys from both maps
	for k, v := range oldFlat {
		if newVal, exists := newFlat[k]; !exists || toString(newVal) != toString(v) {
			dontMatch[k] = difference{
				Old: toString(v),
				New: toString(newVal),
			}
		}
	}

	for k, v := range newFlat {
		if _, exists := oldFlat[k]; !exists {
			dontMatch[k] = difference{
				Old: "",
				New: toString(v),
			}
		}
	}

	// Rebuild nested structure
	output := make(map[string]any, len(dontMatch))
	for k, v := range dontMatch {
		keys := strings.Split(k, ".")
		slices.Reverse(keys)
		addToMap(output, keys, v)
	}

	printColoredln(unchagedColor, "{")
	printDiff(output, nil)
	printColoredln(unchagedColor, "}")
}

func printColoredln(color string, a ...any) {
	if !noColor {
		fmt.Printf(color)
		defer fmt.Printf(colorReset)
	}

	fmt.Println(a...)
}

func printDiff(v any, path []string) {
	switch val := v.(type) {

	case map[string]any:
		for k, vv := range val {
			printColoredln(unchagedColor, strings.Repeat("    ", len(path)+1), `"`+toString(k)+`": {`)
			printDiff(vv, append(path, k))
			printColoredln(unchagedColor, strings.Repeat("    ", len(path)+1), "}")
		}
	case map[string]string:
		for k, vv := range val {
			printColoredln(unchagedColor, strings.Repeat("    ", len(path)+1), `"`+toString(k)+`": {`)
			printDiff(vv, append(path, k))
			printColoredln(unchagedColor, strings.Repeat("    ", len(path)+1), "}")
		}

	case []any:
		for i, vv := range val {
			printColoredln(unchagedColor, strings.Repeat("    ", len(path)+1), `"`+toString(i)+`": {`)
			printDiff(vv, append(path, fmt.Sprintf("%d", i)))
			printColoredln(unchagedColor, strings.Repeat("    ", len(path)+1), "}")
		}

	case difference:
		printColoredln(colorGreen, "[+]", strings.Repeat("    ", len(path)), `"`+toString(val.New)+`",`)
		printColoredln(colorRed, "[-]", strings.Repeat("    ", len(path)), `"`+toString(val.Old)+`",`)
	}
}

func addToMap(obj map[string]any, keys []string, value difference) {
	if len(keys) == 1 {
		obj[keys[0]] = value
		return
	}

	current_key := keys[len(keys)-1]

	_, ok := obj[current_key]
	if !ok {
		obj[keys[len(keys)-1]] = make(map[string]any, 0)
	}

	addToMap(obj[current_key].(map[string]any), keys[:len(keys)-1], value)
}

func walk(v any, path []string, flat map[string]any) {
	switch val := v.(type) {

	case map[string]any:
		for k, vv := range val {
			switch va := vv.(type) {
			case []any:
				for i := range va {
					walk(vv.([]any)[i], append(path, k+"["+toString(i)+"]"), flat)
				}
			default:
				walk(vv, append(path, k), flat)
			}
		}
	case map[string]string:
		for k, vv := range val {
			walk(vv, append(path, k), flat)
		}

	case []any:
		fmt.Println("we shouldn't be here")
		os.Exit(2)

	default:
		flat[strings.Join(path, ".")] = val
	}
}
