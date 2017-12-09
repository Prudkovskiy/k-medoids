# Алгоритм PAM
Рассмотрена и реализована усовершенствованная версия одного из распространенных алгоритмов кластеризации графов – PAM (Partitioning Around Medoids). 
Разбиение на кластеры зависит от входных данных, в отличие от базового метода, где количество кластеров считается изначально известным.

## Формат входных и выходных данных
Входные данные будут содержаться в файле in.txt, выходные – в out.txt.

### Формат входных данных:
В первой строке записано число k_max (int) – верхняя граница кластеризации.
Затем в формате yaml сам граф в виде списков смежности

### Формат выходных данных:
Кластеры списками смежностей в формате yaml (k-шт.)

Примеры смотреть в test_data

### Инструкция пользования
В зависимости от платформы скачать необходимый файл из папки executable_files. 
В консоли вводить: `./your_exefile input_file.txt output_file.txt`
Для запуска тестов в папке с test_coverage.go вводит: `go test -v`
