package main

import (
	"html/template"
)

// TemplateData holds the data for rendering the HTML template
type TemplateData struct {
	Filename string
	Content  template.HTML
	Port     int
}

// htmlTemplate contains the HTML structure with embedded CSS and WebSocket script
const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Filename}}</title>
    <style>
        body {
            margin: 0;
            padding: 20px;
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
            font-size: 16px;
            line-height: 1.6;
            color: #24292e;
            background-color: #fff;
        }
        .markdown-body {
            max-width: 980px;
            margin: 0 auto;
            padding: 45px;
            box-sizing: border-box;
        }
        .markdown-body h1, .markdown-body h2, .markdown-body h3, 
        .markdown-body h4, .markdown-body h5, .markdown-body h6 {
            margin-top: 24px;
            margin-bottom: 16px;
            font-weight: 600;
            line-height: 1.25;
        }
        .markdown-body h1 {
            font-size: 2em;
            border-bottom: 1px solid #eaecef;
            padding-bottom: 0.3em;
        }
        .markdown-body h2 {
            font-size: 1.5em;
            border-bottom: 1px solid #eaecef;
            padding-bottom: 0.3em;
        }
        .markdown-body h3 { font-size: 1.25em; }
        .markdown-body h4 { font-size: 1em; }
        .markdown-body h5 { font-size: 0.875em; }
        .markdown-body h6 { font-size: 0.85em; color: #6a737d; }
        .markdown-body p {
            margin-top: 0;
            margin-bottom: 16px;
        }
        .markdown-body a {
            color: #0366d6;
            text-decoration: none;
        }
        .markdown-body a:hover {
            text-decoration: underline;
        }
        .markdown-body code {
            padding: 0.2em 0.4em;
            margin: 0;
            font-size: 85%;
            background-color: rgba(27,31,35,0.05);
            border-radius: 3px;
            font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, monospace;
        }
        .markdown-body pre {
            padding: 16px;
            overflow: auto;
            font-size: 85%;
            line-height: 1.45;
            background-color: #f6f8fa;
            border-radius: 3px;
            margin-bottom: 16px;
        }
        .markdown-body pre code {
            padding: 0;
            background-color: transparent;
            border-radius: 0;
        }
        .markdown-body blockquote {
            padding: 0 1em;
            color: #6a737d;
            border-left: 0.25em solid #dfe2e5;
            margin: 0 0 16px 0;
        }
        .markdown-body ul, .markdown-body ol {
            padding-left: 2em;
            margin-top: 0;
            margin-bottom: 16px;
        }
        .markdown-body li {
            margin-bottom: 0.25em;
        }
        .markdown-body table {
            border-collapse: collapse;
            margin-bottom: 16px;
            width: 100%;
        }
        .markdown-body table th, .markdown-body table td {
            padding: 6px 13px;
            border: 1px solid #dfe2e5;
        }
        .markdown-body table th {
            font-weight: 600;
            background-color: #f6f8fa;
        }
        .markdown-body table tr {
            background-color: #fff;
            border-top: 1px solid #c6cbd1;
        }
        .markdown-body table tr:nth-child(2n) {
            background-color: #f6f8fa;
        }
        .markdown-body img {
            max-width: 100%;
            box-sizing: border-box;
        }
        .markdown-body hr {
            height: 0.25em;
            padding: 0;
            margin: 24px 0;
            background-color: #e1e4e8;
            border: 0;
        }
    </style>
</head>
<body>
    <article class="markdown-body">
        {{.Content}}
    </article>
    
    <script>
        // WebSocket for live reload
        function connect() {
            const ws = new WebSocket('ws://localhost:{{.Port}}/ws');
            
            ws.onmessage = function(event) {
                if (event.data === 'reload') {
                    location.reload();
                }
            };
            
            ws.onerror = function() {
                console.log('WebSocket error, will retry...');
            };
            
            ws.onclose = function() {
                console.log('WebSocket closed, reconnecting in 1s...');
                setTimeout(connect, 1000);
            };
        }
        
        connect();
    </script>
</body>
</html>`

// GetTemplate returns the parsed HTML template
func GetTemplate() (*template.Template, error) {
	return template.New("html").Parse(htmlTemplate)
}
