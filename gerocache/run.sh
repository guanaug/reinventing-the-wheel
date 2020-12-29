#!/usr/bin/env bash

go build -o gerocache
./gerocache -addr 127.0.0.1:8000 -peers 127.0.0.1:8001,127.0.0.1:8002 &
./gerocache -addr 127.0.0.1:8001 -peers 127.0.0.1:8000,127.0.0.1:8002 &
./gerocache -addr 127.0.0.1:8002 -peers 127.0.0.1:8000,127.0.0.1:8001 &