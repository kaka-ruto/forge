# Edge AI Inference Device

This example demonstrates how to build an edge AI device for real-time computer vision and machine learning inference using Forge OS. The device runs TensorFlow Lite models locally with camera input and publishes results via MQTT.

## Features

- **Real-time Inference**: Sub-100ms processing with TensorFlow Lite
- **Camera Integration**: Multi-camera support with V4L2
- **GPU Acceleration**: Hardware-accelerated ML on compatible devices
- **Model Management**: Dynamic model loading and switching
- **MQTT Publishing**: Real-time results streaming
- **Performance Monitoring**: FPS tracking and latency metrics

## Use Case Scenario

A company needs edge devices for real-time video analysis in manufacturing. The devices should:

1. Capture video from production line cameras
2. Run object detection models locally
3. Identify defects or anomalies in real-time
4. Publish alerts and metrics to central monitoring
5. Operate autonomously without cloud dependency

## Quick Start

```bash
# Create the edge AI device project
forge new edge-ai-inference --template=iot --arch=aarch64

# Navigate to the project
cd edge-ai-inference

# Add AI/ML packages
forge add package opencv
forge add package tensorflow-lite
forge add package python3
forge add package gstreamer

# Add camera and processing packages
forge add package v4l-utils
forge add package ffmpeg

# Add monitoring features
forge add feature monitoring
forge add feature auto-updates

# Configure cameras and models (edit forge.yml)
# - Set camera devices and resolutions
# - Configure ML models and thresholds
# - Set up MQTT publishing

# Build with performance optimizations
forge build --optimize-for=performance

# Test locally
forge test --port 8080

# Deploy to edge hardware
forge deploy remote --host edge-device-01.local
```

## Hardware Requirements

- ARM64/AArch64 device (Jetson Nano, Xavier, Raspberry Pi 4)
- Camera(s) with V4L2 support (USB webcams, CSI cameras)
- 4GB+ RAM for ML processing
- GPU support (optional, for acceleration)
- Ethernet connectivity
- 32GB+ storage

## Configuration

### Camera Setup

The device supports multiple camera configurations:

```yaml
cameras:
  - device: "/dev/video0"
    resolution: "1920x1080"
    fps: 30
    name: "Primary Camera"
```

### ML Model Configuration

Pre-configured models for common tasks:

```yaml
models:
  - name: "object_detection"
    framework: "tflite"
    model_path: "/opt/models/efficientdet_lite0.tflite"
    input_size: "320x320"
    threshold: 0.5
```

### Inference Pipeline

Configurable processing pipeline:

```yaml
inference:
  pipeline:
    - stage: "capture"
      camera: "primary"
      resolution: "640x480"
    - stage: "preprocess"
      operations: ["resize", "normalize"]
    - stage: "detect"
      model: "object_detection"
    - stage: "output"
      format: "mqtt"
```

## Expected Performance

- **Build Time**: < 45 minutes
- **Image Size**: ~800MB
- **Boot Time**: < 30 seconds
- **Inference Latency**: < 100ms per frame
- **Power Consumption**: 5-15W depending on hardware
- **Memory Usage**: 1-2GB during inference

## Supported Models

### Object Detection
- EfficientDet-Lite (fast, accurate)
- SSD MobileNet (balanced performance)
- YOLOv5 (high accuracy)

### Classification
- MobileNetV2 (efficient)
- ResNet (accurate)
- Custom trained models

### Pose Estimation
- MoveNet (real-time)
- PoseNet (legacy)

## Camera Compatibility

### USB Cameras
- Logitech C920, C930e
- Microsoft LifeCam
- Generic UVC cameras

### CSI Cameras (Raspberry Pi/Jetson)
- Raspberry Pi Camera Module
- Jetson Nano camera
- IMX219/IMX477 sensors

### IP Cameras
- RTSP stream support
- ONVIF compatibility
- H.264/H.265 decoding

## GPU Acceleration

### NVIDIA Jetson
- TensorRT optimization
- CUDA acceleration
- 5-10x performance improvement

### Intel Movidius
- Neural Compute Stick
- Myriad X VPU
- Low-power inference

