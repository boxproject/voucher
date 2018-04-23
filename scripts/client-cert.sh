#!/usr/bin/env bash
# call this script with an email address (valid or not).
# must install openssl before run this script.
# like:
# ./makecert.sh demo@random.com

FULLPATH=$1

rm -f ${FULLPATH}/client.pem ${FULLPATH}/client.key

SUBJECT="/C=CN/ST=Shanghai/L=Earth/O=BOX/OU=DC/CN=box.la/emailAddress"

EMAIL=${2:-develop@2se.com}
DAYS=${3:-3650}

openssl req -new -nodes -x509 -out ${FULLPATH}/client.pem -keyout ${FULLPATH}/client.key -days ${DAYS} -subj "${SUBJECT}=${EMAIL}"