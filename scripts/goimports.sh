#!/bin/bash

## This shell script is used with https://github.com/TekWizely/pre-commit-golang
#  You should add it to your .pre-commit-config.yaml file with the options like
#
#  - repo: https://github.com/tekwizely/pre-commit-golang
#    rev: v1.0.0-rc.1
#    hooks:
#      - id: my-cmd
#        name: goimports
#        alias: goimports
#        args: [ scripts/goimports.sh, github.com/syhily/hobbit ]

module="$1"
file="$2"

# Detect the running OS.
COMMAND="sed"
if [[ $OSTYPE == 'darwin'* ]]; then
  # macOS have to use the gsed which can be installed by `brew install gsed`.
  COMMAND="gsed"
fi

# Detect the command.
command -v $COMMAND >/dev/null 2>&1 || { echo >&2 "Require ${COMMAND} but it's not installed. Aborting."; exit 1; }
command -v goimports >/dev/null 2>&1 || { echo >&2 "Require goimports but it's not installed. Aborting."; exit 1; }

# Remove all the import spaces in staging golang files.
REPLACEMENT=$(cat <<-END
'
  /^import (/,/)/ {
    /^$/ d
  }
'
END
)
bash -c "${COMMAND} -i ${REPLACEMENT} ${file}"

# Format the staging golang files.
goimports -l -d -local "${module}" -w "${file}"
