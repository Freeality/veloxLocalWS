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


// implementa GetPath
func (c *ContaController) WebDelete(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r) // obtém os parâmetros de r
	status, bodyOut := c.webDelete(params) // obtém a resposta adequada

	w.WriteHeader(status) // atribui o status ao cabeçalho de resposta
	fmt.Fprintln(w, bodyOut) // responde à solicitação imprimindo bodyOut em w
}

// webDelete recebe param e retorna http.Status e a mensagem para o WebDelete
func (c *ContaController) webDelete(params map[string]string) (int, string) {
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

// WebPost implementa o POST
func (c *ContaController) WebPost(w http.ResponseWriter, r *http.Request) {

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
	status, bodyOut := c.webPost(params, bodyIn) // obtém a resposta adequada
	w.WriteHeader(status) // atribui o status ao cabeçalho de resposta
	fmt.Fprintln(w, bodyOut) // responde à solicitação imprimindo bodyOut em w
}

// implementa o POST
func (c *ContaController) webPost(params map[string]string, body []byte) (int, string) {

	if len(params) != 0 {
		// Sem chaves. Não há suporte.
		return http.StatusMethodNotAllowed, "Método não suportado"
	}

	// desempacotando o registro enviado pelo usuário
	var conta models.Conta
	err := json.Unmarshal(body, &conta)

	if err != nil {
		// não pode desempacotar registro
		return http.StatusBadRequest, "dados inválidos no JSON"
	}

	// se chegou até aqui, tudo certo!
	// acrescenta o registro solicitado pelo usuário
	c.AdicionaRegistro(conta.Email, conta.Senha)

	// agora retorna dizendo que está tudo ok
	return http.StatusOK, "novo registro foi criado"
}

// implementa o GET
func (c *ContaController) WebGet(params map[string]string) (int, string) {

	// Se não houver parâmetros retorna todos os registros
	if len(params) == 0 {
		status, registros := c.obtemRegistrosEStatus()
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

// obtem registros em JSON e http.Status
func (c *ContaController) obtemRegistrosEStatus() (int, string) {

	encodedEntries, err := json.Marshal(c.ObtemRegistros())
	if err != nil {
		// falha codificando os registros
		return http.StatusInternalServerError, "Erro convertendo dados"
	}

	return http.StatusOK, string(encodedEntries)
}