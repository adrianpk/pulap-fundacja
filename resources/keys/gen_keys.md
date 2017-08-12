# Private key
$ openssl genrsa -out fundacja.rsa 1024
# Public key
$ openssl rsa -in fundacja.rsa -pubout > fundacja.rsa.pub
