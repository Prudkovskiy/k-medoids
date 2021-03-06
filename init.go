package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
)

// FileReader работает с файлом входных данных
type FileReader struct {
	FileNameIn  string
	FileNameOut string
	Matrix      [][]float64 // матрица смежности
	Kmax        int         // верхняя граница кластеризации
}

// Read считывает данные из файла
func (r *FileReader) Read() ([]byte, error) {
	file, err := os.Open(r.FileNameIn)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// get the file size
	stat, _ := file.Stat()
	// read the file
	data := make([]byte, stat.Size())
	file.Read(data)

	return data, nil
}

func (r *FileReader) Write(data []byte) error {
	file, err := os.Create(r.FileNameOut)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(data)
	return nil
}

// Pack запаковывает пришедшую мапу кластеров в слайс байт
// при этом он сортирует ее по ключам
func (r *FileReader) Pack(clusters map[int][]int) []byte {

	// Создаем однозначное представление кластеров
	// для сравнения ожидание/реальность в тестах
	sortedKeys := make([]int, len(clusters)) // O(Kmax)
	i := 0
	for key := range clusters {
		sortedKeys[i] = key
		i++
	}
	sort.Slice(sortedKeys, func(i, j int) bool { return sortedKeys[i] < sortedKeys[j] })

	visualClust := make(map[string][]int, len(clusters))
	i = 1
	for _, key := range sortedKeys {
		s := strconv.Itoa(i)
		visualClust[s] = make([]int, 0, len(clusters[key]))
		for _, node := range clusters[key] {
			visualClust[s] = append(visualClust[s], node+1)
		}
		sort.Slice(visualClust[s], func(k, j int) bool { return visualClust[s][k] < visualClust[s][j] })
		i++
	}

	data, _ := json.Marshal(visualClust)
	return data
}

// Unpack распаковывает содержимое файла в матрицу смежности и Kmax
func (r *FileReader) Unpack(data []byte) error {
	// доходим до конца первой строки
	var start int
	for k, i := range data {
		if i == 10 {
			start = k
			break
		}
	}
	// инициализируем границу кластеризации первой строкой файла
	k, err := strconv.Atoi(string(data[:start-1]))
	if err != nil {
		return err
	}
	r.Kmax = k

	// достаем списки смежности из файла в формате yaml
	var info []map[string](map[string]float64)
	err = json.Unmarshal(data[start:], &info)
	if err != nil {
		return err
	}

	// выделяем место под матрицу смежности n*n
	n := len(info)

	if n < r.Kmax {
		return fmt.Errorf("uncorrect: Kmax > number of nodes in your graph")
	}

	r.Matrix = make([][]float64, n)
	for i := 0; i < n; i++ {
		r.Matrix[i] = make([]float64, n)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i != j {
				r.Matrix[i][j] = math.Inf(1)
			} else {
				r.Matrix[i][j] = 0
			}
		}
	}

	// дополняем ее весами ребер
	for _, adjacency := range info {
		for node, spisok := range adjacency {
			for friendnode, weight := range spisok {
				nodeInt, err := strconv.Atoi(node)
				if err != nil {
					return err
				}
				friendNodeInt, err := strconv.Atoi(friendnode)
				if err != nil {
					return err
				}
				r.Matrix[nodeInt-1][friendNodeInt-1] = weight
			}
		}
	}

	return nil
}
