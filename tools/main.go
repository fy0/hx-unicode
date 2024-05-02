package main

import (
	"fmt"
	"os"
	"unicode"
)

var head = `// Auto generated code. DO NOT EDIT.
package hxUnicode;

typedef RangeChar = Array<UInt>; // Lo, Hi, Stride

typedef RangeTable = {
	var R16:Array<RangeChar>;
	var R32:Array<RangeChar>;
	var LatinOffset:Int; // number of entries in R16 with Hi <= MaxLatin1
}

// linearMax is the maximum size table for linear search for non-Latin1 rune.
var linearMax = 18;
var MaxRune = 0x0010FFFF; // Maximum valid Unicode code point.
var ReplacementChar = 0xFFFD; // Represents invalid code points.
var MaxASCII = 0x007F; // maximum ASCII value.
var MaxLatin1 = 0x00FF; // maximum Latin-1 value.

// is16 reports whether r is in the sorted slice of 16-bit ranges.
private function is16(ranges:Array<RangeChar>, r:UInt):Bool {
	if (ranges.length <= linearMax || r <= MaxLatin1) {
		for (i in ranges) {
			if (r < i[0]) {
				return false;
			}
			if (r <= i[1]) {
				return i[2] == 1 || (r - i[0]) % i[2] == 0;
			}
		}
		return false;
	}

	// binary search over ranges
	var lo = 0;
	var hi = ranges.length;
	while (lo < hi) {
		var m = lo + Std.int((hi - lo) / 2);
		var range_ = ranges[m];
		if (range_[0] <= r && r <= range_[1]) {
			return range_[2] == 1 || (r - range_[0]) % range_[2] == 0;
		}
		if (r < range_[0]) {
			hi = m;
		} else {
			lo = m + 1;
		}
	}
	return false;
}

// is32 reports whether r is in the sorted slice of 32-bit ranges.
private function is32(ranges:Array<RangeChar>, r:UInt):Bool {
	if (ranges.length <= linearMax) {
		for (i in ranges) {
			if (r < i[0]) {
				return false;
			}
			if (r <= i[1]) {
				return i[2] == 1 || (r - i[0]) % i[2] == 0;
			}
		}
		return false;
	}

	// binary search over ranges
	var lo = 0;
	var hi = ranges.length;
	while (lo < hi) {
		var m = lo + Std.int((hi - lo) / 2);
		var range_ = ranges[m];
		if (range_[0] <= r && r <= range_[1]) {
			return range_[2] == 1 || (r - range_[0]) % range_[2] == 0;
		}
		if (r < range_[0]) {
			hi = m;
		} else {
			lo = m + 1;
		}
	}
	return false;
}

// Is reports whether the rune is in the specified table of ranges.
function Is(rangeTab:RangeTable, r:UInt):Bool {
	var r16 = rangeTab.R16;
	// Compare as uint32 to correctly handle negative runes.
	if (r16.length > 0 && r < r16[r16.length - 1][1]) {
		return is16(r16, r);
	}

	var r32 = rangeTab.R32;
	if (r32.length > 0 && r >= r32[0][0]) {
		return is32(r32, r);
	}

	return false;
}

`

func rangeTable(class string) *unicode.RangeTable {
	if rt, ok := unicode.Categories[class]; ok {
		return rt
	}
	if rt, ok := unicode.Properties[class]; ok {
		return rt
	}
	if rt, ok := unicode.Scripts[class]; ok {
		return rt
	}

	// cannot happen
	panic(fmt.Sprintf("invalid Unicode class: %s", class))
}

/*
	var XX:RangeTable = {
		R16: [
	        [0x00ad, 0x0600, 3]
	    ],
		R32: [
	        [1, 2, 3]
	    ],
		LatinOffset: 2,
	}
*/

func _writeln(f *os.File, format string, args ...any) {
	_, _ = f.WriteString(fmt.Sprintf(format+"\n", args...))
}

func writeTable(f *os.File, name string) {
	x := rangeTable(name)
	writeln := func(format string, args ...any) {
		_writeln(f, format, args...)
	}

	//write := func(format string, args ...any) {
	//	_, _ = f.WriteString(fmt.Sprintf(format, args))
	//}

	writeln(`private var _%s:RangeTable = {`, name)
	writeln("\tR16: [")
	for _, i := range x.R16 {
		writeln("\t\t[0x%x,0x%x,%d],", i.Lo, i.Hi, i.Stride)
	}
	writeln("\t],")
	writeln("\tR32: [")
	for _, i := range x.R32 {
		writeln("\t\t[0x%x,0x%x,%d],", i.Lo, i.Hi, i.Stride)
	}
	writeln("\t],")
	writeln("\tLatinOffset: 2,")
	writeln("}")
}

func main() {
	f, _ := os.OpenFile("./Unicode.hx", os.O_CREATE|os.O_RDWR, 0644)
	_, _ = f.WriteString(head)

	var m map[string]bool
	var names []string

	for k, _ := range unicode.Categories {
		if !m[k] {
			names = append(names, k)
		}
	}
	for k, _ := range unicode.Properties {
		if !m[k] {
			names = append(names, k)
		}
	}
	for k, _ := range unicode.Scripts {
		if !m[k] {
			names = append(names, k)
		}
	}

	for _, i := range names {
		writeTable(f, i)
	}

	_writeln(f, "")
	_writeln(f, "// to hide long text for haxe vscode auto complete")
	for _, i := range names {
		_writeln(f, `var %s = _%s;`, i, i)
	}

	_ = f.Close()
}
