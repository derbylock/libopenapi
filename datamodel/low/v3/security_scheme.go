// Copyright 2022 Princess B33f Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package v3

import (
	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/index"
	"github.com/pb33f/libopenapi/utils"
	"gopkg.in/yaml.v3"
)

const (
	SecurityLabel        = "security"
	SecuritySchemesLabel = "securitySchemes"
	OAuthFlowsLabel      = "flows"
)

type SecurityScheme struct {
	Type             low.NodeReference[string]
	Description      low.NodeReference[string]
	Name             low.NodeReference[string]
	In               low.NodeReference[string]
	Scheme           low.NodeReference[string]
	BearerFormat     low.NodeReference[string]
	Flows            low.NodeReference[*OAuthFlows]
	OpenIdConnectUrl low.NodeReference[string]
	Extensions       map[low.KeyReference[string]]low.ValueReference[any]
}

type SecurityRequirement struct {
	ValueRequirements []low.ValueReference[map[low.KeyReference[string]][]low.ValueReference[string]]
}

func (ss *SecurityScheme) FindExtension(ext string) *low.ValueReference[any] {
	return low.FindItemInMap[any](ext, ss.Extensions)
}

func (ss *SecurityScheme) Build(root *yaml.Node, idx *index.SpecIndex) error {
	ss.Extensions = low.ExtractExtensions(root)

	oa, oaErr := low.ExtractObject[*OAuthFlows](OAuthFlowsLabel, root, idx)
	if oaErr != nil {
		return oaErr
	}
	if oa.Value != nil {
		ss.Flows = oa
	}

	return nil
}

func (sr *SecurityRequirement) FindRequirement(name string) []low.ValueReference[string] {
	for _, r := range sr.ValueRequirements {
		for k, v := range r.Value {
			if k.Value == name {
				return v
			}
		}
	}
	return nil
}

func (sr *SecurityRequirement) Build(root *yaml.Node, idx *index.SpecIndex) error {
	if utils.IsNodeArray(root) {
		var requirements []low.ValueReference[map[low.KeyReference[string]][]low.ValueReference[string]]
		for _, n := range root.Content {
			var currSec *yaml.Node
			if utils.IsNodeMap(n) {
				res := make(map[low.KeyReference[string]][]low.ValueReference[string])
				var dat []low.ValueReference[string]
				for i, r := range n.Content {
					if i%2 == 0 {
						currSec = r
						continue
					}
					if utils.IsNodeArray(r) {
						// value (should be) an array of strings
						var keyValues []low.ValueReference[string]
						for _, strN := range r.Content {
							keyValues = append(keyValues, low.ValueReference[string]{
								Value:     strN.Value,
								ValueNode: strN,
							})
						}
						dat = keyValues
					}
				}
				res[low.KeyReference[string]{
					Value:   currSec.Value,
					KeyNode: currSec,
				}] = dat
				requirements = append(requirements,
					low.ValueReference[map[low.KeyReference[string]][]low.ValueReference[string]]{
						Value:     res,
						ValueNode: n,
					})
			}
		}
		sr.ValueRequirements = requirements
	}
	return nil
}