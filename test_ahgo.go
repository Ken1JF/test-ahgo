/*
 *	File:       src/gitHub.com/Ken1JF/ahgo/test_ahgo.go
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
 *		mv test_ahgo_new.txt test_ahgo_out.txt
 *	build and run using buildScript.bash
 *		diff test_ahgo_new.txt test_ahgo_out.txt
 */

package main

import (
	"flag"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
	//	"sort"
    "runtime"
	"gitHub.com/Ken1JF/ahgo/ah"
	"gitHub.com/Ken1JF/ahgo/sgf"
	"gitHub.com/Ken1JF/ahgo/sgfdb"
	"fmt"

//	"unsafe"
)

const SGF_GEN_GO_VERSION = "1.0 (update AbstHier working, one level. Built with Go version 1.2 Generate whole board patterns...)"

    // SGF Specification file is copied from a different project: Projects/GenSGFProperties
const defaultSGFSpecFile = "/Users/ken/Projects/abst-hier/src/gitHub.com/Ken1JF/ahgo/sgf_properties_spec.txt"
var SGFSpecFile = defaultSGFSpecFile

    // Some DEBUG (print) controls:
var doAllTests = false

var doPrintSizes = false
var doPrintConstants = false
var doPrintHandicaps = false
var doPrintDirections = false // requires doPrintConstants to be true
var doPrintSGFProperties = false
var doPrintZKeys = false
var doVerifySGFPropertyOrder = false
var doSmallSGFTests = false
var doTransTest = false
var doCountMoves = false
var doReadWriteDatabase = false
var doReadDatabaseAndBuild = false
var doReadTeachingGames = false
var doReadWriteFuseki = false

    // count should be 56 with all tests included. One test currently in "WorkOnLater"
const SmallSGFTestOutputVerified int = 54  // controls the printing of last graph
const SmallSGFTestStringsVerified int = 55 // controls tracing

const defaultSmallSGFTestDir = "/Users/ken/Documents/GO/Tests/sgf_tests/"
var SmallSGFTestDir = defaultSmallSGFTestDir

    // Set to empty string to suppress writing output:
const defaultSmallSGFTestOutDir = "/Users/ken/Documents/GO/Tests/sgf_tests/goTestOut/"
var SmallSGFTestOutDir = defaultSmallSGFTestOutDir


// Test Controls:

// SGF Database and testout dirs:
const defaultDatabaseDir = "/Users/ken/Documents/GO/GoGoD/Go/Database/"
var DatabaseDir = defaultDatabaseDir
const defaultTestOutDir = "/Users/ken/Documents/GO/GoGoD/testout/"
var TestOutDir = defaultTestOutDir
const defaultBoardPatternsDir = "/Users/ken/Projects/OpenLikeAPro/Library/Board/"
var BoardPatternsDir = defaultBoardPatternsDir
const defaultTeachingDir = "/Users/ken/Documents/GO/Games/teaching/"
var TeachingDir = defaultTeachingDir
const defaultTeachingPatternsDir = "/Users/ken/Documents/GO/Games/teaching/patterns/"
var TeachingPatternsDir = defaultTeachingPatternsDir
const defaultFusekiFileName = "/Users/ken/ahgo/Fuseki.sgf"
var fusekiFileName = defaultFusekiFileName
const defaultOutFusekiFileName = "/Users/ken/ahgo/Fuseki2.sgf"
var outFusekiFileName = defaultOutFusekiFileName

var fileLimit int = 0 // no fileLimit
var moveLimit = 0     // no moveLimit
var patternLimit = 0  // no patternLimit
var skipFiles = 0   // no skipping

type ProcessOptions  uint8

const (
       RemoveLabels    ProcessOptions  = 1 << iota
       AddAllLabels    
       )


var sgfProcessOptions ProcessOptions = 0

var removeLabels = false
var allLabels = false

func ReadSmallSGFTests(dir string, outDir string) {
    fmt.Println("Reading Small SGF Tests, dir = ", dir,", outDir = ", outDir);
	dirFiles, err := ioutil.ReadDir(dir)
	if err != nil && err != io.EOF {
		fmt.Println("Can't read test directory: ", dir)
		return
	}
	count := 0
	for _, f := range dirFiles {
		if strings.Index(f.Name(), ".sgf") >= 0 {
			count += 1
			fmt.Println("Processing: ", f.Name())
			if count > SmallSGFTestStringsVerified {
				ah.SetAHTrace(true)
				fmt.Println("Tracing", f.Name())
			}
			fileName := dir + f.Name()
			b, err := ioutil.ReadFile(fileName)
			if err != nil && err != io.EOF {
				fmt.Println("Error reading file: ", fileName, err)
				return
			}
			//			prsr , errL := sgf.ParseFile(fileName, b, sgf.ParseComments, 0)
			prsr, errL := sgf.ParseFile(fileName, b, sgf.ParseComments+sgf.Play, 0)
			if len(errL) != 0 {
				fmt.Println("Error while parsing: ", fileName, ", ", errL.Error())
				return
			}
			if outDir != "" {
				outFileName := outDir + "/" + f.Name()
				err = prsr.GameTree.WriteFile(outFileName, sgfdb.SGFDB_NUM_PER_LINE)
				if err != nil {
					fmt.Printf("Error writing: %s, %s\n", outFileName, err)
				}
			}
			if count > SmallSGFTestOutputVerified {
				prsr.PrintAbstHier(fileName, true)
			}
			ah.SetAHTrace(false)
		}
	}
}

