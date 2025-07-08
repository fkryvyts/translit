# translit
Library for transliterating non-ASCII text into ASCII

Focuses on low memory footprint (which is not easy to achieve because Japanese transliteration dictionaries are huge) and fast execution times

Uses special dictionaries for Japanese to achieve more accurate transliteration for that language.
If no special rules for specific language are found falls back to https://github.com/anyascii/go for generic transliteration

Dictionaries come from
- [go-kakasi](https://github.com/sarumaj/go-kakasi)
- [JMnedict](https://www.edrdg.org/enamdict/enamdict_doc.html)

Currently it is mostly proof of concept so it may crash with OOM errors but you are free to use this and report bugs for it so the library can be improved. 

## Usage:

```go

line = Transliterate(line, Hints{
    Language: Japanese,
})

fmt.Println(line)

```