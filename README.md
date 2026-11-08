# Onyx

Onyx is a Windows executable packer written in Go. It provides a graphical interface for selecting an executable, compressing it, encrypting it, and generating a new wrapped executable.

### Features
- **File Selection:** Select Windows `.exe` files through a graphical file dialog.
- **AES Encryption:** Encrypts file contents using AES-GCM with a randomly generated key.
- **Compression:** Compresses files using XZ before encryption.
- **Key Generation:** Generates a random encryption key for each build.
- **Executable Generation:** Uses an embedded Go source template to compile the encrypted data into a new Windows executable.
- **GUI:** Built with WebView and provides a simple interface for the packing process.

### Technologies
- Go
- AES-GCM
- XZ Compression
- WebView
- Go Compiler

