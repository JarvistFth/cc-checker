BUILD = GOPATH=${PWD}/../..:$(GOPATH) go build

all:
	${BUILD} -o contract_audit audit/audit.go
	${BUILD} -o contract_face face/service.go
	${BUILD} -o contract_identify identify/service.go
	${BUILD} -o contract_location location/service.go
	${BUILD} -o contract_time time/service.go
