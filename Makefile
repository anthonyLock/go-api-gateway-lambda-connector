build: 
	go build -v ./...

test: 
	go test ./...

	
publish-patch:
	./publish.sh -v patch

publish-major:
	./publish.sh -v major

publish-minor:
	./publish.sh -v minor