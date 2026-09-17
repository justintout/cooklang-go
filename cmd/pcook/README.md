# 🦚👨‍🍳 `pcook`
> `cooklang-go` example parsing binary

Pass a `.cook` file, get JSON out.

```
$ go build
$ ./pcook example.cook | jq
```

`example.cook` is written for the latest version of the spec, and
`example-condensed.cook` for version 5, whose `>> key: value` metadata and
one-step-per-line are what `-spec 5` reads. `-spec auto` reads whichever the
file looks like; a recipe with no metadata to go on is read as the latest.

```
$ ./pcook -spec 5 example-condensed.cook | jq
$ ./pcook -spec auto example-condensed.cook | jq
```
