#!/usr/bin/env python

from setuptools import setup, find_packages

setup(
    name='EffectiveCouscous',
    author='Maxime Landon',
    version='1.0',
    author_email='maximelandon@gmail.com',
    description='A structured set of transforms that interface the Metasploit Framework, and other tools.',
    license='GPLv3',
    packages=find_packages('src'),
    package_dir={'': 'src'},
    zip_safe=False,
    package_data={
        '': ['*.gif', '*.png', '*.conf', '*.mtz', '*.machine']  # list of resources
    },
    install_requires=[
        'canari>=3.3.9,<4'
    ],
    dependency_links=[
        # custom links for the install_requires
    ]
)
