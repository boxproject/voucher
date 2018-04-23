#!/usr/bin/env bash
# call this script with an email address (valid or not).
# must install openssl before run this script.
# like:
# ./makecert.sh demo@random.com

set -e
FULLPATH=$1

rm -f ${FULLPATH}/server.pem ${FULLPATH}/server.key
SUBJECT="/C=CN/ST=Shanghai/L=Earth/O=BOX/OU=DC/CN=box.la/emailAddress"

EMAIL=${2:-develop@2se.com}
DAYS=${3:-3650}

openssl req -new -nodes -x509 -out ${FULLPATH}/server.pem -keyout ${FULLPATH}/server.key -days ${DAYS} -subj "${SUBJECT}=${EMAIL}"