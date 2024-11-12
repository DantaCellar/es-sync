package main

func main() {
	InitConfig()
	InitSearch()
	InitDB()

	DumpFloors(IndexNameFloor)
	DumpProject()
	DumpTag()
}
