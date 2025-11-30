## Запустить проект
up:
	docker-compose up --build

## Остановить проект
down:
	docker-compose down
	
## + удалить volume с БД (полная очистка)
drop:
	docker-compose down -v

## Запустить тесты на производительность (k6)
rps:
	k6 run test.js
