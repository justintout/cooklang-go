# cooklang-go
> A Cooklang parser in Go. 

[![CI](https://github.com/justintout/cooklang-go/actions/workflows/workflow.yaml/badge.svg)](https://github.com/justintout/cooklang-go/actions/workflows/workflow.yaml) | [![Canonical Tests](https://github.com/justintout/cooklang-go/actions/workflows/canonical.yaml/badge.svg)](https://github.com/justintout/cooklang-go/actions/workflows/canonical.yaml) | [![Go Reference](https://pkg.go.dev/badge/github.com/justintout/cooklang-go.svg)](https://pkg.go.dev/github.com/justintout/cooklang-go)

## Usage

See the [`pcook` executable](./cmd/pcook/) for usage.

## Development

Issues and pull requests welcome. 

### Testing

`canonical.yaml` is updated manually. It should be in lockstep with the [official canonical tests](https://github.com/cooklang/spec/tree/main/tests).

```bash
go test ./...
```

## External References

- [Cooklang](https://cooklang.org/)
- [aquilax/cooklang-go](https://github.com/aquilax/cooklang-go) - another Cooklang parser in Go, from which I borrowed the test JSON (thank you!)
- [Lexical Scanning in Go, Rob Pike, 2011](https://talks.golang.org/2011/lex.slide#1) - basis for the lexer