# IMPORTANT: edits to this file must be pushed to git 
# for the targeted branch to test. Otherwise the docker 
# container will not be aware of the changes. 

# ======= emulate travis build auth-service =========:
# emulating make docker from auth-service makefile:
bash install_dependencies.sh --gcloud
bash ./deployment/docker/create-images.sh any-version --dev-only --skip-docker

# ======= emulate travis deploy auth-service =========:
# run deployment scripts
bash ./deployment/deploy.sh dev --debug