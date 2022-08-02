package service

import (
	"todo"
	"todo/pkg/repository"
)

type ListService struct {
	repo repository.TodoList
}

func NewListService(repo repository.TodoList) *ListService {

	return &ListService{repo: repo}
}

func (s *ListService) Create(userId int, list todo.TodoList) (int, error) {
	return s.repo.Create(userId, list)
}
