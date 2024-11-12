curl -v -X POST -d "username=user&password=user123" http://localhost:8000/login

curl -v --cookie "token=ТОКЕН" http://localhost:8000/protected


попытка входа под юзером
curl -v --cookie "token=ТОКЕН" http://localhost:8000/admin


теперь под админом
curl -v -X POST -d "username=admin&password=admin123" http://localhost:8000/login

curl -v --cookie "token=ТОКЕН" http://localhost:8000/admin
