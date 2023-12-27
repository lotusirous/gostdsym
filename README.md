# Fuzzy finding go doc symbols

`stdsym` simplifies Go documentation exploration by extracting all exported
symbols from the standard library, enabling fuzzy searching (fzf) through Go
documents.

## Demo

Watch a quick demonstration showcasing the usage of this tool:

![Demo](./demo.gif)

## Installation

Get started quickly with `stdsym`:

```
go install github.com/lotusirous/gostdsym/stdsym@latest
```

Create a handy `gdoc` alias for instant symbol lookups:

```bash
alias gdoc="stdsym |fzf | xargs go doc "
```

If you want to view the results on [pkg.go.dev](https://pkg.go.dev/), use this alias. This example is for macOS, where the open command opens the link in the default browser:

```bash
alias gdocb="stdsym | fzf | awk '{print \"https://pkg.go.dev/\" \$1}' | xargs open"
```
