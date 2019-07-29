#!/bin/bash

function filetime() {
    local FILETIME
    if [[ $(uname) == "Darwin" ]]; then
        FILETIME=$(stat -t %s -f %m $1)
    else
        FILETIME=$(stat $1 -c %Y)
    fi
    echo "$FILETIME"
}

function log() {
    echo -e "[$(date)]" $@
}

function cleanup() {
    if [[ $SKIP_CLEANUP != "true" ]]; then
        rm -rf $HOME_DIR
        rm -rf $CONFIG_DIR
    fi

    if [[ -n $1 ]]; then
        exit $1
    fi
}

function fail() {
    log "\e[31mTest Failure:\e[0m $@"
    log "Failing command:"
    log "\t$LAST_CMD"
    log "Test Environment:"
    log "\tHOME_DIR   = $HOME_DIR"
    log "\tCONFIG_DIR = $CONFIG_DIR"
    log "\tDFM        = $DFM_BIN"
    cleanup 1
}

LAST_CMD=""

function x() {
    LAST_CMD="$DFM_BIN --config-dir $CONFIG_DIR $@"
    $LAST_CMD
    if [[ $? != 0 ]]; then
        FAILED_CODE=$?
        fail "Non-zero exit code"

        if [[ $DEBUG_ON_ERROR == true ]]; then
            rust-lldb -- $@
        fi

        cleanup $FAILED_CODE
    fi
}

##############
# CLONE TEST #
##############
function dfm_clone_test() {
    local DFM=$1
    log "Testing DFM binary: $DFM"
    shift;
    local PROFILE_NAME=$1
    shift;
    local PROFILE_REPO=$1

    HOME_DIR=$(mktemp -d)
    CONFIG_DIR=$(mktemp -d)

    log Using: HOME_DIR $HOME_DIR
    log Using: CONFIG_DIR $CONFIG_DIR

    mkdir -p $HOME_DIR
    export HOME=$HOME_DIR

    log "Running clone tests..."
    log "Retrieving profile from: $PROFILE_REPO"

    log "DFM Version: " $(x --version)
    x clone --name $PROFILE_NAME $PROFILE_REPO
    x link $PROFILE_NAME

    dotfiles=(.dotfile after_link_ran before_link_ran)
    for file in $dotfiles; do
        local full_path=$HOME_DIR/$file
        if [[ ! -e $full_path ]]; then
            fail "$full_path should exist and does not"
        fi
    done

    skipped_dotfiles=(.dfm.yml .git .gitignore LICENSE.md README.md)
    for file in $skipped_dotfiles; do
        if [[ -e $HOME_DIR/$file ]]; then
            fail "Found file that should have been skipped: $HOME_DIR/$file" 
        fi
    done

    CURTIME=$(filetime $HOME_DIR/after_link_ran)
    x run-hook after_link
    NOWTIME=$(filetime $HOME_DIR/after_link_ran)
    TIMEDIFF=$(expr $CURTIME - $NOWTIME)
    if [[ $TIMEDIFF -gt 0 ]]; then
        fail "$HOME_DIR/after_link_ran was expected to be newer than $CURTIME got $NOWTIME"
    fi

    log "clone and link test -- SUCCESS"
    cleanup
}

#############
# INIT TEST #
#############
function dfm_init_test() {
    local DFM=$1
    log "Testing DFM binary: $DFM"
    shift
    local PROFILE_NAME=$1

    HOME_DIR=$(mktemp -d)
    CONFIG_DIR=$(mktemp -d)

    log Using: HOME_DIR $HOME_DIR
    log Using: CONFIG_DIR $CONFIG_DIR

    mkdir -p $HOME_DIR
    export HOME=$HOME_DIR

    log "Running init tests..."

    x init $PROFILE_NAME
    x link $PROFILE_NAME

    LIST_OUTPUT=$(x list)
    if [[ $LIST_OUTPUT != "$PROFILE_NAME" ]]; then
        fail "expected list to show integration profile, got: $LIST_OUTPUT"
    fi

    echo "a dotfile" > $HOME_DIR/.dotfile

    x add --no-git $HOME_DIR/.dotfile

    local fp="$CONFIG_DIR/dfm/profiles/$PROFILE_NAME/.dotfile"
    if [[ ! -f $fp ]]; then
        fail "expected $fp to be a file and is not." 
    fi

    if [[ ! -L $HOME_DIR/.dotfile ]]; then
        fail "expected $HOME_DIR/.dotfile to be a link and is not."
    fi

    OUTPUT=$(x git status --porcelain)
    if [[ $OUTPUT != "?? .dotfile" ]]; then
        fail "expected $DFM_BIN git status --porcelain to succeed and return ?? .dotfile got: $OUTPUT"
    fi

    x remove $PROFILE_NAME 

    log "init, add, and remove test -- SUCCESS"
    cleanup
}

DEBUG_ON_ERROR=""
SKIP_CLEANUP=""
DFM_BIN=""
PROFILE_REPO="https://github.com/chasinglogic/dfm_dotfile_test.git"
PROFILE_NAME="integration"
SKIP_CLEANUP=0

while getopts ":b:r:n:scd" opt; do
    case $opt in
        b) DFM_BIN="$OPTARG"      ;;
        r) PROFILE_REPO="$OPTARG" ;;
        n) PROFILE_NAME="$OPTARG" ;;
        s) SKIP_CLEANUP=true      ;;
        d) DEBUG_ON_ERROR=true    ;;
        c)
            cargo build
            if [[ $? != 0 ]]; then
                exit $?
            fi
            ;;
    esac
done

if [[ -z $DFM_BIN ]]; then
    log "Must provide path to dfm binary via -b flag." 
    exit 1
fi

dfm_clone_test $DFM_BIN $PROFILE_NAME $PROFILE_REPO
dfm_init_test $DFM_BIN $PROFILE_NAME
