# New API proposal for docopts - docopt on shell (bash)

Date: 2015-03-13

See [issue #7](https://github.com/docopt/docopts/issues/7) for the discussion.

Submit pull request for changes.

## Proposed API

Return a string, representing a bash shell associative array for $args.

## Example of usage

Converted from [quick_example.py](https://github.com/docopt/docopt/blob/master/examples/quick_example.py)

~~~bash
#!/bin/bash 
#
# Usage:
#  quick_example.sh tcp <host> <port> [--timeout=<seconds>]
#  quick_example.sh serial <port> [--baud=9600] [--timeout=<seconds>]
#  quick_example.sh -h | --help | --version

libpath=../

# sourcing the API provides some bash functions ie: docopt()
# still using the python parser
source $libpath/docopts.sh

help=$(get_help_string $0)
version='0.1.1rc'

# $@ needs to be at the last param
parsed=$(docopt -A args "$help" $version "$@")
echo "$parsed"                                                                                            
eval "$parsed"

# Evaluating the parsed output, will create $args in the current scope.
# It is an associative array the name is passed from the command line with -A
# (same switch as bash: declare -A assoc)
~~~

## shell helpers

~~~bash
# auto extract the Usage string from the top shell script comment
# ie: help=$(sed -n -e '/^# Usage:/,/^$/ s/^# \?//p' < $0)
help_string=$(get_help_string $0)

# if the option as mulitples values, you can get it into an array
array_opt=( ${args[--multiple-time]} )
~~~

## Examples

### naval_fate.sh

See [python/docopt](https://github.com/docopt/docopt#docopt-creates-beautiful-command-line-interfaces)

~~~bash
naval_fate.sh ship Guardian move 100 150 --speed=15
~~~

returned $parsed string to be evaluated:

~~~
declare -A arguments
arguments['--drifting']=false    
arguments['mine']=false
arguments['--help']=false
arguments['move']=true
arguments['--moored']=false
arguments['new']=false
arguments['--speed']='15'
arguments['remove']=false
arguments['--version']=false
arguments['set']=false
arguments['<name>']='Guardian'
arguments['ship']=true
arguments['<x>']='100'
arguments['shoot']=false
arguments['<y>']='150'
~~~

## Limitations

* repeatable options with filename with space inside, will not be easily split on $IFS