// Transformation test data:
// For these tests, use char values instead of defined PointStatus values
//
var test_5 = []string{
	"1....",
	"2..x.",
	"3.+..",
	"4....",
	"5....",
}

var test_7 = []string{
	"1......",
	"2......",
	"3.+.+..",
	"4...x..",
	"5.+.+..",
	"6......",
	"7......",
}

var test_9 = []string{
	"1........",
	"2........",
	"3.+...+..",
	"4.....x..",
	"5........",
	"6........",
	"7.+...+..",
	"8........",
	"9........",
}

var test_11 = []string{
	"1..........",
	"2..........",
	"3.+.....+..",
	"4.......x..",
	"5..........",
	"6....+.....",
	"7..........",
	"8..........",
	"9.+.....+..",
	"A..........",
	"B..........",
}

var test_13 = []string{
	"1............",
	"2............",
	"3............",
	"4..+.....+X..",
	"5............",
	"6............",
	"7....+.......",
	"8............",
	"9............",
	"A..+.....+...",
	"B............",
	"C............",
	"D............",
}

var test_15 = []string{
	"1..............",
	"2..............",
	"3..............",
	"4..+.......+X..",
	"5..............",
	"6..............",
	"7..............",
	"8.....+........",
	"9..............",
	"A..............",
	"B..............",
	"C..+.......+...",
	"D..............",
	"E..............",
	"F..............",
}

var test_17 = []string{
	"1................",
	"2................",
	"3................",
	"4..+.........+X..",
	"5................",
	"6................",
	"7................",
	"8................",
	"9......+.........",
	"A................",
	"B................",
	"C................",
	"D................",
	"E..+.........+...",
	"F................",
	"G................",
	"H................",
}

var test_19 = []string{
	"1..................",
	"2..................",
	"3..................",
	"4..+.....+.....+X..",
	"5..................",
	"6..................",
	"7..................",
	"8..................",
	"9..................",
	"A..+.....+.....+...",
	"B..................",
	"C..................",
	"D..................",
	"E..................",
	"F..................",
	"G..+.....+.....+...",
	"H..................",
	"I..................",
	"J..................",
}

// printInitBoard prints the PointType values
// after a Board is initialized (via SetSize)
//
func printInitBoard(abhr *ah.AbstHier, title string) {

	//	Black_Occ_Pt:		"◉",	
	//	White_Occ_Pt:		"◎",	

	var c ah.ColValue
	var r ah.RowValue
	nCol, nRow := abhr.GetSize()
	fmt.Println(title, "Board", nCol, "by", nRow)
	for r = 0; r < nRow; r++ {
		for c = 0; c < nCol; c++ {
			bp := abhr.Graphs[ah.PointLevel].GetPoint(c, r)
			hs := bp.GetNodeHighState()
			if hs == uint16(ah.White) {
				fmt.Print("◎")
			} else if hs == uint16(ah.Black) {
				fmt.Print("◉")
			} else {
				fmt.Print(ah.PtTypeNames[bp.GetPointType()])
			}
		}
		fmt.Println(" ")
	}
}

// printInitBoard2 is equivalent to printInitBoard
// but uses the iteration function ah.EachNode
// and a literal func.
//
func printInitBoard2(abhr *ah.AbstHier) {
	var row ah.RowValue = 0
	nCol, nRow := abhr.GetSize()
	fmt.Println("Board", nCol, "by", nRow)
	abhr.EachNode(ah.PointLevel,
		func(brd *ah.Graph, nl ah.NodeLoc) {
			_, r := brd.Nodes[nl].GetPointColRow()
			if r != row {
				fmt.Println(" ")
				row = r
			}
			fmt.Print(ah.PtTypeNames[brd.Nodes[nl].GetPointType()])
		})
	fmt.Println(" ")
}

// differBrds checks the LowStates of the Nodes
// only suitable for special set boards
//
func differBrds(brd1, brd2 *ah.AbstHier) (ret bool) {
	var c ah.ColValue
	var r ah.RowValue
	nCol, nRow := brd1.GetSize()
	nCol2, nRow2 := brd2.GetSize()
	if (nCol != nCol2) || (nRow != nRow2) {
		ret = true
	} else {
		for r = 0; r < nRow; r++ {
			for c = 0; c < nCol; c++ {
				nl := ah.MakeNodeLoc(c, r)
				bp1 := &brd1.Graphs[ah.PointLevel].Nodes[nl]
				bp2 := &brd2.Graphs[ah.PointLevel].Nodes[nl]
				if bp1.GetNodeLowState() != bp2.GetNodeLowState() {
					ret = true
					break
				}
			}
		}
	}
	return ret
}

// checkHandicapBrds checks the LowStates of the Nodes
// only suitable for special set boards
//
func checkHandicapBrds(brd1, brd2 *ah.AbstHier) (ret bool) {
	var c ah.ColValue
	var r ah.RowValue
	nCol, nRow := brd1.GetSize()
	nCol2, nRow2 := brd2.GetSize()
	if (nCol != nCol2) || (nRow != nRow2) {
		ret = true
	} else {
		for r = 0; r < nRow; r++ {
			for c = 0; c < nCol; c++ {
				nl := ah.MakeNodeLoc(c, r)
				bp1 := &brd1.Graphs[ah.PointLevel].Nodes[nl]
				bp2 := &brd2.Graphs[ah.PointLevel].Nodes[nl]
				low1 := bp1.GetNodeLowState()
				low2 := bp2.GetNodeLowState()
				// check that both are occupied or unoccupied
				if ah.IsOccupied(ah.PointStatus(low1)) != ah.IsOccupied(ah.PointStatus(low2)) {
					ret = true
					break
				}
			}
		}
	}
	return ret
}

