# Итоговое задание в рамках курсв "Go - разработчик".

## Выполненные задачи со "звёздочкой":
- [x] Реализуйте возможность определять извне порт при запуске сервера.
- [x] Реализуйте возможность определять путь к файлу базы данных через переменную окружения.
- [ ] Вычисление следующей даты для правил повторения "w" и "m" (недели и месяцы).
- [x] Возможность выбрать задачи через строку поиска.
- [ ] Аутентификация.
- [ ] Docker образ.

## Описание сервиса
Сервис для описания, хранения, изменения и отслеживания датированных записей. К записям можно добавлять комментарий и правила для повторений. Хранение записей осуществляется в базе данных SQLite. Обращение к сервису через web интерфейс. Работа сервиса как в виде Docker контейнера, так и в виде исполняемого бинарного файла.  

## Работа сервиса

* Интерфейс - web.
* Хранение записей - база данных SQLite.
* Порт для подключения - http://localhost:7540/
* Переменное окружение:
  * **```TODO_PORT```** - переменная для определения порда для работы сервиса
  * **```TODO_DBFILE```** - переменная для адреса файла базы данных

## Развертывание сервиса

Если значение глобальных переменных **```TODO_PORT```** и **```TODO_DBFILE```** не будет задано до начала развертывания или тестирования, то порт по умолчанию - **```7540```**, расположение базы данных в том же каталоге, чтои исполняемый файл.
Если глобальная переменная **```TODO_DBFILE```** задана, но по её адресу не найдется файл **```scheduler.db```** с развернутой базой данных, то расположение базы данных будет в том же каталоге, чтои исполняемый файл. 

1. Для запуска сервиса без компиляции:
   * Для запуска используйте команду **```make run```**.
   
2. Для тестирования сервиса:
   * Настройки тестов находятся в файле **```tests/settings.go```**:
     * **```Port = ```** - переменная для номера порта при тестировании;
     * **```DBFile = ```** - переменная для относительного адреса файла базы данных при тестировании;
     * **```FullNextDate = ```** - переменная для параметров тестирования с учетом правил повторения для еженедельных "w" и ежемесячных "m" значений - **```true```**, без учета - **```false```**; 
     * **```Search = ```** - переменная для тестирования возможности поиска задачи - **```true```**, без поиска - **```false```**;
     * **```Token = ```** - токен для тестирования с при реализованной аутентификации.
   * Для запуска всех тестов последовательно используйте команду **```make test```** после успешного выполнения команды **```make run```**.

3. Для запуска сервиса на различных ОС:
   * Среда разработки сервича **```Ubuntu 22.04.4 LTS```**

   * **```Windows 64-bit```** для компиляции сервиса команда **```make build_win64```**. После выполнения команды в каталоге **OS_bin/TODO_win64** будет создан исполняемый файл и необходимые сопутствующие каталоги. Настройки компиляции можно поменять в **```Makefile```**.
   * **```Linux 64-bit```** для компиляции сервиса команда **```make build_lin64```**. После выполнения команды в каталоге **OS_bin/TODO_lin64** будет создан исполняемый файл и необходимые сопутствующие каталоги. Настройки компиляции можно поменять в **```Makefile```**.
   * **```MacOS 64-bit```** для компиляции сервиса команда **```make build_mac64```**. После выполнения команды в каталоге **OS_bin/TODO_mac64** будет создан исполняемый файл и необходимые сопутствующие каталоги. Настройки компиляции можно поменять в **```Makefile```**.

4. Для запуска сервича в Docker контейнере:
   * Перед запуском ознакомьтесь с Dockerfile. Расположение базы данных по умолчанию в контейнере и определяется переменной **```TODO_DBFILE```**. Использование внешнего хранилища описано ниже.

   * Для создания контейнера используйте команду **```make docker_build```**.

   * Для запуска контейнера используйте команду  **```make docker_run```**.

   * Для подключения внешнего хранилища используйте команду **```make docker_run_db```**. При этом на момент запуска контейнера файл внешней базы данных должен существовать и должен быть доступным для записи.

## База данных

* база данных -  SQLite
* файл базы данных - scheduler.db
* файл для создания базы данных - scheduler_creator.sql

#### Поля базы данных:

* id - уникальный идентификатор задачи
* date - дата ближайшего выполнения задачи (обязательное поле)
* title - заголовок задачи (обязательное поле)
* comment - комментарий или описание задачи
* repeat - правило для повторения задачи

## Работа с приложением по API

<!--suppress HtmlDeprecatedAttribute -->
<table>1. Работа с отдельными задачами.
    <tr>
        <th width="70px">Тип</th>
        <th width="100px">Тело</th>
        <th width="200px">Содержание</th>
        <th width="400px">Описание</th>
    </tr>
    <tr>
        <td>POST</td>
        <td>/api/task</td>
        <td>date, title, comment, repeat</td>
        <td>Размешает в базе новую задачу, с учетом правила по её повторению.</td>
    </tr>
    <tr>
        <td>GET</td>
        <td>/api/task</td>
        <td>id</td>
        <td>Возвращает задачу в соответствии с переданным id.</td>
    </tr>
    <tr>
        <td>PUT</td>
        <td>/api/task</td>
        <td>id</td>
        <td>Вносит изменения в задачу, в соответствии с переданным id и параметрами.</td>
    </tr>
    <tr>
        <td>DELETE</td>
        <td>/api/task</td>
        <td>id</td>
        <td>Удаляет задачу, в соответствии с переданным id.</td>
    </tr>
</table>

```HTML
Пример 1: /api/task?id=<число>
```

