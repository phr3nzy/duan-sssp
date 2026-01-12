package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/phr3nzy/duan-sssp/graph"
)

// GraphData for JSON export
type GraphData struct {
	Vertices int      `json:"vertices"`
	Edges    []Edge   `json:"edges"`
	Stats    Stats    `json:"stats"`
	Results  []Result `json:"results"`
}

type Edge struct {
	From   int     `json:"from"`
	To     int     `json:"to"`
	Weight float64 `json:"weight"`
}

type Stats struct {
	Vertices  int     `json:"vertices"`
	Edges     int     `json:"edges"`
	AvgDegree float64 `json:"avgDegree"`
	MaxDegree int     `json:"maxDegree"`
	Density   float64 `json:"density"`
}

type Result struct {
	Algorithm string        `json:"algorithm"`
	Time      time.Duration `json:"time"`
	TimeMS    float64       `json:"timeMs"`
	Speedup   float64       `json:"speedup"`
}

const htmlTemplate = `<!DOCTYPE html>
<html>
<head>
    <title>Duan SSSP Visual Benchmark</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            margin: 0;
            padding: 20px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: #333;
        }
        .container {
            max-width: 1400px;
            margin: 0 auto;
            background: white;
            border-radius: 10px;
            padding: 30px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.3);
        }
        h1 {
            text-align: center;
            color: #667eea;
            margin-bottom: 10px;
        }
        .subtitle {
            text-align: center;
            color: #666;
            margin-bottom: 30px;
            font-style: italic;
        }
        .grid {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 20px;
            margin-bottom: 30px;
        }
        .panel {
            border: 2px solid #e0e0e0;
            border-radius: 8px;
            padding: 20px;
        }
        .panel h2 {
            margin-top: 0;
            color: #667eea;
            border-bottom: 2px solid #667eea;
            padding-bottom: 10px;
        }
        canvas {
            border: 1px solid #ddd;
            border-radius: 5px;
            max-width: 100%;
        }
        .stats-grid {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 15px;
        }
        .stat {
            background: #f5f5f5;
            padding: 15px;
            border-radius: 5px;
            border-left: 4px solid #667eea;
        }
        .stat-label {
            font-size: 12px;
            color: #666;
            text-transform: uppercase;
        }
        .stat-value {
            font-size: 24px;
            font-weight: bold;
            color: #333;
        }
        .benchmark {
            margin: 10px 0;
            padding: 15px;
            background: #f9f9f9;
            border-radius: 5px;
            position: relative;
        }
        .benchmark-name {
            font-weight: bold;
            margin-bottom: 8px;
        }
        .bar-container {
            height: 30px;
            background: #e0e0e0;
            border-radius: 15px;
            overflow: hidden;
            position: relative;
        }
        .bar {
            height: 100%;
            transition: width 1s ease-out;
            display: flex;
            align-items: center;
            padding-left: 10px;
            color: white;
            font-weight: bold;
        }
        .bar.fastest { background: linear-gradient(90deg, #11998e 0%, #38ef7d 100%); }
        .bar.fast { background: linear-gradient(90deg, #f093fb 0%, #f5576c 100%); }
        .bar.slow { background: linear-gradient(90deg, #fa709a 0%, #fee140 100%); }
        .speedup {
            position: absolute;
            right: 10px;
            top: 50%;
            transform: translateY(-50%);
            font-weight: bold;
        }
        #graph-canvas {
            background: #fafafa;
        }
        .winner {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 20px;
            border-radius: 10px;
            text-align: center;
            margin: 20px 0;
            font-size: 20px;
            font-weight: bold;
        }
        .footer {
            text-align: center;
            color: #666;
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #e0e0e0;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🚀 Duan SSSP Visual Benchmark</h1>
        <p class="subtitle">Breaking the O(m + n log n) Sorting Barrier</p>
        
        <div class="grid">
            <div class="panel">
                <h2>📊 Graph Visualization</h2>
                <canvas id="graph-canvas" width="500" height="400"></canvas>
            </div>
            
            <div class="panel">
                <h2>📈 Graph Statistics</h2>
                <div class="stats-grid" id="stats"></div>
            </div>
        </div>
        
        <div class="panel">
            <h2>⚡ Benchmark Results</h2>
            <div id="results"></div>
            <div id="winner"></div>
        </div>
        
        <div class="footer">
            Implementation of Duan et al. (2025) | O(m log^(2/3) n) algorithm
        </div>
    </div>
    
    <script>
        const data = DATA_PLACEHOLDER;
        
        // Render stats
        const statsContainer = document.getElementById('stats');
        const stats = [
            { label: 'Vertices', value: data.stats.vertices.toLocaleString() },
            { label: 'Edges', value: data.stats.edges.toLocaleString() },
            { label: 'Avg Degree', value: data.stats.avgDegree.toFixed(2) },
            { label: 'Max Degree', value: data.stats.maxDegree },
            { label: 'Density', value: (data.stats.density * 100).toFixed(2) + '%' },
            { label: 'Cores Used', value: navigator.hardwareConcurrency || '?' }
        ];
        
        stats.forEach(stat => {
            const div = document.createElement('div');
            div.className = 'stat';
            div.innerHTML = '<div class="stat-label">' + stat.label + '</div>' +
                          '<div class="stat-value">' + stat.value + '</div>';
            statsContainer.appendChild(div);
        });
        
        // Render graph
        const canvas = document.getElementById('graph-canvas');
        const ctx = canvas.getContext('2d');
        renderGraph(ctx, data.edges, data.vertices);
        
        // Render results
        const resultsContainer = document.getElementById('results');
        const maxTime = Math.max(...data.results.map(r => r.time));
        
        data.results.forEach((result, idx) => {
            const div = document.createElement('div');
            div.className = 'benchmark';
            
            const barWidth = (result.time / maxTime) * 100;
            let barClass = 'fastest';
            if (idx > 0) barClass = 'fast';
            if (result.speedup < 0.5) barClass = 'slow';
            
            div.innerHTML = 
                '<div class="benchmark-name">' + result.algorithm + '</div>' +
                '<div class="bar-container">' +
                    '<div class="bar ' + barClass + '" style="width: ' + barWidth + '%">' +
                        result.timeMs.toFixed(3) + ' ms' +
                    '</div>' +
                    '<span class="speedup">' + result.speedup.toFixed(2) + 'x</span>' +
                '</div>';
            
            resultsContainer.appendChild(div);
        });
        
        // Winner announcement
        const winnerDiv = document.getElementById('winner');
        const winner = data.results[0];
        const runner = data.results[1];
        const speedup = (runner.time / winner.time).toFixed(1);
        
        winnerDiv.className = 'winner';
        winnerDiv.innerHTML = '🏆 ' + winner.algorithm + ' wins by ' + speedup + 'x!';
        
        function renderGraph(ctx, edges, vertexCount) {
            const width = canvas.width;
            const height = canvas.height;
            const sampleSize = Math.min(50, vertexCount);
            const padding = 40;
            
            // Generate vertex positions in circular layout
            const positions = [];
            const centerX = width / 2;
            const centerY = height / 2;
            const radius = Math.min(width, height) / 2 - padding;
            
            for (let i = 0; i < sampleSize; i++) {
                const angle = (i / sampleSize) * 2 * Math.PI;
                positions.push({
                    x: centerX + radius * Math.cos(angle),
                    y: centerY + radius * Math.sin(angle)
                });
            }
            
            // Draw edges (sample)
            ctx.strokeStyle = '#ddd';
            ctx.lineWidth = 1;
            
            edges.slice(0, Math.min(100, edges.length)).forEach(edge => {
                if (edge.from < sampleSize && edge.to < sampleSize) {
                    const from = positions[edge.from];
                    const to = positions[edge.to];
                    
                    ctx.beginPath();
                    ctx.moveTo(from.x, from.y);
                    ctx.lineTo(to.x, to.y);
                    ctx.stroke();
                }
            });
            
            // Draw vertices
            positions.forEach((pos, idx) => {
                ctx.fillStyle = idx === 0 ? '#667eea' : '#38ef7d';
                ctx.beginPath();
                ctx.arc(pos.x, pos.y, 6, 0, 2 * Math.PI);
                ctx.fill();
                
                // Label source
                if (idx === 0) {
                    ctx.fillStyle = '#333';
                    ctx.font = 'bold 12px sans-serif';
                    ctx.fillText('Source', pos.x + 10, pos.y);
                }
            });
            
            // Info text
            ctx.fillStyle = '#666';
            ctx.font = '12px sans-serif';
            ctx.fillText('Showing ' + sampleSize + ' of ' + vertexCount + ' vertices', 10, height - 10);
        }
    </script>
</body>
</html>`

