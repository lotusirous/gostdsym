# Fuzzy finding go doc symbols

`stdsym` simplifies Go documentation exploration by extracting all exported
symbols from the standard library, enabling fuzzy searching through Go
documents.

## Demo

Watch a quick demo showcasing the usage of this tool:

![Demo](./demo.gif)

## Installation

Get started quickly with `stdsym`:

```
go install github.com/lotusirous/gostdsym/stdsym@latest
```

Create a handy `gdoc` alias for instant symbol lookups:

```bash
$ stdsym > ~/.gostdsym
$ alias gdoc="cat ~/.gostdsym |fzf | xargs go doc "
```
