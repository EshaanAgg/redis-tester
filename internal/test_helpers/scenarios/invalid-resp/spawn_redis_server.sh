#!/usr/bin/env -S python3 -u
import socket
import time
sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)

sock.bind(('', 6379))

sock.listen(1)

conn, cli_addr = sock.accept()
conn.send(b"SYNTAX_ERROR")
