package dict

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

// Similar to dictionary except that it stores all its nodes in memory
// and allows to dump all its nodes to file
type DictionaryBuilder struct {
	nodes []*Node
	words []Word
}

func NewDictionaryBuilder() *DictionaryBuilder {
	d := &DictionaryBuilder{}
	d.nodes = append(d.nodes, &Node{Fail: -1})
	return d
}

func (d *DictionaryBuilder) AddWord(word Word) {
	node := d.nodes[rootIdx]
	for _, r := range word.Word {
		if _, ok := node.Children[r]; !ok {
			d.nodes = append(d.nodes, &Node{Fail: -1})

			if node.Children == nil {
				node.Children = make(map[rune]NodeIdx)
			}

			node.Children[r] = NodeIdx(len(d.nodes) - 1)
		}
		node = d.nodes[node.Children[r]]
	}

	for _, wi := range node.Output {
		if d.words[wi] == word {
			return
		}
	}

	d.words = append(d.words, word)
	node.Output = append(node.Output, WordIdx(len(d.words)-1))
}

// Build builds the failure links (BFS)
func (d *DictionaryBuilder) Build() {
	queue := []*Node{}
	for _, childIdx := range d.nodes[rootIdx].Children {
		child := d.nodes[childIdx]
		child.Fail = rootIdx
		queue = append(queue, child)
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for r, childIdx := range current.Children {
			child := d.nodes[childIdx]

			failIdx := current.Fail
			for failIdx != -1 {
				if _, ok := d.nodes[failIdx].Children[r]; ok {
					break
				}

				failIdx = d.nodes[failIdx].Fail
			}

			if failIdx != -1 {
				child.Fail = d.nodes[failIdx].Children[r]
			} else {
				child.Fail = rootIdx
			}

			child.Output = append(child.Output, d.nodes[child.Fail].Output...)
			queue = append(queue, child)
		}
	}
}

func (d *DictionaryBuilder) Save(path string) (retErr error) {
	nodesFile, err := d.openFile(path + ".nodes")
	if err != nil {
		return err
	}

	defer func() {
		err := nodesFile.Close()
		if err != nil {
			retErr = fmt.Errorf("failed to close file: %w", err)
		}
	}()

	for _, node := range d.nodes {
		var buff bytes.Buffer
		err := json.NewEncoder(&buff).Encode(node)
		if err != nil {
			return fmt.Errorf("failed to encode node: %w", err)
		}

		_, err = nodesFile.Write(buff.Bytes())
		if err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}

	wordsFile, err := d.openFile(path + ".words")
	if err != nil {
		return err
	}

	defer func() {
		err := wordsFile.Close()
		if err != nil {
			retErr = fmt.Errorf("failed to close file: %w", err)
		}
	}()

	for _, word := range d.words {
		_, err = wordsFile.Write(word.Serialize())
		if err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}

	return nil
}

func (d *DictionaryBuilder) openFile(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}

func (d *DictionaryBuilder) Load(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	lines := bytes.Split(content, []byte("\n"))

	d.nodes = make([]*Node, 0, len(lines))

	for _, content := range lines {
		if len(content) == 0 {
			continue
		}

		node := new(Node)

		err := json.NewDecoder(bytes.NewReader(content)).Decode(node)
		if err != nil {
			return fmt.Errorf("failed to decode node: %w", err)
		}

		d.nodes = append(d.nodes, node)
	}

	return nil
}

func (d *DictionaryBuilder) GetNodes() []*Node {
	return d.nodes
}
