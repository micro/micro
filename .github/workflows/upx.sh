#!/usr/bin/env bash

for FILE in dist/micro_*/*; do
    du -sh ${FILE}
    upx ${FILE}
    du -sh ${FILE}
done
