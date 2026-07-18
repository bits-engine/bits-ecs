package main

// Defining some components and resources
type health struct {
	current int
	max     int
}
type name struct {
	value string
}

// resource
type gameState struct {
	tick          int
	stopAfterTick int
}
