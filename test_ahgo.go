/*
 *	File:       src/github.com/Ken1JF/test_ahgo.go
 *	Project:	abst-hier
 *
 *	Created by Ken Friedenbach on 12/18/2009
 *	Copyright:	2009-2014, Ken Friedenbach
 *	All rights reserved.
 *
 *	Test the ah, sgf, and sgfdb packages.
 *
 *	To build and run:
 *      see the buildScript.bash
 *
 *	To test output:
 *		diff test_ahgo_new.txt test_ahgo_out.txt
 *	Depending on which options are set:
 *	Only differences should be times, and order of directory count completions.
 *
 *	When new functionality is added and output is verified:
 *		rm test_ahgo_out.txt
 *		cp test_ahgo_new.txt test_ahgo_out.txt
 *	build and run using buildScript.bash
 *		diff test_ahgo_new.txt test_ahgo_out.txt
 */

package main

import (
	"flag"
	"fmt"
	"github.com/Ken1JF/ah"
	"github.com/Ken1JF/sgf"
	"github.com/Ken1JF/sgfdb"
	"io"
	"io/ioutil"
	"runtime"
	//	"strings"
	"time"
)

const SGF_GEN_GO_VERSION = "1.0 (update AbstHier working, one level. Built with Go version 1.2 Generate whole board patterns...)"

// SGF Specification file is copied from a different project: Projects/GenSGFProperties
const defaultSGFSpecFile = "src/github.com/Ken1JF/sgf/sgf_properties_spec.txt"

var SGFSpecFile = defaultSGFSpecFile

// Some test and print controls.
// default values are controlled by init()
// final values are set by program arguments

var doReadDatabaseAndBuild bool
var doReadTeachingGames bool
var doReadWriteFuseki bool

// Test Controls:

// SGF Database and testout dirs:
const defaultDatabaseDir = "/usr/local/GoGoD/Go/Database/"

var DatabaseDir = defaultDatabaseDir

const defaultBoardPatternsDir = "/Users/ken/Projects/OpenLikeAPro/Library/Board/"

var BoardPatternsDir = defaultBoardPatternsDir

const defaultTeachingDir = "/Users/ken/Documents/GO/Games/teaching/"

var TeachingDir = defaultTeachingDir

const defaultTeachingPatternsDir = "/Users/ken/Documents/GO/Games/teaching/patterns/"

var TeachingPatternsDir = defaultTeachingPatternsDir

const defaultFusekiFileName = "./Fuseki.sgf"

var fusekiFileName = defaultFusekiFileName

const defaultOutFusekiFileName = "./Fuseki2.sgf"

var outFusekiFileName = defaultOutFusekiFileName

var fileLimit int = 0 // no fileLimit
var moveLimit = 0     // no moveLimit
var patternLimit = 0  // no patternLimit
var skipFiles = 0     // no skipping

type ProcessOptions uint8

const (
	RemoveLabels ProcessOptions = 1 << iota
	AddAllLabels
)

var sgfProcessOptions ProcessOptions = 0

var removeLabels = false
var allLabels = false

func init() {

	flag.IntVar(&moveLimit, "ml", 0, "ml = move limit. limit the number of moves read each .sgf file, 0 means no limit")
	flag.IntVar(&fileLimit, "fl", 0, "fl = file limit. limit the number of .sgf files read from a directory, 0 means no limit")
	flag.IntVar(&patternLimit, "pl", 0, "pl = pattern limit. limit the depth of pattern storing, 0 means no limit")
	flag.IntVar(&skipFiles, "sf", 0, "sf = skip files. skip this number of .sgf files before reading from a directory, 0 means no skip")

	flag.BoolVar(&doReadDatabaseAndBuild, "rdab", false, "rdab = do Read Database And Build patterns, false (default) means do not do Read Database And Build patterns.")
	flag.BoolVar(&doReadTeachingGames, "rtg", false, "rtg = do Read Teaching Games, false (default) means do not do read teaching games.")
	flag.BoolVar(&doReadWriteFuseki, "rwf", false, "rwf = do Read Write Fuseki, false (default) means do not read and write Fuseki file.")
	flag.BoolVar(&removeLabels, "rl", false, "rl = remove labels, false (default) means do not remove labels from Fuseki file.")
	flag.BoolVar(&allLabels, "al", false, "al = all labels, false (default) means do generate all labels in Fuseki file.")

	flag.StringVar(&SGFSpecFile, "ssf", defaultSGFSpecFile, "path to the SGF properties specification file.")
	flag.StringVar(&DatabaseDir, "dbdir", defaultDatabaseDir, "path to the Database directory.")
	flag.StringVar(&BoardPatternsDir, "bpdir", defaultBoardPatternsDir, "path to the Board Patterns directory.")
	flag.StringVar(&TeachingDir, "tdir", defaultTeachingDir, "path to teaching games directory.")
	flag.StringVar(&TeachingPatternsDir, "tpdir", defaultTeachingPatternsDir, "path to the teaching patterns directory.")

	flag.StringVar(&fusekiFileName, "ffn", defaultFusekiFileName, "path to the Fuseki file.")
	flag.StringVar(&outFusekiFileName, "offn", defaultOutFusekiFileName, "path to the output Fuseki file.")
}

