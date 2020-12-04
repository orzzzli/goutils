package graph

import (
	"fmt"
	"testing"
)

func TestGraph(t *testing.T) {
	vertexes := []interface{}{
		"a",
		"b",
		"c",
		"d",
	}
	matrix := New(vertexes)
	fmt.Println(matrix.degreeNumber, matrix.vertexNumber, matrix.degrees, matrix.vertexes)
	err := matrix.AddDegree(0, 1, 1)
	fmt.Println(err, matrix.degreeNumber, matrix.degrees)
	err = matrix.AddDegree(0, 4, 1)
	fmt.Println(err, matrix.degreeNumber, matrix.degrees)
	err = matrix.BatchAddDegrees([]DegreeWeight{
		{
			v0:     1,
			v1:     2,
			weight: 1,
		},
		{
			v0:     2,
			v1:     3,
			weight: 2,
		},
		{
			v0:     5,
			v1:     1,
			weight: 3,
		},
	})
	fmt.Println(err, matrix.degreeNumber, matrix.degrees)
}

func TestOutDegreeVertexes(t *testing.T) {
	vertexes := []interface{}{
		"a",
		"b",
		"c",
		"d",
	}
	matrix := New(vertexes)
	matrix.AddDegree(0, 1, 1)
	matrix.AddDegree(1, 2, 1)
	matrix.AddDegree(2, 1, 1)
	matrix.AddDegree(2, 0, 1)
	matrix.AddDegree(3, 1, 1)
	matrix.AddDegree(3, 2, 1)
	fmt.Println(matrix.OutDegreeVertexes(0))
	fmt.Println(matrix.OutDegreeVertexes(1))
	fmt.Println(matrix.OutDegreeVertexes(2))
	fmt.Println(matrix.OutDegreeVertexes(3))
}

func TestOutDegreeInDegreeAndOutAndInDegree(t *testing.T) {
	vertexes := []interface{}{
		"a",
		"b",
		"c",
		"d",
	}
	matrix := New(vertexes)
	matrix.AddDegree(0, 1, 1)
	matrix.AddDegree(1, 2, 1)
	matrix.AddDegree(2, 1, 1)
	matrix.AddDegree(2, 0, 1)
	matrix.AddDegree(3, 1, 1)
	matrix.AddDegree(3, 2, 1)
	fmt.Println(matrix.OutDegree(0, 1))
	fmt.Println(matrix.OutDegree(1, 0))
	fmt.Println(matrix.InDegree(0, 1))
	fmt.Println(matrix.InDegree(1, 0))
	fmt.Println(matrix.InAndOutDegree(1, 0))
	fmt.Println(matrix.InAndOutDegree(1, 2))
}
