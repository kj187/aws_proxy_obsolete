build_dev:
	docker build --target dev . -t go
run_dev:
	docker run -it -v ${PWD}:/work go sh
