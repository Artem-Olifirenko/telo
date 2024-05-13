#!/bin/sh -e
# ВАЖНО!
# Данный скрипт является стандартным для всех библиотек, его нельзя менять, предварительно не обсудив
# с тимлидом.
UNIT_COVERAGE_MIN=100

GOFLAGS=
CGO_ENABLED=0

# Запуск unit-тестов
unit(){
  echo "run unit tests"
  deps
  go test ./...
}

unit_race() {
  echo "run unit tests with race test"
  deps
  go test -race ./...
}

# Запуск go-lint
lint(){
  echo "run linter"
  go mod vendor
  docker run -v $(pwd):/work:ro -w /work golangci/golangci-lint:v1.40 golangci-lint run -v --modules-download-mode=vendor
  rm -Rf vendor
}

gosec() {
  echo "run gosec"
  go mod vendor
  docker run --rm -it -w /app/ -v $(pwd)/:/app securego/gosec:v2.5.0 -exclude=G307,G304 ./...
  rm -Rf vendor
}

fmt() {
  echo "run go fmt"
  go fmt ./...
}

vet() {
  echo "run go vet"
  go vet ./...
}

unit_coverage() {
  echo "run test coverage"
  go test -coverprofile=cover_profile.out ./...
  CUR_COVERAGE=$( go tool cover -func=cover_profile.out | tail -n 1 | awk '{ print $3 }' | sed -e 's/^\([0-9]*\).*$/\1/g' )
  rm cover_profile.out
  if [ "$CUR_COVERAGE" -lt $UNIT_COVERAGE_MIN ]
  then
    echo "coverage is not enough $CUR_COVERAGE < $UNIT_COVERAGE_MIN"
    return 1
  else
    echo "coverage is enough $CUR_COVERAGE > $UNIT_COVERAGE_MIN"
  fi
}

# Запуск всех тестов
test(){
  fmt
  vet
  unit
  unit_race
  unit_coverage
  lint
  gosec
}

# Подтянуть зависимости
deps(){
  go get ./...
}

# Собрать исполняемый файл
build(){
  deps
  go build
}

# Добавьте сюда список команд
using(){
  echo "Укажите команду при запуске: ./run.sh [command]"
  echo "Список команд:"
  echo "  unit - запустить unit-тесты"
  echo "  unit_race - запуск unit тестов с проверкой на data-race"
  echo "  unit_coverage - запуск unit тестов и проверка покрытия кода тестами"
  echo "  lint - запустить все линтеры"
  echo "  test - запустить все тесты"
  echo "  deps - подтянуть зависимости"
  echo "  build - собрать приложение"
  echo "  fmt - форматирование кода при помощи 'go fmt'"
  echo "  vet - проверка правильности форматирования кода"
  echo "  gosec - запуск статического анализатора безопасности"
}

############### НЕ МЕНЯЙТЕ КОД НИЖЕ ЭТОЙ СТРОКИ #################

command="$1"
if [ -z "$command" ]
then
 using
 exit 0;
else
 $command $@
fi