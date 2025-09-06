package examples

// ReverseArgs 反转字符串的参数
type ReverseArgs struct {
	S string
}

// ReverseReply the result of reversing a string
type ReverseReply struct {
	Reversed string
}

type StringService struct{}

func (s *StringService) Reverse(args *ReverseArgs, reply *ReverseReply) error {
	runes := []rune(args.S)
	for i := 0; i < len(runes)/2; i++ {
		runes[i], runes[len(runes)-1-i] = runes[len(runes)-1-i], runes[i]
	}
	reply.Reversed = string(runes)
	return nil
}
