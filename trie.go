package chomolungma

import "strings"

type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang
	part     string  // 路由中的一部分，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为true
}

func (n *node) matchPart(part string) *node {
	for _, s := range n.children {
		if s.isWild || s.part == part {
			return s
		}
	}

	return nil
}

func (n *node) matchNodes(part string) []*node {
	// 同时匹配到同一层的多个兄弟节点
	nodes := make([]*node, 0, 4)
	for _, node := range n.children {
		if node.part == part || node.isWild {
			nodes = append(nodes, node)
		}
	}

	return nodes
}

func (n *node) insert(pattern string, parts []string, depth int) {
	if len(parts) == depth {
		n.pattern = pattern
		return
	}

	part := parts[depth]
	child := n.matchPart(part)
	if child == nil {
		// part 不可能为空, parts为空的情况已在第一次insert时处理, len(parts) == depth
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}

	child.insert(pattern, parts, depth+1)
}

func (n *node) search(parts []string, depth int) *node {
	// n.part 可能为 nil, 当配置了 / 路由时存在 n.part 为空的情况
	if len(parts) == depth || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}

		return n
	}

	part := parts[depth]
	nodes := n.matchNodes(part)

	for _, node := range nodes {
		result := node.search(parts, depth+1)
		// 首次匹配到就返回结果
		if result != nil {
			return result
		}
	}

	return nil
}
