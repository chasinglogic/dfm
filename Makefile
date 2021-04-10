PYTHON := python3
DFM_BIN := $(shell which dfm)

lint:
	$(PYTHON) -m pydocstyle src
	$(PYTHON) -m pylint src
	$(PYTHON) -m pylint --disable=C0116,W0621 tests

fmt:
	$(PYTHON) -m black src tests
	$(PYTHON) -m isort src tests

clean:
	rm -rf build dist
	rm -rf {} **/*.egg-info
	rm -f **/*.pyc

install:
	$(PYTHON) setup.py install

install-dev:
	pip install --editable .
	pip install -r requirements.dev.txt

test:
	PYTHONPATH="$$PYTHONPATH:src" $(PYTHON) -m pytest -v -s -m 'not slow'

test-all:
	PYTHONPATH="$$PYTHONPATH:src" $(PYTHON) -m pytest --disable-pytest-warnings
	bash ./scripts/integration_tests.sh -b $(DFM_BIN)

publish: clean
	python setup.py sdist bdist_wheel
	twine upload dist/*
