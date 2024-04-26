# go get -d github.com/99designs/gqlgen@v0.17.36
GOTOOLCHAIN=go1.22.2 go run github.com/99designs/gqlgen@latest generate --config ./user.gqlgen.yml

# go get -d github.com/99designs/gqlgen@v0.17.36
GOTOOLCHAIN=go1.22.2 go run github.com/99designs/gqlgen@latest generate --config ./admin.gqlgen.yml