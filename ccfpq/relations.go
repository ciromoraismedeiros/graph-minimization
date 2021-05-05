package ccfpq

import (
	"fmt"
	ds "github.com/ciromdrs/graph-tools/datastructures"
)

/* Interfaces */
type Relation interface {
	Node() ds.Vertex
	Label() ds.Vertex
	Objects() ds.VertexSet
	Show()
	AddObjects(ds.VertexSet)
	IsNested() bool
}

type relationsSet interface {
	set(ds.Vertex, ds.Vertex, Relation)
	get(ds.Vertex, ds.Vertex) Relation
	iterate() <-chan Relation
}

/* Structs */

// A TraceItem is a rule and a list of position sets. They will replace
// Relations in future versions.
type TraceItem struct {
	rule   []ds.Vertex
	posets []ds.VertexSet
}

type mapRelationsSet struct {
	data map[ds.Vertex]map[ds.Vertex]Relation
}

type sliceRelationsSet struct {
	data  []Relation
	ESize int
}

type BaseRelation struct {
	node    ds.Vertex
	label   ds.Vertex
	objects ds.VertexSet
}

type TerminalRelation struct {
	BaseRelation
}

type NonTerminalRelation struct {
	BaseRelation
	rules []*NodeSet
}

type NestedRelation struct {
	BaseRelation
	rule *NodeSet
}

type ExpressionRelation struct {
	BaseRelation
	rule *NodeSet
}

type Symbol struct {
	subjNodeSet *NodeSet
	predicate   ds.Vertex
	objNodeSet  *NodeSet
}

/* Symbol Methods and Functions */
func NewSymbol(predicate ds.Vertex) *Symbol {
	return &Symbol{predicate: predicate}
}

/* BaseRelation Methods and Functions */
func (r *BaseRelation) Node() ds.Vertex {
	return r.node
}

func (r *BaseRelation) Label() ds.Vertex {
	return r.label
}

func (r *BaseRelation) Objects() ds.VertexSet {
	return r.objects
}

func (r *BaseRelation) AddObjects(new ds.VertexSet) {
	r.objects.Update(new)
}

func (r *BaseRelation) Show() {
	fmt.Print("R(", r.Node(), ",", r.Label(), ") = {")
	r.Objects().Show()
	fmt.Println("}")
}

func (r *BaseRelation) IsNested() bool {
	return false
}

func (r *BaseRelation) IsNonTerminal() bool {
	return false
}

func (r *BaseRelation) IsTerminal() bool {
	return false
}

func (r *BaseRelation) IsExpression() bool {
	return false
}

func (r *BaseRelation) IsEmpty() bool {
	return false
}

/* NonTerminalRelation Methods and Functions */
func NewNonTerminalRelation(node, label ds.Vertex) *NonTerminalRelation {
	return &NonTerminalRelation{
		BaseRelation: BaseRelation{
			node:    node,
			label:   label,
			objects: f.NewVertexSet(),
		},
		//rules   : nil,
	}
}

func (r *NonTerminalRelation) Show() {
	r.BaseRelation.Show()
	fmt.Println("len(rules) =", len(r.rules))
	for _, rule := range r.rules {
		fmt.Print(r.label, " -> {")
		rule.nodes.Show()
		rule.new.Show() // ShowNew()
		labelData := rule.next
		for labelData != nil {
			fmt.Print("}   ", labelData.predicate, "   {")
			labelData.objNodeSet.nodes.Show()
			labelData.objNodeSet.new.Show() // ShowNew()
			labelData = labelData.objNodeSet.next
		}
		fmt.Println("}")
	}
}

// TODO: remove G from parameters
func (r *NonTerminalRelation) AddRule(startVertices ds.VertexSet, labels []ds.Vertex, G *Grammar) {
	nodeSet := NewNodeSet()
	nodeSet.new.Update(startVertices)
	nodeSet.relation = r
	AddNew(nodeSet)
	r.rules = append(r.rules, nodeSet)

	for _, label := range labels {
		labelData := NewSymbol(label)
		labelData.subjNodeSet = nodeSet
		nodeSet.next = labelData

		objects := NewNodeSet()
		objects.prev = labelData
		objects.relation = r
		labelData.objNodeSet = objects
		nodeSet = objects
	}
}

func (r *NonTerminalRelation) IsNonTerminal() bool {
	return true
}

// TraceItems returns TraceItem objects that represent the relation.
func (r *NonTerminalRelation) TraceItems() []*TraceItem {
	items := make([]*TraceItem, len(r.rules))
	for i, rule := range r.rules {
		start := f.NewVertexSet()
		start.Add(r.Node())
		var newrule []ds.Vertex
		newrule = append(newrule, r.Label())
		for symbol := rule.next; symbol != nil; symbol = symbol.objNodeSet.next {
			newrule = append(newrule, symbol.predicate)
		}
		items[i] = NewTraceItem(start, newrule)
	}
	return items
}

/*  TerminalRelation Methods and Functions */
func NewTerminalRelation(node, label ds.Vertex, objects ds.VertexSet) *TerminalRelation {
	r := &TerminalRelation{
		BaseRelation: BaseRelation{
			node:    node,
			label:   label,
			objects: objects,
		},
	}
	return r
}

func (r *TerminalRelation) IsTerminal() bool {
	return true
}

/* NestedRelation Methods and Functions */
func NewNestedRelation(node, label ds.Vertex) *NestedRelation {
	r := &NestedRelation{
		BaseRelation: BaseRelation{
			node:    node,
			label:   label,
			objects: f.NewVertexSet(),
		},
	}
	return r
}