// Print the boards, after transformation
//
func printBrds(msg string, brd *ah.AbstHier, newBrd *ah.AbstHier, tName string) {
	var c ah.ColValue
	var r ah.RowValue
	nCol, nRow := brd.GetSize()
	fmt.Println("Board size", nCol, "by", nRow, "after", tName)
	for r = 0; r < nRow; r++ {
		for c = 0; c < nCol; c++ {
			bp := brd.Graphs[ah.PointLevel].GetPoint(c, r)
			ch := bp.GetNodeLowState()
			fmt.Printf("%c", byte(ch))
		}
		fmt.Print("     ")
		for c = 0; c < nCol; c++ {
			nbp := newBrd.Graphs[ah.PointLevel].GetPoint(c, r)
			ch := nbp.GetNodeLowState()
			fmt.Printf("%c", byte(ch))
		}
		fmt.Println(" ")
	}
}

// SetUpTestBoard stores the test data (string characters)
// in the Board as PointStatus information.
//
func SetUpTestBoard(N int, brd *ah.AbstHier, data *[]string) {
	for r := 0; r < N; r++ {
		for c := 0; c < N; c++ {
			brd.SetPoint(ah.MakeNodeLoc(ah.ColValue(c), ah.RowValue(r)), ah.PointStatus((*data)[r][c]))
		}
	}
}

// Eight boards of various sizes.
var brd_5, brd_7, brd_9, brd_11, brd_13, brd_15, brd_17, brd_19 *ah.AbstHier

// and an array to hold them.
var brds [8]*ah.AbstHier

// Test the transformation logic
//
func TestTrans() {
	// Set up the test data boards.
	var col ah.ColValue
	var row ah.RowValue
	for size := 5; size <= 19; size += 2 {
		switch size {
		case 5:
			col = 5
			row = 5
			brd_5 = brd_5.InitAbstHier(col, row, ah.StringLevel, true)
			//				ah.SetAHTrace(false)
			printInitBoard(brd_5, "Initial 5x5 Board")
			brd_5.PrintAbstHier("Initial 5x5 Board", true)
			SetUpTestBoard(size, brd_5, &test_5)
			brds[0] = brd_5
		case 7:
			col = 7
			row = 7
			//				brd_7 = new(ah.AbstHier)
			//				brd_7.SetSize(col, row)
			brd_7 = brd_7.InitAbstHier(col, row, ah.StringLevel, true)
			printInitBoard2(brd_7)
			SetUpTestBoard(size, brd_7, &test_7)
			brds[1] = brd_7
		case 9:
			col = 9
			row = 9
			//				brd_9 = new(ah.AbstHier)
			//				brd_9.SetSize(col, row)
			brd_9 = brd_9.InitAbstHier(col, row, ah.StringLevel, true)
			printInitBoard(brd_9, "Initial 9x9 Board")
			SetUpTestBoard(size, brd_9, &test_9)
			brds[2] = brd_9
		case 11:
			col = 11
			row = 11
			//				brd_11 = new(ah.AbstHier)
			//				brd_11.SetSize(col, row)
			brd_11 = brd_11.InitAbstHier(col, row, ah.StringLevel, true)
			printInitBoard2(brd_11)
			SetUpTestBoard(size, brd_11, &test_11)
			brds[3] = brd_11
		case 13:
			col = 13
			row = 13
			brd_13 = brd_13.InitAbstHier(col, row, ah.StringLevel, true)
			printInitBoard(brd_13, "Initial 13x13 Board")
			SetUpTestBoard(size, brd_13, &test_13)
			brds[4] = brd_13
		case 15:
			col = 15
			row = 15
			brd_15 = brd_15.InitAbstHier(col, row, ah.StringLevel, true)
			printInitBoard2(brd_15)
			SetUpTestBoard(size, brd_15, &test_15)
			brds[5] = brd_15
		case 17:
			col = 17
			row = 17
			brd_17 = brd_17.InitAbstHier(col, row, ah.StringLevel, true)
			printInitBoard(brd_17, "Initial 17x17 Board")
			SetUpTestBoard(size, brd_17, &test_17)
			brds[6] = brd_17
		case 19:
			col = 19
			row = 19
			brd_19 = brd_19.InitAbstHier(col, row, ah.StringLevel, true)
			printInitBoard2(brd_19)
			SetUpTestBoard(size, brd_19, &test_19)
			brds[7] = brd_19
		}
	}
	// Print each board, after applying one of the transformations,
	// and print it (for visual verification)
	//	ah.SetAHTrace(true) // trace first one
	for i, brd := range brds {
		fmt.Println("Checking brds[", i, "]")
		if brd == nil {
			fmt.Println("Error in setup: brd == nil")
		} else {
			newBrd := brd.TransBoard(ah.BoardTrans(i))
			printBrds("Visual Check", brd, newBrd, ah.TransName[i])
		}
		ah.SetAHTrace(false) // turn off after first one
	}
	// Verify that the inverse transformations produce the original
	for i, brd := range brds {
		t := ah.BoardTrans(i)
		inv := ah.InverseTrans[t]
		fmt.Println("Checking", ah.TransName[i], "and its inverse:", ah.TransName[inv])
		newBrd := brd.TransBoard(t)
		newBrdInv := newBrd.TransBoard(inv)
		if differBrds(brd, newBrdInv) {
			printBrds("Error: inverse differs", brd, newBrdInv, ah.TransName[i])
		}
	}
	// Verify the transformation composition table
	nxtBrd := 0 // used to pick the next board
	for A := ah.T_FIRST; A <= ah.T_LAST; A++ {
		for B := ah.T_FIRST; B <= ah.T_LAST; B++ {
			C := ah.ComposeTrans[A][B]
			fmt.Println("Checking", ah.TransName[C], " = ", ah.TransName[A], "*", ah.TransName[B])
			brd := brds[nxtBrd]
			nxtBrd++
			if nxtBrd >= 8 {
				nxtBrd = 0
			}
			brdA := brd.TransBoard(A)
			brdAB := brdA.TransBoard(B)
			brdC := brd.TransBoard(C)
			if differBrds(brdAB, brdC) {
				printBrds("Error: "+ah.TransName[ah.ComposeTrans[A][B]], brdAB, brdC,
					"not equal"+ah.TransName[A]+"*"+ah.TransName[B])
			}
		}
	}
}

