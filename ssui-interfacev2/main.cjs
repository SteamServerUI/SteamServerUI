const { app, BrowserWindow } = require('electron');
const path = require('path');
const fs = require('fs');

function createWindow() {
  const win = new BrowserWindow({
    width: 800,
    height: 600,
    webPreferences: {
      nodeIntegration: false,
      contextIsolation: true
    }
  });

  // Determine the correct path to index.html
  let indexPath;
  
  // In production
  const prodPath = path.join(process.resourcesPath, 'UIMod/v2/index.html');
  
  // Check if path exists
  if (fs.existsSync(prodPath)) {
    indexPath = prodPath;
  } else {
    console.error('Could not find index.html');
    app.quit();
    return;
  }
  
  console.log('Loading UI from:', indexPath);
  win.loadFile(indexPath);
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