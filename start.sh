#!/bin/bash

if [[ $FIREFOX_PROFILE == "" ]]
then
    profile_dir=$(mktemp -d)
    cp /profile_tmpl/user.js "${profile_dir}/"
    export FIREFOX_OPTS="${FIREFOX_OPTS} --profile ${profile_dir}"
fi

X="${FIREFOX:=/firefox/firefox}"
export FIREFOX="$X"

exec /usr/local/bin/prpr
