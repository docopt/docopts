#!/usr/bin/env python
# -*- coding: utf-8 -*-
# vim: set ts=4 sw=4 sts=4 et:

__doc__ = """Shell interface for docopt, the CLI description language.

Usage:
  docopts [options] -h <msg> : [<argv>...]

Options:
  -h <msg>, --help=<msg>        The help message in docopt format.
                                If - is given, read the help message from
                                standard input.
                                If no argument is given, print docopts's own
                                help message and quit.
  -V <msg>, --version=<msg>     A version message.
                                If - is given, read the version message from
                                standard input.  If the help message is also
                                read from standard input, it is read first.
                                If no argument is given, print docopts's own
                                version message and quit.
  -O, --options-first           Disallow interspersing options and positional
                                arguments: all arguments starting from the
                                first one that does not begin with a dash will
                                be treated as positional arguments.
  -H, --no-help                 Don't handle --help and --version specially.
  -A <name>                     Export the arguments as a Bash 4.x associative
                                array called <name>.
  -s <str>, --separator=<str>   The string to use to separate the help message
                                from the version message when both are given
                                via standard input. [default: ----]

"""

__version__ = """docopts 0.6.1+fix
Copyright (C) 2013 Vladimir Keleshev, Lari Rasku.
License MIT <http://opensource.org/licenses/MIT>.
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.

"""

import re
import sys

try:
    from cStringIO import StringIO
except ImportError:
    from io import StringIO

from docopt import docopt, DocoptExit, DocoptLanguageError

# helper functions
def shellquote(s):
    return "'" + s.replace("'", r"'\''") + "'"

def isbashidentifier(s):
    return re.match(r'^([A-Za-z]|[A-Za-z_][0-9A-Za-z_]+)$', s)

def to_bash(obj):
    return {
        type(None): lambda x: '',
        bool:       lambda x: 'true' if x else 'false',
        int:        lambda x: str(x),
        str:        lambda x: shellquote(x),
        list:       lambda x: '(' + ' '.join(map(shellquote, x)) + ')',
    }[type(obj)](obj)

def name_mangle(elem):
    if elem == '-' or elem == '--':
        return None
    elif re.match(r'^<.*>$', elem):
        var = elem[1:-1]
    elif re.match(r'^-[^-]$', elem):
        var = elem[1]
    elif re.match(r'^--.+$', elem):
        var = elem[2:]
    else:
        var = elem
    var = var.replace('-', '_')
    if not isbashidentifier(var):
        raise ValueError(elem)
    else:
        return var

# parse docopts's own arguments
try:
    args = docopt(__doc__, help=False, options_first=True)
except DocoptExit as e:
    message = e.args[0]
    if '-h' == message[0:2] or '--help' == message[0:6]:
        print(__doc__.strip())
        sys.exit()
    if '-V' == message[0:2] or '--version' == message[0:9]:
        print(__version__.strip())
        sys.exit()
    else:
        sys.exit(message)

argv = args['<argv>']
doc = args['--help']
version = args['--version']
options_first = args['--options-first']
help = not args['--no-help']
name = args['-A']
separator = args['--separator']

if doc == '-' and version == '-':
    doc, version = (page.strip() for page in
                    sys.stdin.read().split(separator, 1))
elif doc == '-':
    doc = sys.stdin.read().strip()
elif version == '-':
    version = sys.stdin.read().strip()

# parse options or abort if there is an error in docopt
try:
    # temporarily redirect stdout to a StringIO so we can catch docopt()
    # output on --help and --version
    stdout = sys.stdout
    sys.stdout = StringIO()
    exit_message = None
    args = docopt(doc, argv, help, version, options_first)
except DocoptLanguageError as e:
    # invalid docstring by user
    sys.exit("%s: invalid doc argument: %s" % (sys.argv[0], e))
except DocoptExit as e:
    # invoked with invalid arguments
    exit_message = "echo %s >&2\nexit 64" % (shellquote(str(e)),)
except SystemExit as e:
    # --help or --version found and --no-help was not given
    exit_message = "echo -n %s\nexit 0" % (shellquote(sys.stdout.getvalue()),)
finally:
    # restore stdout to normal and quit if a docopt parse error happened
    sys.stdout.close()
    sys.stdout = stdout
    if exit_message:
        print(exit_message)
        sys.exit()

if name is not None:
    if not isbashidentifier(name):
        sys.exit("%s: not a valid Bash identifier: %s" % (sys.argv[0], name))
    # fake nested Bash arrays for repeatable arguments with values
    arrays = dict((elem, value) for elem, value in args.items() if
                  isinstance(value, list))
    for elem, value in arrays.items():
        del args[elem]
        args[elem+',#'] = len(value)
        args.update(('%s,%d' % (elem, i), v) for i,v in enumerate(value))
    print('declare -A %s' % (name,))
    for elem, value in args.items():
        print('%s[%s]=%s' % (name, shellquote(elem), to_bash(value)))
else:
    try:
        variables = dict(zip(map(name_mangle, args.keys()),
                             map(to_bash, args.values())))
    except ValueError as e:
        sys.exit("%s: name could not be mangled into a valid Bash "
                 "identifier: %s" % (sys.argv[0], e))
    else:
        variables.pop(None, None)
        args.pop('-', None)
        args.pop('--', None)
    if len(variables) < len(args):
        sys.exit("%s: two or more elements have identically mangled names" %
                 (sys.argv[0],))
    for var, value in variables.items():
        print("%s=%s" % (var, value))
