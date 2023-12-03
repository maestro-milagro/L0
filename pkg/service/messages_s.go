package service

import (
	message "awesomeProject"
	"awesomeProject/pkg/cache"
	"awesomeProject/pkg/repository"
)

type MessageService struct {
	repo  repository.MessagesTB
	repoD repository.DeliveriesTB
	repoP repository.PaymentsTB
	repoI repository.ItemsTB
}

func NewMessageService(repo repository.MessagesTB, repoD repository.DeliveriesTB, repoP repository.PaymentsTB, repoI repository.ItemsTB) *MessageService {
	return &MessageService{
		repo:  repo,
		repoD: repoD,
		repoP: repoP,
		repoI: repoI,
	}
}

// Функция для сохранения переданного объекта в бд
func (s *MessageService) Create(message message.Message) (int, error) {
	//Получаем объект хранилища кэша
	c := cache.C
	//Сохраняем данные предназначенную для них в таблицу
	delId, err := s.repoD.Create(message.Delivery)
	if err != nil {
		return 0, err
	}
	payId, err := s.repoP.Create(message.Payment)
	if err != nil {
		return 0, err
	}
	itemsId, err := s.repoI.Create(message.Items)
	if err != nil {
		return 0, err
	}
	//Получив id объекта из каждой таблицы сохраняем их в таблицы связи,
	//а так же сохраняем оставшиеся данные в таблицу message
	id, err := s.repo.Create(message, delId, payId, itemsId)
	//Сохраняем id и данные в кэш
	c.Update(id, message)
	return id, err
}

// Функция для получения всех объектов из бд
func (s *MessageService) GetAll() ([]message.Message, error) {
	mess, err := s.repo.GetAll()
	pay, err := s.repoP.GetAll()
	if err != nil {
		return []message.Message{}, err
	}
	del, err := s.repoD.GetAll()
	if err != nil {
		return []message.Message{}, err
	}
	for i, v := range pay {
		mess[i].Payment = v
		mess[i].Delivery = del[i]
	}
	for i, _ := range mess {
		mess[i].Items, err = s.repoI.GetById(mess[i].MessageId)
		if err != nil {
			return []message.Message{}, err
		}
	}
	return mess, err
}

// Функция для получения объекта по id
func (s *MessageService) GetById(messageId int) (message.Message, error) {
	//Получаем объект хранилища кэша
	c := cache.C
	//Если объект с переданным id есть в кэше передаём данные из кэша
	//Если нет поочереди достаем данные из каждой таблицы, соединяем их сохраняем в кэш и возвращаем
	posAnsw, ok := c.Read(messageId)
	if ok {
		return posAnsw, nil
	}
	del, err := s.repoD.GetById(messageId)
	if err != nil {
		return message.Message{}, err
	}
	pay, err := s.repoP.GetById(messageId)
	if err != nil {
		return message.Message{}, err
	}
	it, err := s.repoI.GetById(messageId)
	if err != nil {
		return message.Message{}, err
	}
	mess, err := s.repo.GetById(messageId)
	if err != nil {
		return message.Message{}, err
	}
	mess.Payment = pay
	mess.Delivery = del
	mess.Items = it
	c.Update(messageId, mess)
	return mess, err
}

// Функция для удаления объекта по id
func (s *MessageService) Delete(messageId int) error {
	return s.repo.Delete(messageId)
}
