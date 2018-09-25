PYTHON := python3

lint:
	$(PYTHON) -m pydocstyle taskforge
	$(PYTHON) -m pylint taskforge tests

fmt:
	$(PYTHON) -m yapf --recursive -i taskforge tests

clean:
	rm -rf build dist
	rm -rf {} **/*.egg-info
	rm -f **/*.pyc

install:
	$(PYTHON) setup.py install

install-dev:
	pip install --editable .
	pip install --editable ".[mongo]"
	pip install -r requirements.dev.txt

test:
	PYTHONPATH="$$PYTHONPATH:src" $(PYTHON) -m pytest -m 'not slow'

test-all:
	PYTHONPATH="$$PYTHONPATH:src" $(PYTHON) -m pytest --disable-pytest-warnings

publish: clean
	python setup.py sdist bdist_wheel
	twine upload dist/*

# You can set these variables from the command line.
SPHINXOPTS    =
SPHINXBUILD   = sphinx-build
SOURCEDIR     = src/docs
BUILDDIR      = build/docs

.PHONY: docs
docs: html
	rm -rf docs/*
	mv build/docs/html/* docs/
	mv docs/_static/* docs/
	rm -rf docs/_static

# Put it first so that "make" without argument is like "make help".
help:
	@$(SPHINXBUILD) -M help "$(SOURCEDIR)" "$(BUILDDIR)" $(SPHINXOPTS) $(O)

.PHONY: help Makefile

livehtml:
	sphinx-autobuild --watch ../taskforge -b html $(SPHINXOPTS) "$(SOURCEDIR)" $(BUILDDIR)/html

# Catch-all target: route all unknown targets to Sphinx using the new
# "make mode" option.  $(O) is meant as a shortcut for $(SPHINXOPTS).
%: Makefile
	@$(SPHINXBUILD) -M $@ "$(SOURCEDIR)" "$(BUILDDIR)" $(SPHINXOPTS) $(O)
