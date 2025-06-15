# Generate TLS 100 years

```bash
# Change DNS name `mail` as your hostname

openssl req -x509 -newkey rsa:4096 \                                                                                                                                                                                                           ─╯
    -nodes -keyout key.pem -out cert.pem \
    -sha256 -days 36500 \
    -addext "subjectAltName = DNS:mail"
```

Note: TLS files [cert.pem](cert.pem) and [key.pem](key.pem) for hostname `mail`
