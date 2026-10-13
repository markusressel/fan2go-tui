#!/usr/bin/env python3
"""Prepares a recording of the demo for agg, which renders it as a GIF.

The shortcut bar separates the key and the name of each entry with a non-breaking space, which agg draws as a
stray glyph with some fonts. It looks the same as a space in a terminal, so it is replaced by one.

usage: prepare_cast.py <input .cast> <output .cast>
"""
import json
import sys


def main():
    source, target = sys.argv[1], sys.argv[2]
    with open(source) as input_file, open(target, "w") as output_file:
        for number, line in enumerate(input_file):
            if number > 0 and line.strip():
                event = json.loads(line)
                if event[1] == "o":
                    event[2] = event[2].replace(" ", " ")
                line = json.dumps(event) + "\n"
            output_file.write(line)


if __name__ == "__main__":
    main()
