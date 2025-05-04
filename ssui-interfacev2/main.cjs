const { app, BrowserWindow } = require('electron');
const path = require('path');
const fs = require('fs');
const express = require('express');
const http = require('http');

// Static file server
let server;
let serverPort = 49969;

function startServer() {
  const expressApp = express();
  
  // Determine the assets directory path
  let assetsPath;
  
  // In production
  const prodPath = path.join(process.resourcesPath, 'UIMod/v2');
  
  // Check if path exists
  if (fs.existsSync(prodPath)) {
    assetsPath = prodPath;
  } else {
    console.error('Could not find assets directory');
    app.quit();
    return false;
  }
  
  console.log('Serving assets from:', assetsPath);
  
  // Serve static files from the assets directory
  expressApp.use(express.static(assetsPath));
  
  // Start the server
  server = http.createServer(expressApp);
  
  return new Promise((resolve, reject) => {
    server.listen(serverPort, () => {
      console.log(`Server running at http://localhost:${serverPort}`);
      resolve(true);
    });
    
    server.on('error', (err) => {
      console.error('Server error:', err);
      if (err.code === 'EADDRINUSE') {
        serverPort++; // Try next port
        console.log(`Port ${serverPort - 1} in use, trying ${serverPort}`);
        server.close();
        startServer().then(resolve).catch(reject);
      } else {
        reject(err);
      }
    });
  });
}

async function createWindow() {
  // Start the local server first
  const serverStarted = await startServer();
  if (!serverStarted) return;

  const win = new BrowserWindow({
    width: 800,
    height: 600,
    webPreferences: {
      nodeIntegration: false,
      contextIsolation: true
    }
  });

  // Load from local server instead of file
  console.log(`Loading UI from: http://localhost:${serverPort}/index.html`);
  win.loadURL(`http://localhost:${serverPort}/index.html`);
  
  // For debugging
  // win.webContents.openDevTools();
}

app.commandLine.appendSwitch('ignore-certificate-errors');

app.whenReady().then(createWindow);

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

app.on('activate', () => {
  if (BrowserWindow.getAllWindows().length === 0) {
    createWindow();
  }
});

app.on('quit', () => {
  if (server) {
    server.close();
  }
});