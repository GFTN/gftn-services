echo "now running travis emulator scripts"
echo "node version" ; node -v
echo "npm version" ; npm -v

# cd into auth-service
cd gftn-services/auth-service

# display git version
git --version

# get newly added branches
git fetch --all

# checkout relevant branch
git checkout ${branch}

# pull changes
git pull

pwd ; ls

sh ./cicd/travis/emulator/scripts/travis-emulation.sh