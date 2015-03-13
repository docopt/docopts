# new API prosal for docopts - docopt on shell (bash)

See [issue #7](https://github.com/docopt/docopts/issues/7) for the discussion.

Submit pull request for changes.

## proposed API

Specifying the associative array for $args and some shell centric behavior have to be considered, too.


~~~bash
#!/bin/bash
# sourcing the API providing some bash functions ie: docopt()
# still using the python parser
source docopts.sh

parsed=$(docopt -A args "$help_string" "$@" "$version" "$options_first")
eval "$parsed"

# it will create in the current scope $args is an associative array the name passed
# from command line with -A (same switch as bash)
declare -A args
args['--option']='value'
~~~

and may be some shell helpers too:

~~~bash
help_string=$(get_help_string $0)

array_opt=( $(get_docopt array "${args[--multiple-time]})" )
~~~
