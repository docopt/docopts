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

class Section(str):
    

def stripindent(string):
    return (sum(1 for s in itertools.takewhile(str.isspace, string)),
            string.lstrip())