* id - уникальный идентификатор задачи
* date - дата ближайшего выполнения задачи (обязательное поле)
* title - заголовок задачи (обязательное поле)
* comment - комментарий или описание задачи
* repeat - правило для повторения задачи

Примеры POST запросов в формате JSON:

```JSON
{
  "date": "20240201",
  "title": "Подвести итог",
  "comment": "Мой комментарий",
  "repeat": "d 5"
}
```
или
```JSON
{
  "date": "20240202",
  "title": "Сходить зоопарк",
  "repeat": ""
}
```
Пример успешного ответа в формате JSON:
```JSON
{"id":"186"}
```
Пример ответа при ошибке в формате JSON:
```JSON
{"error":"Не указан заголовок задачи"}
```
Примеры GET запросов в формате JSON:

```JSON
{
    "id": "185",
    "date": "20240201",
    "title": "Подвести итог",
    "comment": "",
    "repeat": "d 5"
}
```
Пример ответа при ошибке в формате JSON:

```JSON
{"error": "Задача не найдена"}
```
или
```JSON
{"error": "Не указан идентификатор"}
```

При успешном выполнении PUT и DELETE запростов возвращается пустой JSON {}

---

<table>2. Завершение отдельной задачи.
    <tr>
        <th width="70px">Тип</th>
        <th width="100px">Тело</th>
        <th width="200px">Содержание</th>
        <th width="400px">Описание</th>
    </tr>
    <tr>
        <td>POST</td>
        <td>/api/task/done</td>
        <td>id</td>
        <td>При наличии правила по повторению задачи изменяет срок задачи на новый.
            При отсутствии правила для повторения удаляет задачу.</td>
    </tr>

</table>

```HTML

Пример 1: /api/task/done?id=<идентификатор>

```

* id - уникальный идентификатор задачи

В случае успешного завершения возвращается пустой JSON {}

Пример ошибки:

```JSON
{"error": "текст ошибки"}
```

---

<table>3. Работа с группами задач.
    <tr>
        <th width="70px">Тип</th>
        <th width="100px">Тело</th>
        <th width="200px">Содержание</th>
        <th width="400px">Описание</th>
    </tr>
    <tr>
        <td >GET</td>
        <td>/api/tasks</td>
        <td>search</td>
        <td>Возвращает задачи, соответствующие переданным параметрам поиска или ошибку.</td>
    </tr>
</table>

```HTML

Пример 1: /api/tasks?search=бассейн

Пример 2: /api/tasks?search=08.02.2024

```
* search - строка или дата для поиска задач

Пример возвращаемого значения в формате JSON:

```JSON
{
    "tasks": [
    {
        "id": "171",
        "date": "20240131",
        "title": "Заголовок задачи",
        "comment": "",
        "repeat": ""
    },
    {
        "id": "176",
        "date": "20240131",
        "title": "Фитнес",
        "comment": "",
        "repeat": "d 3"
    },
    {
        "id": "185",
        "date": "20240201",
        "title": "Подвести итог",
        "comment": "Мой комментарий",
        "repeat": "d 5"
    },
    ]
}
```
Пример пустого ответа:

```JSON
{"tasks": []}
```

Пример ошибки:

```JSON
{"error": "текст ошибки"}
```

---

<table>4. Дополнительные функции.
    <tr>
        <th width="70px">Тип</th>
        <th width="100px">Тело</th>
        <th width="200px">Содержание</th>
        <th width="400px">Описание</th>
    </tr>
    <tr>
        <td >GET</td>
        <td>/api/nextdate</td>
        <td>now, date, repeat </td>
        <td>Возвращает следующую дату для выполнения задачи или ошибку.</td>
    </tr>
</table>

```HTML

Пример: /api/nextdate?now=<20060102>&date=<20060102>&repeat=<правило>
    
```    
    
* now - текущая дата
* date - дата начала отсчета
* repeat - правило повторения

Правила повторений:

1. Если правило не указано, отмеченная выполненной задача будет удаляться из таблицы после выполнения;
2. d <число> — задача переносится на указанное число дней. Максимально допустимое число равно 400. Примеры:
   * d 1 — каждый день;
   * d 7 — для вычисления следующей даты добавляем семь дней;
   * d 60 — переносим на 60 дней.
3. y — задача выполняется ежегодно. Этот параметр не требует дополнительных уточнений. При выполнении задачи дата перенесётся на год вперёд.
4. (не реализоввано!) w <через запятую от 1 до 7> — задача назначается в указанные дни недели, где 1 — понедельник, 7 — воскресенье. Например:
   * w 7 — задача перенесётся на ближайшее воскресенье;
   * w 1,4,5 — задача перенесётся на ближайший понедельник, четверг или пятницу;
   * w 2,3 — задача перенесётся на ближайший вторник или среду.
5. (не реализоввано!) m <через запятую от 1 до 31,-1,-2> [через запятую от 1 до 12] — задача назначается в указанные дни месяца. При этом вторая последовательность чисел опциональна и указывает на определённые месяцы. Например:
   * m 4 — задача назначается на четвёртое число каждого месяца;
   * m 1,15,25 — задача назначается на 1-е, 15-е и 25-е число каждого месяца;
   * m -1 — задача назначается на последний день месяца;
   * m -2 — задача назначается на предпоследний день месяца;
   * m 3 1,3,6 — задача назначается на 3-е число января, марта и июня;
   * m 1,-1 2,8 — задача назначается на 1-е и последнее число число февраля и авгуcта.

## Примечания:
* В коде присутствуют комментарии "//TODO" с описанрием дальнейших действий по дополнению и изменению кода.
* Дальнейшие изменения сохранят обратную совместимость.
