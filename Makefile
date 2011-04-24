all:
	$(MAKE) -C src/

clean:
	$(MAKE) -C src/ clean

.phony: clean all