// Print the Direction and static PointStatus values:
//
func printConst() {
	if doAllTests || doPrintDirections {
		// Check Direction definitions:
		// They should work as bit masks that can be ORed together
		for d := ah.NoDir; d <= ah.RightDir; d++ {
			switch d {
			case ah.NoDir:
				fmt.Println("NoDir:", d)
			case ah.UpperDir:
				fmt.Println("UpperDir:", d)
			case ah.LeftDir:
				fmt.Println("LeftDir:", d)
			case ah.LowerDir:
				fmt.Println("LowerDir:", d)
			case ah.RightDir:
				fmt.Println("RightDir:", d)
			default:
				//					fmt.Println(" skipping:", d)
			}
		}
	}

	// Check  PointType definitions:
loop:
	for pt := ah.SingletonPt; pt <= ah.Line_7_Pt; pt++ {
		switch pt {
		case ah.SingletonPt:
			fmt.Println("SingletonPt:", pt)

		case ah.LowerEndPt:
			fmt.Println("LowerEndPt:", pt)
		case ah.RightEndPt:
			fmt.Println("RightEndPt:", pt)
		case ah.UpperEndPt:
			fmt.Println("UpperEndPt:", pt)
		case ah.LeftEndPt:
			fmt.Println("LeftEndPt:", pt)

		case ah.UpperLowerBridgePt:
			fmt.Println("UpperLowerBridgePt:", pt)
		case ah.LeftRightBridgePt:
			fmt.Println("LeftRightBridgePt:", pt)

		case ah.UpperLeftCornerPt:
			fmt.Println("UpperLeftCornerPt:", pt)
		case ah.UpperRightCornerPt:
			fmt.Println("UpperRightCornerPt:", pt)
		case ah.LowerRightCornerPt:
			fmt.Println("LowerRightCornerPt:", pt)
		case ah.LowerLeftCornerPt:
			fmt.Println("LowerLeftCornerPt:", pt)

		case ah.UpperEdgePt:
			fmt.Println("UpperEdgePt:", pt)
		case ah.LeftEdgePt:
			fmt.Println("LeftEdgePt:", pt)
		case ah.LowerEdgePt:
			fmt.Println("LowerEdgePt:", pt)
		case ah.RightEdgePt:
			fmt.Println("RightEdgePt:", pt)

		case ah.CenterPt:
			fmt.Println("CenterPt:", pt)

		case ah.HoshiPt:
			fmt.Println("HoshiPt:", pt)

		case ah.Corner_2_2_Pt:
			fmt.Println("Corner_2_2_Pt:", pt)
		case ah.Line_2_Pt:
			fmt.Println("Line_2_Pt:", pt)
		case ah.Corner_3_3_Pt:
			fmt.Println("Corner_3_3_Pt:", pt)
		case ah.Line_3_Pt:
			fmt.Println("Line_3_Pt:", pt)
		case ah.Corner_4_4_Pt:
			fmt.Println("Corner_4_4_Pt:", pt)
		case ah.Line_4_Pt:
			fmt.Println("Line_4_Pt:", pt)
		case ah.Corner_5_5_Pt:
			fmt.Println("Corner_5_5_Pt:", pt)
		case ah.Line_5_Pt:
			fmt.Println("Line_5_Pt:", pt)
		case ah.Corner_6_6_Pt:
			fmt.Println("Corner_6_6_Pt:", pt)
		case ah.Line_6_Pt:
			fmt.Println("Line_6_Pt:", pt)
		case ah.Corner_7_7_Pt:
			fmt.Println("Corner_7_7_Pt:", pt)
		case ah.Line_7_Pt:
			fmt.Println("Line_7_Pt:", pt)
			//			case ah.Black_Occ_Pt:
			//				fmt.Println("Black_Occ_Pt:", pt)
			//			case ah.White_Occ_Pt:
			//				fmt.Println("White_Occ_Pt:", pt)
			break loop // loop will not terminate 255++ => 0

			//			case ah.UninitializedPt:
			//				fmt.Println("UninitializedPt:", pt)

		default:
			//				fmt.Println(" skipping:", pt)
		}
	}

	// Check  PointStatus definitions:
	for ps := ah.UndefinedPointStatus; ps <= ah.LastPointStatus; ps++ {
		switch ps {

		case ah.UndefinedPointStatus:
			fmt.Println("UndefinedPointStatus:", ps)

			// Occupied, stone color
		case ah.Black:
			fmt.Println("Black:", ps)
		case ah.White:
			fmt.Println("White:", ps)

			// Move types, for AB, AW, and AE properties:

		case ah.AB_U:
			fmt.Println("AB_U:", ps)
		case ah.AB_W:
			fmt.Println("AB_W:", ps)
		case ah.AE_B:
			fmt.Println("AE_B:", ps)
		case ah.AE_W:
			fmt.Println("AE_W:", ps)
		case ah.AW_B:
			fmt.Println("AW_B:", ps)
		case ah.AW_U:
			fmt.Println("AW_U:", ps)

			// Unoccupied, generic value
		case ah.Unocc:
			fmt.Println("Unocc:", ps)

			// No Adjacent Stones:
		case ah.B0W0:
			fmt.Println("B0W0:", ps)

			// Single Adjacent Stone:
		case ah.W1:
			fmt.Println("W1:", ps)
		case ah.B1:
			fmt.Println("B1:", ps)

			// Two Adjacent Stones:
		case ah.W2:
			fmt.Println("W2:", ps)
		case ah.B1W1:
			fmt.Println("B1W1:", ps)
		case ah.W1B1:
			fmt.Println("W1B1:", ps)
		case ah.B2:
			fmt.Println("B2:", ps)

			// Three Adjacent Stones:
		case ah.B3:
			fmt.Println("B3:", ps)
		case ah.B2W1:
			fmt.Println("B2W1:", ps)
		case ah.B1W2:
			fmt.Println("B1W2:", ps)
		case ah.W3:
			fmt.Println("W3:", ps)
		case ah.WBB:
			fmt.Println("WBB:", ps)
		case ah.WBW:
			fmt.Println("WBW:", ps)
		case ah.BWB:
			fmt.Println("BWB:", ps)
		case ah.W2B1:
			fmt.Println("W2B1:", ps)

			// Four Adjacent Stones:
		case ah.B4:
			fmt.Println("B4:", ps)
		case ah.B3W1:
			fmt.Println("B3W1:", ps)
		case ah.BWBW:
			fmt.Println("BWBW:", ps)
		case ah.BBWW:
			fmt.Println("BBWW:", ps)
		case ah.B1W3:
			fmt.Println("B1W3:", ps)
		case ah.W4:
			fmt.Println("W4:", ps)

		case ah.LastPointStatus:
			fmt.Println("LastPointStatus:", ps)

		default:
			//				fmt.Println(" skipping:", ps)
		}
	}
}

