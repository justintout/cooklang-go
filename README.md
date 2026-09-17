# cooklang-go
> A Cooklang parser in Go. 

[![CI](https://github.com/justintout/cooklang-go/actions/workflows/ci.yaml/badge.svg)](https://github.com/justintout/cooklang-go/actions/workflows/ci.yaml) | [![Canonical Tests](https://github.com/justintout/cooklang-go/actions/workflows/canonical.yaml/badge.svg)](https://github.com/justintout/cooklang-go/actions/workflows/canonical.yaml) | [![Go Reference](https://pkg.go.dev/badge/github.com/justintout/cooklang-go.svg)](https://pkg.go.dev/github.com/justintout/cooklang-go)

## Usage

See the [`pcook` executable](./cmd/pcook/) for usage. It reads the latest
version of the spec, and takes a version for a recipe written before metadata
moved into YAML front matter and a step became a paragraph:

```bash
go run ./cmd/pcook -spec 5 old-recipe.cook
```

[`cookupgrade`](./cmd/cookupgrade/) rewrites such a recipe for the latest
version instead, if you would rather store the newer spelling.

## Development

Issues and pull requests welcome. 

### Testing

`canonical.yaml` is updated manually. It should be in lockstep with the [official canonical tests](https://github.com/cooklang/spec/tree/main/tests).

`canonical-v5.yaml` and `canonical-v6.yaml` are those tests as they stood while
those versions were current, pinned to their revisions ([9b857a2][v5] and
[8b26235][v6]). Both are read under the version 5 rules, which is the claim
their names make. They are frozen, so they need no updating.

[v5]: https://github.com/cooklang/spec/commit/9b857a24a
[v6]: https://github.com/cooklang/spec/commit/8b2623588

```bash
go test ./...
```

## External References

- [Cooklang](https://cooklang.org/)
- [aquilax/cooklang-go](https://github.com/aquilax/cooklang-go) - another Cooklang parser in Go, from which I borrowed the test JSON (thank you!)
- [Lexical Scanning in Go, Rob Pike, 2011](https://talks.golang.org/2011/lex.slide#1) - basis for the lexer