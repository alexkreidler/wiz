#!/usr/bin/env bash

file=processors.apib

echo Starting aglio to serve $file

# Run a live preview server on http://localhost:3000/
aglio -i $file -s