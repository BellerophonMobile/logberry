#!/bin/sh

##
## Simple example of extracting commit information from git in order
## to generate an automatically up-to-date build statement.
##
## There are two optional parameters:
##   ./update-rev.sh [package] [target]
##
##       package: The go package in which to place the output; default
##                is main.
##
##       target: The file to which to write the output; default is
##               stdout.
##
##

#-- Parameters
target="/dev/stdout"
package=main
if [ $# -ge 1 ]
then
    package=$1

  if [ $# -ge 2 ]
  then
      target=$2
  fi
fi


#-- Get the data
root=$(basename `git rev-parse --show-toplevel`)

branch=$(git rev-parse --abbrev-ref HEAD)

commit=$(git rev-parse HEAD)
modified=
if [ "$(git status -uno | grep modified | wc -l)" -ne "0" ]
then
    modified="*"
fi

host=$(hostname)
user=$(whoami)

date=$(date --rfc-3339=seconds)


#-- Write it out

cat > "$target" <<EOF
$root $branch $commit$modified $host $user $date
EOF

#-- Done!
