package main

import (
	"image"
	"image/jpeg"
	"social-network/pkg/services"
)

func init() {
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
}

func main() {
	services.Server()
}
