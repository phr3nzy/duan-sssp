#!/bin/bash
# Quick start script for Duan SSSP Visual Benchmark
# Uses all CPU cores with visualization

set -e

echo "🚀 Duan SSSP Visual Benchmark - Quick Start"
echo "============================================="
echo ""

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo "❌ Error: Please run this script from the project root directory"
    exit 1
fi

echo "📦 Building visual benchmark tool..."
go build -o visualbench ./cmd/visualbench

if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi

echo "✅ Build successful!"
echo ""
echo "🎯 Running visual benchmark with ALL CPU cores..."
echo ""

# Run with all cores and web visualization
./visualbench \
  -vertices=10000 \
  -edge-factor=3 \
  -iterations=10 \
  -parallel=true \
  -show-graph=true \
  -web=true

echo ""
echo "🎉 Benchmark complete!"
