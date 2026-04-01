package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"

	ort "github.com/yalue/onnxruntime_go"
	"golang.org/x/image/draw"
)

func main() {
	modelPath := flag.String("model", "model.onnx", "Path to the ONNX model file")
	imagePath := flag.String("image", "", "Path to the input image (PNG or JPEG)")
	imageSize := flag.Int("size", 200, "Image size the model expects (e.g. 200 for 200x200)")
	labelsFlag := flag.String("labels", "triangle,quadrilateral,pentagon,hexagon,heptagon,octagon,decagon", "Comma-separated class labels (e.g. circle,square,triangle)")
	flag.Parse()

	if *imagePath == "" {
		fmt.Fprintln(os.Stderr, "error: --image is required")
		flag.Usage()
		os.Exit(1)
	}

	labels := parseLabels(*labelsFlag)

	pixels, err := loadImageTensor(*imagePath, *imageSize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading image: %v\n", err)
		os.Exit(1)
	}

	logits, err := runONNX(*modelPath, pixels, *imageSize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running ONNX inference: %v\n", err)
		os.Exit(1)
	}

	predIdx := argmax(logits)

	predLabel := fmt.Sprintf("class_%d", predIdx)
	if predIdx < len(labels) {
		predLabel = labels[predIdx]
	}

	fmt.Printf("Predicted class: %s (index %d)\n", predLabel, predIdx)

	if len(labels) > 0 {
		fmt.Println("\nScores:")
		scores := softmax(logits)
		for i, score := range scores {
			name := fmt.Sprintf("class_%d", i)
			if i < len(labels) {
				name = labels[i]
			}
			bar := progressBar(score, 20)
			fmt.Printf("  %-14s %s %.4f\n", name, bar, score)
		}
	}
}

func loadImageTensor(path string, size int) ([]float32, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	src, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	dst := image.NewGray(image.Rect(0, 0, size, size))
	draw.BiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Src, nil)

	pixels := make([]float32, size*size)
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			pixels[y*size+x] = float32(dst.GrayAt(x, y).Y) / 255.0
		}
	}
	return pixels, nil
}

func runONNX(modelPath string, pixels []float32, size int) ([]float32, error) {
	ort.SetSharedLibraryPath(ortLibPath())
	if err := ort.InitializeEnvironment(); err != nil {
		return nil, fmt.Errorf("init ORT: %w", err)
	}
	defer ort.DestroyEnvironment()

	inputShape := ort.NewShape(1, 1, int64(size), int64(size))
	inputTensor, err := ort.NewTensor(inputShape, pixels)
	if err != nil {
		return nil, fmt.Errorf("create input tensor: %w", err)
	}
	defer inputTensor.Destroy()

	// DynamicAdvancedSession allocates output tensors itself after Run(),
	// so we don't need to know num_classes ahead of time.
	session, err := ort.NewDynamicAdvancedSession(
		modelPath,
		[]string{"image"},
		[]string{"logits"},
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}
	defer session.Destroy()

	inputs := []ort.Value{inputTensor}
	var outputs []ort.Value
	outputTensors := session.Run(inputs, outputs)
	if len(outputTensors) == 0 {
		return nil, fmt.Errorf("run session: no outputs returned")
	}
	defer func() {
		for _, o := range outputTensors {
			o.Destroy()
		}
	}()

	// Cast the first output to a float32 tensor
	logitTensor, ok := outputTensors[0].(*ort.Tensor[float32])
	if !ok {
		return nil, fmt.Errorf("output tensor is not float32")
	}

	raw := logitTensor.GetData()
	result := make([]float32, len(raw))
	copy(result, raw)
	return result, nil
}

func argmax(vals []float32) int {
	best := 0
	for i := 1; i < len(vals); i++ {
		if vals[i] > vals[best] {
			best = i
		}
	}
	return best
}

func softmax(logits []float32) []float32 {
	maxVal := logits[0]
	for _, v := range logits {
		if v > maxVal {
			maxVal = v
		}
	}
	var sum float64
	exps := make([]float64, len(logits))
	for i, v := range logits {
		exps[i] = math.Exp(float64(v - maxVal))
		sum += exps[i]
	}
	out := make([]float32, len(logits))
	for i, e := range exps {
		out[i] = float32(e / sum)
	}
	return out
}

func progressBar(val float32, width int) string {
	filled := int(val * float32(width))
	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	return bar
}

func parseLabels(s string) []string {
	if s == "" {
		return nil
	}
	var labels []string
	cur := ""
	for _, c := range s {
		if c == ',' {
			if cur != "" {
				labels = append(labels, cur)
				cur = ""
			}
		} else {
			cur += string(c)
		}
	}
	if cur != "" {
		labels = append(labels, cur)
	}
	return labels
}

func ortLibPath() string {
	if p := os.Getenv("ORT_LIB_PATH"); p != "" {
		return p
	}
	return "./onnxruntime.so"
}
