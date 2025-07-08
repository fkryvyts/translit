package resources

import (
	_ "embed"
)

//go:embed dicts/compiled/hepburn.ja.dict
var HepburnJa []byte

//go:embed dicts/compiled/kanwa.ja.dict
var KanwaJa []byte

//go:embed dicts/compiled/normalize.ja.dict
var NormalizeJa []byte
