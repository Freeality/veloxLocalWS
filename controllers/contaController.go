package controllers

import (
	"github.com/freeality/veloxLocalWS/models"
	"sync"
	"fmt"
)

// Representa o controller
type ContaController struct {
	contas []*models.Conta
	Control
}

// contas para teste
func (c *ContaController)CriarContasParaTeste() {
	c.AdicionaRegistro("nome1", "email1", "senha1")
	c.AdicionaRegistro("nome2", "email2", "senha2")
	c.AdicionaRegistro("nome3", "email3", "senha3")
}

// precisamos disso para criar um singleton
var instance *ContaController // contém a instancia única
var once sync.Once // once executa uma função uma única vez

// retorna a instância singleton
func NewContaController() *ContaController {

	once.Do(func() {
		instance = &ContaController{
			make([]*models.Conta, 0),
			Control{new(sync.Mutex)},
		}
	})

	return instance
}

// Adiciona registros
func (c *ContaController) AdicionaRegistro(nome, email, senha string) int {

	// adquire uma tranca e nos certificamos que ela será liberada
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// o novo id é a quantidade de elementos em Contas
	novoId := len(c.contas)

	// cria um novo registro com os dados fornecidos e o novo id
	novoRegistro := &models.Conta{
		novoId,
		nome,
		email,
		senha,
	}

	// adiciona registro
	c.contas = append(c.contas, novoRegistro)
	return novoId
}

// RemoveRegistro remove o registro com o id correspondente. Retorna nil
// em caso de sucesso ou o erro específico.
func (c *ContaController) RemoveRegistro(id int) error {

	// adquire uma tranca e certifica-se que ela será liberada
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// verifica se temos um id válido
	if !idValido(id, c.contas) {
		return fmt.Errorf("id inválido")
	}

	// remove o registro
	RemoveID(id, &c.contas)
	fmt.Printf("id %d removido", id)

	return nil
}

// ObtemRegistro retorna a entrada identificada pelo id fornecido ou um
// erro em caso de falha
func (c *ContaController) ObtemRegistro(id int) (*models.Conta, error) {
	// verifica se o id é valido
	if !idValido(id, c.contas) {
		return nil, fmt.Errorf("id inválido")
	}

	return c.contas[id], nil
}

// ObtemRegistros retorna todos as instâncias registradas
func (c *ContaController) ObtemRegistros() []*models.Conta {

	return c.contas
}

// RemoveTodos remove todos os registros em Contas
func (c *ContaController) RemoveTodos() {
	// adquire a trava e certifica-se de soltá-la no final
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.contas = []*models.Conta{}
}

// Utilidades a partir daqui
// criei esse método aqui pois acho que ficou melhor.
// func RemoveID(slice []*GuestBookEntry, id int) {
func RemoveID(id int, contas *[]*models.Conta) {

	s := *contas
	s = append(s[:id], s[id+1:]...)
	*contas = s
}

// Testa um id
func idValido(id int, contas []*models.Conta) bool {

	if id < 0 || id >= len(contas) {
		return false
	}

	return true
}