package macro

import (
	"errors"
	"fmt"
	"github.com/motoki317/lisp-interpreter/node"
)

/**
Example:
(define-syntax my-cond
  (syntax-rules (else) # list of reserved keywords
    ((_ (else e1 ...)) # branch ... if code matches this
     (begin e1 ...)) # replace to this
    ((_ (e1 e2 ...))
     (when e1 e2 ...))
    ((_ (e1 e2 ...) c1 ...)
     (if e1
	 (begin e2 ...)
	 (cond c1 ...)))))
*/

type Macro struct {
	name     string
	branches []*branch
}

type matcherType int

const (
	keyword matcherType = iota
	identifier
	data
	variadic
	nested
)

type branch struct {
	matcher *matcher
	target  *matcher
}

type builder struct {
	id       map[string]*node.Node
	variadic []*node.Node
}

type matcher struct {
	matcherType matcherType
	str         string
	data        *node.Node
	children    []*matcher
}

// NewMacro creates a new macro instance from the given code.
// Returns error if the code is malformed.
func NewMacro(n *node.Node) (*Macro, error) {
	if n.Type != node.Branch || len(n.Children) != 3 {
		return nil, errors.New("expected macro to be a list of length 3")
	}
	if n.Children[0].Type != node.Keyword || n.Children[0].Str != "define-syntax" {
		return nil, fmt.Errorf("expected define-syntax, but got %v", n.Children[0])
	}

	if n.Children[1].Type != node.Identifier {
		return nil, fmt.Errorf("expected macro identifier, but got %v", n.Children[1].Type)
	}
	macroName := n.Children[1].Str

	syntaxRules := n.Children[2].Children
	if len(syntaxRules) == 0 || syntaxRules[0].Str != "syntax-rules" {
		return nil, fmt.Errorf("expected syntax rules, but got %v", syntaxRules)
	}
	if len(syntaxRules) <= 2 {
		return nil, fmt.Errorf("expected length of syntax rule to be >= 3, but got %v", len(syntaxRules))
	}
	if syntaxRules[1].Type != node.Branch {
		return nil, fmt.Errorf("expected 2nd element of syntax-rules to be a list of keywords, but got %v", syntaxRules[1])
	}

	allowedKeywords := make([]string, len(syntaxRules[1].Children))
	for i, allowedKeyword := range syntaxRules[1].Children {
		if allowedKeyword.Type != node.Keyword {
			return nil, fmt.Errorf("expected keywords in syntax-rules allowed keywords, but got %v", allowedKeyword.Type)
		}
		allowedKeywords[i] = allowedKeyword.Str
	}

	// allowed keywords are checked on branch validation time
	branches := make([]*branch, 0, len(syntaxRules)-2)
	for _, branchCode := range syntaxRules[2:] {
		b, err := newBranch(branchCode, allowedKeywords)
		if err != nil {
			return nil, fmt.Errorf("malformed branch: %w", err)
		}
		branches = append(branches, b)
	}

	return &Macro{
		name:     macroName,
		branches: branches,
	}, nil
}

// Replace checks the given node recursively, and applies the macro (once) if possible.
func (m *Macro) Replace(n *node.Node) (res *node.Node, ok bool) {
	if n.Type != node.Branch {
		return n, false
	}
	// check if the whole node is applicable
	if res, ok = m.replaceOne(n); ok {
		return
	}
	// check if each children is applicable
	for i, child := range n.Children {
		if res, ok = m.Replace(child); ok {
			n.Children[i] = res
			return n, true
		}
	}
	return n, false
}

// replaceOne checks the given node but NOT checking recursively, and applies the macro (once) if possible.
func (m *Macro) replaceOne(n *node.Node) (res *node.Node, ok bool) {
	if n.Type != node.Branch || len(n.Children) == 0 ||
		n.Children[0].Type != node.Identifier || n.Children[0].Str != m.name {
		return n, false
	}
	// drop the first elt in the list (which corresponds to macro name) before checking
	n = &node.Node{
		Type:     node.Branch,
		Children: n.Children[1:],
	}
	for _, branch := range m.branches {
		if res, ok = branch.replace(n); ok {
			return
		}
	}
	return n, false
}

// newBranch creates a new branch from the given node, and allowed keywords.
// Returns error if the code is malformed.
func newBranch(n *node.Node, allowedKeywords []string) (*branch, error) {
	if len(n.Children) != 2 {
		return nil, errors.New("expected branch to be a list of length 2")
	}

	matcherCode := n.Children[0]
	targetCode := n.Children[1]
	if matcherCode.Type != node.Branch || matcherCode.Children[0].Type != node.Keyword ||
		matcherCode.Children[0].Str != "_" {
		return nil, errors.New("expected \"_\" in the first element of the branch matcher")
	}
	// drop the first elt in the list which corresponds to the macro name
	matcherCode.Children = matcherCode.Children[1:]

	matcher, err := newMatcher(matcherCode, allowedKeywords)
	if err != nil {
		return nil, fmt.Errorf("malformed matcher: %w", err)
	}
	target, err := newMatcherTarget(targetCode)
	if err != nil {
		return nil, fmt.Errorf("malformed target: %v", err)
	}

	return &branch{
		matcher: matcher,
		target:  target,
	}, nil
}

// replace checks if the macro can be applied, and returns transformed node if yes.
func (b *branch) replace(n *node.Node) (res *node.Node, ok bool) {
	if !b.matcher.match(n) {
		return nil, false
	}
	builder := builder{
		id:       make(map[string]*node.Node),
		variadic: nil,
	}
	b.matcher.retrieve(n, &builder)
	return builder.buildTarget(b.target), true
}

