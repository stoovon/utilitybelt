.PHONY: clean compile compiled/plugins gomod

compile: gomod compiled/utilitybelt compiled/plugins

clean:
	@rm -rf ./compiled
	@rm -f utilitybelt

gomod:
	@go mod vendor

compiled/utilitybelt:
	@go build -o ./compiled/utilitybelt .
	@ln -s ./compiled/utilitybelt utilitybelt

compiled/plugins:
	@ls ./plugin | xargs -I {} go build -o ./compiled/plugins/{} ./plugin/{}