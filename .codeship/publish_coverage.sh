#!/bin/bash
goveralls -coverprofile=coverage.out -service=codeship -jobid $CI_COMMIT_ID -repotoken $COVERALLS_TOKEN