func startWebVisualization(g *graph.Graph, results []BenchmarkResult) {
	// Prepare data
	edges := make([]Edge, 0)
	for u := 0; u < min(g.V, 100); u++ { // Limit for JSON size
		for _, edge := range g.Adj[u] {
			if edge.To < 100 {
				edges = append(edges, Edge{
					From:   u,
					To:     edge.To,
					Weight: edge.Weight,
				})
			}
		}
	}

	maxDegree := 0
	totalDegree := 0
	for u := 0; u < g.V; u++ {
		deg := len(g.Adj[u])
		totalDegree += deg
		if deg > maxDegree {
			maxDegree = deg
		}
	}

	graphData := GraphData{
		Vertices: g.V,
		Edges:    edges,
		Stats: Stats{
			Vertices:  g.V,
			Edges:     len(edges),
			AvgDegree: float64(totalDegree) / float64(g.V),
			MaxDegree: maxDegree,
			Density:   float64(totalDegree) / float64(g.V*g.V),
		},
		Results: make([]Result, len(results)),
	}

	baseline := results[0].Time
	for i, r := range results {
		graphData.Results[i] = Result{
			Algorithm: r.Algorithm,
			Time:      r.Time,
			TimeMS:    float64(r.Time.Microseconds()) / 1000.0,
			Speedup:   float64(baseline) / float64(r.Time),
		}
	}

	jsonData, _ := json.Marshal(graphData)

	// Create HTML file
	htmlContent := htmlTemplate
	htmlContent = string([]byte(htmlContent))
	htmlContent = replaceString(htmlContent, "DATA_PLACEHOLDER", string(jsonData))

	filename := "benchmark_viz.html"
	err := os.WriteFile(filename, []byte(htmlContent), 0644)
	if err != nil {
		fmt.Printf("Error creating HTML: %v\n", err)
		return
	}

	fmt.Printf("\n%s🌐 Web visualization created: %s%s\n", colorCyan, filename, colorReset)
	fmt.Printf("%sOpening in browser...%s\n", colorCyan, colorReset)

	// Open in browser
	openBrowser("http://localhost:8080/" + filename)

	// Start simple HTTP server
	go func() {
		http.HandleFunc("/"+filename, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			http.ServeFile(w, r, filename)
		})
		http.ListenAndServe(":8080", nil)
	}()

	time.Sleep(2 * time.Second) // Give browser time to open
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		fmt.Printf("Please open %s in your browser\n", url)
		return
	}

	if err != nil {
		fmt.Printf("Please open %s in your browser\n", url)
	}
}

func replaceString(s, old, new string) string {
	result := ""
	remaining := s

	for {
		idx := findString(remaining, old)
		if idx == -1 {
			result += remaining
			break
		}
		result += remaining[:idx] + new
		remaining = remaining[idx+len(old):]
	}

	return result
}

func findString(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
