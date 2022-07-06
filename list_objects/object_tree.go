package list_objects

import (
	"s3-gateway/util"
	"strings"
)

type ObjectTree struct {
	Label       string       `json:"label"`
	FilePath    string       `json:"filePath"`
	IsDirectory bool         `json:"IsDirectory"`
	Children    []ObjectTree `json:"children"`
}

func BuildObjectTree(t *util.Trie[string], prefix, delimiter string) ObjectTree {
	fullPath := prefix + t.Value
	if !t.IsLeaf && t.Value != "" && !strings.HasSuffix(fullPath, delimiter) {
		fullPath += delimiter
	}

	ret := ObjectTree{
		Label:       t.Value,
		FilePath:    fullPath,
		IsDirectory: !t.IsLeaf,
	}

	if len(t.Node) != 0 {
		ch := make([]ObjectTree, 0)
		for _, v := range t.Node {
			node := BuildObjectTree(v, fullPath, delimiter)
			ch = append(ch, node)
		}
		ret.Children = ch
	}

	return ret
}

func FilterFile(obj ObjectTree) ObjectTree {
	ret := ObjectTree{
		Label:       obj.Label,
		FilePath:    obj.FilePath,
		IsDirectory: obj.IsDirectory,
	}

	ch := make([]ObjectTree, 0)
	for _, u := range obj.Children {
		if u.IsDirectory == true {
			ch = append(ch, FilterFile(u))
		}
	}

	if len(ch) != 0 {
		ret.Children = ch
	}
	return ret
}