// checkHandicapCanonical
//
func checkHandicapCanonical() {
	// Verify that the handicap patterns are preserved by transformaions,
	for i, brd := range brds {
		t := ah.BoardTrans(i)
		inv := ah.InverseTrans[t]
		fmt.Println("Checking", ah.TransName[i], "and its inverse:", ah.TransName[inv])
		newBrd := brd.TransBoard(t)
		newBrdInv := newBrd.TransBoard(inv)
		if differBrds(brd, newBrdInv) {
			printBrds("Error: inverse differs", brd, newBrdInv, ah.TransName[i])
		}
	}

}

// printCannonicalHandicap points
//
func printCannonicalHandicap() {
	for ha := 0; ha <= 9; ha++ {
		if ha != 1 {
			var gam *sgf.GameTree = new(sgf.GameTree)
			gam.InitAbstHier(19, 19, ah.StringLevel, true)
			gam.SetHandicap(ha)
			gam.PlaceHandicap(ha, 19)
			for r := 0; r < 19; r++ {
				for c := 0; c < 19; c++ {
					nl := ah.MakeNodeLoc(ah.ColValue(c), ah.RowValue(r))
					bp := &gam.Graphs[ah.PointLevel].Nodes[nl]
					if bp.GetNodeLowState() != uint16(ah.Black) {
						if gam.IsCanonical(nl, ah.BoardHandicapSymmetry[ha]) {
							bp.SetNodeHighState(uint16(ah.White))
							//gam.SetPoint(nl, ah.White)
						}
					}
				}
			}
			str := "Handicap pattern " + strconv.Itoa(ha)
			printInitBoard(&gam.AbstHier, str)
			for trans := ah.T_FIRST; trans <= ah.T_LAST; trans += 1 {
				newBrd := gam.AbstHier.TransBoard(trans)
				if checkHandicapBrds(&gam.AbstHier, newBrd) {
					fmt.Print("false,  /* ", ah.TransName[trans], " */ ")
					// TODO: replace?					printBrds("Error: inverse differs", &gam.AbstHier, newBrd, ah.TransName[trans])
				} else {
					fmt.Print("true, /* ", ah.TransName[trans], " */ ")
				}
			}
			fmt.Println()
		}
	}
    checkHandicapCanonical()
}

