#!/bin/sh

case $1 in
  '1') 
    CODECRAFTERS_CURRENT_STAGE_SLUG='init'
    ;;
  *)
    echo 'Invalid stage'
    exit
    ;;
esac

cd grep-tester
go build -o ../grep-go/test.out ./cmd/tester

cd ../grep-go
CODECRAFTERS_SUBMISSION_DIR=$(pwd) \
CODECRAFTERS_CURRENT_STAGE_SLUG=${CODECRAFTERS_CURRENT_STAGE_SLUG} \
./test.out
rm ./test.out