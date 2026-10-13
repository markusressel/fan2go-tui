#!/usr/bin/env python3
"""Runs fan2go-tui in a pty and types the keys of the demo, for "just demo-record" (asciinema rec --headless).

The output of the application is passed through, so asciinema records it. The application runs with "demo --ui",
so it shows the simulated fans, curves and sensors of the demo server, and no fan2go daemon is needed.

usage: driver.py <binary>
"""
import fcntl
import os
import pty
import select
import struct
import sys
import termios
import time

UP, DOWN, RIGHT, LEFT = "\x1b[A", "\x1b[B", "\x1b[C", "\x1b[D"
TAB, ENTER, ESC = "\t", "\r", "\x1b"
CTRL_Q = "\x11"

# the keys of the demo: (seconds to wait before, keys)
SCRIPT = [
    # a shorter refresh interval, so the graphs fill up quicker
    (3.0, "-"), (0.5, "-"), (0.5, "-"),
    (5.0, DOWN), (3.0, DOWN), (3.0, DOWN),
    (3.0, "2"), (4.0, DOWN), (3.0, DOWN),
    (3.0, "3"), (4.0, DOWN), (3.0, DOWN),
    (3.0, TAB),
    (4.0, CTRL_Q),
]


def window_size():
    try:
        return fcntl.ioctl(sys.stdout.fileno(), termios.TIOCGWINSZ, b"\0" * 8)
    except OSError:
        return struct.pack("HHHH", 42, 179, 0, 0)


def main():
    binary = os.path.abspath(sys.argv[1])
    pid, fd = pty.fork()
    if pid == 0:
        # port 0: any free port, so a running fan2go daemon (or demo server) is not in the way
        os.execv(binary, [binary, "demo", "--ui", "--fresh", "--port", "0"])
    fcntl.ioctl(fd, termios.TIOCSWINSZ, window_size())

    steps = list(SCRIPT)
    next_at = time.monotonic() + steps[0][0]
    out = sys.stdout.buffer
    while True:
        timeout = max(next_at - time.monotonic(), 0) if steps else 0.1
        ready, _, _ = select.select([fd], [], [], timeout)
        if ready:
            try:
                data = os.read(fd, 65536)
            except OSError:
                break
            if not data:
                break
            out.write(data)
            out.flush()
        if steps and time.monotonic() >= next_at:
            _, keys = steps.pop(0)
            os.write(fd, keys.encode())
            if steps:
                next_at = time.monotonic() + steps[0][0]
    os.waitpid(pid, 0)


if __name__ == "__main__":
    main()
