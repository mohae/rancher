#!/usr/bin/env bash
# Update the system and install basics
apt-get update -y -q > /dev/null
apt-get upgrade -y -q > /dev/null
apt-get install -y -q curl wget git vim rsync openssh sudo
apt-get intsall -y gcc make gcc-c++ kernel-devel-`uname-r`
