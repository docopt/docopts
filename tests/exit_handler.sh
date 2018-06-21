#!/usr/bin/env bash
exit_handler() {
    exit_code=$?
    echo "recieved exit $exit_code"
}

trap exit_handler EXIT

