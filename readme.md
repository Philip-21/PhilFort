# PhilFort CLI & Library

PhilFort is a Go-based SSH and SFTP library with a companion CLI tool, designed for secure machine-to-machine communication.  
It supports **modern cryptographic techniques**, key-based authentication, file transfers, and remote command execution, making it ideal for DevOps automation and backend orchestration.

---

## Features

- **SSH Server**: Run a secure SSH server on a VM or machine.
- **Key-Based Authentication**: Authenticate clients using **RSA-4096** or **Ed25519** key pairs.
- **SFTP Support**: Upload and download files securely between machines.
- **Remote Command Execution**: Execute shell commands on the server via SSH.
- **Strong Cryptography**:
  - **RSA and Ed25519 key generation** for secure authentication
  - **Secure file transfers** with SFTP over SSH
  - **In-memory key management** to prevent accidental exposure
- **CLI Tool (`philfort`)**: Easy-to-use command-line interface to generate keys, connect to servers, and transfer files.
- **Dynamic Key Loading**: Supports multiple authorized keys and runtime key reload.
- **Self-Contained Go Library**: Can be used directly in your Go projects.

---

## Installation

Clone the repository:

```bash
git clone https://github.com/Philip-21/PhilFort
cd PhilFort
go mod tidy
go build -o philfort .

```

## CLI Usage

  PhilFort provides a unified CLI interface for key management, connections, and file transfer.
  - Generate SSH Key Pair:
    
       ```bash
            philfort gen-key --output ~/.ssh/philfort --pub ~/.ssh/philfort.pub
      ```

  -  Connect to a Server :
  ```bash
         philfort connect --host 192.168.1.20 --user ubuntu --key ~/.ssh/philfort

  ```


- Download a File
  
 ```bash
     philfort download --host 192.168.1.20 --user ubuntu --key ~/.ssh/philfort --src /home/ubuntu/app.log --dst ./app_copy.log
```



- Upload a File

 ```bash
       philfort upload --host 192.168.1.20 --user ubuntu --key ~/.ssh/philfort --src ./app.log --dst /home/ubuntu/app.log
 ```


- To see all available commands and options, run

   ```bash
        philfort --help
   ```


## Notes

 -  Always protect your private keys (chmod 600 ~/.ssh/philfort)
  
 -  Use Ed25519 for faster and smaller key pairs (support coming soon)
  
 - PhilFort is currently in active development — expect frequent updates and CLI improvements!




<!-- 
## License

PhilFort is licensed under the MIT License.
© 2025 Philip-21 — Open for contributions and improvements. -->