// buildTarget builds macro target from the retrieved id-to-node map.
func (b *builder) buildTarget(target *matcher) *node.Node {
	switch target.matcherType {
	case keyword:
		return &node.Node{
			Type: node.Keyword,
			Str:  target.str,
		}
	case identifier:
		if res, ok := b.id[target.str]; ok {
			return res
		} else {
			return &node.Node{
				Type: node.Identifier,
				Str:  target.str,
			}
		}
	case data:
		return target.data
	case variadic:
		return &node.Node{
			Type:     node.Branch,
			Children: b.variadic,
		}
	case nested:
		res := &node.Node{
			Type:     node.Branch,
			Children: make([]*node.Node, 0, len(target.children)),
		}
		tLength := len(target.children)
		if tLength == 0 {
			return res
		}
		// match each one
		for _, childMatcher := range target.children {
			if childMatcher.matcherType == variadic {
				res.Children = append(res.Children, b.buildTarget(childMatcher).Children...)
			} else {
				res.Children = append(res.Children, b.buildTarget(childMatcher))
			}
		}
		return res
	}
	panic(fmt.Sprintf("type %v not implemented", target.matcherType))
}

// newMatcher creates a new matcher from the given code and allowed keywords.
func newMatcher(n *node.Node, allowedKeywords []string) (*matcher, error) {
	switch n.Type {
	case node.Keyword:
		if n.Str == "..." {
			return &matcher{matcherType: variadic}, nil
		} else {
			if !contains(allowedKeywords, n.Str) {
				return nil, fmt.Errorf("unexpected keyword: %v", n.Str)
			}
			return &matcher{matcherType: keyword, str: n.Str}, nil
		}
	case node.Identifier:
		return &matcher{matcherType: identifier, str: n.Str}, nil
	case node.Branch:
		children := make([]*matcher, len(n.Children))
		for i, childCode := range n.Children {
			child, err := newMatcher(childCode, allowedKeywords)
			if err != nil {
				return nil, err
			}
			// only allow variadic capture at the end of a list
			if child.matcherType == variadic && i != len(n.Children)-1 {
				return nil, errors.New("variadic capture only allowed at the end of a list")
			}
			children[i] = child
		}
		return &matcher{matcherType: nested, children: children}, nil
	}
	return nil, fmt.Errorf("unexpected type: %v", n.Type)
}

// newMatcherTarget creates a new matcher target from the given code.
func newMatcherTarget(n *node.Node) (*matcher, error) {
	switch n.Type {
	case node.Keyword:
		if n.Str == "..." {
			return &matcher{matcherType: variadic}, nil
		} else {
			return &matcher{matcherType: keyword, str: n.Str}, nil
		}
	case node.Identifier:
		return &matcher{matcherType: identifier, str: n.Str}, nil
	case node.Number:
		fallthrough
	case node.Boolean:
		fallthrough
	case node.String:
		return &matcher{matcherType: data, data: n}, nil
	case node.Branch:
		children := make([]*matcher, len(n.Children))
		for i, childCode := range n.Children {
			child, err := newMatcherTarget(childCode)
			if err != nil {
				return nil, err
			}
			// allow variadic target at any position in a list
			children[i] = child
		}
		return &matcher{matcherType: nested, children: children}, nil
	}
	return nil, fmt.Errorf("unexpected type: %v", n.Type)
}

func contains(lst []string, target string) bool {
	for _, elt := range lst {
		if target == elt {
			return true
		}
	}
	return false
}

// match checks if the node matches this macro.
func (m *matcher) match(n *node.Node) bool {
	switch m.matcherType {
	case keyword:
		return n.Type == node.Keyword && n.Str == m.str
	case identifier:
		return true
	case variadic:
		return n.Type == node.Branch
	case nested:
		nLength, mLength := len(n.Children), len(m.children)
		isVariadic := mLength > 0 && m.children[mLength-1].matcherType == variadic
		if n.Type != node.Branch {
			return false
		}
		if nLength == 0 && mLength == 0 {
			return true
		}
		if isVariadic {
			if nLength < mLength-1 {
				return false
			}
			for i, child := range m.children[:mLength-1] {
				if !child.match(n.Children[i]) {
					return false
				}
			}
		} else {
			if nLength != mLength {
				return false
			}
			for i, child := range m.children {
				if !child.match(n.Children[i]) {
					return false
				}
			}
		}
		return true
	}
	panic(fmt.Sprintf("type %v not implemented", m.matcherType))
}

// retrieve retrieves the id-to-node map from the given node.
// Asserts the node matches this macro.
func (m *matcher) retrieve(n *node.Node, b *builder) {
	switch m.matcherType {
	case keyword:
		// nop
	case identifier:
		b.id[m.str] = n
	case data:
		// nop
	case variadic:
		// assert n.Type == node.Branch
		b.variadic = n.Children
	case nested:
		nLength, mLength := len(n.Children), len(m.children)
		if nLength == 0 && mLength == 0 {
			return
		}
		// match each one
		for i, child := range m.children[:mLength-1] {
			child.retrieve(n.Children[i], b)
		}
		// the last one could be variadic, if so, match
		if m.children[mLength-1].matcherType == variadic {
			m.children[mLength-1].retrieve(&node.Node{
				Type:     node.Branch,
				Children: n.Children[mLength-1:],
			}, b)
			return
		}
		// assert nLength == mLength
		m.children[mLength-1].retrieve(n.Children[mLength-1], b)
	}
}
