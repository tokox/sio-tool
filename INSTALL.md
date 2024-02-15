# INSTALLATION


## Linux


### The script
If you don't feel the need to understand what's going on, you can just use this script:
```bash
wget -q https://github.com/Arapak/sio-tool/releases/latest/download/st_linux_x64.zip && \
unzip st_linux_x64.zip && \
sudo mv st /usr/bin && \
rm st_linux_x64.zip
```


### 1. Download the binary
Download the latest binary for linux [here](https://github.com/Arapak/sio-tool/releases/latest/download/st_linux_x64.zip) or just run the command:
```bash
wget https://github.com/Arapak/sio-tool/releases/latest/download/st_linux_x64.zip
```


### 2. Extract the contents of the compressed folder


You can use your default file manager, right-click on the file, and then click extract, or just use this command in your terminal:
```bash
unzip /path/to/file/st_linux_x64.zip
```
In your current location in the terminal, you should have an `st` file; to check this, you can run the `ls` command.


If you get the error "unzip command not found", on debian-based distros (Ubuntu, Linux Mint, etc.), you can download it using
```bash
sudo apt install unzip
```


### 3. Add binary to the path


For the program to work in your terminal, you have to add it to the path. The easiest way to do it in linux is to move it to `/usr/bin` with this command:
```bash
sudo mv /path/to/file/st /usr/bin
```



## Compilation from source


Prerequisite **(go >= 1.18)**:


```bash
git clone https://github.com/Arapak/sio-tool.git
cd sio-tool
go build -ldflags "-s -w" st.go
```