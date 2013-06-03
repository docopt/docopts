# -*- coding: utf-8 -*-

import itertools
from distutils.core import Command

class build(Command):
    
    description = "build man pages"
    
    user_options = []
    
    def initialize_options(self):
        pass
    
    def finalize_options(self):
        pass
    
    def run(self):
        if not self.dry_run:
            print(self.distribution.manpage_sources)

class install(Command):
    
    description = "install man pages"
    
    user_options = []
    
    def initialize_options(self):
        pass
    
    def finalize_options(self):
        pass
    
    def run(self):
        if not self.dry_run:
            pass

def stripindent(string):
    return (sum(1 for s in itertools.takewhile(str.isspace, string)),
            string.lstrip())

def convert(source):
    lines = source.expandtabs().splitlines()
    header = lines.pop(0)
    footer = lines.pop(-1)
    
    iterator = iter(lines)
    indent, line = stripindent(iterator.next())
    for lookahead in iterator:
        lookahead_indent, lookahead_line = stripindent(lookahead)
        if lookahead_indent > indent > 0:
            pass

def setup():
    pass

if __name__ == '__main__':
    import sys
    with open(sys.argv[1]) as f:
        print(convert(f.read()))
