module github.com/dgframe/dg-validation

go 1.25.0

replace (
	github.com/dgframe/core => ../core
	github.com/dgframe/dg-http-validation/adapters/gookit => ./adapters/gookit
)

require github.com/dgframe/core v1.8.0
