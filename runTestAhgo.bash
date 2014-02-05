#!/bin/bash

#  runTestAhgo.bash
#  abst-hier
#
#  Created by Ken Friedenbach on 2/1/14.
#  Copyright 2012-2014. All rights reserved.

cd ../../../../
pwd

#export GOOS=darwin
#export GOARCH=amd64
#export GOROOT=/usr/local/go
#export GOBIN=$GOROOT/bin
export GOPATH=/usr/local/go:/Users/ken/Projects/abst-hier
export PATH=$PATH:/Users/ken/Projects/abst-hier/bin

echo "veting"
go tool vet src/github.com/Ken1JF/test-ahgo/test_ahgo.go

echo "installing"
go install github.com/Ken1JF/ah
go install github.com/Ken1JF/sgf
go install github.com/Ken1JF/sgfdb


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
test_ahgo -al=true -offn="src/github.com/Ken1JF/test-ahgo/Fuseki3.sgf" -ffn="src/github.com/Ken1JF/test-ahgo/Fuseki2.sgf" -rwf=true >& test_ahgo_new.txt -ssf="src/github.com/Ken1JF/sgf/sgf_properties_spec.txt"
diff test_ahgo_new.txt test_ahgo_out.txt

