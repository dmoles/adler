# Adler

Adler is a minimalist Markdown wiki viewer.

## Usage

```
adler <root-dir> [-p <port>]
```

E.g., to serve from Markdown files in `/Users/irene/suda.wiki` on port 8282
(default is 8181):

```
adler /Users/irene/suda.wiki -p 8282
```

## Known issues

```
$ grep TODO adler/server.go | awk '{$1=$1};1'
// TODO: handle subdirectories properly
// TODO: generate directory indexes
// TODO: handle URL paths that already end in `.md`
// TODO: use a template language
// TODO: some CSS
```

## Name

Adler is named for [Ada Adler](https://en.wikipedia.org/wiki/Ada_Adler)
(1878-1946), philologist, classical scholar, and translator of the
[_Suda_](https://en.wikipedia.org/wiki/Suda), arguably the greatest
encyclopedia of the Early Middle Ages.
