package main

import (
    "fmt"
    "sync"
)

func main() {
    // Определим матрицы A и B
    A := [][]float64{
        {1, 2, 3},
        {4, 5, 6},
        {7, 8, 9},
    }
    B := [][]float64{
        {9, 8, 7},
        {6, 5, 4},
        {3, 2, 1},
    }

    // Результирующая матрица C
    C := make([][]float64, len(A))
    for i := range C {
        C[i] = make([]float64, len(B[0]))
    }

    // Количество потоков
    numThreads := 4

    // Канал для передачи строк матрицы C
    rows := make(chan []float64)

    // Запуск потоков
    var wg sync.WaitGroup
    wg.Add(numThreads)
    for i := 0; i < numThreads; i++ {
        go worker(A, B, C, rows, &wg)
    }

    // Распределение строк матрицы C по потокам
    for i := range C {
        rows <- C[i]
    }

    close(rows)

    // Ожидание завершения потоков
    wg.Wait()

    // Вывод результирующей матрицы C
    fmt.Println("Матрица C:")
    for i := range C {
        for j := range C[i] {
            fmt.Printf("%f ", C[i][j])
        }
        fmt.Println()
    }
}

func worker(A [][]float64, B [][]float64, C [][]float64, rows chan []float64, wg *sync.WaitGroup) {
    defer wg.Done()

    for row := range rows {
        // Вычисление элементов строки матрицы C
        for j := range B[0] {
            var sum float64
            for k := range A[0] {
                sum += A[row][k] * B[k][j]
            }
            C[row][j] = sum
        }

        // Отправка следующей строки на обработку
        select {
        case nextRow := <-rows:
            row = nextRow
        default:
            return
        }
    }
}
