# hx-unicode

Provides data and functions to test some properties of Unicode code points.

It's port of golang's unicode package.

### How to use

Use `hxUnicode.U.Is` to check is a character can match a unicode property.

All `RangeTable` from golang/unicode package is available.

```haxe
import hxUnicode.U.Nd;
import hxUnicode.U.Is as UnicodeIs;

function main() {
    // check if is a number \u{Nd}, '1' match, 'A' not match.
    // Note: really want module alaias for haxe.
    trace(UnicodeIs(Nd, '1'.charCodeAt(0)));
    trace(UnicodeIs(Nd, 'A'.charCodeAt(0)));
}
```

### Contribute

Do note edit haxe code. Edit go code.

RangeTables were generated from tools/main.go.

```
go run tools/main.go
```

Then use vscode haxe extension to format generated code.
