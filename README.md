# Shape Infer

A Go CLI for carrying out inference on images trained on a shape dataset for classification.


## Set up

1. Install Onnx Runtime for Linux

```
wget https://github.com/microsoft/onnxruntime/releases/download/v1.24.4/onnxruntime-linux-x64-1.24.4.tgz
tar -xzf onnxruntime-linux-x64-1.24.4.tgz
cp onnxruntime-linux-x64-1.24.4/lib/libonnxruntime.so.1.24.4 .
ln -sf libonnxruntime.so.1.24.4 onnxruntime.so
```

2. Build the Go binary

```
go build
```

## Usage 

```
./go_infer --image images/triangle.png --model shape_classifier.onnx 
```

Example output


```
Predicted class: triangle (index 0)

Scores:
  triangle       ███████████████████░ 0.9997
  quadrilateral  ░░░░░░░░░░░░░░░░░░░░ 0.0001
  pentagon       ░░░░░░░░░░░░░░░░░░░░ 0.0001
  hexagon        ░░░░░░░░░░░░░░░░░░░░ 0.0000
  heptagon       ░░░░░░░░░░░░░░░░░░░░ 0.0000
  octagon        ░░░░░░░░░░░░░░░░░░░░ 0.0000
  decagon        ░░░░░░░░░░░░░░░░░░░░ 0.0000
```
