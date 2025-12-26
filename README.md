# Vimyl 🎵

<img width="903" height="747" alt="20251226_11h58m56s_grim" src="https://github.com/user-attachments/assets/8136558c-8dfa-4a8c-8827-dbae17090da5" />

**Vimyl** is a lightweight, keyboard-centric music player written in Go. It combines the efficiency of **Vim-like navigation** with a modern, dark aesthetic and a real-time audio spectrum visualizer.

Built with **Gio UI**, it runs natively on Linux, Windows, and macOS with high performance.

## ✨ Features

- **Vim-Style Navigation**: Navigate your music library without touching the mouse (`j`/`k` to scroll).
    
- **Dual Pane Layout**: Split view for directory browsing (left) and file selection (right).
    
- **Real-time Visualizer**: Cava-style frequency bars powered by FFT (Fast Fourier Transform).
    
- **Audio Support**: Supports MP3 and WAV formats.
    
- **Modern UI**: Dark theme with Inter-style typography and smooth animations.
    
- **Playback Control**: Seek bar, Next/Prev, and Play/Pause controls.
    

## 🛠 Installation

### Prerequisites

You need **Go 1.21+** installed.

**Linux (Arch/Ubuntu/etc):** You need development headers for audio and graphics (ALSA, Wayland/X11).

```
# Arch Linux
sudo pacman -S alsa-lib libx11 mesa

# Ubuntu/Debian
sudo apt install libasound2-dev libx11-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev libgl1-mesa-dev
```

**Windows:** No special requirements if building without CGO, but GCC (MinGW) is recommended for best audio driver support.

### Build from Source

1. **Clone the repository:**
    
    ```
    git clone https://github.com/berezovskyivalerii/Vimyl.git
    cd vimyl
    ```
    
2. **Build:**
    
    - **Linux:**
        
        ```
        go build -ldflags "-s -w" -o vimyl .
        ```
        
    - **Windows:**
        
        ```
        go build -ldflags "-H=windowsgui -s -w" -o Vimyl.exe .
        ```
        
3. **Run:**
    
    ```
    ./vimyl
    ```
    

## ⌨️ Controls

|Key|Action|
|---|---|
|**J** / **↓**|Move cursor down|
|**K** / **↑**|Move cursor up|
|**Enter**|Play selected song / Enter directory|
|**Space**|Play / Pause|
|**Mouse Click**|Focus panel / Seek slider / Control buttons|

## 🏗 Tech Stack

- **Language**: [Go (Golang)](https://go.dev/ "null")
    
- **GUI Framework**: [Gio UI](https://gioui.org/ "null") - Immediate mode GUI.
    
- **Audio Engine**: [Beep](https://github.com/gopxl/beep "null") - High-level audio playback.
        

## 🤝 Contributing

Pull requests are welcome! If you'd like to implement new features (search or cover art support), feel free to fork the repo.

## 📝 License

This project is open-source and available under the **MIT License**.
