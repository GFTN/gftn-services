# PURPOSE:
# The purpose of this script is to deploy the firebase 
# portal, database secruity rules and trigger functions.
# Only works locally not for travis. See deploy.sh for 
# full deployment through CI/CD pipeline.  


# usage example `$ sh -x deploy_firebase.sh "dev-2-c8774" "1/-mmRFBlvGxxxxxxxxxxxxxxxx"`

# navigate to portal
cd ../gftn-web ;

# get token by running `$ firebase login:ci`
ci_token=$2

# get firebase project
gcloud_project=$1

# deploy firebase database rules
firebase-bolt database.rules.bolt ; firebase deploy --only database --project="$gcloud_project" --token=$ci_token

# deploy firebase functions
firebase deploy --only functions --project="$gcloud_project" --token=$ci_token

# deploy firebase portal (hosting)
firebase deploy --only hosting --project="$gcloud_project" --token=$ci_token
