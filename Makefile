all: go-test td-test test
.PHONY: all

test: go-test td-test

go-test:
	cd ./goAbstractor/ && go test ./...

td-test:
	cd ./techDebtMetrics/ && dotnet test
