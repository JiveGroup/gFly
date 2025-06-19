#!/bin/sh

lsof -i :7789 -sTCP:LISTEN | awk 'NR > 1 {print $2}' | xargs kill -9