func PrintOptionsSet() {
	// Print changes to the options:
	if allLabels {
		fmt.Printf("al, all labels has value %v\n", allLabels)
	}
	if BoardPatternsDir != defaultBoardPatternsDir {
		fmt.Printf("bpdir, board patterns directory has value \"%s\"\n", BoardPatternsDir)
	}
	if DatabaseDir != defaultDatabaseDir {
		fmt.Printf("dbdir, data base directory has value \"%s\"\n", DatabaseDir)
	}
	if fusekiFileName != defaultFusekiFileName {
		fmt.Printf("ffn, Fuseki file name has value \"%s\"\n", fusekiFileName)
	}
	if fileLimit != 0 {
		fmt.Printf("fl, file limit has value %d\n", fileLimit)
	}
	if moveLimit != 0 {
		fmt.Printf("ml, move limit has value %d\n", moveLimit)
	}
	if outFusekiFileName != defaultOutFusekiFileName {
		fmt.Printf("offn, output Fuseki file name has value \"%s\"\n", outFusekiFileName)
	}
	if patternLimit != 0 {
		fmt.Printf("pl, pattern limit has value %d\n", patternLimit)
	}
	if doReadDatabaseAndBuild {
		fmt.Printf("rdab, do read database and build has value %t\n", doReadDatabaseAndBuild)
	}
	if removeLabels {
		fmt.Printf("rl, removeLabels has value %t\n", removeLabels)
	}
	if doReadTeachingGames {
		fmt.Printf("rtg, do read teaching games has value %t\n", doReadTeachingGames)
	}
	if doReadWriteFuseki {
		fmt.Printf("rwf, do read write Fuseki file has value %t\n", doReadWriteFuseki)
	}
	if skipFiles != 0 {
		fmt.Printf("sf, skip files has value %d\n", skipFiles)
	}
	if SGFSpecFile != defaultSGFSpecFile {
		fmt.Printf("ssf, SGF specification file has value \"%s\"\n", SGFSpecFile)
	}
	if TeachingDir != defaultTeachingDir {
		fmt.Printf("tdir, teaching directory has value \"%s\"\n", TeachingDir)
	}
	if TeachingPatternsDir != defaultTeachingPatternsDir {
		fmt.Printf("tpdir, teaching patterns directory has value \"%s\"\n", TeachingPatternsDir)
	}
}

func main() {
	fmt.Printf("Program to generate opening pattern libraries:\n Version %s\n", SGF_GEN_GO_VERSION)
	nCPUs := runtime.NumCPU()
	if sgfdb.TheDBReadReq.DoMultiCPU {
		oldMaxProcs := runtime.GOMAXPROCS(nCPUs)
		fmt.Printf(" num CPUs = %d, default max Procs was %d, now set to num CPUs\n\n", nCPUs, oldMaxProcs)
	} else {
		fmt.Printf(" num CPUs = %d, but multi-processing not enabled.\n\n", nCPUs)
	}
	start := time.Now()

	flag.Parse()

	flag.Usage()

	// Set the sgfProcessOptions based on Boolean Flags
	if removeLabels {
		sgfProcessOptions = sgfProcessOptions + RemoveLabels
	}
	if allLabels {
		sgfProcessOptions = sgfProcessOptions + AddAllLabels
	}

	PrintOptionsSet()

	// do not ask for verification of SGF Specification file,
	// or ask for verbose output. These are done in sgf_test.go
	// If that test is ok, then the file is ok.
	err := sgf.SetupSGFProperties(SGFSpecFile, false, false)

	stop := time.Now()
	fmt.Printf("All tests took %v to run.\n", stop.Sub(start))

	if err == 0 { // don't try these tests if SGF Setup failed.
		if doReadDatabaseAndBuild {
			buildStat := sgfdb.ReadDatabaseAndBuildPatterns(DatabaseDir, BoardPatternsDir, ah.WHOLE_BOARD_PATTERN, fileLimit, moveLimit, skipFiles)

			if buildStat > 0 {
				fmt.Printf("Errors during Build Patterns, status = %d.\n", buildStat)
			}
		}

		if doReadTeachingGames {
			teachStat := sgfdb.ReadTeachingDirectory(TeachingDir, TeachingPatternsDir, fileLimit, moveLimit, patternLimit, skipFiles)

			if teachStat > 0 {
				fmt.Printf("Errors while Reading TeachingDir, status = %d.\n", teachStat)
			}
		}

		if doReadWriteFuseki {
			fusekiFile, err := ioutil.ReadFile(fusekiFileName)
			if err != nil && err != io.EOF {
				fmt.Printf("Error reading teaching Fuseki file: %s, %s\n", fusekiFileName, err)
			} else {
				prsr, errL := sgf.ParseFile(fusekiFileName, fusekiFile,
					sgf.ParseComments+sgf.ParserGoGoD+sgf.ParserPlay, moveLimit)
				if len(errL) != 0 {
					fmt.Printf("Error %s during parsing: %s\n", errL.Error(), fusekiFileName)
				} else {
					//TODO: add error reporting? ErrorList return value?
					if (sgfProcessOptions & RemoveLabels) > 0 {
						fmt.Println("Removing labels from:", fusekiFileName, "to", outFusekiFileName)
						prsr.GameTree.BreadthFirstTraverse(true, sgf.DoRemoveLabels)
						prsr.GameTree.ReportDeletedProperties()
					}
					if (sgfProcessOptions & AddAllLabels) > 0 {
						fmt.Println("Adding labels to:", outFusekiFileName, "from", fusekiFileName)
						prsr.GameTree.DepthFirstTraverse(true, sgf.DoAddLabels)
						fmt.Println("The number of added labels =", sgf.NumberOfAddedLabels)
					}
					err = prsr.GameTree.WriteFile(outFusekiFileName, 1)
					if err != nil {
						fmt.Printf("Error writing: %s, %s\n", outFusekiFileName, err)
					} else {
						// Build ZCode mapping of unique board positions and transformations:
						// TODO: err = prsr.GameTree.BuildFusekiTable()
					}
				}
			}
		}

	}

	finish := time.Now()
	fmt.Printf("Complete run took %v to run.\n", finish.Sub(start))
}
