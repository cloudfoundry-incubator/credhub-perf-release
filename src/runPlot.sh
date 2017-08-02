#!/bin/bash

cd headroomplot
source bin/activate
pip3 install -r requirements.txt
python3 headroomplot.py $1
deactivate
