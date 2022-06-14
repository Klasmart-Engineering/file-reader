set -x
awslocal s3 mb s3://organization
awslocal s3 mb s3://school
awslocal s3 mb s3://user
awslocal s3 mb s3://class
awslocal s3 mb s3://organization.membership
awslocal s3 mb s3://school.membership
set +x