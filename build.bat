set GOOS=freebsd
go build -o bin/freebsd/ntlm_auth_wrapper
set GOOS=linux
go build -o bin/linux/ntlm_auth_wrapper
