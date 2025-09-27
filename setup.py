from setuptools import setup, find_packages

setup(
    name="dsl_bypass_framework",
    version="0.1.0",
    packages=find_packages(),
    install_requires=[
        "httpx",
        "pysqlite3",
    ],
    entry_points={
        'console_scripts': [
            'dslam-scanner=utils.network_scanner:main',
        ],
    },
)