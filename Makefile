NAME := mongo-db-admin

.PHONY: tools build

tools:
	cd ~ && go get -u github.com/aws/aws-lambda-go/cmd/build-lambda-zip && cd -

build:
	mkdir -p ./aws/tf/files
	GOOS=linux GOARCH=amd64 go build -o ./aws/tf/files/mongo-db-admin ./aws/

package: tools
	cd ./aws/tf \
		&& build-lambda-zip --output files/mongo-db-admin.zip files/mongo-db-admin \
		&& tar -cvzf aws-module.tgz *.tf files/mongo-db-admin.zip \
		&& mv aws-module.tgz ../../