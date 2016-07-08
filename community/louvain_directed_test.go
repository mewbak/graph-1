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
	{
		name: "zachary",
		g:    zachary,
		// community structure and modularity calculated by C++ implementation: louvain igraph.
		// Note that louvain igraph returns Q as an unscaled value.
		structures: []structure{
			{
				resolution: 1,
				memberships: []set{
					0: linksTo(0, 1, 2, 3, 7, 11, 12, 13, 17, 19, 21),
					1: linksTo(4, 5, 6, 10, 16),
					2: linksTo(8, 9, 14, 15, 18, 20, 22, 26, 29, 30, 32, 33),
					3: linksTo(23, 24, 25, 27, 28, 31),
				},
				want: 34.3417721519 / 79 /* 5->6 and 6->5 because of co-equal rank */, tol: 1e-4,
			},
		},
		wantLevels: []level{
			{
				q: 0.43470597660631316,
				communities: [][]graph.Node{
					{simple.Node(0), simple.Node(1), simple.Node(2), simple.Node(3), simple.Node(7), simple.Node(11), simple.Node(12), simple.Node(13), simple.Node(17), simple.Node(19), simple.Node(21)},
					{simple.Node(4), simple.Node(5), simple.Node(6), simple.Node(10), simple.Node(16)},
					{simple.Node(8), simple.Node(9), simple.Node(14), simple.Node(15), simple.Node(18), simple.Node(20), simple.Node(22), simple.Node(26), simple.Node(29), simple.Node(30), simple.Node(32), simple.Node(33)},
					{simple.Node(23), simple.Node(24), simple.Node(25), simple.Node(27), simple.Node(28), simple.Node(31)},
				},
			},
			{
				q: 0.35410991828232663,
				communities: [][]graph.Node{
					{simple.Node(0), simple.Node(1), simple.Node(2), simple.Node(3), simple.Node(7), simple.Node(11), simple.Node(12), simple.Node(13), simple.Node(17), simple.Node(19), simple.Node(21)},
					{simple.Node(4), simple.Node(10)},
					{simple.Node(5), simple.Node(6), simple.Node(16)},
					{simple.Node(8), simple.Node(30)},
					{simple.Node(9), simple.Node(14), simple.Node(15), simple.Node(18), simple.Node(20), simple.Node(22), simple.Node(32), simple.Node(33)},
					{simple.Node(23), simple.Node(25)},
					{simple.Node(24), simple.Node(27)},
					{simple.Node(26), simple.Node(29)},
					{simple.Node(28), simple.Node(31)},
				},
			},
			{
				q: -0.014580996635154624,
				communities: [][]graph.Node{
					{simple.Node(0)},
					{simple.Node(1)},
					{simple.Node(2)},
					{simple.Node(3)},
					{simple.Node(4)},
					{simple.Node(5)},
					{simple.Node(6)},
					{simple.Node(7)},
					{simple.Node(8)},
					{simple.Node(9)},
					{simple.Node(10)},
					{simple.Node(11)},
					{simple.Node(12)},
					{simple.Node(13)},
					{simple.Node(14)},
					{simple.Node(15)},
					{simple.Node(16)},
					{simple.Node(17)},
					{simple.Node(18)},
					{simple.Node(19)},
					{simple.Node(20)},
					{simple.Node(21)},
					{simple.Node(22)},
					{simple.Node(23)},
					{simple.Node(24)},
					{simple.Node(25)},
					{simple.Node(26)},
					{simple.Node(27)},
					{simple.Node(28)},
					{simple.Node(29)},
					{simple.Node(30)},
					{simple.Node(31)},
					{simple.Node(32)},
					{simple.Node(33)},
				},
			},
		},
	},
	{
		name: "blondel",
		g:    blondel,
		// community structure and modularity calculated by C++ implementation: louvain igraph.
		// Note that louvain igraph returns Q as an unscaled value.
		structures: []structure{
			{
				resolution: 1,
				memberships: []set{
					0: linksTo(0, 1, 2, 3, 4, 5, 6, 7),
					1: linksTo(8, 9, 10, 11, 12, 13, 14, 15),
				},
				want: 11.1428571429 / 28, tol: 1e-4,
			},
		},
		wantLevels: []level{
			{
				q: 0.3979591836734694,
				communities: [][]graph.Node{
					{simple.Node(0), simple.Node(1), simple.Node(2), simple.Node(3), simple.Node(4), simple.Node(5), simple.Node(6), simple.Node(7)},
					{simple.Node(8), simple.Node(9), simple.Node(10), simple.Node(11), simple.Node(12), simple.Node(13), simple.Node(14), simple.Node(15)},
				},
			},
			{
				q: 0.2971938775510204,
				communities: [][]graph.Node{
					{simple.Node(0), simple.Node(3), simple.Node(5), simple.Node(7)},
					{simple.Node(1), simple.Node(2), simple.Node(4), simple.Node(6)},
					{simple.Node(8), simple.Node(15)},
					{simple.Node(9), simple.Node(12), simple.Node(14)},
					{simple.Node(10), simple.Node(11), simple.Node(13)},
				},
			},
			{
				q: -0.022959183673469385,
				communities: [][]graph.Node{
					{simple.Node(0)},
					{simple.Node(1)},
					{simple.Node(2)},
					{simple.Node(3)},
					{simple.Node(4)},
					{simple.Node(5)},
					{simple.Node(6)},
					{simple.Node(7)},
					{simple.Node(8)},
					{simple.Node(9)},
					{simple.Node(10)},
					{simple.Node(11)},
					{simple.Node(12)},
					{simple.Node(13)},
					{simple.Node(14)},
					{simple.Node(15)},
				},
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
