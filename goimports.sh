#!/bin/bash

##
# parse_file_hook_args
# Creates global vars:
#   OPTIONS: List of options to passed to command
#   FILES  : List of files to process, filtered against ignore_file_pattern_array
#
function parse_file_hook_args {
	OPTIONS=()
	# If arg doesn't pass [ -f ] check, then it is assumed to be an option
	#
	while [ $# -gt 0 ] && [ "$1" != "--" ] && [ ! -f "$1" ]; do
		OPTIONS+=("$1")
		shift
	done

	local all_files
	all_files=()
	# Assume start of file list (may still be options)
	#
	while [ $# -gt 0 ] && [ "$1" != "--" ]; do
		all_files+=("$1")
		shift
	done

	# If '--' next, then files = options
	#
	if [ "$1" == "--" ]; then
		shift
		# Append to previous options
		#
		OPTIONS+=("${all_files[@]}")
		all_files=()
	fi

	# Any remaining arguments are assumed to be files
	#
	all_files+=("$@")

	# Filter out vendor entries and ignore_file_pattern_array
	#
	FILES=()
	local file pattern
	ignore_file_pattern_array+=( "vendor/*" "*/vendor/*" "*/vendor" )
	for file in "${all_files[@]}"; do
		for pattern in "${ignore_file_pattern_array[@]}"; do
			if [[ "${file}" == ${pattern} ]] ; then # pattern => unquoted
				continue 2
			fi
		done
		FILES+=("${file}")
	done
}

parse_file_hook_args "$@"

##
# Remove all blank lines in go 'imports' statements, then sort with default goimports.
# This scripts only works in macOS with `brew install gnu-sed`
# Change gsed to sed in case of you are developing under the Linux.
#
for file in "${FILES[@]}"; do
  gsed -i '
    /^import (/,/)/ {
      /^$/ d
    }
  ' "${file}"
  goimports -l -d -local github.com/bookstairs/bookhunter -w "${file}"
done
