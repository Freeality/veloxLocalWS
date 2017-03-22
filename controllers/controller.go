package controllers

import "sync"

// todo o controller deve possuir esses métodos
type Controller interface {
	RemoveRegistro(id int) error
	RemoveTodos()
}

// os controlles devem poder lidar com muitas requisições
// o mutex tranca a gorotine para evitar que outras
// requisições tentem alterar os dados simultaneamente provocando erros
type Control struct {
	mutex *sync.Mutex
}