func (r *NestedRelation) Objects() ds.VertexSet {
	// if the sub-relation has objects, it means its nested expression
	// succesfully derived a path, so this should return its node.
	o := f.NewVertexSet()
	if r.objects.Size() > 0 {
		o.Add(r.node)
	}
	return o
}

func (r *NestedRelation) SetRule(labels []ds.Vertex) {
	nodeSet := NewNodeSet()
	nodeSet.new.Add(r.node)
	nodeSet.relation = r
	AddNew(nodeSet)
	r.rule = nodeSet

	for _, label := range labels {
		labelData := NewSymbol(label)
		labelData.subjNodeSet = nodeSet
		nodeSet.next = labelData

		objects := NewNodeSet()
		objects.prev = labelData
		objects.relation = r
		labelData.objNodeSet = objects

		//~ if _, isNonTerminal := G.NonTerm[label]; isNonTerminal {
		//~ O[*NewPair(r.node,r.label)] = append(O[*NewPair(r.node,r.label)], labelData.objNodeSet)
		//~ }

		nodeSet = objects
		//~ nodeSet.relation = r
	}
}

func (r *NestedRelation) Show() {
	r.BaseRelation.Show()
	fmt.Print(r.label, " -> {")
	r.rule.nodes.Show()
	r.rule.new.Show() // ShowNew()
	labelData := r.rule.next
	for labelData != nil {
		fmt.Print("}   ", labelData.predicate, "   {")
		labelData.objNodeSet.nodes.Show()
		labelData.objNodeSet.new.Show() // ShowNew()
		labelData = labelData.objNodeSet.next
	}
	fmt.Println("}")
}

func (r *NestedRelation) IsNested() bool {
	return true
}

/* ExpressionRelation Methods and Functions */
func NewExpressionRelation(node, label ds.Vertex) *ExpressionRelation {
	r := &ExpressionRelation{
		BaseRelation: BaseRelation{
			node:    node,
			label:   label,
			objects: f.NewVertexSet(),
		},
	}
	return r
}

func (r *ExpressionRelation) SetRule(startVertices ds.VertexSet, labels []ds.Vertex) {
	nodeSet := NewNodeSet()
	nodeSet.new.Update(startVertices)
	nodeSet.relation = r
	AddNew(nodeSet)
	r.rule = nodeSet

	for _, label := range labels {
		labelData := NewSymbol(label)
		labelData.subjNodeSet = nodeSet
		nodeSet.next = labelData

		objects := NewNodeSet()
		objects.prev = labelData
		objects.relation = r
		labelData.objNodeSet = objects
		nodeSet = objects
	}
}

func (r *ExpressionRelation) IsExpression() bool {
	return true
}

/* mapRelationsSet Functions and Methods */
func newMapRelationsSet(VSize, ESize int) *mapRelationsSet {
	r := &mapRelationsSet{data: make(map[ds.Vertex]map[ds.Vertex]Relation, VSize)}
	for k := range r.data {
		r.data[k] = make(map[ds.Vertex]Relation, ESize)
	}
	return r
}

func (m *mapRelationsSet) get(node, symbol ds.Vertex) Relation {
	if _, ok := m.data[node]; !ok {
		return nil
	}
	return m.data[node][symbol]
}

func (m *mapRelationsSet) set(node, symbol ds.Vertex, r Relation) {
	if _, ok := m.data[node]; !ok {
		m.data[node] = make(map[ds.Vertex]Relation)
	}
	m.data[node][symbol] = r
}

func (m *mapRelationsSet) iterate() <-chan Relation {
	ch := make(chan Relation)
	go func() {
		for _, aux := range m.data {
			for _, r := range aux {
				ch <- r
			}
		}
		defer close(ch)
	}()
	return ch
}

/* sliceRelationsSet Functions and Methods */
func newSliceRelationsSet(VSize, ESize int) *sliceRelationsSet {
	return &sliceRelationsSet{
		data:  make([]Relation, VSize*ESize),
		ESize: ESize,
	}
}

func (m *sliceRelationsSet) get(node, symbol ds.Vertex) Relation {
	v := node.(ds.BitVertex)
	s := symbol.(ds.BitVertex)
	i := (m.ESize * v.IndexInSlice()) + s.IndexInSlice()
	return m.data[i]
}

func (m *sliceRelationsSet) set(node, symbol ds.Vertex, r Relation) {
	v := node.(ds.BitVertex)
	s := symbol.(ds.BitVertex)
	i := (m.ESize * v.IndexInSlice()) + s.IndexInSlice()
	m.data[i] = r
}

func (m *sliceRelationsSet) iterate() <-chan Relation {
	ch := make(chan Relation)
	go func() {
		for _, r := range m.data {
			if r != nil {
				ch <- r
			}
		}
		defer close(ch)
	}()
	return ch
}

/* TraceItem methods and functions */

// NewTraceItem returns a new TraceItem object
func NewTraceItem(start ds.VertexSet, rule []ds.Vertex) *TraceItem {
	ti := &TraceItem{
		rule:   rule,
		posets: make([]ds.VertexSet, len(rule)),
	}
	ti.posets[0] = start
	for i := 1; i < len(ti.posets); i++ {
		ti.posets[i] = f.NewVertexSet()
	}
	return ti
}

// Equals returns true if the trace items are equal. It checks the rule and the
// elements in the position sets.
func (ti *TraceItem) Equals(other *TraceItem) bool {
	if len(ti.rule) != len(other.rule) {
		return false
	}
	for i := range ti.rule {
		if !ti.rule[i].Equals(other.rule[i]) ||
			!ti.posets[i].Equals(other.posets[i]) {
			return false
		}
	}
	return true
}
