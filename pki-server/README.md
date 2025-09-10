
```bash

ls -l demo-pki/root/root.cert.pem
ls -l demo-pki/certs/client.cert.pem
ls -l demo-pki/private/client.key.pem

curl https://localhost:4433/hello \
  --cacert demo-pki/root/root.cert.pem \
  --cert demo-pki/client/client.cert.pem \
  --key demo-pki/client/client.key.pem


  curl https://localhost:4433 \
  --cacert demo-pki/root/root.cert.pem \
  --cert demo-pki/client/client.cert.pem \
  --key demo-pki/client/client.key.pem
```

