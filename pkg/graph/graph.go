package graph

// Graph представляет граф с узлами и ребрами.
type Graph[K comparable] struct {
	nodes map[K]*Node[K]
}

// NewGraph создает и возвращает новый пустой граф.
func NewGraph[K comparable]() *Graph[K] {
	return &Graph[K]{
		nodes: make(map[K]*Node[K]),
	}
}

// AddNode добавляет узел в граф.
func (g *Graph[K]) AddNode(key K) {
	if !g.ContainsNode(key) {
		g.nodes[key] = NewNode[K]()
	}
}

// getOrAddNode возвращает узел, если он существует, или добавляет и возвращает новый узел.
func (g *Graph[K]) getOrAddNode(node K) *Node[K] {
	n, ok := g.nodes[node]
	if !ok {
		n = NewNode[K]()
		g.nodes[node] = n
	}
	return n
}

// AddEdge добавляет ребро между двумя узлами.
func (g *Graph[K]) AddEdge(from K, to K) {
	f := g.getOrAddNode(from)
	g.AddNode(to)
	f.AddEdge(to)
}

// ContainsNode проверяет, существует ли узел в графе.
func (g *Graph[K]) ContainsNode(key K) bool {
	_, ok := g.nodes[key]
	return ok
}

// TopSort выполняет топологическую сортировку графа и возвращает отсортированный список узлов.
func (g *Graph[K]) TopSort() ([]K, error) {
	sortedNodes := newOrderedSet[K]()
	visitedNodes := newOrderedSet[K]()

	for key := range g.nodes {
		if !visitedNodes.Contains(key) {
			err := g.visit(key, sortedNodes, visitedNodes)
			if err != nil {
				return nil, err
			}
		}
	}
	return sortedNodes.items, nil
}

// visit посещает узел и выполняет рекурсивную проверку на циклы.
func (g *Graph[K]) visit(key K, sortedNodes *orderedSet[K], visitedNodes *orderedSet[K]) error {
	added := visitedNodes.Add(key)
	if !added {
		return nil
	}
	n := g.nodes[key]
	for _, edge := range n.Edges() {
		if !sortedNodes.Contains(edge) {
			err := g.visit(edge, sortedNodes, visitedNodes)
			if err != nil {
				return err
			}
		}
	}

	sortedNodes.Add(key)
	visitedNodes.Remove(key)
	return nil
}

// Node представляет узел графа, содержащий ребра к другим узлам.
type Node[K comparable] struct {
	edges map[K]bool
}

// NewNode создает и возвращает новый узел.
func NewNode[K comparable]() *Node[K] {
	return &Node[K]{
		edges: make(map[K]bool),
	}
}

// AddEdge добавляет ребро к узлу.
func (n *Node[K]) AddEdge(key K) {
	n.edges[key] = true
}

// Edges возвращает список узлов, к которым ведут ребра.
func (n *Node[K]) Edges() []K {
	var keys []K
	for k := range n.edges {
		keys = append(keys, k)
	}
	return keys
}

// orderedSet представляет упорядоченное множество узлов.
type orderedSet[K comparable] struct {
	indexes map[K]int
	items   []K
	length  int
}

// newOrderedSet создает и возвращает новое упорядоченное множество.
func newOrderedSet[K comparable]() *orderedSet[K] {
	return &orderedSet[K]{
		indexes: make(map[K]int),
		length:  0,
	}
}

// Add добавляет элемент в упорядоченное множество.
func (s *orderedSet[K]) Add(item K) bool {
	_, ok := s.indexes[item]
	if !ok {
		s.indexes[item] = s.length
		s.items = append(s.items, item)
		s.length++
	}
	return !ok
}

// Contains проверяет, содержится ли элемент в упорядоченном множестве.
func (s *orderedSet[K]) Contains(item K) bool {
	_, ok := s.indexes[item]
	return ok
}

// Remove удаляет элемент из упорядоченного множества.
func (s *orderedSet[K]) Remove(item K) {
	if index, ok := s.indexes[item]; ok {
		delete(s.indexes, item)
		s.items = append(s.items[:index], s.items[index+1:]...)
		s.length--
		// Обновляем индексы
		for i := index; i < s.length; i++ {
			s.indexes[s.items[i]] = i
		}
	}
}

// Index возвращает индекс элемента в упорядоченном множестве.
func (s *orderedSet[K]) Index(item K) int {
	if index, ok := s.indexes[item]; ok {
		return index
	}
	return -1
}
