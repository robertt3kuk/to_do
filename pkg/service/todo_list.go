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
func (s *ListService) GetAll(userId int) ([]todo.TodoList, error) {
	return s.repo.GetAll(userId)
}
func (s *ListService) GetById(userId, listId int) (todo.TodoList, error) {
	return s.repo.GetById(userId, listId)
}
func (s *ListService) Delete(userId, listId int) error {
	return s.repo.Delete(userId, listId)
}
func (s *ListService) Update(userId, listId int, input todo.UpdateListInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(userId, listId, input)
}
