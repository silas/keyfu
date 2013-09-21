#!/usr/bin/env bash

ulimit -n 8192

rootdir='/opt/keyfu'
bindir="$rootdir/bin"
interface=${1-:8000}

exec $bindir/keyfu \
  -interface=$interface \
  -static=$rootdir/static \
  -templates=$rootdir/templates
