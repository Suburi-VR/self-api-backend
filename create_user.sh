#!/bin/bash -eu

USER_POOL_ID="ap-northeast-1_Kjb4vUZPh"
TABLE_NAME="UserTable"

export USER_EMAIL=${1}
export USER_PASS=${2}
export USER_GIP=${3}
export USER_NAME=${4}
export USER_KANA=${5}
export USER_COMPANY=${6}
export USER_DEPT=${7}

echo "Create item file........."

echo '{' > item.json
echo "  \"username\": {\"S\": \"${USER_EMAIL}\"}," >> item.json
echo "  \"secret\": {\"S\": \"${USER_PASS}\"}," >> item.json
echo "  \"orgid\": {\"N\": \"${USER_GIP}\"}," >> item.json
echo "  \"nickname\": {\"S\": \"${USER_NAME}\"}," >> item.json
echo "  \"kana\": {\"S\": \"${USER_KANA}\"}," >> item.json
echo "  \"company\": {\"S\": \"${USER_COMPANY}\"}," >> item.json
echo "  \"department\": {\"S\": \"${USER_DEPT}\"}," >> item.json
echo "  \"anonflg\": {\"BOOL\": false}" >> item.json
echo '}' >> item.json

aws cognito-idp admin-create-user --user-pool-id ${USER_POOL_ID} --username ${USER_EMAIL} --user-attributes Name="custom:supporter",Value="true"

aws cognito-idp admin-set-user-password --user-pool-id ${USER_POOL_ID} --username ${USER_EMAIL} --password ${USER_PASS} --permanent



aws --endpoint http://localhost:8000 dynamodb put-item --table-name ${TABLE_NAME} --item file://item.json

exit 0