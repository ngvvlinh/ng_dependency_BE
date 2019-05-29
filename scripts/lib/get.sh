#!/usr/bin/env bash
this=$0
prog=$(basename $0)
current=$prog
quickusage="Run '$prog help' for usage."


quickusage() {
    >&2 echo ""
    >&2 echo "$current: unknown command"
    >&2 echo "$quickusage"
    exit 126
}

expect() {
    local inlineusage=$1
    shift
    if [[ -z "$1" ]]; then
        >&2 echo "Usage: $current $inlineusage"
        >&2 echo "$quickusage"
        exit 126
    fi
    current="$current $1"
}

cmd_help() {
  usage="$prog is a tool for retrieving environment information for scripting.

Usage:

    $prog <command> [arguments]

The commands are:

    cmd <cmd>           get the command
    dir                 print some predefined directory or error
        backend         print '\$ETOPDIR/backend'
        config          print '\$ETOPDIR/etop-config'
        etop            print '\$ETOPDIR'
    gopath              print the current gopath or default to ~/go
    mod
        path <import>   print the location of 'import' from \$GOPATH/pkg/mod
"
  echo "$usage"
}

cmd_gopath() {
    if [[ -n "$GOPATH" ]]; then
        echo $GOPATH
    else
        echo ~/go
    fi
}

cmd_dir() {
    set -eo pipefail
    if [[ -z "$ETOPDIR" ]]; then
        >&2 echo "ETOPDIR is not set"
        exit 1
    fi

    dir=$1 ; expect "<dir>" $@ ; shift
    case $dir in
    etop)
        echo $ETOPDIR
        ;;
    backend)
        echo $ETOPDIR/backend
        ;;
    config)
        echo $ETOPDIR/config
        ;;
    *)
        quickusage
    esac
}

cmd_mod() {
    set -eo pipefail

    subcmd=$1 ; expect "path [arguments]" $@ ; shift
    case $subcmd in
    path)
        import=$1 ; expect "<import>" $@ ; shift
        path=$(go list -m $import)
        echo $($this gopath)/pkg/mod/${path/ /@}
        ;;
    *)
        quickusage
    esac
}

cmd_cmd() {
	path=$1 ; expect "<cmd>" $@ ; shift
	if [[ "$path" =~ ^etop.vn ]]; then
		go install "$path" || exit 1
	fi

	name=$(basename $path)
	cmdpath=$(which $name)
	if [[ $? == 0 ]]; then echo $cmdpath; exit 0; fi

	# try download it
	if [[ "$path" =~ ^k8s.io ]]; then
		GO111MODULE=off go get "$path"
		cmdpath=$(which $name)
		if [[ $? == 0 ]]; then echo $cmdpath; exit 0; fi
	fi

	>&2 echo "Command $name not found. Install it at $path"
	exit 125
}

command=$1 ; expect "<command>" $@ ; shift
case $command in
"-h" | "--help")
    cmd_help
    ;;
*)
    cmd_${command} $@
    if [[ $? == 127 ]]; then quickusage ; fi
esac