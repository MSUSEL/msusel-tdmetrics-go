 .PHONY: help go-test j-test td-test test go-clean j-clean td-clean clean

help:
	@echo Help for MSU SEL TD Metrics - Go!
	@echo.
	@echo `make clean` will clean all the parts. To run individual parts use go-clean, j-clean, or td-clean.
	@echo.
	@echo `make test` will run all the tests for all parts. To run individual parts use go-test, j-test, or td-test.
	@echo.

test: go-test j-test td-test

go-test:
	cd ./goAbstractor/ && go test -count=1 ./...

j-test:
	cd ./javaAbstractor/ && mvn
# Add `-Dtest=TestClass#testMethod` to run specific tests

td-test:
	cd ./techDebtMetrics/ && dotnet test
# Add `--filter StubTest0007` to run specific tests

clean: go-clean j-clean td-clean

go-clean:
	cd ./goAbstractor/ && go clean -cache

j-clean:
	cd ./javaAbstractor/ && mvn clean

td-clean:
	cd ./techDebtMetrics/ && dotnet clean
