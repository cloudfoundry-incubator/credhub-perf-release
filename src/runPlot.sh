#!/bin/bash

cd headroomplot
source bin/activate
pip install -r requirements.txt
python headroomplot.py $1
deactivate
