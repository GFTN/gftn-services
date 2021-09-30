// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package utility

import (
	"github.com/golang/geo/r3"
	"github.com/golang/geo/s2"
	"github.com/GFTN/gftn-services/gftn-models/model"
)

func ContainsPoint(ids []string, candidatePoints map[string][]model.Coordinate, requestCoordinate model.Coordinate) ([]string, error) {

	targetPoint := s2.Point{r3.Vector{*requestCoordinate.Long, *requestCoordinate.Lat, 0}}
	var newIds []string

	for index, geos := range candidatePoints {
		var p s2.Point
		query := s2.NewConvexHullQuery()

		for _, geo := range geos {
			p = s2.Point{r3.Vector{*geo.Long, *geo.Lat, 0}}
			query.AddPoint(p)
		}
		//after complete building convex hull
		//detect if the hull contains the point
		// if not, then we will remove the location id from the candidate array
		hull := query.ConvexHull()
		if hull.ContainsPoint(targetPoint) {
			newIds = append(newIds, index)
		}
	}
	return newIds, nil
}
