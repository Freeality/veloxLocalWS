package controllers

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"strconv"
	"io/ioutil"
	"encoding/json"
	"github.com/freeality/veloxLocalWS/models"
)

// Obtem parametros de http.Request
// Obtem a resposta de delete com os parametros
// Preenche o cabeçalho e o body com a resposta em http.ResponseWrite
func (c *ContaController) processaRequestDelete(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r) // obtém os parâmetros de r
	status, bodyOut := c.apagaRegistroComIDEmMap(params) // obtém a resposta adequada

	w.WriteHeader(status) // atribui o status ao cabeçalho de resposta
	fmt.Fprintln(w, bodyOut) // responde à solicitação imprimindo bodyOut em w
}

// Obtem o body de http.Request
// Se estiver OK, obtem parametros
// processa as informações julgando sua validade
// Preenche o cabeçalho e o body com a resposta em http.ResponseWrite
func (c *ContaController) processaRequestPost(w http.ResponseWriter, r *http.Request) {

	// Tentar certeza que o Body está fechado quando terminado
	defer r.Body.Close()

	// Obtém o body de r. ioutil.ReadAll retorna []byte, error
	bodyIn, err := ioutil.ReadAll(r.Body)

	// Problemas com o bodyIn, para por aqui
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // atribui o status ao cabeçalho de resposta
		fmt.Fprintln(w, []byte("Erro interno")) // responde à solicitação imprimindo bodyOut em w
		return
	}

	params := mux.Vars(r) // obtém os parâmetros

	// Quando houver um id válido em param o método fará uma atualização de um
	// registro existente. Caso não id, mas um body válido o registro será
	// adicionado.
	if len(params) >= 1 {
		// Alteração de um ou mais registros: Ainda não há suporte.
		return http.StatusMethodNotAllowed, "Método não suportado"
	}

	status, bodyOut := c.adicionaUmRegistroEmBody(bodyIn) // obtém a resposta adequada
	w.WriteHeader(status) // atribui o status ao cabeçalho de resposta
	fmt.Fprintln(w, bodyOut) // responde à solicitação imprimindo bodyOut em w
}

// Obtem os parametros de http.Request
// Processa as informações julgando sua validade
// Preenche o cabeçalho e o body com a resposta em http.ResponseWrite
func (c *ContaController) processaRequestGet(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req) // obtém os parâmetros de r
	status, bodyOut := c.obtemRegistroComIDEmMap(params) // obtém a resposta adequada

	res.WriteHeader(status) // atribui o status ao cabeçalho de resposta
	fmt.Fprintln(res, bodyOut) // responde à solicitação imprimindo bodyOut em w
}

// Busca um id válido em params e retorna
// O código de status do servidor e a informação string
// Caso o parâmetro seja um id válido retornará o registro
// em JSON e o código de Status OK
func (c *ContaController) obtemRegistroComIDEmMap(params map[string]string) (int, string) {

	// Se não houver parâmetros retorna todos os registros
	if len(params) == 0 {
		status, registros := c.obtemRegistrosJSON()
		return status, registros
	}

	// Se há parametros continua aqui
	// Converte id para inteiro
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		// Id não é um número
		return http.StatusBadRequest, "id inválido"
	}

	// O id é um inteiro válido. Tente obter o registro
	registro, err := c.ObtemRegistro(id)
	if err != nil {
		// Registro não encontrado
		return http.StatusNotFound, "registro não encontrado"
	}

	// O registro existe. Converta-o para JSON
	registroCodificado, err := json.Marshal(registro)
	if err != nil {
		// o registro não pode ser convertido
		return http.StatusInternalServerError, "Não pude converter p/JSON"
	}

	// Registro encontrado e convertido para JSON com sucesso
	return http.StatusOK, string(registroCodificado)
}

// Busca um id válido em params
// Apaga o registro, caso o id seja válido
// Retorna o código de Status do servidor OK
// e a mensagem de sucesso. O método julga as informações
// e retorna a mensagem adequada para cada outras situações.
func (c *ContaController) apagaRegistroComIDEmMap(params map[string]string) (int, string) {
	// Sem parâmetros. Bad request. entrada inválido
	if len(params) == 0 {
		return http.StatusBadRequest, "entrada inválida"
	}

	// id == 'ALL'. Deleta tudo!
	if params["id"] == "ALL" {
		c.RemoveTodos()
		return http.StatusOK, "lista apagada"
	}

	// O id deve ser um inteiro
	id, err := strconv.Atoi(params["id"])

	// não é um número
	if err != nil {
		return http.StatusBadRequest, "id inválido"
	}

	// se o número é valido
	err = c.RemoveRegistro(id)

	// não encontrou o registro
	if err != nil {
		return http.StatusNotFound, "registro não encontrado"
	}

	// se chegou aqui, está tudo certo!
	return http.StatusOK, "registro apagado"
}

// Adiciona um registro novo
// Caso o body seja uma conta em JSON válida, acrescenta o registro,
// retorna status ok e a mensagem de sucesso
func (c *ContaController) adicionaUmRegistroEmBody(body []byte) (int, string) {

	// desempacotando o registro enviado pelo usuário
	var conta models.Conta
	err := json.Unmarshal(body, &conta)

	if err != nil {
		// não pode desempacotar registro
		return http.StatusBadRequest, "dados inválidos no JSON"
	}

	// se chegou até aqui, tudo certo!
	// acrescenta o registro solicitado pelo usuário
	c.AdicionaRegistro(conta.Nome, conta.Email, conta.Senha)

	// agora retorna dizendo que está tudo ok
	return http.StatusOK, "novo registro foi criado"
}

// obtem registros em JSON e http.Status
// utiliza json.Marshal para empacotar os registros,
// retorna status OK e os dados se estiver tudo certo.
func (c *ContaController) obtemRegistrosJSON() (int, string) {

	encodedEntries, err := json.Marshal(c.ObtemRegistros())
	if err != nil {
		// falha codificando os registros
		return http.StatusInternalServerError, "Erro convertendo dados"
	}

	return http.StatusOK, string(encodedEntries)
}