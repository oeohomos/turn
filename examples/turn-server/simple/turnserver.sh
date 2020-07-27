#!/bin/bash

public_ip=192.168.0.108
users='everguard:coguard'
min_port=65000
max_port=65535

./turnserver \
  -public-ip=${public_ip} \
  -users=${users} \
  -min-port=${min_port} \
  -max-port=${max_port}
