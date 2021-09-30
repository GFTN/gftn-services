# NOTE: install dependencies first via $  sh install_dependencies.sh

checkerrors(){

    errors=$?

    # helps to alert failure if code does not build properly
    # and helps vscode debuger exit on typscript/lint error
    if [ $errors -ne 0 ]
    then
        echo $(tput setaf 1)"======= errors: "$errors" ======="
        echo $(tput setaf 1)"ERROR: build process error, see run-tests.sh"
        exit 1
    fi  
    
}

# run application secret decryption tests 
node secret-mgmt/node_modules/mocha/bin/_mocha -r ts-node/register --project secret-mgmt/tsconfig.tests.json --colors --timeout 60000 secret-mgmt/src/index.spec.ts ; checkerrors
# tsc secret-mgmt/src/index.spec.ts ; checkerrors