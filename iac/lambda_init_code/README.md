# Base Lambda

This is a base code for Lambdas using `Go 1.x`.
This code is mainly required to bootstrap a Lambda function with runtime `Go 1.x`, and you don't want to use Terraform to deploy your code.

To generate the excutable, just run:

Linux:

```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/
```

Windows (Powershell):

```powershell
$env:GOOS='linux'; $env:GOARCH='amd64'; $env:CGO_ENABLED=0; go build -o main
```
