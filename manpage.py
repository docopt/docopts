# -*- coding: utf-8 -*-

from distutils.core import Command

class build_manpage(Command):
    
    description = "build manpage"
    
    user_options = []
    
    def initialize_options(self):
        pass
    
    def finalize_options(self):
        pass
    
    def run(self):
        if not self.dry_run:
            pass

class install_manpage(Command):
    
    description = "install manpage"
    
    user_options = []
    
    def initialize_options(self):
        pass
    
    def finalize_options(self):
        pass
    
    def run(self):
        if not self.dry_run:
            pass
