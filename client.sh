#!/bin/bash

# configurations
#
# Set CURL to curl binary path, empty string to disable curl
# Set WGET to wget binary path, empty string to disable wget
# One of CURL and WGET must be set. CURL is used if both are set.

remote="http://localhost:9801"
if [[ $SERVER != "" ]]
then
    remote="$SERVER"
fi

function urlencode() {
    local length="${#1}"
    for (( i = 0; i < length; i++ )); do
        local c="${1:i:1}"
        case $c in
            [a-zA-Z0-9.~_-])
		printf "$c"
		;;
            *)
		printf "$c" | xxd -p -c1 | while read x
		do
		    printf "%%%s" "$x"
		done
		;;
        esac
    done
}

# with curl
function c {
    uri="$1"
    wait="$2"
    secret="$3"
    opts="-sL"

    if [[ $DEBUG != "" ]]
    then
	opts="-SL"
    fi

    $CURL -X POST \
	  --data-urlencode "uri=${uri}" \
	  --data-urlencode "wait=${wait}" \
	  --data-urlencode "secret=${secret}" \
	  $opts \
	  "${remote}/grab"
}

# with wget
function w {
    uri=$(urlencode "$1")
    wait=$(urlencode "$2")
    secret=$(urlencode "$3")
    opts="-O -"

    if [[ $DEBUG == "" ]]
    then
	opts="-q -O -"
    fi

    $WGET \
	--post-data="uri=${uri}&wait=${wait}&secret=${secret}" \
	$opts \
	"${remote}/grab"	
}

function help {
    echo "Usage: ${1} uri wait [secret]"
    echo ''
    echo 'Example:'
    echo "  ${1} http://google.com 'div#xfoot > script'"
    echo ''
    echo 'Envvars:'
    echo '  SERVER: Specify location of your server'
    echo "          SERVER=http://example.com ${1} http://google.com 'div#xfoot > script'"
    echo '  DEBUG:  Make curl/wget show basic messages'
    echo "          DEBUG=1 ${1} http://google.com 'div#xfoot > script'"
}

if [[ $1 == "" || $2 == "" ]]
then
    help "$0"
    exit 1
fi

# validate configurations
if [[ $CURL == "" && $WGET == "" ]]
then
    # invalid, try auto detect
    which curl > /dev/null 2>& 1
    if [[ $? == 0 ]]
    then
	CURL=curl
    else
	which wget > /dev/null 2>&1
	if [[ $? == 0 ]]
	then
	    WGET=wget
	else
	    echo 'None of CURL and WGET is set, and we cannot determine path to any of it.'
	    exit 1
	fi
    fi
fi

if [[ $CURL == "" ]]
then
    w "$1" "$2" "$3"
else
    c "$1" "$2" "$3"
fi
