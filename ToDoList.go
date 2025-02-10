package main

import (
	"fmt"
	"time"
)

type Person struct {
	name    string
	age     int
	address string
}

type Employee struct {
	name     string
	position string
	salary   float64
	bonus    float64
}

type Student struct {
	name            string
	solvedProblems  int
	scoreForOneTask float64
	passingScore    float64
}

type Task struct {
	summary     string
	description string
	deadline    time.Time
	priority    int
}

type Note struct {
	title string
	text  string
}

type ToDoList struct {
	name  string
	tasks []Task
	notes []Note
}

func (p Person) Print() {
	fmt.Printf("Name:%s\n\nAge:%d\n\nAdress:%s", p.name, p.age, p.address)
}

func (e Employee) CalculateTotalSalary() {
	fmt.Printf("Employee:%s\n\nPosition:%s\n\nTotal Salary:%.2f\n\n", e.name, e.position, e.salary+e.bonus)
}

func (s Student) IsExcellentStudent() bool {
	return s.scoreForOneTask*float64(s.solvedProblems) >= s.passingScore
}

func (t Task) IsOverdue() bool {
	return time.Now().After(t.deadline)
}

func (t Task) IsTopPriority() bool {
	return t.priority >= 4
}

func (todo ToDoList) TasksCount() int {
	return len(todo.tasks)
}

func (todo ToDoList) NotesCount() int {
	return len(todo.notes)
}

func (todo ToDoList) CountTopPrioritiesTasks() int {
	c := 0
	for _, item := range todo.tasks {
		if item.IsTopPriority() {
			c ++
		}
	}
	return c
}

func (todo ToDoList) CountOverdueTasks() int {
	c := 0
	for _, item := range todo.tasks {
		if item.IsOverdue() {
			c ++
		}
	}
	return c
}
