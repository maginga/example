package main

func main() {
	parser, _ := NewParser()
	parser.ReadFile("/Users/hansonjang/Downloads/10.log")
	parser.WriteFile("sample3.log")
}
