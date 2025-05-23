# Subpub-project
Это система подписки и публикации (Pub/Sub), реализованная с использованием gRPC. Она поддерживает подписку на темы, публикацию сообщений и возможность закрывать соединения.

## Стек 
- [Go](https://golang.org/) — язык программирования для сервера
- [gRPC](https://grpc.io/) — для общения между сервисами
- [Protocol Buffers](https://developers.google.com/protocol-buffers) — для описания сообщений и сервисов
- [Testify](https://github.com/stretchr/testify) — для тестирования
- [Protoc](https://grpc.io/docs/protoc-installation/) — компилятор для Protocol Buffers


## Установка и запуск проекта

### 1. Клонирование репозитория
```bash
git clone https://github.com/voronkov44/subpub-project.git
```

### 2. Переход в корневую директорию 
```bash
cd subpub-project
```

### 3. Запуск проекта
#### 3.1. Открываем первый терминал и запускаем gRPC сервер
```bash
go run cmd/main.go
```

#### 3.2. Открываем второй терминал и подписываемся на топик(слушаем порт)
*Требуется установка [grpcurl](https://github.com/fullstorydev/grpcurl), если не установлен смотрите [зависимости.](https://github.com/voronkov44/subpub-project?tab=readme-ov-file#%D1%83%D1%81%D1%82%D0%B0%D0%BD%D0%BE%D0%B2%D0%BA%D0%B0-%D0%BF%D0%B0%D0%BA%D0%B5%D1%82%D0%B0-grpcurl)*

Подписка — это стрим, он ждёт сообщения в консоли:

```bash
grpcurl -plaintext -d '{"key":"test-topic"}' localhost:50051 subpub.PubSub/Subscribe
```
*После этой команды терминал "зависает" и будет выводить все новые события.*

#### 3.3. Открываем третий терминал и публикуем сообщение в топик
Отправляем события в топик:
```bash
grpcurl -plaintext -d '{"key":"test-topic", "data":"Hello from grpcurl!"}' localhost:50051 subpub.PubSub/Publish
```
*Результат должен быть `{}`, а на втором терминале вывод нашего сообщения*

### 4. Запуск проекта через docker
Для удобства был написан `dockerfile`, через который можно запустить проект

*Требуется установка [docker](https://www.docker.com/products/docker-desktop/), если не установлен, смотрите [зависимости.](https://github.com/voronkov44/subpub-project?tab=readme-ov-file#%D1%83%D1%81%D1%82%D0%B0%D0%BD%D0%BE%D0%B2%D0%BA%D0%B0-docker)*

#### 4.1 Сборка образа:
```bash
docker build -t subpub:v1 .
```

#### 4.2 Запуск проекта в контейнере:
```bash
docker run -d --name subpub-server -p 50051:50051 subpub:v5
```

#### 4.3 Проверка что наш контейнер действительно создался:
```bash
docker ps
```

Так же можно проверить логи
```bash
docker logs subpub-server
```

#### 4.4 Открываем второй терминал и подписываемся на топик(слушаем порт):
```bash
grpcurl -plaintext -d '{"key":"test-topic"}' localhost:50051 subpub.PubSub/Subscribe
```
*После этой команды терминал "зависает" и будет выводить все новые события.*

#### 4.5. Открываем третий терминал и публикуем сообщение в топик
Отправляем события в топик:
```bash
grpcurl -plaintext -d '{"key":"test-topic", "data":"Hello from grpcurl!"}' localhost:50051 subpub.PubSub/Publish
```
*Результат должен быть `{}`, а на втором терминале вывод нашего сообщения*


#### Общие команды для управления контейнерами:
```
# Просмотр запущенных контейнеров
docker ps

#Просмотр всех контейнеров, включая остановленные
docker ps -a

# Остановка контейнера
docker stop <container_id>

# Удлание контейнера
docker rm <container_id>

# Удаление образа
docker rmi <image_id>

# Очистка системы
docker system prune
```



## 5. Тесты
В проекте реализовано 4 unit-теста:
| Тест                                 | Что проверяет                                                             |
| :----------------------------------- | :------------------------------------------------------------------------ |
| `TestMultipleIndependentSubscribers` | Независимость подписчиков и получение всех сообщений в правильном порядке |
| `TestPublishSubscribe`               | Корректную доставку сообщения подписчику                                  |
| `TestUnsubscribe`                    | Корректную работу отписки от топика                                       |
| `TestCloseWithContextCancel`         | Обработку отмены контекста при закрытии                                   |


### Запуск тестов:
```bash
go test -v ./internal/subpub
```

### Описание тестов:

1. `TestMultipleIndependentSubscribers`

**Что проверяет:**

Что несколько независимых подписчиков на один топик получают все опубликованные сообщения в правильном порядке, независимо от скорости обработки сообщений.

**Как работает:**

- Создаются 3 подписчика на топик test-topic.

- Каждый подписчик сохраняет принятые сообщения в свой слайс.

- Второй подписчик искусственно «медленный» (пауза в 200 мс).

- Публикуются 3 сообщения.

- Проверяется, что все подписчики получили все 3 сообщения в правильном порядке.

- Если хотя бы один не успеет — тест упадёт по таймауту 3 сек.

**Этот тест гарантирует, что подписчики работают независимо и получают все сообщения без потерь.**



2. `TestPublishSubscribe`


**Что проверяет:**

Что публикация сообщения успешно доставляет его подписчику.

**Как работает:**

- Создаётся подписчик на test-topic.

- Публикуется сообщение "hello".

- Подписчик проверяет, что получил именно "hello".

- Завершается после получения.

- Подписка отписывается.

**Минимальный sanity check, что публикация и подписка работают.**

3. `TestUnsubscribe`


**Что проверяет:**

Что отписка от подписки действительно предотвращает получение сообщений.

**Как работает:**

- Создаётся подписчик на topic.

- Подписчик сразу отписывается.

- Публикуется сообщение.

- Проверяется, что обработчик не был вызван (флаг called остался false).

**Тестирует корректность работы метода Unsubscribe.**

4. `TestCloseWithContextCancel`

**Что проверяет:**

Что метод Close корректно обрабатывает отменённый контекст.

**Как работает:**

- Создаётся контекст с отменой.

- Вызывается sp.Close(ctx).

- Проверяется, что метод вернул context.Canceled.

**Это проверка корректной обработки отмены контекста при закрытии pubsub'a.**



















## Зависимости
### Установка пакета [grpcurl](https://github.com/fullstorydev/grpcurl)

MacOS (через Homebrew):
```bash
brew install grpcurl
```

Linux:

Для установки на Linux воспользуйтесь установкой через пакет [Go](https://go.dev/doc/install):

```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```
Или скачайте бинарник с релизов на [GitHub](https://github.com/fullstorydev/grpcurl/releases)

Windows:

Для установки на Windows воспользуйтесь установкой через пакет [Go](https://go.dev/doc/install):

```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```
Или скачайте exe-файл с релизов на [GitHub](https://github.com/fullstorydev/grpcurl/releases)
*После установки через go install убедитесь, что ваш $GOPATH/bin или $HOME/go/bin добавлен в PATH.*

### Установка docker
Установка пакета [Docker Engine](https://docs.docker.com/engine/install/)
