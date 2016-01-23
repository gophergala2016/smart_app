#!/bin/bash

while read -r line
do
  export line
done < <(cat .env)

go run -- *.go
