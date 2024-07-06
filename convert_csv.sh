#!/bin/bash

filename=$1

iconv -f utf-16le -t utf-8 "$filename" > result.csv
rm "$filename"
mv result.csv "$filename"