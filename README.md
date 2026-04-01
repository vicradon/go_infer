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