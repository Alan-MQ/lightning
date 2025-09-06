package examples

// AddArgs 加法运算的参数
type AddArgs struct {
	A, B int
}

// AddReply 加法运算的结果
type AddReply struct {
	Sum int
}

// MultiplyArgs 乘法运算的参数
type MultiplyArgs struct {
	A, B int
}

// MultiplyReply 乘法运算的结果
type MultiplyReply struct {
	Product int
}

// Calculator 计算器服务
type Calculator struct{}

// Add 加法运算
func (c *Calculator) Add(args *AddArgs, reply *AddReply) error {
	reply.Sum = args.A + args.B
	return nil
}

// Multiply 乘法运算
func (c *Calculator) Multiply(args *MultiplyArgs, reply *MultiplyReply) error {
	reply.Product = args.A * args.B
	return nil
}

// 这个方法不符合RPC规范（参数数量不对），应该被忽略
func (c *Calculator) InvalidMethod(a int) int {
	return a * 2
}

// 这个方法不是导出的，应该被忽略
func (c *Calculator) privateMethod(args *AddArgs, reply *AddReply) error {
	reply.Sum = args.A + args.B
	return nil
}