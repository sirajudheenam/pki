#!/bin/env bash
echo "CURRENT : $(pwd)"
echo 
ls scripts/*
ls /home/runner/work/pki/pki/*
ls /home/runner/work/pki/pki/pki-go/*
ls /home/runner/work/pki/pki/pki-go/certs/*
ls /home/runner/work/pki/pki/pki-go/certs/localhost/*
ls /home/runner/work/pki/pki/demo-pki/localhost/*
ls /home/runner/work/pki/pki/demo-pki/localhost/server/*
ls /home/runner/work/pki/pki/demo-pki/localhost/root/*
mkdir -p $(pwd)/pki-go/certs/localhost/server
mkdir -p $(pwd)/pki-go/certs/localhost/client
cp $(pwd)/demo-pki/localhost/client/* $(pwd)/pki-go/certs/localhost/client/
cp $(pwd)/demo-pki/localhost/server/* $(pwd)/pki-go/certs/localhost/server/
cp $(pwd)/demo-pki/localhost/root/* $(pwd)/pki-go/certs/localhost/server/
ls /home/runner/work/pki/pki/*
ls /home/runner/work/pki/pki/pki-go/*
ls /home/runner/work/pki/pki/pki-go/certs/localhost/*
ls /home/runner/work/pki/pki/pki-go/certs/localhost/server/*
ls /home/runner/work/pki/pki/pki-go/certs/localhost/client/*