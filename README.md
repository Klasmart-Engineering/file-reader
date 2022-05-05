# file-reader

## localstack
### awscli 
#### mac
1. Install homebrew `/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`
2. Install awscli `brew install awscli`

Start up localstack via docker-compose
`docker-compose up -d`

Ensure the env variables have been exported via `direnv`.

Exposed port is `4566`, you're ready to use localstack e.g.
`aws --endpoint-url=http://localhost:4566 s3 mb s3://mytestbucket`