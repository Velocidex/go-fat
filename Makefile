all:
	go build -o fat bin/*.go


windows:
	GOOS=windows GOARCH=amd64 \
            go build -ldflags="-s -w" \
	    -o fat.exe ./bin/*.go

generate:
	cd parser/ && binparsegen conversion.spec.yaml > fat_gen.go


test:
	go test -v ./...
