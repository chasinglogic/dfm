#!/bin/bash

function log() {
    echo -e "[$(date)]" $@
}

function cleanup() {
    if [[ $SKIP_CLEANUP != "true" ]]; then
        rm -rf $HOME_DIR
        rm -rf $CONFIG_DIR
    fi

    if [[ -n $1 ]]; then
        log "\e[31mTEST FAILURE\e[0m"
        log "Test Environment:"
        log "\tHOME_DIR   = $HOME_DIR"
        log "\tCONFIG_DIR = $CONFIG_DIR"
        log "\tDFM        = $DFM_BIN"
        exit $1
    fi
}

function x() {
    $DFM_BIN --config-dir $CONFIG_DIR $@
    if [[ $? != 0 ]]; then
        FAILED_CODE=$?
        log "Failed to run:"
        log "\tHOME=$HOME_DIR $DFM_BIN --config-dir $CONFIG_DIR $@"

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

    dotfiles=(.dotfile)
    for file in $dotfiles; do
        if [[ ! -f $HOME_DIR/$file ]]; then
            log "Missing file: $HOME_DIR/$file" 
            log "Failed command:"
            log "\tHOME=$HOME_DIR $DFM_BIN --config-dir $CONFIG_DIR link $PROFILE_NAME"
            cleanup 1
        fi
    done

    dotfiles=(.dotfile after_link_ran before_link_ran)
    for file in $dotfiles; do
        local full_path=$HOME_DIR/$file
        if [[ ! -e $full_path ]]; then
            log "$full_path should exist and does not"
            exit 1
        fi
    done

    skipped_dotfiles=(.dfm.yml .git .gitignore LICENSE.md README.md)
    for file in $skipped_dotfiles; do
        if [[ -e $HOME_DIR/$file ]]; then
            log "Found file that should have been skipped: $HOME_DIR/$file" 
            log "Failed command:"
            log "\tHOME=$HOME_DIR $DFM_BIN --config-dir $CONFIG_DIR link $PROFILE_NAME"
            cleanup 1
        fi
    done

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

    if [[ $(x list) != "integration\n" ]]; then
        log "expected list to show integration profile, got:"
        log $(x list)
        cleanup 1
    fi

    echo "a dotfile" > $HOME_DIR/.dotfile

    x add $HOME_DIR/.dotfile

    local fp="$CONFIG_DIR/dfm/profiles/$PROFILE_NAME/.dotfile"
    if [[ ! -f $fp ]]; then
        log "expected $fp to be a file and is not." 
        cleanup 1
    fi

    if [[ ! -L $HOME_DIR/.dotfile ]]; then
        log "expected $HOME_DIR/.dotfile to be a link and is not."
        cleanup 1
    fi

    log "init and add test -- SUCCESS"
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