### CPU-only Fallback
- NEON optimization on ARM
- Multi-threading support
- Quantized models

## MQTT Integration

Results published in structured format:

```json
{
  "device": "edge-device-01",
  "timestamp": "2024-01-01T12:00:00Z",
  "detections": [
    {
      "label": "person",
      "confidence": 0.87,
      "bbox": [100, 150, 200, 300]
    }
  ],
  "fps": 12.5,
  "latency_ms": 85
}
```

## Monitoring & Metrics

Comprehensive performance tracking:

- **Inference FPS**: Real-time processing rate
- **Detection Count**: Objects detected per minute
- **Processing Latency**: End-to-end delay
- **Model Load Time**: Model switching performance
- **Memory Usage**: RAM consumption tracking

## Development API

REST API for development and debugging:

```bash
# Get device status
curl http://edge-device:8080/status

# List available models
curl http://edge-device:8080/models

# Switch models
curl -X POST http://edge-device:8080/model \
  -d '{"name": "pose_estimation"}'

# Get inference results
curl http://edge-device:8080/results
```

## Model Management

### Pre-trained Models
Included models for common tasks:
- COCO object detection (80 classes)
- Face detection and recognition
- Pose estimation
- Hand tracking

### Custom Models
Support for custom trained models:
- TensorFlow Lite format
- Quantized for edge deployment
- Optimized for target hardware

### Model Updates
- OTA model deployment
- A/B model switching
- Fallback to previous version
- Performance validation

## Power Optimization

### CPU/GPU Scaling
- Dynamic frequency scaling
- Workload-based adjustment
- Thermal throttling protection

### Model Optimization
- INT8 quantization
- Model pruning
- Knowledge distillation

### Processing Optimization
- Frame skipping for low-power modes
- Resolution scaling
- Batch processing

## Security Features

### Model Protection
- Encrypted model storage
- Runtime integrity checks
- Tamper detection

### Network Security
- TLS MQTT encryption
- Authentication tokens
- Access control lists

### Device Security
- Secure boot
- Encrypted storage
- Remote attestation

## Troubleshooting

### Camera Issues
- Check V4L2 device nodes: `v4l2-ctl --list-devices`
- Verify camera permissions
- Test with `gst-launch` pipeline
- Check USB power delivery

### Model Loading Problems
- Verify model file paths
- Check model compatibility
- Review TensorFlow Lite logs
- Test with sample images

### Performance Issues
- Monitor CPU/GPU usage
- Check thermal throttling
- Profile inference pipeline
- Optimize model parameters

### MQTT Connection Problems
- Verify broker connectivity
- Check authentication credentials
- Review network configuration
- Test with MQTT client tools

## Integration Examples

### Manufacturing Quality Control
- Defect detection on assembly lines
- Product counting and verification
- Safety gear compliance monitoring

### Retail Analytics
- Customer counting and tracking
- Shelf inventory monitoring
- Queue length analysis

### Smart City Applications
- Traffic monitoring and analysis
- Parking space detection
- Waste level monitoring

## Performance Benchmarks

### Jetson Nano
- Object Detection: 15-25 FPS
- Pose Estimation: 10-15 FPS
- Face Recognition: 8-12 FPS

### Raspberry Pi 4
- Object Detection: 3-5 FPS
- Pose Estimation: 2-3 FPS
- Classification: 5-8 FPS

### x86_64 with GPU
- Object Detection: 50-100 FPS
- All models: Hardware dependent

## Expansion Options

### Additional Sensors
- LiDAR integration
- Thermal cameras
- Multi-spectral imaging
- Audio processing

### Advanced Features
- Multi-model pipelines
- Federated learning
- Edge training
- Model compression

### Cloud Integration
- Model updates from cloud
- Result aggregation
- Remote monitoring
- Firmware updates

## Community Resources

- **TensorFlow Lite**: Model optimization guides
- **OpenCV**: Computer vision tutorials
- **MQTT**: Protocol documentation
- **Forge OS**: Framework documentation

## Next Steps

- Train custom models for specific use cases
- Implement multi-camera synchronization
- Add edge training capabilities
- Integrate with industrial protocols
- Develop custom inference pipelines