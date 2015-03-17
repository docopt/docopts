#!/bin/bash 
# vim: set et sw=4 ts=4 sts=4:
#
# docopt for bash
#
# Usage: See API_proposal.md
#   source path/to/docopt.sh
#   docopt -A args -h "$help" -v $version : "$@"
#
# the prefix docopt_* is used to export globals and functions

# compute this file dirpath:
docopt_sh_me=$(readlink -f "${BASH_SOURCE[0]}")
docopt_sh_dir="$(dirname "$docopt_sh_me")"

# fetch Usage: from the given filname
# usually $0 in the main level script
docopt_get_help_string() {
    local myfname=$1
    # filter the block (all blocks) starting at a "# Usage:" and ending 
    # at an empty line, one level of comment markup is removed
    sed -n -e '/^# Usage:/,/^$/ s/^# \?//p' < $myfname
}

# function wrapper
# Usage: same as docopts.py
docopt() {
    #   docopts [options] -h <msg> : [<argv>...]
    # call python parser on embedded code
    python3 <(sed -n -e '/^### EMBEDDED/,$ s/^#> // p' "$docopt_sh_me") "$@"
}

docopt_find_docopts() {
    # docopts is the python wrapper using docopt.py
    # it is now embedded in docopts
    echo ../docopts
    # will do: docopt_sh_dir="$.dirname $.readlink -f "${BASH_SOURCE[0]}"))"
}

# convert a repeatable option parsed by docopts.py into a bash ARRAY
# Usage: myarray=( $(docopt_get_values args --repeatable-option"
docopt_get_values() {
    local opt=$2
    local ref="\${$1[$opt,#]}"
    local nb_val=$(eval echo "$ref")
    local i=0
    local vars=""
    while [[ $i -lt $nb_val ]] ; do
        ref="\${$1[$opt,$i]}"
        eval "vars+=\" $ref\""
        i=$(($i + 1))
    done
    echo $vars
}

# echo evaluable code to get alls the values into a bash array
# Usage: eval "$(docopt_get_eval_array args FILE myarray)"
docopt_get_eval_array() {
    local ref="\${$1[$2,#]}"
    local nb_val=$(eval echo "$ref")
    local i=0
    local vars=""
    echo "declare -a $3"
    while [[ $i -lt $nb_val ]] ; do
        ref="\${$1[$2,$i]}"
        eval "echo \"$3+=( '$ref' )\""
        i=$(($i + 1))
    done
}

# Auto parser for same docopts usage over script, or lazyness.
# use this convention:
#  - help string in: $help
#  - Usage parse by docopt_get_help_string at beginning of the script
#  - arguments evaluated at global level in: $args
#  - no version information
#  
docopt_auto_parse() {
    local script_fname=$1
    shift
    help="$(docopt_get_help_string "$script_fname")"
    # $args assoc array must be declared outside on this function
    # or it's scope will be local
    docopt -A args -h "$help" : "$@" | grep -v -- 'declare -A args'
}

## main code
# --auto : don't forget to pass $@
# Usage: source docopts.sh --auto "$@"
if [[ "$1" == "--auto" ]] ; then
    shift
    declare -A args
    eval "$(docopt_auto_parse "${BASH_SOURCE[1]}" "$@")"
fi


