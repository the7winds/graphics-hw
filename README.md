# README

A simple application to draw a fractal (Julia set).
There are some settings:
* number of iterations
* complex seed

## Build
```
git submodule --init --recursive
cmake .
make -j 4
```
## Run
Before running the application export path to shaders directory via SHADERS_DIR variable.