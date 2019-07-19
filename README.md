# Adler

Adler is a minimalist Markdown wiki viewer.

It follows the [Gollum](https://github.com/gollum/gollum) / [GitHub
wiki](https://help.github.com/en/articles/about-wikis) convention of
linking based on the page title without extension, so a link of the
form

```
[some link text](Page)
```

becomes

```
<a href="Page">some link text</a>
```

and is served from the Markdown file `Page.md`.

<!-- TODO: Just use real links -->

## Usage

```
adler start <root-dir> [-p <port>]
```

E.g., to serve from Markdown files in `/Users/irene/suda.wiki` on port 8181
(default is 8080):

```
adler start /Users/irene/suda.wiki -p 8181
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
