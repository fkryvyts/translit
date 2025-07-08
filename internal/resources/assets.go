package resources

import (
	_ "embed"
)

//go:embed dicts/compiled/hepburn.ja.dict.nodes
var HepburnJaNodes []byte

//go:embed dicts/compiled/hepburn.ja.dict.words
var HepburnJaWords []byte

//go:embed dicts/compiled/kanwa.ja.dict.nodes
var KanwaJaNodes []byte

//go:embed dicts/compiled/kanwa.ja.dict.words
var KanwaJaWords []byte

//go:embed dicts/compiled/normalize.ja.dict.nodes
var NormalizeJaNodes []byte

//go:embed dicts/compiled/normalize.ja.dict.words
var NormalizeJaWords []byte
