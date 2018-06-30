#!/usr/bin/env python3

import os
import sys
import subprocess


def color2hex(col):
    # col=0 is darkest, col=29 is brightest
    return '{:02x}{:02x}{:02x}'.format(((col + 3) // 3) * 24, ((col + 2) // 3) * 24, ((col + 1) // 3) * 24)

def svg_to_png(infp, outfp):
	proc = subprocess.Popen(['inkscape', '-z', '--export-png', outfp, infp], stderr=subprocess.PIPE, stdout=subprocess.PIPE)
	err, out = proc.communicate()
	err, out = err.decode('utf-8'), out.decode('utf-8')
	if out:
		print(out)
	if err and proc.returncode != 0:
		print(err)


def get_rules(t):
    rules = {}
    for line in t.splitlines():
        if '->' not in line:
            continue
        key, v = line.split('->')
        rows = v.strip().split(',')
        cells = [[0] * 3 for _ in range(3)]
        for r, row in enumerate(rows):
            for c, col in enumerate(row.strip().split()):
                cells[r][c] = color2hex(int(col.strip()))

        rules[color2hex(int(key.strip()))] = cells
    return rules

def main(infile, basefile):
    with open(infile) as fd:
        rules = get_rules(fd.read())

    with open(basefile) as fd:
        base = fd.read()

    m = {0: 0, 1: 127, 2: 255}
    for match, cells in rules.items():
        src = base
        for r, row in enumerate(cells):
            for c, col in enumerate(row):
                template = ':#{:02x}{:02x}00'.format(m[r], m[c])
                repl = ':#' + col
                src = src.replace(template, repl)

        filename = 'rule-{}'.format(match.upper())
        with open(filename + 'FF.svg', 'w') as fd:
            fd.write(src)

        svg_to_png(filename + 'FF.svg', filename + 'FF.png')


if __name__ == '__main__':
    if len(sys.argv) != 3:
        print('usage: {} <rulesfile.txt> <basefile.svg>'.format(sys.argv[0]))
        sys.exit(1)

    main(sys.argv[1], sys.argv[2])
