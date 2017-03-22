package models

import (
	"testing"
	"time"
	"fmt"
)

// Testa o formato do JSON
func TestJSONTime_MarshalJSON(t *testing.T) {
	// data de hoje
	data := JSONTime(time.Date(2017, time.March, 16, 0,0,0,0, time.Local))
	conta := Conta{Id:1, Email:"email", Senha:"1234"}
	dataString := "16-3-2017"

	L := Licenca{Id:1, Conta:conta, Validade:data}
	JSONdata, err := L.Validade.MarshalJSON()

	if err != nil {
		t.Errorf("não pode desempacotar: %v", err)
		return
	}

	fmt.Printf("A Data da licenca %s", JSONdata)

	if fmt.Sprintf("%s", JSONdata) == dataString {
		t.Log("\nTudo certo!\n")
		return
	}

	t.Errorf("A data não ficou certa %v", JSONdata)
}