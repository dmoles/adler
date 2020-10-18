# Adler

[![Build Status](https://travis-ci.org/dmoles/adler.svg?branch=main)](https://travis-ci.org/dmoles/adler)

Adler is a minimalist Markdown wiki viewer.

## Building

Out of the box, Adler builds with `go build` / `go install`.

A [`magefile`](https://magefile.org) is provided to support compiling SCSS
and embedding static assets; run `mage -l` for the list of tasks. Note that
the first `mage` invocation may take some time as [golibsass](https://github.com/bep/golibsass) 

## Usage

```
adler <root-dir> [-p <port>]
```

E.g., to serve from Markdown files in `/Users/irene/suda.wiki` on port 8282
(default is 8181):

```sh
adler /Users/irene/suda.wiki -p 8282
```

## Name

Adler is named for [Ada Adler](https://en.wikipedia.org/wiki/Ada_Adler)
(1878-1946), philologist, classical scholar, and translator of the
[_Suda_](https://en.wikipedia.org/wiki/Suda), arguably the greatest
encyclopedia of the Early Middle Ages.
