File:       README.txt
Project:    abst-hier

Created by Ken Friedenbach on Sept. 2, 2012
Updated for abst-hier workspace Jan. 1, 2014
Last updated Jan 31, 3014

This project is a re-engineering of AHGoGo,
a C++ project.

For overall status and objectives, see:
    TODO: add document ...
in the Projects directory.

The $GOPATH should be set to:
    abst-hier/
which has subdirectories
    src/
    bin/
    pkg/
    test/
and so on.

From here, the path to the sources,
following Go recommendations is:
    src/github.com/Ken1JF

The source is divided into three directories:
    /ah - abstraction hierarchy updating
    /sgf - reading and writing .sgf files
        SGF properties are defined in: 
            /sgf/sgf_properties_spec.txt
    /sgfdb - read and write the GoGoD and other databases
These are built and installed as three packages.

There are now test files in each package, which use
the Go test technology. They can be used to verify
functionality, as well as getting cover analysis.

The old test program is in:
    /test-ahgo/test_ahgo.go
It is gradually having test functionality removed,
until only code under development will remain: the building
of pattern files.

There is a script to build and run the tests:
    test-ahgo/buildScript.bash
    
Some major tests:
    read write database
        after this test, run "./checkout.txt" in
            TODO: add to abst-hier
            TODO: get a golden output, then stop using
            
Working on:
    Processing: Fuseki.sgf - a file of Fuseki openings.
        Based on the book: Fuseki Small Encyclopedia Vol. 2 with 3 Pro games transformed to normalize the first move, so it is in the Left-Upper Octant.
        Note: 99 InteriorNodes are caused by 99 C[...] comments.
    Fuseki2.sgf has all labels removed, by test_ahgo.
    Fuseki2B.sgf has labels added by CGoban 3.2 for "all moves".
        2464 "LB[" and 374 "][" 
    Fuseki2C.sgf has labesl added by CGoban 3.2 for "variations" only.
        239 "LB[" and 374 "]["
    Fuseki3.sgf has labels added by test_ahgo.
    /usr/local/bin/diffsgf -l option to ignore LB properties.
    BreadFirst and DepthFirst Search are working, with Pre or Post processing options. 
    DoRemoveLabels is working.
    DoAddLabels is working: only generates labels when there is more than one choice.
    Does no put labels on pass moves.
    Does not count pass moves when deciding if to put labels; differs from CGoban 3.2. (CGoban does not put a label on a pass move, but counts it, and ends up putting a label on the only other, as a singleton.)
    Fuseki3 and Fuseki2C have the same 8 labels, dp:A -> ec:H, for the second layer, but in different orders on output...
        Fuseki3 has order A -> H
        Fuseki2 has order cc, dc, ec, cd, dd, cp, dp, cq, i.e. col-row, sorted by row, then col.
TODO:
    read the GoGoD databases, record frequency and results.
    construct other databases for handicaps
    split records based on Komi and Rules.
    add frequency and results for self and opponents.
    build an opening player that uses Fuseki files.
    add files for half boards, Joseki, sides of boards, and quadrants.
    modify player to transition through using these.
    offer opening player to David Doshey's SlugGo.
    