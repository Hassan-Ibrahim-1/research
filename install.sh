#!/bin/bash
git clone https://github.com/ggerganov/llama.cpp\ncd llama.cpp
cd llama.cpp
mkdir build
cd build
cmake .. -DLLAMA_METAL=on
cmake --build . --config Release
