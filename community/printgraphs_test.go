// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build printgraphs

// printgraphs allows us to generate a consistent directed view of
// a set of edges that follows a reasonably real-world-meaningful
// graph. The interpretation of the links in the resulting directed
// graphs are either "suggests" in the context of a Page Ranking or
// possibly "looks up to" in the Zachary graph.
//
// The []set values are printed when go test -tags printgraphs is run.

package community

import (
	"fmt"
	"sort"

	"github.com/gonum/graph"
	"github.com/gonum/graph/internal/ordered"
	"github.com/gonum/graph/network"
	"github.com/gonum/graph/simple"
)

func init() {
	for _, raw := range []struct {
		name string
		set  []set
	}{
		{"zachary", zachary},
		{"blondel", blondel},
	} {
		g := simple.NewUndirectedGraph(0, 0)
		for u, e := range raw.set {
			// Add nodes that are not defined by an edge.
			if !g.Has(simple.Node(u)) {
				g.AddNode(simple.Node(u))
			}
			for v := range e {
				g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v), W: 1})
			}
		}

		nodes := g.Nodes()
		sort.Sort(ordered.ByID(nodes))

		fmt.Printf("%s = []set{\n", raw.name)
		rank := network.PageRank(asDirected{g}, 0.85, 1e-8)
		for _, u := range nodes {
			to := g.From(nodes[u.ID()])
			sort.Sort(ordered.ByID(to))
			var links []int
			for _, v := range to {
				if rank[u.ID()] <= rank[v.ID()] {
					links = append(links, v.ID())
				}
			}

			if links == nil {
				fmt.Printf("\t%d: nil, // rank=%.4v\n", u.ID(), rank[u.ID()])
				continue
			}

			fmt.Printf("\t%d: linksTo(", u.ID())
			for i, v := range links {
				if i != 0 {
					fmt.Print(", ")
				}
				fmt.Print(v)
			}
			fmt.Printf("), // rank=%.4v\n", rank[u.ID()])
		}
		fmt.Println("}")
	}
}

type asDirected struct{ *simple.UndirectedGraph }

func (g asDirected) HasEdgeFromTo(u, v graph.Node) bool {
	return g.UndirectedGraph.HasEdgeBetween(u, v)
}
func (g asDirected) To(v graph.Node) []graph.Node { return g.From(v) }
