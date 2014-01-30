#!/bin/bash

#  buildScript.sh
#  abst-hier
#
#  Created by Ken Friedenbach on 7/17/12.
#  Copyright 2012-2014. All rights reserved.

cd ../../../../
pwd

export GOOS=darwin
export GOARCH=amd64
export GOROOT=/usr/local/go
export GOBIN=$GOROOT/bin
export GOPATH=/usr/local/go:/Users/ken/Projects/abst-hier
export PATH=$PATH:/Users/ken/Projects/abst-hier/bin

echo "cleaning"
go clean github.com/Ken1JF/ah
go clean github.com/Ken1JF/sgf
go clean github.com/Ken1JF/sgfdb
echo "looking for bin/test_ahgo to remove"
if [[ -e bin/test_ahgo ]] ; then echo "remove bin/test_ahgo" ; rm bin/test_ahgo ; fi

echo "veting"
go tool vet src/github.com/Ken1JF/ah
go tool vet src/github.com/Ken1JF/sgf
go tool vet src/github.com/Ken1JF/sgfdb
go tool vet src/github.com/Ken1JF/test-ahgo/test_ahgo.go

echo "installing"
go install github.com/Ken1JF/ah
go install github.com/Ken1JF/sgf
go install github.com/Ken1JF/sgfdb

echo "testing"
go test github.com/Ken1JF/ah
go test github.com/Ken1JF/sgf
go test github.com/Ken1JF/sgfdb

echo "testing individual package test coverage"
go test -cover github.com/Ken1JF/ah
go test -cover github.com/Ken1JF/sgf
go test -cover github.com/Ken1JF/sgfdb

echo "testing combined coverage together"
if [[ -e test/ ]] ; then echo "test/ exists" ; else echo "making test/" ; mkdir test ; fi
go test -cover github.com/Ken1JF/ah -coverprofile test/ahCover.out
go test -cover github.com/Ken1JF/sgf -coverpkg github.com/Ken1JF/ah,github.com/Ken1JF/sgf -coverprofile test/sgfCover.out
go test -cover github.com/Ken1JF/sgfdb -coverpkg github.com/Ken1JF/ah,github.com/Ken1JF/sgf,github.com/Ken1JF/sgfdb -coverprofile test/sgfdbCover.out

#-coverprofile cover.out

#echo "NOT formatting"
# Only do this when everything is clean!
# and XCode has been shut down.
echo "formatting"
go fmt github.com/Ken1JF/ah
go fmt github.com/Ken1JF/sgf
go fmt github.com/Ken1JF/sgfdb
go fmt github.com/Ken1JF/test-ahgo

echo "building test_ahgo"
go build -o bin/test_ahgo src/github.com/Ken1JF/test-ahgo/test_ahgo.go

echo "running test_ahgo"
test_ahgo -at=true -al=true -offn="src/github.com/Ken1JF/test-ahgo/Fuseki3.sgf" -ffn="src/github.com/Ken1JF/test-ahgo/Fuseki2.sgf" -rwf=true >& test_ahgo_new.txt -ssf="src/github.com/Ken1JF/sgf/sgf_properties_spec.txt"
diff test_ahgo_new.txt test_ahgo_out.txt

