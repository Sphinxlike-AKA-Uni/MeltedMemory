package main
import (
	"fmt"
	"strings"
	"image"
	"image/png"
	"image/jpeg"
	"image/gif"
	"image/color"
	"os"
	"time"
	"flag"
)

func maxHeapify(tosort []int, position int) {
	size := len(tosort)
	maximum := position
	leftChild := 2*position + 1
	rightChild := leftChild + 1
	if leftChild < size && tosort[leftChild] > tosort[position] {
		maximum = leftChild
	}
		
	if rightChild < size && tosort[rightChild] > tosort[maximum] {
		maximum = rightChild
	}
	
	if position != maximum {
		tosort[position], tosort[maximum] = tosort[maximum], tosort[position]
		maxHeapify(tosort, maximum)
	}
}

func buildMaxHeap(tosort []int) {
	// from http://en.wikipedia.org/wiki/Heapsort
	// iParent := floor((i-1)/2)
	for i := (len(tosort)-1)/2; i >= 0; i-- {
		maxHeapify(tosort, i)
	}
}

func heapSort(tosort []int) {
	//buildMaxHeap(tosort) // To make the heapsort incomplete
	for i := len(tosort)-1; i >= 1; i-- {
		tosort[i], tosort[0] = tosort[0], tosort[i]
		maxHeapify(tosort[:i-1], 0)
	}
	
}

var filename string = ""

func main() {
	flag.Parse() // Parse command line arguments
	if len(flag.Args()) == 0 { // No arguments
		fmt.Println("No file specified")
		time.Sleep(5*time.Second)
		os.Exit(1)
	}
	
	filename = flag.Args()[0] // Get first provided argument
	
	// Check image format
	if strings.HasSuffix(filename, ".png") {
		image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	} else if strings.HasSuffix(filename, ".jpg") || strings.HasSuffix(filename, ".jpeg") {
		image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	} else if strings.HasSuffix(filename, ".gif") {
		WorksGIF()
		return
	} else {
		fmt.Println("filename does not end with .png, .gif, .jpg, or .jpeg")
		time.Sleep(5*time.Second)
		os.Exit(2)
	}
	WorksImage()
}

func WorksImage() {
	// Open file
	imgfile, err := os.Open(filename)
	if err != nil { // Can't opens
		fmt.Println("Error opening provided image file: ", err)
		time.Sleep(5*time.Second)
		os.Exit(3)
	}
	
	defer imgfile.Close() // close when finished using this
	// Decode the image data
	imgcfg, _, err := image.Decode(imgfile)
	if err != nil {
		fmt.Println("Error decoding image: ", err)
		time.Sleep(5*time.Second)
		os.Exit(4)
	}
	
	// Load pixels width and height
	w, h := imgcfg.Bounds().Max.X, imgcfg.Bounds().Max.Y // width and height
	mi := image.NewRGBA(image.Rect(0, 0, w, h)) // create replicated image
	fmt.Println("Width: ", w, "\nHeight: ", h)
	time.Sleep(1500*time.Millisecond)
	
	// Heapsort by up to down
	for y := 0; y < h; y++ {
		fmt.Println("y: ", y)
		pa := []int{}
		for x := 0; x < w; x++ {
			r, g, b, a := imgcfg.At(x, y).RGBA()
			r = r / 255
			g = g / 255
			b = b / 255
			if a == 65535 {
				a = 255
			} else {
				a = a / 255
			}
			p := int(r<<24 | g<<16 | b<<8 | a)
			pa = append(pa, p)
		}
		heapSort(pa) // rearrange
		for x := 0; x < w; x++ {
			p := pa[x]
			mi.SetRGBA(x, y, color.RGBA{uint8(p>>24), uint8((p<<8)>>24), uint8((p<<16)>>24), uint8((p<<24)>>24)})
		}
	}
	
	
	// Heapsort by left to right
	for x := 0; x < w; x++ {
		fmt.Println("x: ", x)
		pa := []int{}
		for y := 0; y < h; y++ {
			r, g, b, a := imgcfg.At(x, y).RGBA()
			r = r / 255
			g = g / 255
			b = b / 255
			if a == 65535 {
				a = 255
			} else {
				a = a / 255
			}
			p := int(r<<24 | g<<16 | b<<8 | a)
			pa = append(pa, p)
		}
		heapSort(pa) // rearrange
		for y := 0; y < h; y++ {
			p := pa[y]
			mi.SetRGBA(x, y, color.RGBA{uint8(p>>24), uint8((p<<8)>>24), uint8((p<<16)>>24), uint8((p<<24)>>24)})
		}
	}

	// Write the results into jpeg image
	f, _ := os.OpenFile("out.jpeg", os.O_WRONLY|os.O_CREATE, 436) // Opens the file
	defer f.Close() // Closes when finished
	var opt jpeg.Options
	opt.Quality = 99 // Set JPEG quality to 99%
	jpeg.Encode(f, mi, &opt)
	
	// Completes succuessfully

}

func WorksGIF() {
	giffile, err := os.Open(filename)
	if err != nil { // Can't opens
		fmt.Println("Error opening provided image file: ", err)
		time.Sleep(5*time.Second)
		os.Exit(3)
	}
	
	g, err := gif.DecodeAll(giffile)
	if err != nil {
		fmt.Println("DecodeAll Error: ", err)
		time.Sleep(5*time.Second)
		os.Exit(5)
	}
	
	mi := len(g.Image)
	
	for i, ii := range g.Image {
		w, h := ii.Bounds().Max.X, ii.Bounds().Max.Y
		// Heapsort by up to down
		for y := 0; y < h; y++ {
			fmt.Printf("Frame(%d/%d) y: %d\n", i+1, mi, y)
			pa := []int{}
			for x := 0; x < w; x++ {
				r, g, b, a := ii.At(x, y).RGBA()
				r = r / 255
				g = g / 255
				b = b / 255
				if a == 65535 {
					a = 255
				} else {
					a = a / 255
				}
				p := int(r<<24 | g<<16 | b<<8 | a)
				pa = append(pa, p)
			}
			heapSort(pa) // rearrange
			for x := 0; x < w; x++ {
				p := pa[x]
				ii.Set(x, y, color.RGBA{uint8(p>>24), uint8((p<<8)>>24), uint8((p<<16)>>24), uint8((p<<24)>>24)})
			}
		}
		// Heapsort by left to right
		for x := 0; x < w; x++ {
			fmt.Printf("Frame(%d/%d) x: %d\n", i+1, mi, x)
			pa := []int{}
			for y := 0; y < h; y++ {
				r, g, b, a := ii.At(x, y).RGBA()
				r = r / 255
				g = g / 255
				b = b / 255
				if a == 65535 {
					a = 255
				} else {
					a = a / 255
				}
				p := int(r<<24 | g<<16 | b<<8 | a)
				pa = append(pa, p)
			}
			heapSort(pa) // rearrange
			for y := 0; y < h; y++ {
				p := pa[y]
				ii.Set(x, y, color.RGBA{uint8(p>>24), uint8((p<<8)>>24), uint8((p<<16)>>24), uint8((p<<24)>>24)})
			}
		}
	}
	
	f, _ := os.OpenFile("out.gif", os.O_WRONLY|os.O_CREATE, 436) // Opens the file
	defer f.Close() // Closes when finished
	gif.EncodeAll(f, g)
	
}
