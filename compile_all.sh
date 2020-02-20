GOOS=darwin go build -o solrcli
GOOS=linux go build -o solrcli_linux
GOOS=windows GOARCH=386 go build -o solrcli.exe