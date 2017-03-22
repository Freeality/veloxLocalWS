package models

import (
	"time"
	"fmt"
)

type Licenca struct {
	Id int
	Validade JSONTime
	Conta Conta
}

type JSONTime time.Time

func (t JSONTime)MarshalJSON() ([]byte, error) {

	dia := time.Time(t).Day()
	mes := time.Time(t).Month()
	ano := time.Time(t).Year()

	saida := fmt.Sprintf("%v-%d-%v", dia, mes, ano)
	return []byte(saida), nil
}