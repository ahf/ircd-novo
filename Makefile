all:
	$(MAKE) -C src/
	$(MAKE) -C utils/

clean:
	$(MAKE) -C src/ clean
	$(MAKE) -C utils/ clean

.phony: clean all
