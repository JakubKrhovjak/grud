# University Management Microservices

Mikroservisní architektura pro správu univerzity s PostgreSQL databázemi a Bun ORM.

## Architektura

Projekt se skládá ze dvou samostatných mikroservisů:

1. **student-service** - správa studentů
2. **project-service** - správa projektů

Každý mikroservis má:
- Vlastní PostgreSQL databázi
- Vlastní API server
- Nezávislý deployment
- Domain-Driven Design (DDD) strukturu

## Technologie

- Go 1.25
- PostgreSQL 16
- Bun ORM
- Gorilla Mux (HTTP router)
- Zap/Slog Logger (strukturované logování)
- Docker & Docker Compose

## Služby

### Student Service
- Port: `8080`
- Databáze: `university` (port `5439`)
- Endpointy: `/api/students`

### Project Service
- Port: `8081`
- Databáze: `projects` (port `5440`)
- Endpointy: `/api/projects`

## Struktura projektu

```
grud/
├── student-service/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── db/
│   │   ├── logger/
│   │   ├── student/
│   │   │   ├── model.go
│   │   │   ├── repository.go
│   │   │   ├── service.go
│   │   │   └── http.go
│   │   └── app/
│   ├── go.mod
│   └── Dockerfile
│
├── project-service/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── db/
│   │   ├── logger/
│   │   ├── project/
│   │   │   ├── model.go
│   │   │   ├── repository.go
│   │   │   ├── service.go
│   │   │   └── http.go
│   │   └── app/
│   ├── go.mod
│   └── Dockerfile
│
├── docker-compose.yml
├── go.work
└── README.md
```

## Instalace a spuštění

### Předpoklady
- Docker
- Docker Compose

### Spuštění všech služeb

```bash
docker-compose up -d
```

Tento příkaz spustí:
- PostgreSQL databázi pro studenty (port 5439)
- PostgreSQL databázi pro projekty (port 5440)
- Student API (port 8080)
- Project API (port 8081)

### Kontrola běžících služeb

```bash
docker-compose ps
```

### Zobrazení logů

```bash
# Všechny služby
docker-compose logs -f

# Pouze student-service
docker-compose logs -f student_api

# Pouze project-service
docker-compose logs -f project_api

# Databáze
docker-compose logs -f postgres
docker-compose logs -f postgres_projects
```

### Zastavení služeb

```bash
docker-compose down
```

### Zastavení a smazání dat

```bash
docker-compose down -v
```

## Student Service API

### Student Model
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

### Endpointy

#### Vytvořit studenta
```bash
POST http://localhost:8080/api/students
Content-Type: application/json

{
  "first_name": "Jan",
  "last_name": "Novák",
  "email": "jan.novak@university.cz",
  "major": "Computer Science",
  "year": 2
}
```

#### Získat všechny studenty
```bash
GET http://localhost:8080/api/students
```

#### Získat studenta podle ID
```bash
GET http://localhost:8080/api/students/{id}
```

#### Aktualizovat studenta
```bash
PUT http://localhost:8080/api/students/{id}
Content-Type: application/json

{
  "first_name": "Jan",
  "last_name": "Novák",
  "email": "jan.novak@university.cz",
  "major": "Software Engineering",
  "year": 3
}
```

#### Smazat studenta
```bash
DELETE http://localhost:8080/api/students/{id}
```

## Project Service API

### Project Model
```json
{
  "id": 1,
  "name": "Web Application",
  "description": "Modern web app with Go backend",
  "status": "in_progress",
  "start_date": "2024-01-15T00:00:00Z",
  "end_date": null
}
```

### Endpointy

#### Vytvořit projekt
```bash
POST http://localhost:8081/api/projects
Content-Type: application/json

{
  "name": "Web Application",
  "description": "Modern web app with Go backend",
  "status": "in_progress",
  "start_date": "2024-01-15T00:00:00Z"
}
```

#### Získat všechny projekty
```bash
GET http://localhost:8081/api/projects
```

#### Získat projekt podle ID
```bash
GET http://localhost:8081/api/projects/{id}
```

#### Aktualizovat projekt
```bash
PUT http://localhost:8081/api/projects/{id}
Content-Type: application/json

{
  "name": "Web Application",
  "description": "Updated description",
  "status": "completed",
  "start_date": "2024-01-15T00:00:00Z",
  "end_date": "2024-03-20T00:00:00Z"
}
```

#### Smazat projekt
```bash
DELETE http://localhost:8081/api/projects/{id}
```

## Lokální vývoj (bez Dockeru)

### Student Service

```bash
cd student-service
go mod download
export DB_HOST=localhost
export DB_PORT=5439
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=university
export PORT=8080
go run cmd/server/main.go
```

### Project Service

```bash
cd project-service
go mod download
export DB_HOST=localhost
export DB_PORT=5440
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=projects
export PORT=8081
go run cmd/server/main.go
```

## Go Workspace

Projekt používá Go workspace pro práci s více moduly:

```bash
# Aktualizace workspace
go work sync

# Sestavení všech služeb
go work use ./student-service ./project-service
```

## Výhody mikroservisní architektury

- **Nezávislý deployment** - každá služba může být nasazena samostatně
- **Škálovatelnost** - služby lze škálovat nezávisle dle potřeby
- **Technologická svoboda** - každá služba může použít jiné technologie
- **Izolace chyb** - selhání jedné služby neovlivní ostatní
- **Týmová autonomie** - různé týmy mohou pracovat na různých službách
- **Database per Service** - každá služba má vlastní databázi

## Přímý přístup k databázím

### Student Database
```bash
docker exec -it university_db psql -U postgres -d university
```

### Project Database
```bash
docker exec -it projects_db psql -U postgres -d projects
```

## Domain-Driven Design

Každý mikroservis používá DDD strukturu:

### Model Layer
- Definice entit
- Bun tagy pro ORM mapping
- JSON tagy pro API response

### Repository Layer
- Interface pro DB operace
- CRUD metody s Bun ORM
- Vrací Go errors

### Service Layer
- Business logika
- Validace vstupů
- Error handling
- Strukturované logování

### HTTP Layer
- REST handlers
- Request/Response mapping
- HTTP status codes

## Best Practices

1. **Dependency Injection** - dependencies jsou injectované přes konstruktory
2. **Interface segregation** - každá vrstva definuje své interface
3. **Error wrapping** - použití `fmt.Errorf` s `%w`
4. **Context propagation** - context.Context v každé metodě
5. **Validation** - validace na service vrstvě
6. **Structured logging** - strukturované logování pro monitoring
7. **Graceful shutdown** - bezpečné ukončení služeb

## Troubleshooting

### Problém s připojením k databázi

```bash
docker-compose logs postgres
docker-compose logs postgres_projects
```

### Port již používán

Změň porty v `docker-compose.yml`

### Rebuild Docker images

```bash
docker-compose up -d --build
```

### Reset databází

```bash
docker-compose down -v
docker-compose up -d
```