// test EachNode and EachAdjNode
func printBoard() {
	var thePt *ah.GraphNode

	printPoint := func(nl ah.NodeLoc) {
		pp := &brds[3].Graphs[ah.PointLevel].Nodes[nl]
		if pp == thePt {
			fmt.Println("")
		}
		c, r := pp.GetPointColRow()
		fmt.Printf("[%d,%d]", c, r)
		if pp == thePt {
			fmt.Print(": ")
		} else {
			fmt.Print(", ")
		}
	}

	brds[3].EachNode(ah.PointLevel,
		func(brd *ah.Graph, nl ah.NodeLoc) {
			bp := &brd.Nodes[nl]
			thePt = bp
			printPoint(nl)
			brds[3].EachAdjNode(ah.PointLevel, nl, printPoint)
		})
	fmt.Println("")
}

// Print the sizes of struct and type declarations, by package.
//
func PrintSizes() {
	ah.PrintAhStructSizes()
	sgf.PrintSGFStructSizes()
	sgfdb.PrintSgfDbStructSizes()
}

// TODO: sort by second field (last name) if present
func Gtr(a []byte, b []byte) bool {
	idx := 0
	for (idx < len(a)) && (idx < len(b)) {
		if a[idx] > b[idx] {
			return true
		} else if a[idx] < b[idx] {
			return false
		}
		idx += 1
	}
	if len(a) > len(b) {
		return true
	}
	return false
}

func ReportSGFCounts() {
	for i, c := range sgf.ID_Counts {
		if c > 0 {
			fmt.Printf("Property %s used %d times.\n", string(sgf.GetProperty(sgf.PropertyDefIdx(i)).ID), c)
		}
	}
	if sgf.Unkn_Count > 0 {
		fmt.Printf("Property ?Unkn? used %d times.\n", sgf.Unkn_Count)
	}

	// report the HA map
	sum := 0
	for s, n := range sgf.HA_map {
		fmt.Printf("Handicap %s occurred %d times.\n", s, n)
		sum += n
	}
	fmt.Printf("Total Handicap games %d with %d different handicaps\n", sum, len(sgf.HA_map))

	// report the OH map
	sum = 0
	for s, n := range sgf.OH_map {
		fmt.Printf("Old Handicap %s occurred %d times.\n", s, n)
		sum += n
	}
	fmt.Printf("Total Old Handicap games %d with %d different settings\n", sum, len(sgf.OH_map))

	// report the RE map
	sum = 0
	for s, n := range sgf.RE_map {
		fmt.Printf("Result %s occurred %d times.\n", s, n)
		sum += n
	}
	fmt.Printf("Total games with Results %d among %d different settings\n", sum, len(sgf.RE_map))

	// report the RC (result comments)
	sum = 0
	for s, n := range sgf.RC_map {
		fmt.Printf("Result comment %s occurred %d times.\n", s, n)
		sum += n
	}
	fmt.Printf("Total Result comments %d with %d different comments\n", sum, len(sgf.RC_map))

	// report the RU map
	sum = 0
	for s, n := range sgf.RU_map {
		fmt.Printf("Rules %s occurred %d times.\n", s, n)
		sum += n
	}
	fmt.Printf("Total games with Rules %d with %d different settings\n", sum, len(sgf.RU_map))

	// report the BWRank map
	sum = 0
	for s, n := range sgf.BWRank_map {
		fmt.Printf("Rank %s occurred %d times.\n", s, n)
		sum += n
	}
	fmt.Printf("Total players with Ranks %d among %d different settings\n", sum, len(sgf.BWRank_map))

	// report the BWPlayer map
	//	sum = 0
	//	for s, n := range sgf.BWPlayer_map {
	//		fmt.Printf("Player %s occurred %d times, first %s, %s last %s, %s.\n", s, n.NGames, n.FirstGame, n.FirstRank, n.LastGame, n.LastRank)
	//		sum += n.NGames
	//	}
	//	fmt.Printf("Total players %d with %d different names\n", sum, len(sgf.BWPlayer_map))

	// sort the Player names, with counts:
	nPlayers := len(sgf.BWPlayer_map)
	var playerNames [][]byte
	var playerCount []int
	playerNames = make([][]byte, nPlayers)
	playerCount = make([]int, nPlayers)
	idx := 0
	for s, n := range sgf.BWPlayer_map {
		playerNames[idx] = make([]byte, len(s))
		_ = copy(playerNames[idx], s)
		playerCount[idx] = n.NGames
		idx += 1
	}
	// Sort them alphabetically:
	for ix := 0; ix < nPlayers; ix++ {
		for iy := ix; iy < nPlayers; iy++ {
			if Gtr(playerNames[ix], playerNames[iy]) {
				playerNames[ix], playerNames[iy] = playerNames[iy], playerNames[ix]
				playerCount[ix], playerCount[iy] = playerCount[iy], playerCount[ix]
			}
		}
	}
	for i, s := range playerNames {
		n, _ := sgf.BWPlayer_map[string(s)]
		fmt.Printf("Player %s: %d, first: %s, %s, last: %s, %s\n", s, playerCount[i], n.FirstGame, n.FirstRank, n.LastGame, n.LastRank)
	}

	// Sort them numerically:
	for ix := 0; ix < nPlayers; ix++ {
		for iy := ix; iy < nPlayers; iy++ {
			if playerCount[ix] < playerCount[iy] {
				playerNames[ix], playerNames[iy] = playerNames[iy], playerNames[ix]
				playerCount[ix], playerCount[iy] = playerCount[iy], playerCount[ix]
			}
		}
	}
	for i, s := range playerCount {
		n, _ := sgf.BWPlayer_map[string(playerNames[i])]
		fmt.Printf(" %d : %s, first:  %s, %s, last: %s, %s\n", s, playerNames[i], n.FirstGame, n.FirstRank, n.LastGame, n.LastRank)
	}
}

