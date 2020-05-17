run: clean build
	./buggy

build:
    go build ./Models/MSG.go
    go build mailconnector.go
	go build buggy.go

clean:
	rm buggy || true
