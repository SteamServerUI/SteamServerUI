const { app, BrowserWindow, dialog } = require('electron');
const path = require('path');
const fs = require('fs');
const express = require('express');
const http = require('http');

// Static file server
let server;
let currentPortIndex = 0;
const DEFAULT_PORTS = [28080, 28888, 29090, 27070, 26060, 35050, 34040, 30303, 34320, 34899];

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
    dialog.showErrorBox('Error', 'Could not find assets directory: ' + prodPath);
    app.exit(1);
    return false;
  }
  
  console.log('Serving assets from:', assetsPath);
  
  // Serve static files from the assets directory
  expressApp.use(express.static(assetsPath));
  
  // Start the server
  server = http.createServer(expressApp);
  
  return new Promise((resolve, reject) => {
    server.listen(DEFAULT_PORTS[currentPortIndex], () => {
      console.log(`Server running at http://localhost:${DEFAULT_PORTS[currentPortIndex]}`);
      resolve(true);
    });
    
    server.on('error', (err) => {
      console.error('Server error:', err);
      if (err.code === 'EADDRINUSE' && currentPortIndex < DEFAULT_PORTS.length - 1) {
        currentPortIndex++;
        console.log(`Port ${DEFAULT_PORTS[currentPortIndex - 1]} in use, trying ${DEFAULT_PORTS[currentPortIndex]}`);
        server.close();
        startServer().then(resolve).catch(reject);
      } else {
        let errorMsg = err.code === 'EADDRINUSE' ? 'All default ports are in use. Please free up a port and try again.' : err.message;
        dialog.showErrorBox('Server Error', errorMsg);
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
  console.log(`Loading UI from: http://localhost:${DEFAULT_PORTS[currentPortIndex]}/index.html`);
  win.loadURL(`http://localhost:${DEFAULT_PORTS[currentPortIndex]}/index.html`);
  
  // For debugging
  // win.webContents.openDevTools();
}

app.commandLine.appendSwitch('ignore-certificate-errors');

app.whenReady().then(createWindow);

app.on('window-all-closed', () => {
  app.quit();
});

app.on('activate', () => {
  if (BrowserWindow.getAllWindows().length === 0) {
    createWindow();
  }
});

app.on('quit', () => {
  if (server) {
    server.closeAllConnections();
    server.close(() => {
      console.log('Server closed');
    });
    // Force quit after 5 seconds if server doesn't close
    setTimeout(() => {
      console.log('Forcing app exit after timeout');
      app.exit(0);
    }, 5000);
  }
});

// Minimal error handling
process.on('uncaughtException', (err) => {
  console.error('Uncaught Exception:', err);
  dialog.showErrorBox('Unexpected Error', 'An unexpected error occurred: ' + err.message);
  app.exit(1);
});

process.on('unhandledRejection', (err) => {
  console.error('Unhandled Rejection:', err);
  dialog.showErrorBox('Unexpected Error', 'An unexpected error occurred: ' + (err.message || err));
  app.exit(1);
});