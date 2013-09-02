# -*- coding: utf-8 -*-

try:
    import json
except ImportError:
    import simplejson as json

import shutil
import unittest
import os.path

from subprocess import Popen, PIPE
from distutils.core import Command

class run(Command):
    
    description = "run unit tests"
    
    user_options = []
    
    def initialize_options(self):
        pass
    
    def finalize_options(self):
        pass
    
    def suite(self):
        suite = unittest.TestSuite()
        with open('testcases.docopt', 'r') as infile:
            testcases = infile.read().split('r"""')[1:]
        for testcase in testcases:
            doc, tests = testcase.split('"""', 1)
            tests = tests.split('$ ')[1:]
            for test in tests:
                argv, expect = test.split('\n', 1)
                argv = argv.split()[1:]
                expect = json.loads(expect[:expect.find('#')])
                suite.addTest(DocoptsTest(doc, argv, expect))
        return suite
    
    def run(self):
        if self.dry_run:
            return
        suite = self.suite()
        runner = unittest.runner.TextTestRunner(verbosity=self.verbose)
        runner.run(suite)

class DocoptsTest(unittest.TestCase):
    
    def __init__(self, doc, argv, expect):
        super(DocoptsTest, self).__init__()
        self.doc = doc
        self.argv = argv
        self.expect = expect
    
    def runTest(self):
        p = Popen(['./docopts'] + self.argv,
                  stdin=PIPE,
                  stdout=PIPE,
                  stderr=PIPE)
        stdout, stderr = p.communicate(self.doc.encode())
        dirname = stdout.decode().rstrip('\n')
        if not dirname:
            self.assertEqual(self.expect, 'user-error')
            return
        elif self.assertEqual(type(self.expect), dict):
            for arg, value in self.expect.items():
                path = os.path.join(dirname, arg)
                if value is None:
                    self.assertFalse(os.path.exists(path))
                elif isinstance(value, bool):
                    self.assertEqual(value, os.path.exists(path))
                elif isinstance(value, int):
                    with open(path, 'r') as f:
                        self.assertEqual(value, int(f.read()))
                elif isinstance(value, str):
                    with open(path, 'r') as f:
                        self.assertEqual(value, f.read())
                elif isinstance(value, list):
                    for i,v in enumerate(value):
                        subpath = os.path.join(path, str(i+1))
                        with open(subpath, 'r') as f:
                            self.assertEqual(v, f.read())
        shutil.rmtree(dirname)
