#!/bin/sh

mockgen -package=mocks -mock_names=Logger=Logger \
        -destination=mocks/logger.go github.com/masonhubco/rebar/v2 Logger

mockgen -package=mocks -mock_names=TxWrapper=TxWrapper \
        -destination=mocks/tx_wrapper.go github.com/masonhubco/rebar/v2/middleware TxWrapper
