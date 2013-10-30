package main

type User struct {
	id         int32
	characters []Character
}

type Character struct {
	Entity // anonomous field gives character all entity fields
}

type Entity struct {
	id       int32      // Uniquely id this entity in space
	position [2]float32 // coords x,y of entity
	movement [2]float32 // speed
	rotation float32    // speed of rotation around the Z axis (negative is counter clockwise)
	mass     float32    // mass effects physics!
}

type CelestialBody struct {
	Entity
	body_type string // 'star' 'planet' 'asteroid'
}

type Ship struct {
	Entity
	hull string // something something darkside
}
