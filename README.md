# Balance Microservice

![Build Status](https://github.com/phlx/balance-go/workflows/validate-markdown/badge.svg)

## Что это?

Это микросервис для управления балансом пользователей, 
сделанный по [тестовому заданию](https://github.com/avito-tech/job-backend-trainee-assignment) 
ради фана.  
Написан на Go 1.15.  

## Что нужно, чтобы запустить?

Git, чтобы скачать этот репозиторий.  
Docker и Docker Compose, чтобы запустить всё нужное
в виртуальной среде.  
Этого достаточно для запуска приложения.  
Установленный Go нужен только для локального запуска тестов.  
Для других вещей желательно иметь Unix-подобное 
окружение с поддержкой Makefiles.  
На Windows теперь есть [WSL2](https://docs.microsoft.com/en-us/windows/wsl/compare-versions#whats-new-in-wsl-2),
а значит в Unix-подобное окружение на нём можно зайти командой `wsl`.  
Разумеется, предварительно [установив WSL](https://docs.microsoft.com/ru-ru/windows/wsl/install-win10).  

## Как запустить?
Стандартно для Docker Compose:
```bash
docker-compose up -d
```
То же самое делает команда
```bash
make
```

## Как потыкаться?
После запуска на главной странице по адресу [http://localhost](http://localhost)
доступен Swagger UI с загруженной [схемой микросервиса](api/swagger.json).  
Также, описание API есть в виде [опубликованной коллекции Postman](https://documenter.getpostman.com/view/6261504/TVKD2czt).  
Сама коллекция, которую можно импортировать в свой Postman
находится в каталоге [api/balance.postman_collection.json](api/balance.postman_collection.json).  

## А зачем Graphite и Grafana?
По адресу [http://localhost:3030](http://localhost:3030)
можно обнаружить Graphite, а по адресу [http://localhost:3080](http://localhost:3080)
находится Grafana.  
Graphite собирает данные с помощью Statsd. В приложении есть
немного мест, где отправляются простейшие метрики.  
Соответственно, метрики можно посмотреть через сам Graphite.  
Либо настроить их просмотр через Grafana, где можно 
настроить богатые и красивые дашборды.

## А почему дашборды в Grafana не настроены?
Мне было лень разбираться с ними. 🙂  
Так и быть — `// TODO: настроить импорт дашбордов в Grafana`  
Самостоятельно можно подключить к Grafana нужный source 
на странице [http://localhost:3080/datasources](http://localhost:3080/datasources).  
Нажать `Add data source`, выбрать `Graphite` и в разделе `HTTP`
в поле `URL` прописать адрес `http://graphite:80`.  
И далее уже поиграться с дашбордами.

## Какие ещё команды предусмотрены?
Быстрый обзор команд можно посмотреть с помощью
```bash
make help
```
Команда выведет команды и их описание:
```bash
help             This help overview
run              Run application in Docker Compose
build            Build application via Docker Compose
rerun            Stop Docker Compose, build application container, run Docker Compose
test             Run all tests inside app container
fmt              Run fmt with latest docker image of golang
lint             Run lint with latest image of golangci/golangci-lint
concurrent-test  Run test for concurrent transactions in PostgreSQL, see internal/postgres/isolation.go
```

## Запуск локального дебага — это как?
В запущенных Docker контейнерах расположено приложение
и его зависимости в виде PostgreSQL, Redis и StatsD.  
При этом само приложение запущено в контейнере app и выведено
наружу через порт 3000.  

Приложение можно дебажить локально через JetBrains Goland.  
Для этого в IDE нужно [поставить breakpoint в нужном месте и запустить
main/service/main.go](https://blog.jetbrains.com/go/2019/02/06/debugging-with-goland-getting-started/#debugging-an-application).
При этом в [созданной конфигурации дебага](https://www.jetbrains.com/help/go/creating-run-debug-configuration-for-tests.html#test-configuration-for-a-package)
нужно установить в поле Program arguments значение `--debug=true`.  

Запуск приложения со включенным флагом `debug` заставит его
использовать файл окружения `configs/.env.debug`, благодаря чему само
приложение будет запущено на текущей (хостовой) машине,
но подключаться оно будет к зависимостям уже запущенным в контейнерах
Docker Compose.

## Как можно вкратце описать устройство микросервиса?
Роутинг и механизм middleware поддерживает gin.  
Каждый запрос проходит через middleware cors, idempotency и logger:
- CORS позволяет для Swagger UI делать запросы от localhost:80 до localhost:3000.  
- Idempotency срабатывает на POST-запросах и возвращает закешированный ответ
для конкретных idempotence key.  
- Logger позволяет логировать запросы по правилам общего логгера 
(используется Logrus) и выводит сообщения нужного уровня и в нужном формате.  

В приложении создаётся контейнер, но он не имеет отношение к терминологии
Dependency Injection. Контейнер в данном случае — условное название
для объекта, собирающего зависимости и явным образом прокидывающего
их в нужные методы.

Бизнес-логика приложения собрана в одной структуре — это
`internal/services/core/service`. Методам этой структуры на вход поступают
уже провалидированные и подготовленные данные.

Валидацией и подготовкой данных занимаются контроллеры из каталога `internal/handlers`.
В качестве валидатора используется `go-playground/validator`.

CoreService обеспечивает логику работы всех нужных методов по работе с балансом.
В методе Get() используется прокинутый в coreService клиент курсов валют exchangeRatesService.
Все методы coreService делают запросы в Postgres с помощью пакета go-pg, используя 
модели каталога `internal/models`.
Соответствующие моделям таблицы создаются (если их нет) автоматически при старте сервиса. 

В качестве клиента Redis используется пакет go-redis для обеспечения 
идемпотентности (`internal/middlewares/idempotency.go`) и для кеширования 
курсов валют (`internal/exchangerates/service.go`).

Для финансовых вычислений используется библиотека shopspring/decimal.  
Она устраняет недостатки арифметики стандарта IEEE 754 (Floating-Point Arithmetic).  
Пример: прибавление к 0 числа 0.1 тысячу раз вместо ожидаемого 
числа 10 [выдаст 9.999999999999831](https://play.golang.org/p/TQBd4yJe6B).  
Также, для округлений используется банковское округление (до ближайшего чётного) 
для избежания накопления ошибки округления.

## А что с транзакциями и MVCC?
Ошибка ["грязного" чтения](https://postgrespro.ru/docs/postgrespro/12/transaction-iso) исключена благодаря 
уровню изоляции транзакций по умолчанию `READ COMMITTED` и запросу `SELECT ... FOR UPDATE`.
Запрос с `FOR UPDATE` в более строгих уровнях изоляции 
[будет вызывать ошибки](https://postgrespro.ru/docs/postgrespro/12/explicit-locking#LOCKING-ROWS).  

Эта ошибка применима к ситуации, когда много параллельных запросов на списание средств могут 
увести баланс пользователя в минус, в то время как он не может быть меньше нуля.  
Эмулировать такую ситуацию можно с помощью скрипта
```bash
make concurrent-test
```
Этот скрипт:
- Зачисляет пользователю с id 99999999 100 рублей
- Выводит баланс
- Десять раз подряд с разницей в 100 мс пытается списать по 50 рублей, 
совершая в транзакции задержку на секунду между чтением баланса и сохранением 
нового баланса и действия (модели transaction)
- Вновь выводит баланс

В результате баланс пользователя равен нулю и создана одна транзакция 
зачисления на 100 рублей и две транзакции списания по 50 рублей.  
В случае ошибки "грязного" чтения в такой ситуации баланс мог бы остаться нулевым, но при этом 
создались ли бы лишние банковские транзакции, по которым сумма не была бы равна балансу.

## Почему не используются очереди?
Слишком сложно 🙂  
Если принимать запрос и возвращать на него ответ о том, что задача добавлена в очередь, то каким образом
и куда уведомлять об успешном завершении задачи?  

Ответ на этот вопрос неочевиден.
Для этого можно было бы использовать websockets, принимать и уведомлять о завершении по дуплексному каналу,
но это уже не HTTP API, а HTTP + WS API.  
Такое усложнение добавляет кучу работы к реализации такого сервиса.

В текущем виде микросервис в единственном экземпляре в Докере обеспечивает ~3-3,5k запросов в секунду
на эндпоинт списка транзакций по существующему юзеру с множеством транзакций.
Сомнительно, что в боевом окружении такой микросервис может столкнуться с проблемой высоких нагрузок.

Сетевые проблемы, при которых запрос может уйти несколько раз решены с помощью ключа идемпотентности.
В этом случае даже при таймаутах не возникнет проблем при повторных отправках запроса.

## А где unit и интеграционные тесты?
Юнит-тестов в классическом понимании нет.
Есть пара тестов, проверяющих, что в целом сервис готов запуститься.

И есть множество функциональных тестов в каталоге `test/functional`, 
использующих библиотеку `gavv/httpexpect`,
а вместе с ней и `stretchr/testify`.

Функциональные тесты проверяют кейсы и со стороны API, 
и со стороны записей в базе данных (и в Redis, в том числе).
