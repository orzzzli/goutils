package graph

import (
	"errors"
	"strconv"
)

const (
	INFINITY int = 65536 // 不可能值
)

//邻接矩阵
type AdjacencyMatrix struct {
	vertexNumber int           //顶点数量
	degreeNumber int           //边数量
	vertexes     []interface{} //顶点数组
	degrees      [][]int       //矩阵
}

type DegreeWeight struct {
	v0     int //顶点a index，起点
	v1     int //顶点b index，终点
	weight int //权重
}

func New(vertexes []interface{}) *AdjacencyMatrix {
	vertexNumber := len(vertexes)
	adMatrix := &AdjacencyMatrix{
		vertexNumber: vertexNumber,
		degreeNumber: 0,
		vertexes:     vertexes,
	}
	tempDegrees := make([][]int, vertexNumber)
	//初始化矩阵，全部边值为不可达
	for i := 0; i < vertexNumber; i++ {
		tempDegrees[i] = make([]int, vertexNumber)
		for index, _ := range tempDegrees[i] {
			tempDegrees[i][index] = INFINITY
		}
	}
	adMatrix.degrees = tempDegrees
	return adMatrix
}

//新增边
func (am *AdjacencyMatrix) AddDegree(v0 int, v1 int, weight int) error {
	if v0 >= am.vertexNumber {
		return errors.New("v0 out of range. v0 is " + strconv.Itoa(v0) + ". max range is " + strconv.Itoa(am.vertexNumber-1))
	}
	if v1 >= am.vertexNumber {
		return errors.New("v1 out of range. v1 is " + strconv.Itoa(v1) + ". max range is " + strconv.Itoa(am.vertexNumber-1))
	}
	am.degrees[v0][v1] = weight
	am.degreeNumber++
	return nil
}

//批量新增边
func (am *AdjacencyMatrix) BatchAddDegrees(degrees []DegreeWeight) error {
	for _, info := range degrees {
		err := am.AddDegree(info.v0, info.v1, info.weight)
		if err != nil {
			return err
		}
	}
	return nil
}

//获取指定顶点出度的全部顶点IndexList
func (am *AdjacencyMatrix) OutDegreeVertexes(vertexIndex int) ([]int, error) {
	if vertexIndex >= am.vertexNumber {
		return nil, errors.New("vertexIndex out of range. vertexIndex is " + strconv.Itoa(vertexIndex) + ". max range is " + strconv.Itoa(am.vertexNumber-1))
	}
	var output []int
	degrees := am.degrees[vertexIndex]
	for index, weight := range degrees {
		if weight != INFINITY {
			output = append(output, index)
		}
	}
	return output, nil
}

//指定顶点0和顶点1之间是否有出度
func (am *AdjacencyMatrix) OutDegree(vertexIndex0 int, vertexIndex1 int) (bool, error) {
	if vertexIndex0 >= am.vertexNumber {
		return false, errors.New("vertexIndex0 out of range. vertexIndex0 is " + strconv.Itoa(vertexIndex0) + ". max range is " + strconv.Itoa(am.vertexNumber-1))
	}
	if vertexIndex1 >= am.vertexNumber {
		return false, errors.New("vertexIndex1 out of range. vertexIndex1 is " + strconv.Itoa(vertexIndex1) + ". max range is " + strconv.Itoa(am.vertexNumber-1))
	}
	if am.degrees[vertexIndex0][vertexIndex1] != INFINITY {
		return true, nil
	}
	return false, nil
}

//指定顶点0和顶点1之间是否有入度
func (am *AdjacencyMatrix) InDegree(vertexIndex0 int, vertexIndex1 int) (bool, error) {
	if vertexIndex0 >= am.vertexNumber {
		return false, errors.New("vertexIndex0 out of range. vertexIndex0 is " + strconv.Itoa(vertexIndex0) + ". max range is " + strconv.Itoa(am.vertexNumber-1))
	}
	if vertexIndex1 >= am.vertexNumber {
		return false, errors.New("vertexIndex1 out of range. vertexIndex1 is " + strconv.Itoa(vertexIndex1) + ". max range is " + strconv.Itoa(am.vertexNumber-1))
	}
	if am.degrees[vertexIndex1][vertexIndex0] != INFINITY {
		return true, nil
	}
	return false, nil
}

//指定顶点0和顶点1之间是否有出入度
func (am *AdjacencyMatrix) InAndOutDegree(vertexIndex0 int, vertexIndex1 int) (bool, error) {
	if vertexIndex0 >= am.vertexNumber {
		return false, errors.New("vertexIndex0 out of range. vertexIndex0 is " + strconv.Itoa(vertexIndex0) + ". max range is " + strconv.Itoa(am.vertexNumber-1))
	}
	if vertexIndex1 >= am.vertexNumber {
		return false, errors.New("vertexIndex1 out of range. vertexIndex1 is " + strconv.Itoa(vertexIndex1) + ". max range is " + strconv.Itoa(am.vertexNumber-1))
	}
	inDegree, err := am.InDegree(vertexIndex0, vertexIndex1)
	if err != nil {
		return false, err
	}
	outDegree, err := am.OutDegree(vertexIndex0, vertexIndex1)
	if err != nil {
		return false, err
	}
	return inDegree && outDegree, nil
}
