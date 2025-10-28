// src/theme/theme.js

// Available Vars
//"--bg-primary"
//"--bg-secondary"
//"--bg-tertiary"
//"--bg-hover"
//"--bg-active"
//"--text-primary"
//"--text-secondary"
//"--text-accent"
//"--text-warning"
//"--border-color"
//"--accent-primary"
//"--accent-secondary"
//"--accent-tertiary"
//"--shadow-light"
//"--shadow-medium"
//"--transition-speed"
//"--top-nav-height"


const themes = {
    skogsgrön: {
      name: "Skogsgrön",
      properties: {
        "--bg-primary": "#1e1e1e",
        "--bg-secondary": "#252526",
        "--bg-tertiary": "#2d2d2d",
        "--bg-hover": "#3c3c3c",
        "--bg-active": "#3e4033",
        "--text-primary": "#d4d4d4",
        "--text-secondary": "#a9a9a9",
        "--text-accent": "#6a9955",
        "--text-warning": "#ce9178",
        "--border-color": "#3e3e3e",
        "--accent-primary": "#6a9955",
        "--accent-secondary": "#4d7240",
        "--accent-tertiary": "#5f7e52",
        "--shadow-light": "0 2px 8px rgba(0, 0, 0, 0.3)",
        "--shadow-medium": "0 4px 12px rgba(0, 0, 0, 0.4)",
        "--transition-speed": "250ms",
        "--top-nav-height": "3rem",
        "--sidebar-width": "150px",
        "--sidebar-collapsed-width": "60px",
      },
    },
  
    vaxholmMörk: {
      name: "vaxholmMörk",
      properties: {
        "--bg-primary": "#121a12",
        "--bg-secondary": "#1b2a1b",
        "--bg-tertiary": "#243224",
        "--bg-hover": "#2e3a2e",
        "--bg-active": "#3a4a3a",
        "--text-primary": "#d9e6d9",
        "--text-secondary": "#a3b3a3",
        "--text-accent": "#7a9a7a",
        "--text-warning": "#c9a67a",
        "--border-color": "#2a3a2a",
        "--accent-primary": "#7a9a7a",
        "--accent-secondary": "#5f7a5f",
        "--accent-tertiary": "#6a8a6a",
        "--shadow-light": "0 2px 8px rgba(0, 0, 0, 0.5)",
        "--shadow-medium": "0 4px 12px rgba(0, 0, 0, 0.6)",
        "--transition-speed": "250ms",
        "--top-nav-height": "3rem",
        "--sidebar-width": "150px",
        "--sidebar-collapsed-width": "60px",
      },
    },
  
    skärgårdPastell: {
      name: "skärgårdPastell",
      properties: {
        "--bg-primary": "#2e3b3e",
        "--bg-secondary": "#3e4b4e",
        "--bg-tertiary": "#4e5b5e",
        "--bg-hover": "#5e6b6e",
        "--bg-active": "#6e7b7e",
        "--text-primary": "#dce7e7",
        "--text-secondary": "#b0c0c0",
        "--text-accent": "#a3c1ad",
        "--text-warning": "#d9bba3",
        "--border-color": "#4a5a5a",
        "--accent-primary": "#a3c1ad",
        "--accent-secondary": "#8bb394",
        "--accent-tertiary": "#7aa383",
        "--shadow-light": "0 2px 8px rgba(0, 0, 0, 0.25)",
        "--shadow-medium": "0 4px 12px rgba(0, 0, 0, 0.35)",
        "--transition-speed": "250ms",
        "--top-nav-height": "3rem",
        "--sidebar-width": "150px",
        "--sidebar-collapsed-width": "60px",
      },
    },
  
    färgblindVänlig: {
      name: "färgblindVänlig",
      properties: {
        "--bg-primary": "#121212",
        "--bg-secondary": "#1e1e1e",
        "--bg-tertiary": "#2a2a2a",
        "--bg-hover": "#383838",
        "--bg-active": "#454545",
        "--text-primary": "#ffffff",
        "--text-secondary": "#bfbfbf",
        "--text-accent": "#ffb300", // bright yellow for visibility
        "--text-warning": "#ff3b3b", // bright red
        "--border-color": "#666666",
        "--accent-primary": "#ffb300",
        "--accent-secondary": "#ffaa00",
        "--accent-tertiary": "#cc8800",
        "--shadow-light": "0 2px 8px rgba(0, 0, 0, 0.7)",
        "--shadow-medium": "0 4px 12px rgba(0, 0, 0, 0.8)",
        "--transition-speed": "250ms",
        "--top-nav-height": "3rem",
        "--sidebar-width": "150px",
        "--sidebar-collapsed-width": "60px",
      },
    },
  
    cyberpunkGlow: {
      name: "Cyberpunk Glow",
      properties: {
        "--bg-primary": "#0a0a23",
        "--bg-secondary": "#1a1a3a",
        "--bg-tertiary": "#2a2a5a",
        "--bg-hover": "#3a3a7a",
        "--bg-active": "#4a4a9a",
        "--text-primary": "#e0e0ff",
        "--text-secondary": "#a0a0ff",
        "--text-accent": "#ff00ff",
        "--text-warning": "#ff4d4d",
        "--border-color": "#660066",
        "--accent-primary": "#ff00ff",
        "--accent-secondary": "#cc00cc",
        "--accent-tertiary": "#990099",
        "--shadow-light": "0 0 8px #ff00ff",
        "--shadow-medium": "0 0 16px #ff00ff",
        "--transition-speed": "250ms",
        "--top-nav-height": "3rem",
        "--sidebar-width": "150px",
        "--sidebar-collapsed-width": "60px",
      },
    },

    stationeersServerUI: {
        name: "Stationeers Server UI (Classic)",
        properties: {
            "--bg-primary": "#0a0a14",
            "--bg-secondary": "#1b1b2f8f",
            "--bg-tertiary": "#2a2a5a",
            "--bg-hover": "#2a2a5a",
            "--bg-active": "#4a4a9a",
            "--text-primary": "#00fca9",
            "--text-secondary": "#00fca9",
            "--text-accent": "#00fca9",
            "--text-warning": "#ff4d4d",
            "--border-color": "#660066",
            "--accent-primary": "#0eefa9",
            "--accent-secondary": "#cc00cc",
            "--accent-tertiary": "#990099",
            "--shadow-light": "0 0 8px #0df2aa",
            "--shadow-medium": "0 0 16px #0df2aa",
            "--transition-speed": "250ms",
            "--top-nav-height": "3rem",
            "--sidebar-width": "150px",
            "--sidebar-collapsed-width": "60px",
        },
      },
  
    sjöwind: {
      name: "sjöwind",
      properties: {
        "--bg-primary": "#1a2a38",
        "--bg-secondary": "#253545",
        "--bg-tertiary": "#2f4055",
        "--bg-hover": "#3a4c66",
        "--bg-active": "#4a5c76",
        "--text-primary": "#e0eaf0",
        "--text-secondary": "#b0c0d0",
        "--text-accent": "#68c1e8",
        "--text-warning": "#f0ad4e",
        "--border-color": "#456277",
        "--accent-primary": "#68c1e8",
        "--accent-secondary": "#4fa3ca",
        "--accent-tertiary": "#3a89b0",
        "--shadow-light": "0 2px 8px rgba(0, 0, 0, 0.3)",
        "--shadow-medium": "0 4px 12px rgba(0, 0, 0, 0.4)",
        "--transition-speed": "250ms",
        "--top-nav-height": "3rem",
        "--sidebar-width": "150px",
        "--sidebar-collapsed-width": "60px",
      },
    },
  
    solnedgång: {
      name: "solnedgång",
      properties: {
        "--bg-primary": "#272133",
        "--bg-secondary": "#332940",
        "--bg-tertiary": "#3e304d",
        "--bg-hover": "#4b3a5d",
        "--bg-active": "#57446d",
        "--text-primary": "#f5e6ff",
        "--text-secondary": "#d1b6e1",
        "--text-accent": "#ff9e7a",
        "--text-warning": "#ffcc66",
        "--border-color": "#5d4970",
        "--accent-primary": "#ff9e7a",
        "--accent-secondary": "#e68a6a",
        "--accent-tertiary": "#cc775a",
        "--shadow-light": "0 2px 8px rgba(0, 0, 0, 0.35)",
        "--shadow-medium": "0 4px 12px rgba(0, 0, 0, 0.45)",
        "--transition-speed": "250ms",
        "--top-nav-height": "3rem",
        "--sidebar-width": "150px",
        "--sidebar-collapsed-width": "60px",
      },
    },
  
    mintChoklad: {
      name: "mintChoklad",
      properties: {
        "--bg-primary": "#1e2721",
        "--bg-secondary": "#26322a",
        "--bg-tertiary": "#2e3d33",
        "--bg-hover": "#38493e",
        "--bg-active": "#425548",
        "--text-primary": "#e0f0e8",
        "--text-secondary": "#b0c5b8",
        "--text-accent": "#7fe0c3",
        "--text-warning": "#d9b382",
        "--border-color": "#3d4940",
        "--accent-primary": "#7fe0c3",
        "--accent-secondary": "#58c4a3",
        "--accent-tertiary": "#3ba483",
        "--shadow-light": "0 2px 8px rgba(0, 0, 0, 0.3)",
        "--shadow-medium": "0 4px 12px rgba(0, 0, 0, 0.4)",
        "--transition-speed": "250ms",
        "--top-nav-height": "3rem",
        "--sidebar-width": "150px",
        "--sidebar-collapsed-width": "60px",
      },
    },
  
    lavendelFält: {
      name: "lavendelFält",
      properties: {
        "--bg-primary": "#2b2440",
        "--bg-secondary": "#352e4e",
        "--bg-tertiary": "#3f385c",
        "--bg-hover": "#4a426a",
        "--bg-active": "#554c78",
        "--text-primary": "#ece8ff",
        "--text-secondary": "#c7c0e3",
        "--text-accent": "#b28dff",
        "--text-warning": "#ffad9c",
        "--border-color": "#4d4566",
        "--accent-primary": "#b28dff",
        "--accent-secondary": "#9a77e0",
        "--accent-tertiary": "#8360c6",
        "--shadow-light": "0 2px 8px rgba(0, 0, 0, 0.35)",
        "--shadow-medium": "0 4px 12px rgba(0, 0, 0, 0.45)",
        "--transition-speed": "250ms",
        "--top-nav-height": "3rem",
        "--sidebar-width": "150px",
        "--sidebar-collapsed-width": "60px",
      },
    },
  
    kirkeness: {
      name: "Kirkeness",
      properties: {
        "--bg-primary": "#2e3440",
        "--bg-secondary": "#3b4252",
        "--bg-tertiary": "#434c5e",
        "--bg-hover": "#4c566a",
        "--bg-active": "#5e6779",
        "--text-primary": "#eceff4",
        "--text-secondary": "#d8dee9",
        "--text-accent": "#88c0d0",
        "--text-warning": "#ebcb8b",
        "--border-color": "#4c566a",
        "--accent-primary": "#88c0d0",
        "--accent-secondary": "#81a1c1",
        "--accent-tertiary": "#5e81ac",
        "--shadow-light": "0 2px 8px rgba(0, 0, 0, 0.25)",
        "--shadow-medium": "0 4px 12px rgba(0, 0, 0, 0.35)",
        "--transition-speed": "250ms",
        "--top-nav-height": "3rem",
        "--sidebar-width": "150px",
        "--sidebar-collapsed-width": "60px",
      },
    },

    midsommar: {
      name: "Midsommar",
      properties: {
        "--bg-primary": "#1a1a2e",
        "--bg-secondary": "#232342",
        "--bg-tertiary": "#2d2d56",
        "--bg-hover": "#3a3a6a",
        "--bg-active": "#4a4a7e",
        "--text-primary": "#fff4d6",
        "--text-secondary": "#e6d5a8",
        "--text-accent": "#ffd700",
        "--text-warning": "#ff6b9d",
        "--border-color": "#4a4a70",
        "--accent-primary": "#ffd700",
        "--accent-secondary": "#ffb700",
        "--accent-tertiary": "#ff9500",
        "--shadow-light": "0 2px 8px rgba(255, 215, 0, 0.2)",
        "--shadow-medium": "0 4px 12px rgba(255, 215, 0, 0.3)",
        "--transition-speed": "250ms",
        "--top-nav-height": "3rem",
        "--sidebar-width": "150px",
        "--sidebar-collapsed-width": "60px",
      },
    },

    polarljus: {
      name: "Polarljus",
      properties: {
        "--bg-primary": "#0d1b2a",
        "--bg-secondary": "#1b263b",
        "--bg-tertiary": "#253347",
        "--bg-hover": "#2f3e53",
        "--bg-active": "#3d4e63",
        "--text-primary": "#e0fbfc",
        "--text-secondary": "#98c1d9",
        "--text-accent": "#3ddbd9",
        "--text-warning": "#ee6c4d",
        "--border-color": "#415a77",
        "--accent-primary": "#3ddbd9",
        "--accent-secondary": "#00b4d8",
        "--accent-tertiary": "#0096c7",
        "--shadow-light": "0 0 10px rgba(61, 219, 217, 0.3)",
        "--shadow-medium": "0 0 20px rgba(61, 219, 217, 0.4)",
        "--transition-speed": "250ms",
        "--top-nav-height": "3rem",
        "--sidebar-width": "150px",
        "--sidebar-collapsed-width": "60px",
      },
    },

    lingon: {
      name: "Lingon",
      properties: {
        "--bg-primary": "#2d1b1e",
        "--bg-secondary": "#3d252a",
        "--bg-tertiary": "#4d2f36",
        "--bg-hover": "#5d3942",
        "--bg-active": "#6d434e",
        "--text-primary": "#ffe8ed",
        "--text-secondary": "#deb8c4",
        "--text-accent": "#ff3864",
        "--text-warning": "#ffa07a",
        "--border-color": "#5a3a42",
        "--accent-primary": "#ff3864",
        "--accent-secondary": "#e6194b",
        "--accent-tertiary": "#c41e3a",
        "--shadow-light": "0 2px 8px rgba(255, 56, 100, 0.25)",
        "--shadow-medium": "0 4px 12px rgba(255, 56, 100, 0.35)",
        "--transition-speed": "250ms",
        "--top-nav-height": "3rem",
        "--sidebar-width": "150px",
        "--sidebar-collapsed-width": "60px",
      },
    },

    saffran: {
      name: "Saffran",
      properties: {
        "--bg-primary": "#2a1f15",
        "--bg-secondary": "#3a2a1e",
        "--bg-tertiary": "#4a3527",
        "--bg-hover": "#5a4030",
        "--bg-active": "#6a4b39",
        "--text-primary": "#fff5e6",
        "--text-secondary": "#e6d4b8",
        "--text-accent": "#f4a261",
        "--text-warning": "#e76f51",
        "--border-color": "#5a4a3a",
        "--accent-primary": "#f4a261",
        "--accent-secondary": "#e09f5f",
        "--accent-tertiary": "#d68c45",
        "--shadow-light": "0 2px 8px rgba(244, 162, 97, 0.2)",
        "--shadow-medium": "0 4px 12px rgba(244, 162, 97, 0.3)",
        "--transition-speed": "250ms",
        "--top-nav-height": "3rem",
        "--sidebar-width": "150px",
        "--sidebar-collapsed-width": "60px",
      },
    },

    nejönköping: {
      name: "nejönköping",
      properties: {
        "--bg-primary": "#0f0320",
        "--bg-secondary": "#1a0835",
        "--bg-tertiary": "#250d4a",
        "--bg-hover": "#30125f",
        "--bg-active": "#3b1774",
        "--text-primary": "#f0e6ff",
        "--text-secondary": "#c8b6e2",
        "--text-accent": "#ff10f0",
        "--text-warning": "#ffff00",
        "--border-color": "#4d1f87",
        "--accent-primary": "#ff10f0",
        "--accent-secondary": "#00f0ff",
        "--accent-tertiary": "#39ff14",
        "--shadow-light": "0 0 15px rgba(255, 16, 240, 0.5)",
        "--shadow-medium": "0 0 30px rgba(255, 16, 240, 0.6)",
        "--transition-speed": "250ms",
        "--top-nav-height": "3rem",
        "--sidebar-width": "150px",
        "--sidebar-collapsed-width": "60px",
      },
    },

    blåbär: {
      name: "Blåbär",
      properties: {
        "--bg-primary": "#1a1f3a",
        "--bg-secondary": "#242c4a",
        "--bg-tertiary": "#2e395a",
        "--bg-hover": "#38466a",
        "--bg-active": "#42537a",
        "--text-primary": "#e8eeff",
        "--text-secondary": "#b8c8e8",
        "--text-accent": "#6b88ff",
        "--text-warning": "#c77dff",
        "--border-color": "#3d4a6a",
        "--accent-primary": "#6b88ff",
        "--accent-secondary": "#5577e6",
        "--accent-tertiary": "#4466cc",
        "--shadow-light": "0 2px 8px rgba(107, 136, 255, 0.2)",
        "--shadow-medium": "0 4px 12px rgba(107, 136, 255, 0.3)",
        "--transition-speed": "250ms",
        "--top-nav-height": "3rem",
        "--sidebar-width": "150px",
        "--sidebar-collapsed-width": "60px",
      },
    },

    rabarber: {
      name: "Rabarber",
      properties: {
        "--bg-primary": "#1f1a1d",
        "--bg-secondary": "#2d242a",
        "--bg-tertiary": "#3b2e37",
        "--bg-hover": "#493844",
        "--bg-active": "#574251",
        "--text-primary": "#ffe8f5",
        "--text-secondary": "#e6b8d8",
        "--text-accent": "#ff6b9d",
        "--text-warning": "#ffa07a",
        "--border-color": "#4d3a47",
        "--accent-primary": "#ff6b9d",
        "--accent-secondary": "#e6558a",
        "--accent-tertiary": "#cc4077",
        "--shadow-light": "0 2px 8px rgba(255, 107, 157, 0.25)",
        "--shadow-medium": "0 4px 12px rgba(255, 107, 157, 0.35)",
        "--transition-speed": "250ms",
        "--top-nav-height": "3rem",
        "--sidebar-width": "150px",
        "--sidebar-collapsed-width": "60px",
      },
    },
  };
  
  

// Get theme names as an array
const themeNames = Object.keys(themes);

// Current theme state
let currentTheme = 'skogsgrön';

// Apply theme to document
function applyTheme(themeName) {
    const theme = themes[themeName];
    if (!theme) return;

    currentTheme = themeName;
    
    // Apply each CSS variable
    const root = document.documentElement;
    Object.entries(theme.properties).forEach(([property, value]) => {
        root.style.setProperty(property, value);
    });
    
    // Save to localStorage
    localStorage.setItem('theme', themeName);
}

// Rotate to next theme
function nextTheme() {
    const currentIndex = themeNames.indexOf(currentTheme);
    const nextIndex = (currentIndex + 1) % themeNames.length;
    applyTheme(themeNames[nextIndex]);
}

// Get current theme name
function getCurrentTheme() {
    return currentTheme;
}

// Get all theme names
function getThemes() {
    return themeNames;
}

// Initialize theme from localStorage or default
function initTheme() {
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme && themes[savedTheme]) {
        applyTheme(savedTheme);
    } else {
        applyTheme(currentTheme);
    }
}

// Export the service functions
export default {
    initTheme,
    applyTheme,
    nextTheme,
    getCurrentTheme,
    getThemes
};