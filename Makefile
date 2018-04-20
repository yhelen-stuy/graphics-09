all:
	go run draw.go image.go main.go matrix.go parser.go vector.go stack.go

clean:
	rm *.ppm
	rm *.png
