package parse

import (
	"strings"

	"o.o/backend/com/main/location"
	"o.o/backend/pkg/common/validate"
)

func ParseAddress(input string) (best *Result, results []Result) {
	norm := validate.NormalizeSearch(input)
	results = parse(norm)
	if len(results) == 0 {
		return nil, nil
	}

	idx, min := 0, len(results[0].Address)
	for i, r := range results {
		if len(r.Address) < min {
			idx, min = i, len(r.Address)
		}
	}
	r := results[idx]
	return &r, results
}

func parse(norm string) resultList {
	sts := make(stateList, 0, 16)
	rts := make(resultList, 0, 16)
	out := &state{stateList: &sts, resultList: &rts}
	nd := &nodes[nodeZero]
	nd.consume(nd, norm, out)
	return rts
}

type Result struct {
	ProvinceCode string
	DistrictCode string
	WardCode     string
	Address      string
}

func (r *Result) String() string {
	if r == nil {
		return "nil"
	}
	var b strings.Builder
	prov := location.ProvinceIndexCode[r.ProvinceCode]
	if prov != nil {
		b.WriteString(prov.Name)
	}
	b.WriteString("|")
	dist := location.DistrictIndexCode[r.DistrictCode]
	if dist != nil {
		b.WriteString(dist.Name)
	}
	b.WriteString("|")
	ward := location.WardIndexCode[r.WardCode]
	if ward != nil {
		b.WriteString(ward.Name)
	}
	b.WriteString("|")
	b.WriteString(r.Address)
	return b.String()
}

type resultList []Result

func (rs *resultList) collect(r Result) {
	*rs = append(*rs, r)
}

type state struct {
	Result

	key     int
	matched string

	*node
	*stateList
	*resultList
}

type stateList []state

func (n stateList) current() state {
	if len(n) > 0 {
		return n[len(n)-1]
	}
	return state{}
}

func (n stateList) last() state {
	if len(n) > 1 {
		return n[len(n)-2]
	}
	return state{}
}

func (n *stateList) pop() {
	*n = (*n)[:len(*n)-1]
}

func (s *state) push(r state) *state {
	r.stateList = s.stateList
	r.resultList = s.resultList
	*s.stateList = append(*s.stateList, r)
	return &(*s.stateList)[len(*s.stateList)-1]
}

type stream struct {
	input  string
	i, max int
}

func newStreamer(input string, max int) stream {
	s := stream{
		input: input,
		max:   max,
		i:     len(input),
	}
	return s
}

func (s *stream) next() string {
	for i := s.i; i >= 0; i-- {
		if i == 0 || s.input[i-1] == ' ' {
			s.i = i - 1
			return s.input[i:]
		}
	}
	return ""
}

func (s *stream) remain() string {
	if s.i < 0 {
		return ""
	}
	return s.input[:s.i]
}

type node struct {
	key       int
	consume   func(nd *node, input string, out *state)
	optNodes  []int
	nextNodes []int
}

func (nd *node) next(matched, input string, out *state) {
	nd.nextWith(matched, input, out, nd.optNodes)
	nd.nextWith(matched, input, out, nd.nextNodes)
}

func (nd *node) nextWith(matched, input string, out *state, nextNodes []int) {
	last := out.current()
	for _, ni := range nextNodes {
		nextNode := &nodes[ni]
		out = out.push(state{Result: last.Result, key: ni, matched: matched, node: nextNode})
		nextNode.consume(nextNode, input, out)
		out.pop()
	}
}

const (
	nodeZero = iota
	nodeProv
	nodeDist
	nodeWard
	nodePProv
	nodePDist
	nodePWard
	nodeAddr
	nodeLen
)

var nodes [nodeLen]node

func init() {
	nodes = [nodeLen]node{
		nodeZero: {
			consume: func(nd *node, input string, out *state) {
				nd.next("", input, out)
			},
			nextNodes: []int{nodeProv, nodeDist},
		},
		nodePProv: {
			consume: func(nd *node, input string, out *state) {
				for _, s := range location.PrefixProvince {
					if strings.HasSuffix(input, s) {
						nd.nextWith(s, trimRight(input, s), out, out.last().nextNodes)
					}
				}
			},
		},
		nodePDist: {
			consume: func(nd *node, input string, out *state) {
				for _, s := range location.PrefixDistrict {
					if strings.HasSuffix(input, s) {
						nd.nextWith(s, trimRight(input, s), out, out.last().nextNodes)
					}
				}
			},
		},
		nodePWard: {
			consume: func(nd *node, input string, out *state) {
				for _, s := range location.PrefixWard {
					if strings.HasSuffix(input, s) {
						nd.nextWith(s, trimRight(input, s), out, out.last().nextNodes)
					}
				}
			},
		},
		nodeProv: {
			consume: func(nd *node, input string, out *state) {
				st := newStreamer(input, maps.maxWordProvince)
				for part := st.next(); part != ""; part = st.next() {
					remain := st.remain()
					codes := maps.provByName[part]
					for _, code := range codes {
						out.ProvinceCode = code
						nd.next(part, remain, out)
						out.ProvinceCode = ""
					}
				}
			},
			optNodes:  []int{nodePProv},
			nextNodes: []int{nodeDist, nodeWard},
		},
		nodeDist: {
			consume: func(nd *node, input string, out *state) {
				st := newStreamer(input, maps.maxWordDistrict)
				for part := st.next(); part != ""; part = st.next() {
					for _, code := range maps.distByName[part] {
						dist := location.DistrictIndexCode[code]
						if out.ProvinceCode != "" && dist.ProvinceCode != out.ProvinceCode {
							continue
						}
						out.DistrictCode = code
						nd.next(part, st.remain(), out)
						out.ProvinceCode = ""
					}
				}
			},
			optNodes:  []int{nodePDist},
			nextNodes: []int{nodeWard, nodeAddr, nodeProv},
		},
		nodeWard: {
			consume: func(nd *node, input string, out *state) {
				st := newStreamer(input, maps.maxWordWard)
				for part := st.next(); part != ""; part = st.next() {
					for _, code := range maps.wardByName[part] {
						ward := location.WardIndexCode[code]
						if out.DistrictCode != "" && ward.DistrictCode != out.DistrictCode {
							continue
						}
						if out.ProvinceCode != "" && ward.ProvinceCode != out.ProvinceCode {
							continue
						}
						out.WardCode = code
						nd.next(part, st.remain(), out)
						out.WardCode = ""
					}
				}
			},
			optNodes:  []int{nodePWard},
			nextNodes: []int{nodeAddr},
		},
		nodeAddr: {
			consume: func(nd *node, input string, out *state) {
				out.Address = input
				out.collect(out.Result)
				out.Address = ""
			},
		},
	}
	for k := range nodes {
		nodes[k].key = k
	}
}

func trimRight(s string, cut string) string {
	return s[:len(s)-len(cut)-1]
}
