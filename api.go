package dicescript

import "errors"

// EvalOption 用于在执行前自定义 VM，例如调整配置或注册自定义骰子。
type EvalOption func(*Context) error

// WithContext 自定义 VM，便于注入自定义配置、变量或扩展能力。
func WithContext(fn func(*Context)) EvalOption {
	if fn == nil {
		return nil
	}
	return func(ctx *Context) error {
		fn(ctx)
		return nil
	}
}

// WithRollConfig 允许在执行前直接修改 VM 的 RollConfig。
func WithRollConfig(fn func(*RollConfig)) EvalOption {
	if fn == nil {
		return nil
	}
	return func(ctx *Context) error {
		fn(&ctx.Config)
		return nil
	}
}

// EvalResult 描述一次掷骰的计算过程以及最终结果。
type EvalResult struct {
	// Expression 为用户输入的完整表达式
	Expression string
	// Detail 为一步步的计算过程，可用于展示
	Detail string
	// Matched 表示已经被 DiceScript 消耗的文本
	Matched string
	// Rest 表示剩余未被解析的输入
	Rest string
	// Value 为最终返回的 VMValue，调用者可自行继续处理
	Value *VMValue
	// ValueText 为最终结果的字符串形式，便于直接输出
	ValueText string
	// StackTop, Depth, OpCount 提供执行期内部指标，可用于调试或监控
	StackTop int
	Depth    int
	OpCount  IntType
}

// EvalDice 以最简方式执行一条掷骰表达式，并返回计算过程与结果。
func EvalDice(expr string, opts ...EvalOption) (*EvalResult, error) {
	if expr == "" {
		return nil, errors.New("表达式不能为空")
	}

	vm := NewVM()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(vm); err != nil {
			return nil, err
		}
	}

	if err := vm.Run(expr); err != nil {
		return nil, err
	}

	value := NewNullVal()
	if vm.Ret != nil {
		value = vm.Ret.Clone()
	}

	return &EvalResult{
		Expression: expr,
		Detail:     vm.GetDetailText(),
		Matched:    vm.Matched,
		Rest:       vm.RestInput,
		Value:      value,
		ValueText:  value.ToString(),
		StackTop:   vm.StackTop(),
		Depth:      vm.Depth(),
		OpCount:    vm.NumOpCount,
	}, nil
}
