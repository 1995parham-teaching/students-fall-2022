<h1 align="center">Students</h1>
<h6 align="center">Based on a Fall-2022 Internet Engineering Course Project at Amirkabir University of Tech.</h6>

<p align="center">
  <img alt="GitHub Workflow Status" src="https://img.shields.io/github/actions/workflow/status/1995parham-teaching/students-fall-2022/test.yaml?logo=github&style=for-the-badge">
  <img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/1995parham-teaching/students-fall-2022?logo=go&style=for-the-badge">
</p>

## Introduction

This review discusses how to write a web application using the [Echo](https://echo.labstack.com/) HTTP Framework and [GORM](https://gorm.io/) ORM.
The application is designed to store information about students and their courses in a SQLite database.
The relationship between the courses and students is many-to-many,
meaning that students can take multiple courses, and each course can have multiple students.

To ensure code simplicity and maintainability, best practices were used. The code structure is compatible with the popular
[project-layout](https://github.com/golang-standards/project-layout).

The application uses two models, `Student` and `Course`, for in-application communication.
The models use request/responses to serialize data over HTTP and store structures to serialize data from/to the database.
To generate a student ID, a random number is assigned to each student.
There is no authentication over the APIs, and anyone can use CRUD over students and courses.

GraphQL can improve the structure of your APIs, in case of having lots of data using it can reduce the duplicate codes.
Here, I am going to implement it using [99designs/gqlgen](https://github.com/99designs/gqlgen).

## SQLite is not enough?

However, using SQLite has its limitations.
GORM cannot easily switch to PostgreSQL,
and implementing this change would require a complete structure overhaul.
Changing the connection is not enough, and running the migration on store creation is not recommended.

## Up and Running (GraphQL)

You can open [GraphiQL](http://127.0.0.1:1373/v2/graphiql) in your browser and then sending
your queries and mutations.

Following mutation creates a student:

```graphql
mutation {
  createStudent(name: "Elahe Dastan") {
    name
  }
}
```

Then you retrieve them by name:

```graphql
query {
  studentsByName(name: "Elahe Dastan") {
    name
    id
  }
}
```

```json
{
  "data": {
    "studentsByName": [
      {
        "name": "Elahe Dastan",
        "id": "10368677"
      },
      {
        "name": "Elahe Dastan",
        "id": "34092594"
      }
    ]
  }
}
```

## Up and Running (HTTP)

Build and run the students' server:

```bash
go build
./students
```

Student creation request:

```bash
curl 127.0.0.1:1373/v1/students -X POST -H 'Content-Type: application/json' -d '{ "name": "Parham Alvani" }'
```

```json
{ "name": "Parham Alvani", "id": "89846857", "courses": null }
```

Student list request:

```bash
curl 127.0.0.1:1373/v1/students
```

```json
[{ "name": "Parham Alvani", "id": "89846857", "courses": [] }]
```

Course creation request:

```bash
curl 127.0.0.1:1373/v1/courses -X POST -H 'Content-Type: application/json' -d '{ "name": "Internet Engineering" }'
```

```json
{ "Name": "Internet Engineering", "ID": "00000007" }
```

Register student into a course:

```bash
curl 127.0.0.1:1373/v1/students/89846857/register/00000007
```

```json
null
```

And then we have the course into the student course list:

```bash
curl 127.0.0.1:1373/v1/students/89846857
```

```json
{
  "name": "Parham Alvani",
  "id": "89846857",
  "courses": [{ "Name": "Internet Engineering", "ID": "00000007" }]
}
```

Then we can even add new course and register our student into that course too:

```bash
curl 127.0.0.1:1373/v1/courses -X POST -H 'Content-Type: application/json' -d '{ "name": "C Programming" }'
```

```json
{ "Name": "C Programming", "ID": "00000000" }
```

```bash
curl 127.0.0.1:1373/v1/students/89846857/register/00000000
```

```json
null
```

```bash
curl 127.0.0.1:1373/v1/students/89846857
```

```json
{
  "name": "Parham Alvani",
  "id": "89846857",
  "courses": [
    { "Name": "C Programming", "ID": "00000000" },
    { "Name": "Internet Engineering", "ID": "00000007" }
  ]
}
```

## Preload

When you have a relation in your database, you can use `gorm.Preload` to fetch the related information within your
query, under the hood GORM run `N+1` queries as follows:

```text
2024/09/22 04:09:21 /Users/parham/Documents/Git/parham/1995parham-teaching/students-fall-2022/internal/store/student/sql.go:102
[0.081ms] [rows:0] SELECT * FROM `students_courses` WHERE `students_courses`.`sql_item_id` = "27849651"

2024/09/22 04:09:21 /Users/parham/Documents/Git/parham/1995parham-teaching/students-fall-2022/internal/store/student/sql.go:102
[0.691ms] [rows:1] SELECT * FROM `students` WHERE `students`.`id` = "27849651" LIMIT 1
```
