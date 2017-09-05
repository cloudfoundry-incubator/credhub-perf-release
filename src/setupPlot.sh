#!/bin/bash

sudo pip install virtualenv
virtualenv --no-site-packages headroomplot
cd headroomplot
source bin/activate
pip install -r requirements.txt
