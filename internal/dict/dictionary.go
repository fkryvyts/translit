package dict

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"
)

type NodeIdx int

const rootIdx NodeIdx = 0

type Node struct {
	Children map[rune]NodeIdx `json:"c,omitempty"`
	Fail     NodeIdx          `json:"f,omitempty"`
	Output   []Word           `json:"o,omitempty"`
}

type Word struct {
	Word     string
	Translit string
}

func (w Word) MarshalJSON() ([]byte, error) {
	return json.Marshal([2]string{w.Word, w.Translit})
}

func (w *Word) UnmarshalJSON(data []byte) error {
	var arr [2]string
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}

	w.Word = arr[0]
	w.Translit = arr[1]
	return nil
}

type Match struct {
	Word     Word
	Position int
}

type MatchesResult struct {
	Text    string
	Matches []Match
}

func (m MatchesResult) Replace(wordSep string) string {
	res := []rune(m.Text)
	wordSepRunes := []rune(wordSep)
	buff := make([]rune, 0, len(res))

	shift := 0
	nextPos := 0

	for _, match := range m.Matches {
		if match.Position < nextPos {
			continue
		}

		w := []rune(match.Word.Word)

		pos := match.Position + shift
		if string(res[pos:pos+len(w)]) != match.Word.Word {
			continue
		}

		tr := []rune(match.Word.Translit)

		buff = append(buff[:0], res[:pos]...)

		if s := string(res[:pos]); len(s) > 0 && !strings.HasSuffix(s, wordSep) {
			buff = append(buff, wordSepRunes...)
			shift += len(wordSepRunes)
		}

		buff = append(buff, tr...)

		if s := string(res[pos+len(w):]); len(s) > 0 && !strings.HasPrefix(s, wordSep) {
			buff = append(buff, wordSepRunes...)
			shift += len(wordSepRunes)
		}

		buff = append(buff, res[pos+len(w):]...)

		shift += len(tr) - len(w)
		nextPos = match.Position + len(w)
		res, buff = buff, res
	}

	return string(res)
}

type Dictionary struct {
	sync.RWMutex
	nodes          map[NodeIdx]*Node
	nodesQueue     []NodeIdx
	nodesCacheSize int
	lines          [][]byte
}

func NewDictionary(nodesCacheSize int) *Dictionary {
	d := &Dictionary{
		nodes:          make(map[NodeIdx]*Node),
		nodesCacheSize: nodesCacheSize,
	}

	if nodesCacheSize > 0 {
		d.nodesQueue = make([]NodeIdx, 0, nodesCacheSize+1)
	}

	return d
}

// Search that uses Aho-Corasick automaton
func (d *Dictionary) Search(text string) *MatchesResult {
	nodeIdx := rootIdx
	var results []Match
	runes := []rune(text)

	for i := range len(runes) {
		r := runes[i]
		for nodeIdx != rootIdx {
			if _, ok := d.loadNode(nodeIdx).Children[r]; ok {
				break
			}

			nodeIdx = d.loadNode(nodeIdx).Fail
		}

		if child, ok := d.loadNode(nodeIdx).Children[r]; ok {
			nodeIdx = child
		}

		node := d.loadNode(nodeIdx)
		for _, word := range node.Output {
			start := i - len([]rune(word.Word)) + 1
			results = append(results, Match{Word: word, Position: start})
		}
	}

	slices.SortStableFunc(results, func(a, b Match) int {
		if a.Position != b.Position {
			return a.Position - b.Position
		}

		if l1, l2 := len([]rune(a.Word.Word)), len([]rune(b.Word.Word)); l1 != l2 {
			return l2 - l1
		}

		return 0
	})

	return &MatchesResult{
		Text:    text,
		Matches: results,
	}
}

func (d *Dictionary) unloadOldestNodes() {
	if len(d.nodesQueue) <= d.nodesCacheSize || d.nodesCacheSize <= 0 {
		return
	}

	oldestNodeIdx := d.nodesQueue[0]
	delete(d.nodes, oldestNodeIdx)

	// Resize queue without memory reallocation
	for i := range d.nodesCacheSize {
		d.nodesQueue[i] = d.nodesQueue[i+1]
	}

	d.nodesQueue = d.nodesQueue[:d.nodesCacheSize]

}

func (d *Dictionary) loadNode(nodeIdx NodeIdx) *Node {
	if n, ok := d.getNode(nodeIdx); ok {
		return n
	}

	d.Lock()
	defer d.Unlock()

	if d.nodesCacheSize > 0 {
		d.nodesQueue = append(d.nodesQueue, nodeIdx)
		d.unloadOldestNodes()
	}

	// Should never happen when we actually return this node but still
	defNode := &Node{Fail: -1}

	if len(d.lines) <= int(nodeIdx) {
		return defNode
	}

	content := d.lines[nodeIdx]

	node := new(Node)

	err := json.NewDecoder(bytes.NewReader(content)).Decode(node)
	if err != nil {
		return defNode
	}

	d.nodes[nodeIdx] = node

	return node
}

func (d *Dictionary) getNode(nodeIdx NodeIdx) (*Node, bool) {
	d.RLock()
	defer d.RUnlock()

	n, ok := d.nodes[nodeIdx]
	return n, ok
}

func (d *Dictionary) Load(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	return d.LoadFromBytes(content)
}

func (d *Dictionary) LoadFromBytes(content []byte) error {
	d.lines = bytes.Split(content, []byte("\n"))

	return nil
}
