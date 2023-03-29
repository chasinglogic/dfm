"""
Collection of DFM specific exceptions.

This will be raised by DFM internals whenever there is a validation error or
something dfm-specific goes wrong.
"""


class DFMException(Exception):
    """Represents a DFM internal exception."""


class MappingException(DFMException):
    """
    Raised during the Mapping process.

    Largely raised by linking directories instead of files since there is more
    that can go wrong there. This is basically a validation error.
    """
