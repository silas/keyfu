#!/usr/bin/env bash

basedir='/opt/tarsnap'
datadir="$basedir/data"
cachedir="$basedir/cache"

mongoexport -h 127.0.0.1 -d keyfu -c keyword > $datadir/keyfu.keyword.json
mongoexport -h 127.0.0.1 -d keyfu -c user > $datadir/keyfu.user.json

if [[ $1 != "skip" ]]; then
  tarsnap \
    -f "keyfu-$( date +%Y-%m-%d )" \
    -c \
    --keyfile $basedir/tarsnap.key \
    --cachedir $cachedir \
    $datadir
fi
