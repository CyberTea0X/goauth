npx --yes @redocly/cli build-docs -o docs/openapi.html docs/openapi.yaml &
npx --yes @redocly/cli build-docs -o docs/external.html docs/external.yaml &
wait
echo "Docs are be built"
