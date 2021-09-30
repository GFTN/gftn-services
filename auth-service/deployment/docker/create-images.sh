#!/bin/bash

# defaults: 
TRAVIS="false" 
VERSION='non-official-version'

show_help()
{ 
cat << EOF

Create docker image for authentication-service:

--version (alias: -e) <local|dev|qa|st|prod> OPTIONAL
  Target auth-service environment to build deploymet files  
  DEFAULT: 'non-official-version'

EOF

}

# flags usage - https://archive.is/5jGpl#selection-709.0-715.1 or https://archive.is/TRzn4
while :; do
    case $1 in
        -h|-\?|--help)   # Call a "show_help" function to display a synopsis, then exit.
            show_help
            exit
            ;;
        --version) # Takes an option argument, ensuring it has been specified.
            if [ -n "$2" ]; then
                VERSION=$2
                shift
            else
                printf "$(tput setaf 1)ERROR: \"--version\" requires a non-empty option argument.\n" >&2
                exit 1
            fi
            ;;
        -v|--verbose)
            verbose=$((verbose + 1)) # Each -v argument adds 1 to verbosity.
            ;;
        --)              # End of all options.
            shift
            break
            ;;
        -?*)
            printf 'WARN: Unknown option (ignored): %s\n' "$1" >&2
            ;;
        *)               # Default case: If no more options then break out of the loop.
            break
    esac
    shift
done

checkerrors(){

    errors=$?

    # helps to alert failure if code does not build properly
    # and helps vscode debuger exit on typscript/lint error
    if [ $errors -ne 0 ]
    then
        echo $(tput setaf 1)"======= errors: "$errors" ======="
        echo $(tput setaf 1)"ERROR: build process error, see create-images.sh"
        exit 1
    fi  
    
}

# delete out existing build folder, if exists
rm -rf ./authentication/build || true

# NOTE: --env need not be specified as it can be anything since not required for docker build
bash ./deployment/build.sh --debug --cloud docker
checkerrors

echo "$(tput setaf 4)Building docker image version - $VERSION....$(tput bold)$(tput sgr0)"

# build docker image
docker build --build-arg BUILD_VERSION="$VERSION" -f deployment/docker/Dockerfile -t gftn/auth-service .
checkerrors

echo "Done."
