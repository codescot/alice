#!/bin/sh
echo '{"from": "A", "to": "B" }' | sam local invoke --event - HugFunction
