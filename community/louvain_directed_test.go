// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package community

import (
	"math"
	"math/rand"
	"reflect"
	"sort"
	"testing"

	"github.com/gonum/graph/internal/ordered"

	"github.com/gonum/graph"
	"github.com/gonum/graph/simple"
)

var communityDirectedQTests = []struct {
	name       string
	g          []set
	structures []structure

	wantLevels []level
}{
	{
		name: "simple_directed",
		g: []set{
			0: linksTo(1),
			1: linksTo(0, 4),
			2: linksTo(1),
			3: linksTo(0, 4),
			4: linksTo(2),
		},
		// community structure and modularity calculated by C++ implementation: louvain igraph.
		// Note that louvain igraph returns Q as an unscaled value.
		structures: []structure{
			{
				resolution: 1,
				memberships: []set{
					0: linksTo(0, 1),
					1: linksTo(2, 3, 4),
				},
				want: 0.5714285714285716 / 7,
			},
		},
		wantLevels: []level{
			{
				communities: [][]graph.Node{
					{simple.Node(0), simple.Node(1)},
					{simple.Node(2), simple.Node(3), simple.Node(4)},
				},
				q: 0.5714285714285716 / 7,
			},
			{
				communities: [][]graph.Node{
					{simple.Node(0)},
					{simple.Node(1)},
					{simple.Node(2)},
					{simple.Node(3)},
					{simple.Node(4)},
				},
				q: -1.2857142857142856 / 7,
			},
		},
	},
}

func TestLouvainDirected(t *testing.T) {
	const louvainIterations = 20

	for _, test := range communityDirectedQTests {
		g := simple.NewDirectedGraph(0, 0)
		for u, e := range test.g {
			// Add nodes that are not defined by an edge.
			if !g.Has(simple.Node(u)) {
				g.AddNode(simple.Node(u))
			}
			for v := range e {
				g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v), W: 1})
			}
		}

		if test.structures[0].resolution != 1 {
			panic("bad test: expect resolution=1")
		}
		want := make([][]graph.Node, len(test.structures[0].memberships))
		for i, c := range test.structures[0].memberships {
			for n := range c {
				want[i] = append(want[i], simple.Node(n))
			}
			sort.Sort(ordered.ByID(want[i]))
		}
		sort.Sort(ordered.BySliceIDs(want))

		var (
			got   *ReducedDirected
			bestQ = math.Inf(-1)
		)
		// Louvain is randomised so we do this to
		// ensure the level tests are consistent.
		src := rand.New(rand.NewSource(1))
		for i := 0; i < louvainIterations; i++ {
			r := LouvainDirected(g, 1, src)
			if q := Q(r, nil, 1); q > bestQ || math.IsNaN(q) {
				bestQ = q
				got = r

				if math.IsNaN(q) {
					// Don't try again for non-connected case.
					break
				}
			}

			var qs []float64
			for p := r; p != nil; p = p.Expanded() {
				qs = append(qs, Q(p, nil, 1))
			}

			// Recovery of Q values is reversed.
			if reverse(qs); !sort.Float64sAreSorted(qs) {
				t.Errorf("Q values not monotonically increasing: %.5v", qs)
			}
		}

		gotCommunities := got.Communities()
		for _, c := range gotCommunities {
			sort.Sort(ordered.ByID(c))
		}
		sort.Sort(ordered.BySliceIDs(gotCommunities))
		if !reflect.DeepEqual(gotCommunities, want) {
			t.Errorf("unexpected community membership for %s Q=%.4v:\n\tgot: %v\n\twant:%v",
				test.name, bestQ, gotCommunities, want)
			continue
		}

		var levels []level
		for p := got; p != nil; p = p.Expanded() {
			var communities [][]graph.Node
			if p.parent != nil {
				communities = p.parent.Communities()
				for _, c := range communities {
					sort.Sort(ordered.ByID(c))
				}
				sort.Sort(ordered.BySliceIDs(communities))
			} else {
				communities = reduceDirected(g, nil).Communities()
			}
			q := Q(p, nil, 1)
			if math.IsNaN(q) {
				// Use an equalable flag value in place of NaN.
				q = math.Inf(-1)
			}
			levels = append(levels, level{q: q, communities: communities})
		}
		if !reflect.DeepEqual(levels, test.wantLevels) {
			t.Errorf("unexpected level structure:\n\tgot: %v\n\twant:%v", levels, test.wantLevels)
		}
	}
}
