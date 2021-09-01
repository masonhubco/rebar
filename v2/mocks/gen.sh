#!/bin/sh

mockgen -package=mocks -mock_names=Logger=Logger \
        -destination=mocks/logger.go github.com/masonhubco/rebar/v2 Logger
