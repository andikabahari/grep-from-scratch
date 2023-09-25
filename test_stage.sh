#!/bin/sh

case $1 in
  '1') 
    CODECRAFTERS_CURRENT_STAGE_SLUG='init'
    ;;
  '2') 
    CODECRAFTERS_CURRENT_STAGE_SLUG='match_digit'
    ;;
  '3') 
    CODECRAFTERS_CURRENT_STAGE_SLUG='match_alphanumeric'
    ;;
  '4') 
    CODECRAFTERS_CURRENT_STAGE_SLUG='positive_character_groups'
    ;;
  '5') 
    CODECRAFTERS_CURRENT_STAGE_SLUG='negative_character_groups'
    ;;
  '6') 
    CODECRAFTERS_CURRENT_STAGE_SLUG='combining_character_classes'
    ;;
  '7') 
    CODECRAFTERS_CURRENT_STAGE_SLUG='start_of_string_anchor'
    ;;
  '8') 
    CODECRAFTERS_CURRENT_STAGE_SLUG='end_of_string_anchor'
    ;;
  '9') 
    CODECRAFTERS_CURRENT_STAGE_SLUG='one_or_more_quantifier'
    ;;
  '10') 
    CODECRAFTERS_CURRENT_STAGE_SLUG='zero_or_one_quantifier'
    ;;
  '11') 
    CODECRAFTERS_CURRENT_STAGE_SLUG='wildcard'
    ;;
  '12') 
    CODECRAFTERS_CURRENT_STAGE_SLUG='alternation'
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