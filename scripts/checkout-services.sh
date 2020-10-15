# This file checks out the `micro/services` repo to the './test/services'
# folder, and makes the same modifications as what the github
# workflow files do so test can run those services.

rm -rf test/services
cd test
git clone https://github.com/micro/services
cd services
rm go.mod
rm go.sum
grep -rl github.com/micro/services . | xargs sed -i 's/github.com\/micro\/services/github.com\/micro\/micro\/v3\/test\/services/g'
rm -rf .git