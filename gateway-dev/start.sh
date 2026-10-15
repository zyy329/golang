#!/bin/bash
rm -rf run
mkdir run run/logs
nginx -p $PWD/run -c ../conf/nginx.conf
tail -f run/logs/error.log