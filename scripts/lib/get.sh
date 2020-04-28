#!/usr/bin/env bash
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
        backend         print '\$PROJECT_DIR/backend'
        config          print '\$PROJECT_DIR/etop-config'
        etop            print '\$PROJECT_DIR'
    gopath              print the current gopath or default to ~/go
    mod
        path <import>   print the location of 'import' from \$GOPATH/pkg/mod
"
  echo "$usage"
}

if [[ -n "$GOPATH" ]]; then
    gopath="$GOPATH"
else
    gopath="$HOME/go"
fi

cmd_gopath() { echo "$gopath" ; }

cmd_dir() {
    set -eo pipefail
    if [[ -z "$PROJECT_DIR" ]]; then
        >&2 echo "PROJECT_DIR is not set"
        exit 1
    fi

    dir=$1 ; expect "<dir>" "$@" ; shift
    case $dir in
    etop)
        echo $PROJECT_DIR
        ;;
    backend)
        echo $PROJECT_DIR/backend
        ;;
    config)
        echo $PROJECT_DIR/config
        ;;
    *)
        quickusage
    esac
}

cmd_mod() {
    set -eo pipefail

    subcmd=$1 ; expect "path [arguments]" "$@" ; shift
    case $subcmd in
    path)
        import=$1 ; expect "<import>" "$@" ; shift
        path=$(go list -m $import)
        echo "${gopath}/pkg/mod/${path/ /@}"
        ;;
    *)
        quickusage
    esac
}

cmd_cmd() {
	path=$1 ; expect "<cmd>" "$@" ; shift
	name=$(basename $path)

	if [[ "$path" =~ ^o.o ]]; then
		go install "$path" || exit 1
	  cmdpath="${gopath}/bin/${name}"
	  echo $cmdpath
	  exit 0
	fi

	cmdpath=$(which $name)
	if [[ $? == 0 ]]; then echo $cmdpath; exit 0; fi

	# try download it
	if [[ "$path" =~ ^k8s.io || "$path" =~ ^github.com ]]; then
		GO111MODULE=off go get "$path"
		cmdpath=$(which $name)
		if [[ $? == 0 ]]; then echo $cmdpath; exit 0; fi
	fi

	>&2 echo "Command $name not found. Install it at $path"
	exit 125
}

command=$1 ; expect "<command>" "$@" ; shift
case $command in
"-h" | "--help")
    cmd_help
    ;;
*)
    cmd_${command} "$@"
    if [[ $? == 127 ]]; then quickusage ; fi
esac
