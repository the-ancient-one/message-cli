# Project Name

A MSc Dissertation work on the PQC messaging system using the Circl(Cloudflare) Digital Signature(Dilithium) and Key Encapsulation Mechanism (KEM)(Kyber) for acheving the cryptographic security. 


### $${\color{red}(Read \space the \space following \space file \space carefully!)}$$ 

## Table of Contents

- [Folder Structure](#folder)
- [Installation](#installation)
- [Usage](#usage)
- [License](#license)
- [Cryptography Notice](#cryptography-notice)


## Folder

|File/Folder | Description |
| ------------------ | ------------------------------------------ |
| cmd/| The cli application logic folder |
| common/| Common func and features used by the application logic| 
| config/| The configuration load and read for the application by the .env and other files. |
| logs/| Logging directory for the application acctiivties |
| storage/| All the application related files and keys stored | 
| msgcrypto/ | Cryptographic alogroithm and logic |
| .env | Environmental file to store and declare the variables |

### Folder Structure Tree 

```
.
├── LICENSE
├── README
├── cmd
│   ├── listUsers.go
│   ├── readMsg.go
│   ├── root.go
│   ├── sendMsg.go
│   └── userID.go
├── common
│   └── common.go
├── config
│   └── env.go
├── go.mod
├── go.sum
├── logs
│   └── logfile_2024-08-09
├── main.go
├── msgcrypto
│   └── crypto.go
└── storage
    ├── 123
    │   ├── keys
    │   │   ├── kem
    │   │   │   ├── privateKeyKEM
    │   │   │   └── publicKeyKEM
    │   │   └── sign
    │   │       ├── privateKeySK
    │   │       └── publicKeySK
    │   └── messages
    │       ├── counter.txt
    │       ├── encryptedMsg-1.json
    ├── 345
    │   ├── keys
    │   │   ├── kem
    │   │   │   ├── privateKeyKEM
    │   │   │   └── publicKeyKEM
    │   │   └── sign
    │   │       ├── privateKeySK
    │   │       └── publicKeySK
    │   └── messages
    │       ├── counter.txt
    │       ├── encryptedMsg-1.json
    │       ├── encryptedMsg-2.json
    └── self
        └── keys
            ├── kem
            │   ├── privateKeyKEM
            │   └── publicKeyKEM
            └── sign
                ├── privateKeySK
                └── publicKeySK

```

### User/Contact Folder 
```
└── storage --main root folder.
    ├── 123 --contact id.
    │   ├── keys --root key folder for the contact.
    │   │   ├── kem --KEM keys used for deriving shared key pair and key encapsulation mechanism.
    │   │   │   ├── privateKeyKEM 
    │   │   │   └── publicKeyKEM
    │   │   └── sign -Digital Signature key pair used for signing the messages.
    │   │       ├── privateKeySK
    │   │       └── publicKeySK
    │   └── messages --Direcorty to store the messages/conversation.
    │       ├── counter.txt --counter file to keep track of the number of messages in the conversation.
    │       ├── encryptedMsg-1.json --copy of the send messages to the contact.

```

## Installation

The application works a stand alone application which has been completed a proof of concept and not recommended to be used for other than academic purposes.
To install the application, follow these steps:

1. Clone the repository:
    ```
    git clone https://github.com/the-ancient-one/message-cli
    ```

2. Change into the project directory:
    ```
    cd message-cli
    ```

3. Install the dependencies:
    ```
    go mod download
    ```

4. Build the application:
    ```
    go build
    ```

5. Run the application:
    ```
    ./message-cli
    ```

Make sure you have Go installed on your system before proceeding with the installation.

## Usage
To use the application, follow these steps:

1. Open a terminal and navigate to the project directory.

2. Run the application:
    ```
    ./message-cli
    ```

3. Follow the on-screen instructions to interact with the messaging system.

### Flags 

```
Sending message using the flags.
./message-cli sendMsg -u=345 -m="dqwdwd grtgrtg Hellow"

Reading the message with flags. 
./message-cli readMsg -u=345
```

## License

Most code is released under the New BSD (3 Clause) License. If subdirectories
include a different license, that license applies instead.

## Cryptography Notice

This distribution includes cryptographic software. The country in which you
currently reside may have restrictions on the import, possession, use, and/or
re-export to another country, of encryption software. BEFORE using any
encryption software, please check your country's laws, regulations and policies
concerning the import, possession, or use, and re-export of encryption
software, to see if this is permitted. See http://www.wassenaar.org/ for more
information.
