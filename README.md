# Распределённый вычислитель арифметических выражений

Этот проект представляет собой распределённый вычислитель арифметических выражений, который позволяет пользователям отправлять арифметические выражения для вычисления и получать результаты асинхронно. Проект состоит из двух основных компонентов: оркестратора и агента.


## Запуск приложения

Для запуска приложения выполните следующие шаги:

1. Убедитесь, что у вас установлен Go (версия 1.16 или выше).
2. Склонируйте репозиторий:
bash
   git clone 
   cd <имя_папки_репозитория>
   

3. Установите необходимые зависимости:
bash
   go mod tidy
   

4. Запустите приложение:
bash
   go run cmd/main.go
   

## API

### Оркестратор

Оркестратор предоставляет следующие конечные точки:

#### Добавление вычисления арифметического выражения

- **URL:** `/api/v1/calculate`
- **Метод:** `POST`
- **Тело запроса:**
json
  {
      "expression": "<строка с выражением>"
  }
  

- **Ответ:**
  - Код 201: выражение принято для вычисления
  - Код 422: невалидные данные
  - Код 500: что-то пошло не так
- **Тело ответа:**
json
  {
      "id": "<уникальный идентификатор выражения>"
  }
  

#### Получение списка выражений

- **URL:** `/api/v1/expressions`
- **Метод:** `GET`
- **Ответ:**
  - Код 200: успешно получен список выражений
  - Код 500: что-то пошло не так
- **Тело ответа:**
json
  {
      "expressions": [
          {
              "id": "<идентификатор выражения>",
              "status": "<статус вычисления выражения>",
              "result": "<результат выражения>"
          }
      ]
  }
  

#### Получение выражения по его идентификатору

- **URL:** `/api/v1/expressions/:id`
- **Метод:** `GET`
- **Ответ:**
  - Код 200: успешно получено выражение
  - Код 404: нет такого выражения
  - Код 500: что-то пошло не так
- **Тело ответа:**
json
  {
      "expression": {
          "id": "<идентификатор выражения>",
          "status": "<статус вычисления выражения>",
          "result": "<результат выражения>"
      }
  }
  

### Агент

Агент получает выражение для вычисления с сервера, вычисляет его и отправляет на сервер результат выражения.

#### Получение задачи для выполнения

- **URL:** `/internal/task`
- **Метод:** `GET`
- **Ответ:**
  - Код 200: успешно получена задача
  - Код 404: нет задачи
  - Код 500: что-то пошло не так
- **Тело ответа:**
json
  {
      "task": {
          "id": "<идентификатор задачи>",
          "arg1": "<имя первого аргумента>",
          "arg2": "<имя второго аргумента>",
          "operation": "<операция>",
          "operation_time": "<время выполнения операции>"
      }
  }
  

#### Прием результата обработки данных

- **URL:** `/internal/task`
- **Метод:** `POST`
- **Тело запроса:**
json
  {
      "id": 1,
      "result": 2.5
  }

Ответ:

Код 200: успешно записан результат

Код 404: нет такой задачи

Код 422: невалидные данные

Код 500: что-то пошло не так

# Переменные окружения
Время выполнения операций задается переменными среды в миллисекундах