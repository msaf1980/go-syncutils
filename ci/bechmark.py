#!/usr/bin/env python

import argparse
import subprocess
import re
import os
import sys


def parse_cmdline():
    parser = argparse.ArgumentParser(
        description='Detect source changes and run benchmark for changed part (or all, if base package changed)')

    parser.add_argument('-b', '--base', dest='base', action='store', type=str,
                        help='base commit')

    parser.add_argument('-u', '--head', dest='head', action='store', type=str, default='HEAD',
                        help='head commit')

    parser.add_argument('-f', '--files', dest='files', action='store', type=str,
                        help='file with changes (relative paths) list')

    parser.add_argument('-d', '--dir', dest='dirs', action='append', type=str,
                        help='dirs for force add to tests')

    parser.add_argument('-c', '--count', dest='count', action='store', type=int, default=6,
                        help='run count')

    return parser.parse_args()


def addDir(file, dirs, forceDirs, forceFiles, sourcesRegexp):
    if file != '':
        if file in forceFiles:
            return True

        dir = os.path.dirname(file)
        if os.path.isdir(dir):
            for s in sourcesRegexp:
                if s.match(file):
                    dirs.add('./' + dir)
                    for b in forceDirs:
                        if file.startswith(b):
                            # this will trigger all tests
                            return True

    return False


def main():
    args = parse_cmdline()

    # regexp for detect files, which will trigger run benchmark
    sources = [
        re.compile('^.*\.go$'),
        re.compile('^go.mod$'),
        re.compile('^go.sum$')
    ]
    # dirs, which will be trigger all benchmark (not changed dir)
    forceDirs = []
    # files, which will be trigger all benchmark
    forceFiles = {'go.mod', 'go.sum'}
    # dirs for all benchmark
    baseTests = []

    dirs = set()

    isBase = False

    if args.files is None:
        if args.base is None and args.dirs is None:
            sys.exit("base commit not set")
        if args.head is None and args.dirs is None:
            sys.exit("head commit not set")

        if not args.head is None and not args.base is None:
            p = subprocess.run(
                ['git', 'diff', '--name-only', args.base, args.head],
                stdout=subprocess.PIPE,
            )
            for file in p.stdout.decode().split("\n"):
                if addDir(file, dirs, forceDirs, sources):
                    isBase = True
    else:
        with open(args.files) as f:
            for file in f:
                if addDir(file, dirs, forceDirs, forceFiles, sources):
                    isBase = True

    # prerare test dirs list
    tests = []
    if isBase:
        tests = baseTests
    else:
        for d in dirs:
            tests.append(d)
        if not args.dirs is None:
            for d in args.dirs:
                tests.append('./' + d)
        tests.sort()

    ranges = []
    for i in range(1, args.count+1):
        ranges.append(str(i))


    if len(tests) > 0:
        command = \
            "for i in %s; do"\
            "  echo STEP ${i} ; go test -benchmem -run=^$ -bench '^Benchmark' %s || exit 1; "\
            "done" % (' '.join(ranges), ' '.join(tests))

        sys.stderr.write(command+"\n")
        p = subprocess.Popen([command], shell=True)
        try:
            p.wait()
        except KeyboardInterrupt:
            try:
                p.terminate()
            except OSError:
                pass
            p.wait()


if __name__ == "__main__":
    main()