func init() {
    
	flag.IntVar(&moveLimit, "ml", 0, "ml = move limit. limit the number of moves read each .sgf file, 0 means no limit")
	flag.IntVar(&fileLimit, "fl", 0, "fl = file limit. limit the number of .sgf files read from a directory, 0 means no limit")
	flag.IntVar(&patternLimit, "pl", 0, "pl = pattern limit. limit the depth of pattern storing, 0 means no limit")
	flag.IntVar(&skipFiles, "sf", 0, "sf = skip files. skip this number of .sgf files before reading from a directory, 0 means no skip")

    flag.BoolVar(&doAllTests, "at", false, "at = all tests. do all tests, false (default) means not to do all tests, but can still do individual tests.")
    flag.BoolVar(&doPrintSizes, "ps", false, "ps = print sizes. print the sizes of data types, false (default) means not to print sizes.")
    flag.BoolVar(&doPrintConstants, "pc", false, "pc = print constants. print the values of constants, false (default) means not to print constants.")
    flag.BoolVar(&doPrintDirections, "pd", false, "pd = print directions. print the values of directions, false (default) means not to print directions.")
    flag.BoolVar(&doPrintHandicaps, "ph", false, "ph = print handicaps. print the canonical placement of handicaps, false (default) means not to print handicaps.")
    flag.BoolVar(&doPrintSGFProperties, "pp", false, "pp = print SGF properties. print the SGF porperties, false (default) means not to print SGF properties.")
    flag.BoolVar(&doPrintZKeys, "pzk", false, "pzk = print Z Keys. print the Zobrist Keys, false (default) means not to the Z Keys.")
    flag.BoolVar(&doVerifySGFPropertyOrder, "vpo", false, "vpo = verify SGF property order, after reading SGF_Properties_Spec.txt file, false (default) means do not verify.")
    flag.BoolVar(&doSmallSGFTests, "sst", false, "sst = do Small SGF Tests, false (default) means do not do these tests.")
    flag.BoolVar(&doTransTest, "tt", false, "tt = do Trans Test, false (default) means do not do Trans test.")
    flag.BoolVar(&doCountMoves, "cm", false, "cm = do Count Moves, false (default) means do not do count moves.")
    flag.BoolVar(&doReadWriteDatabase, "rwd", false, "rwd = do Read and Write Database, false (default) means do not read and write database.")
    flag.BoolVar(&doReadDatabaseAndBuild, "rdab", false, "rdab = do Read Database And Build patterns, false (default) means do not do Read Database And Build patterns.")
    flag.BoolVar(&doReadTeachingGames, "rtg", false, "rtg = do Read Teaching Games, false (default) means do not do read teaching games.")
    flag.BoolVar(&doReadWriteFuseki, "rwf", false, "rwf = do Read Write Fuseki, false (default) means do not read and write Fuseki file.")
    flag.BoolVar(&removeLabels, "rl", false, "rl = remove labels, false (default) means do not remove labels from Fuseki file.")
    flag.BoolVar(&allLabels, "al", false, "al = all labels, false (default) means do generate all labels in Fuseki file.")
    
    flag.StringVar(&SGFSpecFile, "ssf", defaultSGFSpecFile, "path to the SGF properties specification file.")
    flag.StringVar(&SmallSGFTestDir, "sstdir", defaultSmallSGFTestDir, "path to the Small SGF test directory.")
    flag.StringVar(&SmallSGFTestOutDir, "sstodir", defaultSmallSGFTestOutDir, "path to the Small SGF Tests output directory.")
    flag.StringVar(&DatabaseDir, "dbdir", defaultDatabaseDir, "path to the Database directory.")
    flag.StringVar(&TestOutDir, "todir", defaultTestOutDir, "path to the test output directory.")
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
    if doAllTests {
        fmt.Printf("at, all test has value %t\n", doAllTests)
    }
    if BoardPatternsDir != defaultBoardPatternsDir {
        fmt.Printf("bpdir, board patterns directory has value \"%s\"\n", BoardPatternsDir);
    }
    if doCountMoves {
        fmt.Printf("cm, do count moves has value %t\n", doCountMoves)
    }
    if DatabaseDir != defaultDatabaseDir {
        fmt.Printf("dbdir, data base directory has value \"%s\"\n", DatabaseDir);
    }
    if fusekiFileName != defaultFusekiFileName {
        fmt.Printf("ffn, Fuseki file name has value \"%s\"\n", fusekiFileName);
    }
    if fileLimit != 0 {
        fmt.Printf("fl, file limit has value %d\n", fileLimit)
    }
    if moveLimit != 0 {
        fmt.Printf("ml, move limit has value %d\n", moveLimit)
    }
    if outFusekiFileName != defaultOutFusekiFileName {
        fmt.Printf("offn, output Fuseki file name has value \"%s\"\n", outFusekiFileName);
    }
    if doPrintConstants {
        fmt.Printf("pc, print constants has value %t\n", doPrintConstants)
    }
    if doPrintDirections {
        fmt.Printf("pd, print directions has value %t\n", doPrintDirections)
    }
    if doPrintHandicaps {
        fmt.Printf("ph, print handicaps has value %t\n", doPrintHandicaps)
    }
    if patternLimit != 0 {
        fmt.Printf("pl, pattern limit has value %d\n", patternLimit)
    }
    if doPrintSGFProperties {
        fmt.Printf("pp, print SGF properties has value %t\n", doPrintSGFProperties)
    }
    if doPrintSizes {
        fmt.Printf("ps, print sizes has value %t\n", doPrintSizes)
    }
    if doPrintZKeys {
        fmt.Printf("pzk, print the Z Keys has value %t\n", doPrintZKeys)
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
    if doReadWriteDatabase {
        fmt.Printf("rwd, do read write database has value %t\n", doReadWriteDatabase)
    }
    if doReadWriteFuseki {
        fmt.Printf("rwf, do read write Fuseki file has value %t\n", doReadWriteFuseki)
    }
    if skipFiles != 0 {
        fmt.Printf("sf, skip files has value %d\n", skipFiles)
    }
    if SGFSpecFile != defaultSGFSpecFile {
        fmt.Printf("ssf, SGF specification file has value \"%s\"\n", SGFSpecFile);
    }
    if doSmallSGFTests {
        fmt.Printf("sst, do small SGF tests has value %t\n", doSmallSGFTests)
    }
    if SmallSGFTestDir != defaultSmallSGFTestDir {
        fmt.Printf("sstdir, small SGF test directory has value \"%s\"\n", SmallSGFTestDir);
    }
    if SmallSGFTestOutDir != defaultSmallSGFTestOutDir {
        fmt.Printf("sstodir, small SGF test output directory has value \"%s\"\n", SmallSGFTestOutDir);
    }
    if TeachingDir != defaultTeachingDir {
        fmt.Printf("tdir, teaching directory has value \"%s\"\n", TeachingDir);
    }
    if TestOutDir != defaultTestOutDir {
        fmt.Printf("todir, test output directory has value \"%s\"\n", TestOutDir);
    }
    if TeachingPatternsDir != defaultTeachingPatternsDir {
        fmt.Printf("tpdir, teaching patterns directory has value \"%s\"\n", TeachingPatternsDir);
    }
    if doTransTest {
        fmt.Printf("tt, do trans test has value %t\n", doTransTest)
    }
    if doVerifySGFPropertyOrder {
        fmt.Printf("vpo, verify SGF property order has value %t\n", doVerifySGFPropertyOrder)
    }
}

const DO_MULTI_CPU  = false

func main() {
    fmt.Printf("Program to generate opening pattern libraries:\n Version %s\n", SGF_GEN_GO_VERSION)
    nCPUs := runtime.NumCPU()
    if DO_MULTI_CPU {
        oldMaxProcs :=	runtime.GOMAXPROCS(nCPUs)
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
    
	if doAllTests || doPrintSizes {
		PrintSizes()
	}

	if doAllTests || doPrintConstants || doPrintDirections {
		printConst()
	}

	if doAllTests || doPrintHandicaps {
		printCannonicalHandicap()
	}

	if sgf.Setup(SGFSpecFile, doAllTests || doVerifySGFPropertyOrder, doAllTests || doPrintSGFProperties) == 0 {

		if doAllTests || doSmallSGFTests {
			ReadSmallSGFTests(SmallSGFTestDir, SmallSGFTestOutDir)
			ah.SetAHTrace(false)
		}

		if doAllTests || doTransTest {
			TestTrans()
			printBoard()
		}

		if doAllTests || doCountMoves {
			ah.SetAHTrace(false)
			sgfdb.CountFilesAndMoves(DatabaseDir, fileLimit)
		}
	}

	setup_and_count := time.Now()
    
	fmt.Printf("Setup and CountFilesAndMoves took %v to run.\n", setup_and_count.Sub(start))

	if doAllTests || doReadWriteDatabase {
		stat := sgfdb.ReadAndWriteDatabase(DatabaseDir, TestOutDir, fileLimit, moveLimit, skipFiles)
		if stat > 0 {
			fmt.Printf("Errors during reading and writing database: %d\n", stat)
		}
		ReportSGFCounts()
	}

	if doAllTests || doPrintZKeys {
		ah.PrintZKeys()
	}
    
	stop := time.Now()
    fmt.Printf("All tests took %v to run.\n", stop.Sub(start))

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
                sgf.ParseComments+sgf.GoGoD+sgf.Play, moveLimit)
            if len(errL) != 0 {
                fmt.Printf("Error %s during parsing: %s\n", errL.Error(), fusekiFileName)
            } else {
                    //TODO: add error reporting? ErrorList return value?
                if (sgfProcessOptions & RemoveLabels) > 0 {
                    fmt.Println("Removing labels from: ",fusekiFileName, " to ", outFusekiFileName);
                    prsr.GameTree.BreadthFirstTraverse(true, sgf.DoRemoveLabels)
                    prsr.GameTree.ReportDeletedProperties()
                }
                if (sgfProcessOptions & AddAllLabels) > 0 {
                    fmt.Println("Adding labels to: ", outFusekiFileName, " from ", fusekiFileName);
                    prsr.GameTree.DepthFirstTraverse(true, sgf.DoAddLabels)
                    fmt.Println("The number of added labels = ",sgf.NumberOfAddedLabels)
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
    
	finish := time.Now()
	fmt.Printf("Complete run took %v to run.\n", finish.Sub(start))
}
