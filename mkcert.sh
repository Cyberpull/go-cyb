set -e

# Make SSL Certificate
[[ ! -d certs ]] && mkdir certs

# Create Cert File
echo "$CRT_DATA" > certs/crt.pem

# Create Key File
echo "$KEY_DATA" > certs/key.pem
