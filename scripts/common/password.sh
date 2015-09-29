#!/usr/bin/env bash

function generate_password() {
  echo $(LC_CTYPE=C tr -dc A-Za-z0-9 < /dev/urandom | fold -w 32 | head -n 1)
}
