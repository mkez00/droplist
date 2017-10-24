# droplist

Go application that configures firewall to block IP addresses from Spamhaus DROP list

# Outstanding items
1) Clear previously run entries before running next batch (this will make utility idempotent)
2) General code cleanup

# Build Instructions

1) Linux AMD64: `env GOOS=linux GOARCH=amd64 go build`

# Download

Binary download for Linux AMD64: [a relative link](/bin/droplist.zip)