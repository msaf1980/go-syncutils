#!/usr/bin/env python

import argparse
import subprocess
import re
import os



def parse_cmdline():
    parser = argparse.ArgumentParser(
        description='Detect source changes')

    parser.add_argument('-b', '--base', dest='base', action='store', type=str, required=True,
                        help='base commit')

    parser.add_argument('-u', '--head', dest='head', action='store', type=str, default='HEAD',
                        help='head commit')

    return parser.parse_args()


def main():
    args = parse_cmdline()

    # regexp for detect files, which will trigger run benchmark
    sources = [
        re.compile('^.*\.go$'),
        re.compile('^go.mod$'),
        re.compile('^go.sum$')
    ]

    p = subprocess.run(
        ['git', 'diff', '--name-only', args.base, args.head],
        stdout=subprocess.PIPE,
    )
    for file in p.stdout.decode().split("\n"):
        for s in sources:
            if s.match(file):
                print(file)
                break


if __name__ == "__main__":
    main()
