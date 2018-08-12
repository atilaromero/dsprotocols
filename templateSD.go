// -----------------------------------------
// Tentativa de um Template para o projeto de SD.
//
// Disciplina Topicos PPGCC - PUCRS - Escola Politecnica
// Atila, Dotti, Eliã, Tarcisio (em ordem alfabetica)
// -----------------------------------------

package main

import (
	"fmt"
	// "sync"
	//"time"
)

// -----------------------------------------
// Definição de Template -------------------
// -----------------------------------------
// Template aqui nao significa uma estrutura herdada, pois acho que Go nao tem heranca
// Mas sim um formato que seguiriamos convencionando a estrutura de todos os modulos
// para suportar as definicoes do livro que seguiremos
// este formato teria que ter mapeamento o mais direto possivel com as construcoes
// do livro

type Template struct {
	req, ind chan struct{} // canais oferecidos na api
	uT1      Template2     // assim como declaramos os modulos utilizados por este modulo
	// os canais oferecidos por este estao na sua definicao
	state int // teriamos ainda aqui tudo referente ao estado deste modulo
}

type Template2 struct {
	req, ind chan struct{} // canais oferecidos na api
}

func NewTemplate(x Template2) *Template {
	// used template deveria vir com referencia aos modulos que vai utilizar

	my := &Template{ // aqui inicializamos canais, referenciamos modulos utilizados, inicializamos estado
		req: make(chan struct{}),
		ind: make(chan struct{}),
		uT1: x, // assim como declaramos os modulos utilizados por este modulo
		// os canais oferecidos por este estao na sua definicao
		state: 0, // teriamos ainda aqui tudo referente ao estado deste modulo
	}

	templateEvents := func() { // aqui escrevemos o tratamento de eventos - funcao dentro de NewTemplate
		select {
		case <-my.req:

		case <-my.uT1.ind: // aqui pode mandar msg para my.ind:     my.ind <- msgind
			// e/ou mandar msg para uT1.req:           uT1.req <- msgreq
			// e/ou mudar estado interno.
			// noi caso de mudar estado interno teria que ver
			// se guardas internas seriam automaticamente reavaliadas pelas construcoes
			// da linguagem
			// case - enfim, aqui teriamos que ver se modelamos os casos necessarios de guardas
		}
	}

	go func() {
		for {
			templateEvents()
		} // codigo fixo indicando que fica para sempre tratando eventos
	}()
	return my
}

// -----------------------------------------
// Fim Definição de Template ---------------
// -----------------------------------------

// Usa template

func main() {
	fmt.Println("inicio")
	perfectLink := NewTemplate()
	perfectLink.req <- mensagemNoLink
}
