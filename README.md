# Securely serve a file over http

- serves the given file over http
- server prints qr code for random filename using https://github.com/mdp/qrterminal
- this should point to a random url with sufficient entropy (either use a uuid, or a 4-chunk id)
- scan in using phone
- server shuts down after serving exactly 1 request

## Usage:

```
$ qrserve Passwords.kdbx
```
