bold=$(tput bold)
normal=$(tput sgr0)

get_email() {
    if  [[ ! -z "$1" ]]
    then
        curl -s -H "Authorization: Bearer $1" https://www.googleapis.com/oauth2/v2/userinfo\?fields\=email | python -c 'import sys, json; print json.load(sys.stdin)["email"]'
    fi
}

# Check google project 
GOOGLE_PROJECT=$(gcloud config get-value core/project 2> /dev/null)
if  [[ ! -z "$GOOGLE_PROJECT" ]]
then
    echo "Project: [${bold}${GOOGLE_PROJECT}${normal}]"
else
    echo "" 
    echo "Project not set, run \
${bold}gcloud config set project [PROJECT_ID]${normal}"
exit 1
fi

# Check run region
RUN_REGION=$(gcloud config get-value run/region 2> /dev/null)
if  [[ ! -z "$RUN_REGION" ]]
then
    echo "Cloud Run region: [${bold}${RUN_REGION}${normal}]"
else
    echo "" 
    echo "Cloud Run region not set, \
run ${bold}gcloud config set run/region us-central1${normal} \
or pick a region closer to you"
exit 1
fi


# Check gcloud login
GC_TOKEN=$(gcloud auth print-access-token 2> /dev/null)
if  [[ ! -z "$GC_TOKEN" ]]
then
    echo "gcloud logged in as: [${bold}$(get_email $GC_TOKEN)${normal}]"
else
    echo "" 
    echo "gcloud not logged as: \
run ${bold}gcloud auth login${normal}"
exit 1
fi

# Check application-default login
ADC_TOKEN=$(gcloud auth application-default print-access-token 2> /dev/null)
if  [[ ! -z "$ADC_TOKEN" ]]
then
    echo "gcloud application-default logged in as: [${bold}$(get_email $ADC_TOKEN)${normal}]"
else
    echo "" 
    echo "gcloud application-default not logged in: \
run ${bold}gcloud auth application-default login${normal}"
exit 1
fi

echo ""
echo "Do these setting look right? Run ${bold}make deploy${normal} to deploy the project." 



