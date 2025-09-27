# dsl_bypass_ultra/utils/logger.py

import logging
import sys

def setup_logger(name, level=logging.INFO):
    """
    Sets up a configured logger.

    This utility function creates a standardized logger to be used
    across the entire project.

    Tüm modüller tarafından kullanılacak.
    """

    logger = logging.getLogger(name)
    if logger.hasHandlers():
        return logger # Logger already configured

    logger.setLevel(level)

    # Create a handler
    handler = logging.StreamHandler(sys.stdout)

    # Create a formatter and set it for the handler
    formatter = logging.Formatter(
        '%(asctime)s - %(name)s - %(levelname)s - %(message)s'
    )
    handler.setFormatter(formatter)

    # Add the handler to the logger
    logger.addHandler(handler)

    return logger

# Example usage:
# from .logger import setup_logger
# log = setup_logger(__name__)
# log.info("This is an informational message.")