  GNU nano 5.4                                                                     distil.sh                                                                              
#!/usr/bin/env bash
set -eux

MOONSHINE="bitflip"
DIR="release"
mkdir -p "$DIR"
export GOOS=windows
export GOARCH=386
go build
mv "$MOONSHINE".exe "${DIR}/${MOONSHINE}"-win386.exe
export GOOS=windows
export GOARCH=amd64
go build
mv "$MOONSHINE".exe "${DIR}/${MOONSHINE}"-win64.exe
export GOOS=linux
export GOARCH=amd64
go build
mv "$MOONSHINE" "${DIR}/${MOONSHINE}"-linux64
export GOOS=
export GOARCH=386
go build
mv "$MOONSHINE" "${DIR}/${MOONSHINE}"-darwin386
export GOOS=darwin
export GOARCH=amd64
go build
mv "$MOONSHINE" "${DIR}/${MOONSHINE}"-darwinAmd64
export GOOS=
export GOARCH=
