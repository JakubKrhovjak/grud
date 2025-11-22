# Student Management REST API

REST API pro správu studentů univerzity s PostgreSQL databází.

## Technologie

- Go 1.25
- PostgreSQL 16
- Docker & Docker Compose
- Gorilla Mux (HTTP router)

## Struktura projektu

```
grud/
├── main.go              # Hlavní aplikační soubor s REST API
├── go.mod               # Go module definice
├── go.sum               # Go dependencies checksums
├── Dockerfile           # Docker image pro API
├── docker-compose.yml   # Docker Compose konfigurace
└── README.md           # Dokumentace
```

## Student Model

```json
{
  "id": 1,
  "first_name": "Jan",
  "last_name": "Novák",
  "email": "jan.novak@university.cz",
  "major": "Computer Science",
  "year": 2
}
```

## API Endpointy

### Vytvořit studenta
```bash
POST /api/students
Content-Type: application/json

{
  "first_name": "Jan",
  "last_name": "Novák",
  "email": "jan.novak@university.cz",
  "major": "Computer Science",
  "year": 2
}
```

### Získat všechny studenty
```bash
GET /api/students
```

### Získat studenta podle ID
```bash
GET /api/students/{id}
```

### Aktualizovat studenta
```bash
PUT /api/students/{id}
Content-Type: application/json

{
  "first_name": "Jan",
  "last_name": "Novák",
  "email": "jan.novak@university.cz",
  "major": "Software Engineering",
  "year": 3
}
```

### Smazat studenta
```bash
DELETE /api/students/{id}
```

## Instalace a spuštění

### Předpoklady
- Docker
- Docker Compose

### Spuštění s Docker Compose

1. Naklonujte nebo stáhněte projekt
2. Přejděte do projektového adresáře
3. Spusťte Docker Compose:

```bash
docker-compose up -d
```

Tento příkaz:
- Stáhne PostgreSQL 16 Alpine image
- Vytvoří databázi `university` na portu 5439
- Sestaví a spustí Go API na portu 8080
- Automaticky vytvoří tabulku `students`

### Kontrola běžících služeb

```bash
docker-compose ps
```

### Zobrazení logů

```bash
# Všechny služby
docker-compose logs -f

# Pouze API
docker-compose logs -f api

# Pouze databáze
docker-compose logs -f postgres
```

### Zastavení služeb

```bash
docker-compose down
```

### Zastavení a smazání dat

```bash
docker-compose down -v
```

## Lokální vývoj (bez Dockeru)

### Předpoklady
- Go 1.25+
- PostgreSQL

### Instalace závislostí

```bash
go mod init grud
go get github.com/gorilla/mux
go get github.com/lib/pq
```

### Nastavení proměnných prostředí

```bash
export DB_HOST=localhost
export DB_PORT=5439
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=university
export PORT=8080
```

### Spuštění PostgreSQL

```bash
docker run --name postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=university -p 5439:5432 -d postgres:16-alpine
```

### Spuštění aplikace

```bash
go run main.go
```

## Testování API

### Pomocí curl

#### Vytvořit studenta
```bash
curl -X POST http://localhost:8080/api/students \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Jan",
    "last_name": "Novák",
    "email": "jan.novak@university.cz",
    "major": "Computer Science",
    "year": 2
  }'
```

#### Získat všechny studenty
```bash
curl http://localhost:8080/api/students
```

#### Získat studenta podle ID
```bash
curl http://localhost:8080/api/students/1
```

#### Aktualizovat studenta
```bash
curl -X PUT http://localhost:8080/api/students/1 \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Jan",
    "last_name": "Novák",
    "email": "jan.novak@university.cz",
    "major": "Software Engineering",
    "year": 3
  }'
```

#### Smazat studenta
```bash
curl -X DELETE http://localhost:8080/api/students/1
```

### Pomocí HTTPie

```bash
# Vytvořit studenta
http POST localhost:8080/api/students first_name=Jan last_name=Novák email=jan.novak@university.cz major="Computer Science" year:=2

# Získat všechny studenty
http localhost:8080/api/students

# Získat studenta
http localhost:8080/api/students/1

# Aktualizovat studenta
http PUT localhost:8080/api/students/1 first_name=Jan last_name=Novák email=jan.novak@university.cz major="Software Engineering" year:=3

# Smazat studenta
http DELETE localhost:8080/api/students/1
```

## Konfigurace

Aplikace používá proměnné prostředí pro konfiguraci:

| Proměnná | Výchozí hodnota | Popis |
|----------|----------------|-------|
| DB_HOST | localhost | Hostname PostgreSQL |
| DB_PORT | 5439 | Port PostgreSQL |
| DB_USER | postgres | Databázový uživatel |
| DB_PASSWORD | postgres | Databázové heslo |
| DB_NAME | university | Název databáze |
| PORT | 8080 | Port API serveru |

## Databázová struktura

### Tabulka: students

| Sloupec | Typ | Popis |
|---------|-----|-------|
| id | SERIAL | Primární klíč (auto-increment) |
| first_name | VARCHAR(100) | Křestní jméno |
| last_name | VARCHAR(100) | Příjmení |
| email | VARCHAR(255) | Email (unikátní) |
| major | VARCHAR(100) | Studijní obor |
| year | INTEGER | Ročník |

## Troubleshooting

### Problém s připojením k databázi

Zkontrolujte, že PostgreSQL běží:
```bash
docker-compose logs postgres
```

### Port již používán

Pokud je port 5439 nebo 8080 již používán, změňte port v `docker-compose.yml`:
```yaml
ports:
  - "5440:5432"  # Pro PostgreSQL
  - "8081:8080"  # Pro API
```

### Rebuild Docker image

Pokud jste změnili kód:
```bash
docker-compose up -d --build
```

## Přímý přístup k databázi

```bash
docker exec -it university_db psql -U postgres -d university
```

SQL příkazy:
```sql
-- Zobrazit všechny studenty
SELECT * FROM students;

-- Vytvořit studenta
INSERT INTO students (first_name, last_name, email, major, year)
VALUES ('Jan', 'Novák', 'jan@example.com', 'CS', 2);

-- Smazat všechny studenty
TRUNCATE TABLE students RESTART IDENTITY;
```
