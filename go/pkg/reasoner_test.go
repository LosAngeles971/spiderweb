package pkg

import (
	"fmt"
	"testing"
	"it/losangeles971/joshua/internal/io"
	ctx "it/losangeles971/joshua/internal/context"
	kkk "it/losangeles971/joshua/internal/knowledge"
)

func TestLogicReasoning(t *testing.T) {
	k, err := io.Load("../../resources/k_contadino.yml")
	if err != nil {
		fmt.Println("Knowledge not loaded due to error ", err)
		t.FailNow()
	}
	init := ctx.Create()
	init.Add(&ctx.Variable{Name: "FarmerA", Value: 1.0, Defined: true, })
	init.Add(&ctx.Variable{Name: "FarmerB", Value: 0.0, Defined: true, })
	init.Add(&ctx.Variable{Name: "WolfA", Value: 1.0, Defined: true, })
	init.Add(&ctx.Variable{Name: "WolfB", Value: 0.0, Defined: true, })
	init.Add(&ctx.Variable{Name: "GoatA", Value: 1.0, Defined: true, })
	init.Add(&ctx.Variable{Name: "GoatB", Value: 0.0, Defined: true, })
	init.Add(&ctx.Variable{Name: "CabbageA", Value: 1.0, Defined: true, })
	init.Add(&ctx.Variable{Name: "CabbageB", Value: 0.0, Defined: true, })
	e_name := "The farmer, the wolf, the goat and the cabbage are on the bank B of the river"
	success, ok := k.GetEvent(e_name)
	if !ok {
		fmt.Println("Knowledge does not contain event: ", e_name)
		t.FailNow()
	}
	outcome, queue, err := Verify(k, init, success, 50)
	fmt.Println(outcome)
	fmt.Println("Cycles: ", queue.GetCycles())
	fmt.Println(err)
	queue.SortByCycle()
	fmt.Println(queue.Print())
	if outcome != kkk.CE_OUTCOME_TRUE {
		fmt.Println("Outcome is not ", kkk.CE_OUTCOME_TRUE)
		t.FailNow()
	}
}