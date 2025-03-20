package main

import (
	generator "github.com/lcarilla/sqlc-plugin-php-dbal/internal"
	"github.com/sqlc-dev/plugin-sdk-go/codegen"
)

func main() {
	codegen.Run(generator.Generate)
}
