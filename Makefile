PYTHON := python
DFM_BIN := $(shell which dfm)

lint:
	$(PYTHON) -m pydocstyle src
	$(PYTHON) -m pylint src
	$(PYTHON) -m pylint --disable=C0116,W0621,R1732 tests

fmt:
	$(PYTHON) -m black src tests
	$(PYTHON) -m isort src tests

clean:
	rm -rf build dist
	rm -rf {} **/*.egg-info
	rm -f **/*.pyc

install:
	pipenv install --editable .

setup: install
	pipenv install --dev

test:
	PYTHONPATH="$$PYTHONPATH:src" $(PYTHON) pytest -v -s -m 'not slow'

test-all:
	PYTHONPATH="$$PYTHONPATH:src" $(PYTHON) pytest --disable-pytest-warnings
	bash ./scripts/integration_tests.sh -b $(DFM_BIN)

publish: clean
	pipenv run python setup.py sdist bdist_wheel
	pipenv run twine upload dist/*
