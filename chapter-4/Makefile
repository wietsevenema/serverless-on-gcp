FUNCTIONS = background frontend

.PHONY: config deploy $(FUNCTIONS)

config:
	./configcheck.sh 

deploy: config $(FUNCTIONS)

$(FUNCTIONS):
	cd $@ && make deploy
