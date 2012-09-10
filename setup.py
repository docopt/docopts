# -*- coding: utf-8 -*-

"""``docopts`` is a Python script that provides a shell interface for the
`reference Python implementation <https://github.com/docopt/docopt>`_ of
`docopt <http://docopt.org>`_, the command-line interface description
language.

To install ``docopts`` from source, extract the release archive and execute
the following command in the release directory::

    python setup.py install

You may need superuser permissions; alternatively, you can install it locally
by appending ``--user`` to the above command.

For full documentation, see the ``docopts`` `homepage
<http://github.com/docopt/docopts>`_.

"""

import os
from setuptools import setup

def read(name):
    with open(os.path.join(os.path.dirname(__file__), name)) as f:
        return f.read()

setup(name             = "docopts",
      version          = "0.5.0",
      author           = "Lari Rasku",
      author_email     = "raskug@lavabit.com",
      url              = "https://github.com/docopt/docopts",
      license          = "MIT",
      description      = "Shell interface for docopt, the command-line "
                         "interface description language.",
      keywords         = "shell bash docopt command-line",
      long_description = __doc__,
      scripts          = ["docopts"],
      install_requires = ["docopt == 0.5.0"],
      classifiers      = ["Development Status :: 3 - Alpha",
                          "Environment :: Console",
                          "Intended Audience :: Developers",
                          "License :: OSI Approved :: MIT License",
                          "Programming Language :: Python :: 2.5",
                          "Programming Language :: Python :: 2.6",
                          "Programming Language :: Python :: 2.7",
                          "Programming Language :: Python :: 3.1",
                          "Programming Language :: Python :: 3.2",
                          "Topic :: Software Development :: User Interfaces"])
