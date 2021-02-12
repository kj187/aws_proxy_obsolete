build_dev:
	docker build --target dev . -t go
run_dev:
	docker run -it -v ${PWD}:/work go sh

gr_install:
	brew install goreleaser
gr_test:
	goreleaser --snapshot --skip-publish --rm-dist
gr_release:
	goreleaser release --rm-dist