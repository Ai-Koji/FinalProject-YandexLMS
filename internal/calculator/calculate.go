package calculator

import (
	"errors"
	"strings"
    "strconv"
	"fmt"
	"regexp"
)

// stack realization
type stack []string

func (s stack) Push(v string) stack {
	return append(s, v)
}

func (s stack) Pop() (stack, string) {
	l := len(s)
	return s[:l-1], s[l-1]
}

func (s stack) Get() string {
	l := len(s)
	return s[l-1]
}


func CheckExpression(expression string) (error) {
	matched, err := regexp.MatchString(`()^[\d+\-*/()]+$`, expression)
	if err != nil || !matched {
		err = fmt.Errorf("Incorrect expression")
	}
	return err
}

// translate string to array(slice)
func ArrayNotation(expression string) ([]string, error) {
	expression = strings.ReplaceAll(expression, " ", "") // delete spaces

	var ArrayOperation []string = []string{""}
	ArrayIndex := 0
	var Error error = nil

	for _, Symbol := range expression {
		if Symbol >= 48 && Symbol <= 57 || Symbol == 46 || Symbol == 44 {
			ArrayOperation[ArrayIndex] = ArrayOperation[ArrayIndex] + string(Symbol)
		} else if Symbol >= 40 && Symbol <= 47 && ArrayIndex != -1 {
			if ArrayOperation[ArrayIndex] == "" {
				ArrayOperation[ArrayIndex] = string(Symbol)
			} else {
				ArrayOperation = append(ArrayOperation, string(Symbol))
				ArrayIndex++
			}
			ArrayOperation = append(ArrayOperation, "")
			ArrayIndex++
		} else {
			Error = errors.New("unknown symbol")
			break
		}
	}

	// delete "" on end
	if ArrayOperation[len(ArrayOperation)-1] == "" {
		ArrayOperation = ArrayOperation[: len(ArrayOperation)-1]
	}

	return ArrayOperation, Error
}

// priority for postfix notation
func Priority(operation string) int {
	PriorityRes := 0
	switch operation {
	case "+":
		PriorityRes = 1
		break
	case "-":
		PriorityRes = 1
		break
	case "*":
		PriorityRes = 2
		break
	case "/":
		PriorityRes = 2
		break
	case ")":
		PriorityRes = 4
		break
	case "(":
		PriorityRes = 5
		break
	case "s":
		PriorityRes = 6
		break
	}
	return PriorityRes
}

// translate infix to postfix notation for future calculations
func PostfixNotationArr(expression string) ([]string, error) {
	var Error error = nil
	var NotationArr []string
	NotationArr, Error = ArrayNotation(expression)
	var Result []string


	var NumbersCount int = 0
	var OperationsCount int = 0
	var BracketCount int = 0

	if Error == nil {
		StackOperations := make(stack, 0)
		var PriorityIndex int
		var Operation string

		for _, Symbol := range NotationArr {
			PriorityIndex = Priority(Symbol)
			if PriorityIndex == 0 {
				Result = append(Result, Symbol)	
				NumbersCount++
			} else if PriorityIndex == 5 {
				StackOperations = StackOperations.Push(Symbol);
				BracketCount++
			} else if PriorityIndex == 4 {
				for ; len(StackOperations) > 0 && StackOperations.Get() != "(" ; {
					StackOperations, Operation = StackOperations.Pop()
					Result = append(Result, Operation)
				}
				StackOperations, Operation = StackOperations.Pop()
				BracketCount--
			} else {
				OperationsCount++
				if len(StackOperations) == 0 {
					StackOperations = StackOperations.Push(Symbol)
				} else if PriorityIndex < Priority(StackOperations.Get()) && StackOperations.Get() != "(" {
					for ;len(StackOperations) > 0 && Priority(StackOperations.Get()) >= PriorityIndex; {						
						StackOperations, Operation = StackOperations.Pop()
						Result = append(Result, Operation)
					}
					StackOperations = StackOperations.Push(Symbol) 
				} else {
					StackOperations = StackOperations.Push(Symbol)
				}
			}
		}

		for ;len(StackOperations)>0; {
			StackOperations, Operation = StackOperations.Pop()
			Result = append(Result, Operation)
		}
	}

	if NumbersCount != OperationsCount + 1 || BracketCount != 0 {
		Error = errors.New("incorrect equation")
		Result = []string{}
	}

	return Result, Error
}

//calculate value from postfix notion
func CalcFromPostfix(PostfixArr []string) (float64, error) {
	var err error
	var Result float64 = 0
	StackNumbers := make(stack, 0)
	var Number1, Number2 float64
	var Temp string

	for _, element := range PostfixArr {
		if Priority(element) == 0 {
			StackNumbers = StackNumbers.Push(element);
		} else {
			StackNumbers, Temp = StackNumbers.Pop()
			Number1, _ = strconv.ParseFloat(Temp,64)
			StackNumbers, Temp = StackNumbers.Pop()
			Number2, _ = strconv.ParseFloat(Temp,64)
			switch element {
				case "+":
					StackNumbers = StackNumbers.Push(fmt.Sprintf("%f", Number1+Number2))
					break
				case "-":
					StackNumbers = StackNumbers.Push(fmt.Sprintf("%f", Number2-Number1))
					break
				case "*":
					StackNumbers =StackNumbers.Push(fmt.Sprintf("%f", Number1*Number2))
					break
				case "/":
					StackNumbers =StackNumbers.Push(fmt.Sprintf("%f", Number2/Number1))
					break
				}
		}
	}
	Result, err = strconv.ParseFloat(StackNumbers[0],64)
	return Result, err;
}

// main function for calculate
func Calc(expression string) (float64, error) {
	var Result float64 = 0
	var PostfixArr []string 
	var Error error;
	Error = CheckExpression(expression)

	if (Error == nil) {
		PostfixArr, Error = PostfixNotationArr(expression)
	}
	if (Error == nil) {
		Result, Error = CalcFromPostfix(PostfixArr)
	}

	return Result, Error
}