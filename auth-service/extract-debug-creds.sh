# These are development ONLY credentials
# NOTE: these decryption keys are visible to
# aid developers with **development** environment only. 
# these are not sensitive becuase these can only be used 
# for accessing development
# other environment credentials are encrypted via cicd-cred-vN.tgz.enc
# and these decryption values are only known to the "buildEnv" such as travis
# as they are input manually via the secret manager for the cicd tool 
# (ie: travis.com env var secrets console)
CICD_CRED_KEY_V29=412c55132641886729d815e68cbdb53a294b49a186f56be83ecd16d0314adcfe
CICD_CRED_IV_V29=c1d607cfad8d34398aa345ee7937f1d0

# Run decryption using: 
sh ./secret-mgmt/cicd/decrypt.sh ./cicd-cred-debug-v29.tgz.enc 412c55132641886729d815e68cbdb53a294b49a186f56be83ecd16d0314adcfe c1d607cfad8d34398aa345ee7937f1d0 # end
