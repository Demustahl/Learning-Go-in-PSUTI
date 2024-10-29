Создание mod файла и подключение зависимостей:
go mod init lab8
go get github.com/gorilla/mux
go mod tidy

Установка MongoDB и подключение:
1. установка
yay -S mongodb-bin
2. запуск службы
sudo systemctl start mongodb
3. включение
sudo systemctl enable mongodb
4. проверка удачного запуска
systemctl status mongodb
5. подключение к бд с помощью стандартного клиентаю Теперь можно работать
mongosh

2. Работа с бд:
добавить юзера
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name": "Maksim", "age": 20}'
получить всех юзеров
curl -X GET http://localhost:8080/users
получить по id
curl -X GET http://localhost:8080/users/64b1f2c5e3b3c4d1f0a6e8b9
изменение юзера
curl -X PUT http://localhost:8080/users/64b1f2c5e3b3c4d1f0a6e8b9 -H "Content-Type: application/json" -d '{"name": "Alice Smith", "age": 31}'
удаление юзера
curl -X DELETE http://localhost:8080/users/67211173c42e97bdae2ceaa7

Демонстрация в mongosh:
запуск mongosh
mongosh
заход в нашу бд
use dbLab8
вывести всех юзеров
db.users.find()

3. Обработка ошибок
ошибка "пустое имя"
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name": "", "age": 25}'
ошибка "отрицательный возраст"
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name": "Ameba", "age": -5}'
ошибка "неверный формат id"
curl -X PUT http://localhost:8080/users/chiki-briki -H "Content-Type: application/json" -d '{"name": "Galadriel", "age": 200}'
ошибка "пустое имя для замены"
curl -X PUT http://localhost:8080/users/64b1f2c5e3b3c4d1f0a6e8b9 -H "Content-Type: application/json" -d '{"name": "", "age": 5}'
ошибка "удаление несуществующего юзера"
curl -X DELETE http://localhost:8080/users/64b1f2c5e3b3c4d1f0a6e8aa

4. Пагинация и фильтрация
куууууууча юзеров
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name": "Elven Starfire", "age": 105}'
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name": "Grumpy Dwarf", "age": 45}'
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name": "Shadow Whisperer", "age": 30}'
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name": "Mystic Phoenix", "age": 29}'
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name": "Aurora Flame", "age": 25}'
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name": "Crystal Moon", "age": 21}'
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name": "Thunder Beast", "age": 40}'
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name": "Silent Breeze", "age": 33}'
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name": "Galactic Rover", "age": 19}'
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name": "Draco the Brave", "age": 27}'
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name": "Celestial Star", "age": 22}'
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name": "Sapphire Dragon", "age": 60}'
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name": "Forest Keeper", "age": 31}'
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name": "Ocean Wave", "age": 28}'
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name": "Fiery Hawk", "age": 26}'
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name": "Iron Heart", "age": 37}'
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name": "Golden Saber", "age": 23}'
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name": "Storm Chaser", "age": 35}'
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name": "Twilight Walker", "age": 34}'
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name": "Silver Fox", "age": 50}'
1 страница с 5 юзерами
curl -X GET "http://localhost:8080/users?page=1&pageSize=5"
фильтрация по имени
curl -X GET "http://localhost:8080/users?name=st"
фильтрация по возрасту
curl -X GET "http://localhost:8080/users?age=21"
фильтрация по имени, 2 страница 3 юзера
curl -X GET "http://localhost:8080/users?page=2&pageSize=3&name=al"




ДОБАВИТЬ
1. Обработку удачного удаления: типо сообщение о том, что такой юзер был удален. Сейчас просто ничего не выводится
2. Как происходит проверка на правильный формат ID? Какой формат считается правильным в программе?
3. 