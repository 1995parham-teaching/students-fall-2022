<h1 align="center">Students</h1>
<h6 align="center">Based on a Fall-2022 Internet Engineering Course Project at Amirkabir University of Tech.</h6>

<p align="center">
  <img alt="GitHub Workflow Status" src="https://img.shields.io/github/actions/workflow/status/1995parham-teaching/students-fall-2022/test.yaml?logo=github&style=for-the-badge">
  <img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/1995parham-teaching/students-fall-2022?logo=go&style=for-the-badge">
</p>

## Introduction

Review on how we can write a web application with HTTP Framework named [Echo](https://echo.labstack.com/) and
ORM named [GORM](https://gorm.io/).
This application stores students and their courses into a SQLite database. There is a many-to-many
relationship between course and student which means each student can have multiple courses
and each course may be taken by multiple students.

I tried to use best practices that reduce the code complexity and increase maintainability.
Code structure is somewhat compatible with the famous [project-layout](https://github.com/golang-standards/project-layout).

There are two models named `Student` and `Course`. Models are used for in-application communication
and use request/responses for serializing models over HTTP and use store structures for serializing models
from/to database.
For each student, it generates the student ID randomly and then stores it.
There is no authentication over APIs and anybody can use CRUD over students and courses.
