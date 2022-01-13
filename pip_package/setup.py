import os
import re

from setuptools import setup


setup(
    name='tfacon',
    version="1.0.9",
    description="tfacon",
    author="Red Hat Inc",
)
os.system('source scripts/install_tfacon.sh')
