NAME := mongo-db-admin

.PHONY: tools build

tools:
	cd ~ && go install github.com/aws/aws-lambda-go/cmd/build-lambda-zip@latest && cd -

build:
	mkdir -p ./aws/tf/files
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o ./aws/tf/files/bootstrap ./aws/

package: tools
	cd ./aws/tf \
		&& build-lambda-zip --output files/mongo-db-admin.zip files/bootstrap \
		&& tar -cvzf aws-module.tgz *.tf files/mongo-db-admin.zip \
		&& mv aws-module.tgz ../../

acc: acc-up acc-run acc-down

acc-up:
	cd acc && docker-compose -p mongo-db-admin-acc up -d db

acc-run:
	ACC=1 gotestsum ./acc/...

acc-down:
	cd acc && docker-compose -p mongo-db-admin-acc down
