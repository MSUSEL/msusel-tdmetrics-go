.PHONY: help go-test js-test td-test test go-clean js-clean td-clean clean

help:
	@echo Help for MSU SEL TD Metrics - Go!
	@echo.
	@echo `make clean` will clean all the parts. To run individual parts use go-clean, js-clean, or td-clean.
	@echo.
	@echo `make test` will run all the tests for all parts. To run individual parts use go-test, js-test, or td-test.
	@echo.

test: go-test js-test td-test

go-test:
	cd ./goAbstractor/ && go test -count=1 ./...

td-test:
	cd ./techDebtMetrics/ && dotnet test --filter StubTest0011
# Add `--filter StubTest0007` to run specific tests

clean: go-clean js-clean td-clean

go-clean:
	cd ./goAbstractor/ && go clean -cache

td-clean:
	cd ./techDebtMetrics/ && dotnet clean
