package calculator

import (
	"strconv"
	"time"
	"unicode"

	"github.com/Solmorn/Calculator2/internal/config"
)

var GlobalError = 0
var cfg = config.Load()

type Node struct {
	data  string
	left  *Node
	right *Node
}

func NewNode(data string) *Node {
	return &Node{data: data, left: nil, right: nil}
}

type TAE struct {
	root *Node
}

func NewTAE() *TAE {
	return &TAE{root: nil}
}

func (t *TAE) calcTree(node *Node) float64 {

	if node.data != "+" && node.data != "-" && node.data != "*" && node.data != "/" {
		value, _ := strconv.Atoi(node.data)
		return float64(value)
	}

	switch node.data {
	case "+":
		time.Sleep(cfg.TIME_ADDITION_MS)
		return t.calcTree(node.left) + t.calcTree(node.right)
	case "-":
		time.Sleep(cfg.TIME_SUBTRACTION_MS)
		return t.calcTree(node.left) - t.calcTree(node.right)
	case "*":
		time.Sleep(cfg.TIME_MULTIPLICATIONS_MS)
		return t.calcTree(node.left) * t.calcTree(node.right)
	case "/":
		time.Sleep(cfg.TIME_DIVISIONS_MS)
		return t.calcTree(node.left) / t.calcTree(node.right)
	}
	return 0
}

func (t *TAE) setRoot(node *Node) {
	t.root = node
}

func (t *TAE) CalcTree() float64 {
	return t.calcTree(t.root)
}

func lex(inputString string) []string {

	var result []string
	current := string(inputString[0])
	for i := 1; i < len(inputString); i++ {
		if inputString[i] >= '0' && inputString[i] <= '9' && inputString[i-1] >= '0' && inputString[i-1] <= '9' {
			current += string(inputString[i])
		} else {
			result = append(result, current)
			current = ""
			current += string(inputString[i])
		}
	}
	if current != "" {
		result = append(result, current)
	}

	return result
}

func toPostfix(inputString string) []string {

	curr := lex(inputString)
	var postfix []string
	var ans []string

	for _, item := range curr {
		if item == "(" {
			postfix = append(postfix, item)
		} else if item == ")" {
			for len(postfix) > 0 && postfix[len(postfix)-1] != "(" {
				ans = append(ans, postfix[len(postfix)-1])
				postfix = postfix[:len(postfix)-1]
			}
			postfix = postfix[:len(postfix)-1]
		} else if item == "+" || item == "-" {
			for len(postfix) > 0 && (postfix[len(postfix)-1] == "+" || postfix[len(postfix)-1] == "-" || postfix[len(postfix)-1] == "/" || postfix[len(postfix)-1] == "*") {
				ans = append(ans, postfix[len(postfix)-1])
				postfix = postfix[:len(postfix)-1]
			}
			postfix = append(postfix, item)
		} else if item == "*" || item == "/" {
			for len(postfix) > 0 && (postfix[len(postfix)-1] == "/" || postfix[len(postfix)-1] == "*") {
				ans = append(ans, postfix[len(postfix)-1])
				postfix = postfix[:len(postfix)-1]
			}
			postfix = append(postfix, item)
		} else {
			ans = append(ans, item)
		}
	}

	for len(postfix) > 0 {
		ans = append(ans, postfix[len(postfix)-1])
		postfix = postfix[:len(postfix)-1]
	}

	return ans
}

func ECalc(expression string) float64 {

	tae := NewTAE()
	curr := toPostfix(expression)
	var queue []*Node
	for _, item := range curr {
		if item == "+" || item == "-" || item == "/" || item == "*" {
			rightChild := queue[len(queue)-1]
			queue = queue[:len(queue)-1]
			leftChild := queue[len(queue)-1]
			queue = queue[:len(queue)-1]

			newNode := NewNode(item)
			newNode.left = leftChild
			newNode.right = rightChild
			queue = append(queue, newNode)
		} else {
			newNode := NewNode(item)
			queue = append(queue, newNode)
		}
	}

	root := queue[len(queue)-1]
	tae.setRoot(root)

	return tae.CalcTree()

}

func Calc(expression string) (float64, error) {
	if checker(expression) {
		ans := ECalc(expression)
		if ans != ans+1 {
			return ans, nil
		} else {
			return 0, ErrDivisionByZero
		}
	} else {
		return 0, ErrInvalidExpression
	}
}

func checker(expression string) bool {
	openBrackets := 0
	isLastCharOperator := true

	for _, char := range expression {
		if unicode.IsSpace(char) {
			continue
		}

		if char == '(' {
			openBrackets++
			isLastCharOperator = true
		} else if char == ')' {
			openBrackets--
			if openBrackets < 0 {
				return false
			}
			isLastCharOperator = false
		} else if isOperator(char) {
			if isLastCharOperator {
				return false
			}
			isLastCharOperator = true
		} else if unicode.IsDigit(char) || unicode.IsLetter(char) {
			isLastCharOperator = false
		} else {
			return false
		}

	}

	return openBrackets == 0 && !isLastCharOperator
}

func isOperator(char rune) bool {
	operators := "+-*/"
	for _, op := range operators {
		if char == op {
			return true
		}
	}
	return false
}
