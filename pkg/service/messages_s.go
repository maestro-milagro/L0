package service

import (
	message "awesomeProject"
	"awesomeProject/pkg/repository"
)

type MessageService struct {
	repo repository.MessagesTB
}

func NewMessageService(repo repository.MessagesTB) *MessageService {
	return &MessageService{repo: repo}
}
func (s *MessageService) Create(message message.Message) (int, error) {
	return s.repo.Create(message)
}
func (s *MessageService) GetAll() ([]message.Message, error) {
	return s.repo.GetAll()
}
func (s *MessageService) GetById(messageId int) (message.Message, error) {
	return s.repo.GetById(messageId)
}
func (s *MessageService) Delete(messageId int) error {
	return s.repo.Delete(messageId)
}