# don't alter this code, this is the original python docopts + docopt.py
### EMBEDDED
#> #!/usr/bin/env python
#> # -*- coding: utf-8 -*-
#> # vim: set ts=4 sw=4 sts=4 et:
#> 
#> __doc__ = """Shell interface for docopt, the CLI description language.
#> 
#> Usage:
#>   docopts [options] -h <msg> : [<argv>...]
#> 
#> Options:
#>   -h <msg>, --help=<msg>        The help message in docopt format.
#>                                 If - is given, read the help message from
#>                                 standard input.
#>                                 If no argument is given, print docopts's own
#>                                 help message and quit.
#>   -V <msg>, --version=<msg>     A version message.
#>                                 If - is given, read the version message from
#>                                 standard input.  If the help message is also
#>                                 read from standard input, it is read first.
#>                                 If no argument is given, print docopts's own
#>                                 version message and quit.
#>   -O, --options-first           Disallow interspersing options and positional
#>                                 arguments: all arguments starting from the
#>                                 first one that does not begin with a dash will
#>                                 be treated as positional arguments.
#>   -H, --no-help                 Don't handle --help and --version specially.
#>   -A <name>                     Export the arguments as a Bash 4.x associative
#>                                 array called <name>.
#>   -s <str>, --separator=<str>   The string to use to separate the help message
#>                                 from the version message when both are given
#>                                 via standard input. [default: ----]
#> 
#> """
#> 
#> __version__ = """docopts 0.6.1+fix
#> Copyright (C) 2013 Vladimir Keleshev, Lari Rasku.
#> License MIT <http://opensource.org/licenses/MIT>.
#> This is free software: you are free to change and redistribute it.
#> There is NO WARRANTY, to the extent permitted by law.
#> 
#> """
#> 
#> import re
#> import sys
#> 
#> try:
#>     from cStringIO import StringIO
#> except ImportError:
#>     from io import StringIO
#> 
#> ### embbeded: from docopt import docopt, DocoptExit, DocoptLanguageError
#> #> """Pythonic command-line interface parser that will make you smile.
#> #> 
#> #>  * http://docopt.org
#> #>  * Repository and issue-tracker: https://github.com/docopt/docopt
#> #>  * Licensed under terms of MIT license (see LICENSE-MIT)
#> #>  * Copyright (c) 2013 Vladimir Keleshev, vladimir@keleshev.com
#> #> 
#> #> """
#> #> import sys
#> #> import re
#> #> 
#> #> 
#> #> __all__ = ['docopt']
#> #> __version__ = '0.6.1'
#> 
#> 
#> class DocoptLanguageError(Exception):
#> 
#>     """Error in construction of usage-message by developer."""
#> 
#> 
#> class DocoptExit(SystemExit):
#> 
#>     """Exit in case user invoked program with incorrect arguments."""
#> 
#>     usage = ''
#> 
#>     def __init__(self, message=''):
#>         SystemExit.__init__(self, (message + '\n' + self.usage).strip())
#> 
#> 
#> class Pattern(object):
#> 
#>     def __eq__(self, other):
#>         return repr(self) == repr(other)
#> 
#>     def __hash__(self):
#>         return hash(repr(self))
#> 
#>     def fix(self):
#>         self.fix_identities()
#>         self.fix_repeating_arguments()
#>         return self
#> 
#>     def fix_identities(self, uniq=None):
#>         """Make pattern-tree tips point to same object if they are equal."""
#>         if not hasattr(self, 'children'):
#>             return self
#>         uniq = list(set(self.flat())) if uniq is None else uniq
#>         for i, child in enumerate(self.children):
#>             if not hasattr(child, 'children'):
#>                 assert child in uniq
#>                 self.children[i] = uniq[uniq.index(child)]
#>             else:
#>                 child.fix_identities(uniq)
#> 
#>     def fix_repeating_arguments(self):
#>         """Fix elements that should accumulate/increment values."""
#>         either = [list(child.children) for child in transform(self).children]
#>         for case in either:
#>             for e in [child for child in case if case.count(child) > 1]:
#>                 if type(e) is Argument or type(e) is Option and e.argcount:
#>                     if e.value is None:
#>                         e.value = []
#>                     elif type(e.value) is not list:
#>                         e.value = e.value.split()
#>                 if type(e) is Command or type(e) is Option and e.argcount == 0:
#>                     e.value = 0
#>         return self
#> 
#> 
#> def transform(pattern):
#>     """Expand pattern into an (almost) equivalent one, but with single Either.
#> 
#>     Example: ((-a | -b) (-c | -d)) => (-a -c | -a -d | -b -c | -b -d)
#>     Quirks: [-a] => (-a), (-a...) => (-a -a)
#> 
#>     """
#>     result = []
#>     groups = [[pattern]]
#>     while groups:
#>         children = groups.pop(0)
#>         parents = [Required, Optional, OptionsShortcut, Either, OneOrMore]
#>         if any(t in map(type, children) for t in parents):
#>             child = [c for c in children if type(c) in parents][0]
#>             children.remove(child)
#>             if type(child) is Either:
#>                 for c in child.children:
#>                     groups.append([c] + children)
#>             elif type(child) is OneOrMore:
#>                 groups.append(child.children * 2 + children)
#>             else:
#>                 groups.append(child.children + children)
#>         else:
#>             result.append(children)
#>     return Either(*[Required(*e) for e in result])
#> 
#> 
#> class LeafPattern(Pattern):
#> 
#>     """Leaf/terminal node of a pattern tree."""
#> 
#>     def __init__(self, name, value=None):
#>         self.name, self.value = name, value
#> 
#>     def __repr__(self):
#>         return '%s(%r, %r)' % (self.__class__.__name__, self.name, self.value)
#> 
#>     def flat(self, *types):
#>         return [self] if not types or type(self) in types else []
#> 
#>     def match(self, left, collected=None):
#>         collected = [] if collected is None else collected
#>         pos, match = self.single_match(left)
#>         if match is None:
#>             return False, left, collected
#>         left_ = left[:pos] + left[pos + 1:]
#>         same_name = [a for a in collected if a.name == self.name]
#>         if type(self.value) in (int, list):
#>             if type(self.value) is int:
#>                 increment = 1
#>             else:
#>                 increment = ([match.value] if type(match.value) is str
#>                              else match.value)
#>             if not same_name:
#>                 match.value = increment
#>                 return True, left_, collected + [match]
#>             same_name[0].value += increment
#>             return True, left_, collected
#>         return True, left_, collected + [match]
#> 
#> 
#> class BranchPattern(Pattern):
#> 
#>     """Branch/inner node of a pattern tree."""
#> 
#>     def __init__(self, *children):
#>         self.children = list(children)
#> 
#>     def __repr__(self):
#>         return '%s(%s)' % (self.__class__.__name__,
#>                            ', '.join(repr(a) for a in self.children))
#> 
#>     def flat(self, *types):
#>         if type(self) in types:
#>             return [self]
#>         return sum([child.flat(*types) for child in self.children], [])
#> 
#> 
#> class Argument(LeafPattern):
#> 
#>     def single_match(self, left):
#>         for n, pattern in enumerate(left):
#>             if type(pattern) is Argument:
#>                 return n, Argument(self.name, pattern.value)
#>         return None, None
#> 
#>     @classmethod
#>     def parse(class_, source):
#>         name = re.findall('(<\S*?>)', source)[0]
#>         value = re.findall('\[default: (.*)\]', source, flags=re.I)
#>         return class_(name, value[0] if value else None)
#> 
#> 
#> class Command(Argument):
#> 
#>     def __init__(self, name, value=False):
#>         self.name, self.value = name, value
#> 
#>     def single_match(self, left):
#>         for n, pattern in enumerate(left):
#>             if type(pattern) is Argument:
#>                 if pattern.value == self.name:
#>                     return n, Command(self.name, True)
#>                 else:
#>                     break
#>         return None, None
#> 
#> 
#> class Option(LeafPattern):
#> 
#>     def __init__(self, short=None, long=None, argcount=0, value=False):
#>         assert argcount in (0, 1)
#>         self.short, self.long, self.argcount = short, long, argcount
#>         self.value = None if value is False and argcount else value
#> 
#>     @classmethod
#>     def parse(class_, option_description):
#>         short, long, argcount, value = None, None, 0, False
#>         options, _, description = option_description.strip().partition('  ')
#>         options = options.replace(',', ' ').replace('=', ' ')
#>         for s in options.split():
#>             if s.startswith('--'):
#>                 long = s
#>             elif s.startswith('-'):
#>                 short = s
#>             else:
#>                 argcount = 1
#>         if argcount:
#>             matched = re.findall('\[default: (.*)\]', description, flags=re.I)
#>             value = matched[0] if matched else None
#>         return class_(short, long, argcount, value)
#> 
#>     def single_match(self, left):
#>         for n, pattern in enumerate(left):
#>             if self.name == pattern.name:
#>                 return n, pattern
#>         return None, None
#> 
#>     @property
#>     def name(self):
#>         return self.long or self.short
#> 
#>     def __repr__(self):
#>         return 'Option(%r, %r, %r, %r)' % (self.short, self.long,
#>                                            self.argcount, self.value)
#> 
#> 
#> class Required(BranchPattern):
#> 
#>     def match(self, left, collected=None):
#>         collected = [] if collected is None else collected
#>         l = left
#>         c = collected
#>         for pattern in self.children:
#>             matched, l, c = pattern.match(l, c)
#>             if not matched:
#>                 return False, left, collected
#>         return True, l, c
#> 
#> 
#> class Optional(BranchPattern):
#> 
#>     def match(self, left, collected=None):
#>         collected = [] if collected is None else collected
#>         for pattern in self.children:
#>             m, left, collected = pattern.match(left, collected)
#>         return True, left, collected
#> 
#> 
#> class OptionsShortcut(Optional):
#> 
#>     """Marker/placeholder for [options] shortcut."""
#> 
#> 
#> class OneOrMore(BranchPattern):
#> 
#>     def match(self, left, collected=None):
#>         assert len(self.children) == 1
#>         collected = [] if collected is None else collected
#>         l = left
#>         c = collected
#>         l_ = None
#>         matched = True
#>         times = 0
#>         while matched:
#>             # could it be that something didn't match but changed l or c?
#>             matched, l, c = self.children[0].match(l, c)
#>             times += 1 if matched else 0
#>             if l_ == l:
#>                 break
#>             l_ = l
#>         if times >= 1:
#>             return True, l, c
#>         return False, left, collected
#> 
#> 
#> class Either(BranchPattern):
#> 
#>     def match(self, left, collected=None):
#>         collected = [] if collected is None else collected
#>         outcomes = []
#>         for pattern in self.children:
#>             matched, _, _ = outcome = pattern.match(left, collected)
#>             if matched:
#>                 outcomes.append(outcome)
#>         if outcomes:
#>             return min(outcomes, key=lambda outcome: len(outcome[1]))
#>         return False, left, collected
#> 
#> 
#> class Tokens(list):
#> 
#>     def __init__(self, source, error=DocoptExit):
#>         self += source.split() if hasattr(source, 'split') else source
#>         self.error = error
#> 
#>     @staticmethod
#>     def from_pattern(source):
#>         source = re.sub(r'([\[\]\(\)\|]|\.\.\.)', r' \1 ', source)
#>         source = [s for s in re.split('\s+|(\S*<.*?>)', source) if s]
#>         return Tokens(source, error=DocoptLanguageError)
#> 
#>     def move(self):
#>         return self.pop(0) if len(self) else None
#> 
#>     def current(self):
#>         return self[0] if len(self) else None
#> 
#> 
#> def parse_long(tokens, options):
#>     """long ::= '--' chars [ ( ' ' | '=' ) chars ] ;"""
#>     long, eq, value = tokens.move().partition('=')
#>     assert long.startswith('--')
#>     value = None if eq == value == '' else value
#>     similar = [o for o in options if o.long == long]
#>     if tokens.error is DocoptExit and similar == []:  # if no exact match
#>         similar = [o for o in options if o.long and o.long.startswith(long)]
#>     if len(similar) > 1:  # might be simply specified ambiguously 2+ times?
#>         raise tokens.error('%s is not a unique prefix: %s?' %
#>                            (long, ', '.join(o.long for o in similar)))
#>     elif len(similar) < 1:
#>         argcount = 1 if eq == '=' else 0
#>         o = Option(None, long, argcount)
#>         options.append(o)
#>         if tokens.error is DocoptExit:
#>             o = Option(None, long, argcount, value if argcount else True)
#>     else:
#>         o = Option(similar[0].short, similar[0].long,
#>                    similar[0].argcount, similar[0].value)
#>         if o.argcount == 0:
#>             if value is not None:
#>                 raise tokens.error('%s must not have an argument' % o.long)
#>         else:
#>             if value is None:
#>                 if tokens.current() in [None, '--']:
#>                     raise tokens.error('%s requires argument' % o.long)
#>                 value = tokens.move()
#>         if tokens.error is DocoptExit:
#>             o.value = value if value is not None else True
#>     return [o]
#> 
#> 
#> def parse_shorts(tokens, options):
#>     """shorts ::= '-' ( chars )* [ [ ' ' ] chars ] ;"""
#>     token = tokens.move()
#>     assert token.startswith('-') and not token.startswith('--')
#>     left = token.lstrip('-')
#>     parsed = []
#>     while left != '':
#>         short, left = '-' + left[0], left[1:]
#>         similar = [o for o in options if o.short == short]
#>         if len(similar) > 1:
#>             raise tokens.error('%s is specified ambiguously %d times' %
#>                                (short, len(similar)))
#>         elif len(similar) < 1:
#>             o = Option(short, None, 0)
#>             options.append(o)
#>             if tokens.error is DocoptExit:
#>                 o = Option(short, None, 0, True)
#>         else:  # why copying is necessary here?
#>             o = Option(short, similar[0].long,
#>                        similar[0].argcount, similar[0].value)
#>             value = None
#>             if o.argcount != 0:
#>                 if left == '':
#>                     if tokens.current() in [None, '--']:
#>                         raise tokens.error('%s requires argument' % short)
#>                     value = tokens.move()
#>                 else:
#>                     value = left
#>                     left = ''
#>             if tokens.error is DocoptExit:
#>                 o.value = value if value is not None else True
#>         parsed.append(o)
#>     return parsed
#> 
#> 
#> def parse_pattern(source, options):
#>     tokens = Tokens.from_pattern(source)
#>     result = parse_expr(tokens, options)
#>     if tokens.current() is not None:
#>         raise tokens.error('unexpected ending: %r' % ' '.join(tokens))
#>     return Required(*result)
#> 
#> 
#> def parse_expr(tokens, options):
#>     """expr ::= seq ( '|' seq )* ;"""
#>     seq = parse_seq(tokens, options)
#>     if tokens.current() != '|':
#>         return seq
#>     result = [Required(*seq)] if len(seq) > 1 else seq
#>     while tokens.current() == '|':
#>         tokens.move()
#>         seq = parse_seq(tokens, options)
#>         result += [Required(*seq)] if len(seq) > 1 else seq
#>     return [Either(*result)] if len(result) > 1 else result
#> 
#> 
#> def parse_seq(tokens, options):
#>     """seq ::= ( atom [ '...' ] )* ;"""
#>     result = []
#>     while tokens.current() not in [None, ']', ')', '|']:
#>         atom = parse_atom(tokens, options)
#>         if tokens.current() == '...':
#>             atom = [OneOrMore(*atom)]
#>             tokens.move()
#>         result += atom
#>     return result
#> 
#> 
#> def parse_atom(tokens, options):
#>     """atom ::= '(' expr ')' | '[' expr ']' | 'options'
#>              | long | shorts | argument | command ;
#>     """
#>     token = tokens.current()
#>     result = []
#>     if token in '([':
#>         tokens.move()
#>         matching, pattern = {'(': [')', Required], '[': [']', Optional]}[token]
#>         result = pattern(*parse_expr(tokens, options))
#>         if tokens.move() != matching:
#>             raise tokens.error("unmatched '%s'" % token)
#>         return [result]
#>     elif token == 'options':
#>         tokens.move()
#>         return [OptionsShortcut()]
#>     elif token.startswith('--') and token != '--':
#>         return parse_long(tokens, options)
#>     elif token.startswith('-') and token not in ('-', '--'):
#>         return parse_shorts(tokens, options)
#>     elif token.startswith('<') and token.endswith('>') or token.isupper():
#>         return [Argument(tokens.move())]
#>     else:
#>         return [Command(tokens.move())]
#> 
#> 
#> def parse_argv(tokens, options, options_first=False):
#>     """Parse command-line argument vector.
#> 
#>     If options_first:
#>         argv ::= [ long | shorts ]* [ argument ]* [ '--' [ argument ]* ] ;
#>     else:
#>         argv ::= [ long | shorts | argument ]* [ '--' [ argument ]* ] ;
#> 
#>     """
#>     parsed = []
#>     while tokens.current() is not None:
#>         if tokens.current() == '--':
#>             return parsed + [Argument(None, v) for v in tokens]
#>         elif tokens.current().startswith('--'):
#>             parsed += parse_long(tokens, options)
#>         elif tokens.current().startswith('-') and tokens.current() != '-':
#>             parsed += parse_shorts(tokens, options)
#>         elif options_first:
#>             return parsed + [Argument(None, v) for v in tokens]
#>         else:
#>             parsed.append(Argument(None, tokens.move()))
#>     return parsed
#> 
#> 
#> def parse_defaults(doc):
#>     defaults = []
#>     for s in parse_section('options:', doc):
#>         # FIXME corner case "bla: options: --foo"
#>         _, _, s = s.partition(':')  # get rid of "options:"
#>         split = re.split('\n[ \t]*(-\S+?)', '\n' + s)[1:]
#>         split = [s1 + s2 for s1, s2 in zip(split[::2], split[1::2])]
#>         options = [Option.parse(s) for s in split if s.startswith('-')]
#>         defaults += options
#>     return defaults
#> 
#> 
#> def parse_section(name, source):
#>     pattern = re.compile('^([^\n]*' + name + '[^\n]*\n?(?:[ \t].*?(?:\n|$))*)',
#>                          re.IGNORECASE | re.MULTILINE)
#>     return [s.strip() for s in pattern.findall(source)]
#> 
#> 
#> def formal_usage(section):
#>     _, _, section = section.partition(':')  # drop "usage:"
#>     pu = section.split()
#>     return '( ' + ' '.join(') | (' if s == pu[0] else s for s in pu[1:]) + ' )'
#> 
#> 
#> def extras(help, version, options, doc):
#>     if help and any((o.name in ('-h', '--help')) and o.value for o in options):
#>         print(doc.strip("\n"))
#>         sys.exit()
#>     if version and any(o.name == '--version' and o.value for o in options):
#>         print(version)
#>         sys.exit()
#> 
#> 
#> class Dict(dict):
#>     def __repr__(self):
#>         return '{%s}' % ',\n '.join('%r: %r' % i for i in sorted(self.items()))
#> 
#> 
#> def docopt(doc, argv=None, help=True, version=None, options_first=False):
#>     """Parse `argv` based on command-line interface described in `doc`.
#> 
#>     `docopt` creates your command-line interface based on its
#>     description that you pass as `doc`. Such description can contain
#>     --options, <positional-argument>, commands, which could be
#>     [optional], (required), (mutually | exclusive) or repeated...
#> 
#>     Parameters
#>     ----------
#>     doc : str
#>         Description of your command-line interface.
#>     argv : list of str, optional
#>         Argument vector to be parsed. sys.argv[1:] is used if not
#>         provided.
#>     help : bool (default: True)
#>         Set to False to disable automatic help on -h or --help
#>         options.
#>     version : any object
#>         If passed, the object will be printed if --version is in
#>         `argv`.
#>     options_first : bool (default: False)
#>         Set to True to require options precede positional arguments,
#>         i.e. to forbid options and positional arguments intermix.
#> 
#>     Returns
#>     -------
#>     args : dict
#>         A dictionary, where keys are names of command-line elements
#>         such as e.g. "--verbose" and "<path>", and values are the
#>         parsed values of those elements.
#> 
#>     Example
#>     -------
#>     >>> from docopt import docopt
#>     >>> doc = '''
#>     ... Usage:
#>     ...     my_program tcp <host> <port> [--timeout=<seconds>]
#>     ...     my_program serial <port> [--baud=<n>] [--timeout=<seconds>]
#>     ...     my_program (-h | --help | --version)
#>     ...
#>     ... Options:
#>     ...     -h, --help  Show this screen and exit.
#>     ...     --baud=<n>  Baudrate [default: 9600]
#>     ... '''
#>     >>> argv = ['tcp', '127.0.0.1', '80', '--timeout', '30']
#>     >>> docopt(doc, argv)
#>     {'--baud': '9600',
#>      '--help': False,
#>      '--timeout': '30',
#>      '--version': False,
#>      '<host>': '127.0.0.1',
#>      '<port>': '80',
#>      'serial': False,
#>      'tcp': True}
#> 
#>     See also
#>     --------
#>     * For video introduction see http://docopt.org
#>     * Full documentation is available in README.rst as well as online
#>       at https://github.com/docopt/docopt#readme
#> 
#>     """
#>     argv = sys.argv[1:] if argv is None else argv
#> 
#>     usage_sections = parse_section('usage:', doc)
#>     if len(usage_sections) == 0:
#>         raise DocoptLanguageError('"usage:" (case-insensitive) not found.')
#>     if len(usage_sections) > 1:
#>         raise DocoptLanguageError('More than one "usage:" (case-insensitive).')
#>     DocoptExit.usage = usage_sections[0]
#> 
#>     options = parse_defaults(doc)
#>     pattern = parse_pattern(formal_usage(DocoptExit.usage), options)
#>     # [default] syntax for argument is disabled
#>     #for a in pattern.flat(Argument):
#>     #    same_name = [d for d in arguments if d.name == a.name]
#>     #    if same_name:
#>     #        a.value = same_name[0].value
#>     argv = parse_argv(Tokens(argv), list(options), options_first)
#>     pattern_options = set(pattern.flat(Option))
#>     for options_shortcut in pattern.flat(OptionsShortcut):
#>         doc_options = parse_defaults(doc)
#>         options_shortcut.children = list(set(doc_options) - pattern_options)
#>         #if any_options:
#>         #    options_shortcut.children += [Option(o.short, o.long, o.argcount)
#>         #                    for o in argv if type(o) is Option]
#>     extras(help, version, argv, doc)
#>     matched, left, collected = pattern.fix().match(argv)
#>     if matched and left == []:  # better error message if left?
#>         return Dict((a.name, a.value) for a in (pattern.flat() + collected))
#>     raise DocoptExit()
#> 
#> # ----------------------- end docopt.py ---------------------------
#> 
#> # helper functions
#> def shellquote(s):
#>     return "'" + s.replace("'", r"'\''") + "'"
#> 
#> def isbashidentifier(s):
#>     return re.match(r'^([A-Za-z]|[A-Za-z_][0-9A-Za-z_]+)$', s)
#> 
#> def to_bash(obj):
#>     return {
#>         type(None): lambda x: '',
#>         bool:       lambda x: 'true' if x else 'false',
#>         int:        lambda x: str(x),
#>         str:        lambda x: shellquote(x),
#>         list:       lambda x: '(' + ' '.join(map(shellquote, x)) + ')',
#>     }[type(obj)](obj)
#> 
#> def name_mangle(elem):
#>     if elem == '-' or elem == '--':
#>         return None
#>     elif re.match(r'^<.*>$', elem):
#>         var = elem[1:-1]
#>     elif re.match(r'^-[^-]$', elem):
#>         var = elem[1]
#>     elif re.match(r'^--.+$', elem):
#>         var = elem[2:]
#>     else:
#>         var = elem
#>     var = var.replace('-', '_')
#>     if not isbashidentifier(var):
#>         raise ValueError(elem)
#>     else:
#>         return var
#> 
#> # ----------------------- MAIN -------------
#> 
#> # parse docopts's own arguments
#> try:
#>     args = docopt(__doc__, help=False, options_first=True)
#> except DocoptExit as e:
#>     message = e.args[0]
#>     if '-h' == message[0:2] or '--help' == message[0:6]:
#>         print(__doc__.strip())
#>         sys.exit()
#>     if '-V' == message[0:2] or '--version' == message[0:9]:
#>         print(__version__.strip())
#>         sys.exit()
#>     else:
#>         sys.exit(message)
#> 
#> argv = args['<argv>']
#> doc = args['--help']
#> version = args['--version']
#> options_first = args['--options-first']
#> help = not args['--no-help']
#> name = args['-A']
#> separator = args['--separator']
#> 
#> if doc == '-' and version == '-':
#>     doc, version = (page.strip() for page in
#>                     sys.stdin.read().split(separator, 1))
#> elif doc == '-':
#>     doc = sys.stdin.read().strip()
#> elif version == '-':
#>     version = sys.stdin.read().strip()
#> 
#> # parse options or abort if there is an error in docopt
#> try:
#>     # temporarily redirect stdout to a StringIO so we can catch docopt()
#>     # output on --help and --version
#>     stdout = sys.stdout
#>     sys.stdout = StringIO()
#>     exit_message = None
#>     args = docopt(doc, argv, help, version, options_first)
#> except DocoptLanguageError as e:
#>     # invalid docstring by user
#>     sys.exit("%s: invalid doc argument: %s" % (sys.argv[0], e))
#> except DocoptExit as e:
#>     # invoked with invalid arguments
#>     exit_message = "echo %s >&2\nexit 64" % (shellquote(str(e)),)
#> except SystemExit as e:
#>     # --help or --version found and --no-help was not given
#>     exit_message = "echo -n %s\nexit 0" % (shellquote(sys.stdout.getvalue()),)
#> finally:
#>     # restore stdout to normal and quit if a docopt parse error happened
#>     sys.stdout.close()
#>     sys.stdout = stdout
#>     if exit_message:
#>         print(exit_message)
#>         sys.exit()
#> 
#> if name is not None:
#>     if not isbashidentifier(name):
#>         sys.exit("%s: not a valid Bash identifier: %s" % (sys.argv[0], name))
#>     # fake nested Bash arrays for repeatable arguments with values
#>     arrays = dict((elem, value) for elem, value in args.items() if
#>                   isinstance(value, list))
#>     for elem, value in arrays.items():
#>         del args[elem]
#>         args[elem+',#'] = len(value)
#>         args.update(('%s,%d' % (elem, i), v) for i,v in enumerate(value))
#>     print('declare -A %s' % (name,))
#>     for elem, value in args.items():
#>         print('%s[%s]=%s' % (name, shellquote(elem), to_bash(value)))
#> else:
#>     try:
#>         variables = dict(zip(map(name_mangle, args.keys()),
#>                              map(to_bash, args.values())))
#>     except ValueError as e:
#>         sys.exit("%s: name could not be mangled into a valid Bash "
#>                  "identifier: %s" % (sys.argv[0], e))
#>     else:
#>         variables.pop(None, None)
#>         args.pop('-', None)
#>         args.pop('--', None)
#>     if len(variables) < len(args):
#>         sys.exit("%s: two or more elements have identically mangled names" %
#>                  (sys.argv[0],))
#>     for var, value in variables.items():
#>         print("%s=%s" % (var, value))
