#! /usr/bin/env python3
import sys, json, re, os
from subprocess import Popen, PIPE, STDOUT

def compare_json(expect, result):
    try:
        py_result = json.loads(result)
        py_expect = json.loads(expect)
    except:
        return "BAD_JSON"

    # string compare
    if type(py_expect) is str:
        if type(py_result) is str:
            if py_expect == py_result:
                return "OK"
            else:
                return "ERROR_STRING_DIFFER"
        else:
            return "ERROR_TYPE_DIFFER: string expected got %s" % type(py_result)

    #print(sorted(py_expect.keys()))
    #print(sorted(py_result.keys()))

    # compare json key by key
    return_val = "OK"
    for k in py_expect.keys():
        if k in py_result:
            if py_result[k] == py_expect[k]:
                del py_result[k]
                continue
            else:
                return_val = "ERROR_VALUE_DIFFER: for key %s" % k
                break
        else:
            return_val = "ERROR_KEY_MISSING: %s" % k
            break

    if return_val == "OK" and len(py_result) > 0:
        return_val = "ERROR_TOO_MANY_KEYS: %s" % ",".join(py_result.keys())

    return return_val


## main

fixtures = open(os.path.join(os.path.dirname(__file__), 'testcases.docopt'), 'r').read()

# remove comments
fixtures = re.sub('#.*$', '', fixtures, flags=re.M)

testee = (sys.argv[1] if len(sys.argv) >= 2 else
        exit('Usage: language_agnostic_tester.py ./path/to/executable/testee [ID ...]'))
ids = [int(x) for x in sys.argv[2:]] if len(sys.argv) > 2 else None
summary = ''

index = 0
for fixture in fixtures.split('r"""'):
    doc, _, body = fixture.partition('"""')
    for case in body.split('$')[1:]:
        index += 1
        if ids is not None and index not in ids:
            continue
        argv, _, expect = case.strip().partition('\n')
        prog, _, argv = argv.strip().partition(' ')
        assert prog == 'prog', repr(prog)
        p = Popen(testee + ' ' + argv,
                  stdout=PIPE, stdin=PIPE, stderr=STDOUT, shell=True)
        result = p.communicate(input=doc.encode('utf-8'))[0]

        compare = compare_json(expect, result)

        if compare == 'OK':
            summary += '.'
        elif compare.startswith('ERROR'):
            print((' %d: FAILED ' % index).center(79, '='))
            print(compare)
            print('r"""%s"""' % doc)
            print('$ prog %s\n' % argv)
            print('result>', result)
            print('expect>', expect)
            summary += 'F'
        elif compare == 'BAD_JSON':
            summary += 'J'
            print ( (' %d: BAD JSON ' % index).center(79, '=') )
            print ('result>', result)
            print ('expect>', expect)


print ( (' %d / %d ' % (summary.count('.'), len(summary))).center(79, '=') )
print (summary)
