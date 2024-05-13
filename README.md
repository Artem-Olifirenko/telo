# Organism
Абстракция для управления readiness, liveness пробами в самом микросервисе

# Описание абстракции
Данная маленькая библиотека обладает 2 сущностями: организм (Organism), конечность (Limb). Организм и конечности могут 
быть живы/мертвы (liveness) и готовы к работе или нет (readiness). Если хотя бы одна из конечностей не готова к работе, 
то весь организм считается не готовым к работе. Если хотя бы одна из конечностей умерла, весь организм считается 
мертвым.

# Примеры
Возьмем схематичный http-сервер для ответа на liveness, readiness пробы
```go
package http

func NewServer(organism *Organism) *Organism {
	return &Server{organism: organism}
}

type Server struct {
	organism *Organism
}

func (s *Server) Run() error {
	// ...
}

func (s *Server) HandleLiveness() string {
	if s.organism.IsAlive() {
		return "OK"
    }
    
    return ""
}

func (s *Server) HandleReadiness() string {
    if s.organism.IsReady() {
    	return "OK"
    }
    
    return ""
}
```

Дальше возьмем какой-нибудь схематичный main.go, в котором надо запустить один сервисный http-сервер, и один клиентский 
http rest-api сервер. Логика следующая:
- создаем организм
- для каждого сервиса
    - отрастить конечность
    - в горутине, по запуску сервиса:
        - отложенно умертвить конечность, так как сервис должен жить постоянно и если он вышел, значит данный микросервис 
        нельзя считать больше эивым
        - во время запуска сервиса необходимо проверять его готовность к работе и как проверили, помечаем конечность как 
        готовую к работе 
- в конце помечаем организм как готовый к работе (это нужно чтобы мы успели запустить все сервисы)

```go
package main

import (
	"code.citik.ru/gobase/organism"
	"context"
	"fmt"
	"sync"
)

func main() {
	ctx := context.Background()
	// создаем сам контекст
	org := organism.NewOrganism()
	// создаем waitgroup для блокировки выполнения метода main, пока всео операции не будут завершены. Можно применять 
	// любой способ для блокировки
	wg := sync.WaitGroup{}

	// Запускаем сервисный сервер. Если данный сервер запустится молниеносно и к нему через 1нс обратится kubernetes, то 
	// сервисное апи отдаст OK на liveness пробу и ничего на readiness
	wg.Add(1)
	// отращиваем конечность для сервисного апи
	serviceServerLimb := org.GrowLimb()
	go func() {
		defer wg.Done()
		// данный сервер должен работать всю жизнь приложения, если он вдруг вышел, то помечаем конечность как мертвую
		defer serviceServerLimb.Die()
		server := http.NewServer()
		// запускам метод по проверку готовности http-сервера принимать запросы и отвечать на них
		go checkForReadyHttpServer(server, serviceServerLimb)
		err := server.Run()
		if err != nil {
			fmt.Printf("service-server error: %s", err)
		}
	}()
	
	wg.Add(1)
	// отращиваем конечность для клиентского апи
	apiServerLimb := org.GrowLimb()
	go func() {
		defer wg.Done()
		// данный сервер должен работать всю жизнь приложения, если он вдруг вышел, то помечаем конечность как мертвую
		defer apiServerLimb.Die()
		server := restpai.NewServer()
		// запускам метод по проверку готовности http-сервера принимать запросы и отвечать на них
		go checkForReadyHttpServer(server, apiServerLimb)
		err := server.Run()
		if err != nil {
			fmt.Printf("api-server error: %s", err)
		}
	}()
	
	go func() {
		<-ctx.Done()
		org.Die()
	}()
	org.Ready()
	wg.Wait()
}

// checkForReadyHttpServer примерный метод по проверке готовности http-сервера
func checkForReadyHttpServer(server http.Server, limb *organism.Limb) {
	// ...
	if serviceLaunchedProperly {
		limb.Ready()
    }
    // ...
}
```

Contribute
==========
1. запустите тесты
    ```shell script
    docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v1.40 golangci-lint run -v --deadline=5m \
    && go test --race -coverprofile=cover_profile.out ./... \
    && test "$( go tool cover -func=cover_profile.out | tail -n 1 | awk '{ print $3 }' | sed -e 's/^\([0-9]*\).*$/\1/g' )" -eq 100
    ```
2. Пушните, создайте pull-request и назначьте ревьювером мейнтейнера

Maintainer
====
Nikita Sapogov <sapogov.n@citilink.ru> 