go install github.com/swaggo/swag/cmd/swag@latest
swag init -ot yaml,json
npx --yes @redocly/cli build-docs -o docs/swagger.html docs/swagger.yaml
