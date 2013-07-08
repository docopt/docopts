# -*- coding: utf-8 -*-

from distutils.core import setup
import manpage

setup(name             = "docopts",
      version          = "0.7.0-dev",
      author           = "Lari Rasku",
      author_email     = "lari.o.rasku@gmail.com",
      url              = "https://github.com/docopt/docopts",
      description      = "Interpreter for docopt, the command-line "
                         "interface description language.",
      long_description = open('README').read(),
      scripts          = ["docopts"],
      requires         = ["docopt (==0.7.0)"],
      classifiers      = ["Development Status :: 3 - Alpha",
                          "Environment :: Console",
                          "Intended Audience :: Developers",
                          "License :: OSI Approved :: MIT License",
                          "Programming Language :: Python :: 2.6",
                          "Programming Language :: Python :: 2.7",
                          "Programming Language :: Python :: 3.1",
                          "Programming Language :: Python :: 3.2",
                          "Programming Language :: Python :: 3.3",
                          "Operating System :: POSIX",
                          "Topic :: Utilities"],
      platforms        = ["POSIX"],
      license          = "MIT License",
      cmdclass         = {'build_manpages':   manpage.build,
                          'install_manpages': manpage.install})
