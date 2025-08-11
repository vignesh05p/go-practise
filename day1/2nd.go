switch v := i.(type) {
case string:
	fmt.Println("It's a string:", v)
case int:
	fmt.Println("It's an int:", v)
case bool:
	fmt.Println("It's a bool:", v)
default:
	fmt.Println("Unknown type")
}
