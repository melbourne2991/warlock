# Warlock

Easy locking/unlocking of directories and files via command line.

- Group directories and/or files into "vaults"
- Lock/unlock vaults easily using a password or key
- Locked vaults are encrypted with AES256 GCM

## Security (please read before using)

While I have tried to follow security best practices I am not and do not claim to be a security expert, and as such you should be sure to vet and understand the code thoroughly yourself before using it - particularly if you are using this to store any kind of sensitive data (which I imagine you would be, because that's kind of the whole point of this tool).

This goes for any library but especially new libraries like such as this one which are not yet battle hardened.

I am not personally using this for anything serious (I wrote this mainly for fun and to brush up on my Golang) so you probably shouldn't either. That said, if there are vulnerabilities / attack vectors I have overlooked I will do my best to address them.

## Warlock vs Filevault / Bitlocker

Filevault and bitlocker are provided by OSX and Windows, respectively, and are volume-level encryption tools.
In contrast, Warlock operates at a file/directory level.

## Usage

Add directories you want to lock

```
// Adding to "personal" vault
warlock add personal /home/Bob/Documents/Passports

// Adding to "work" vault
warlock add work /home/Work/Projects
warlock add work /home/Bob/Contacts/Work
```

Lock 'em

```
warlock lock personal // will prompt for passphrase
warlock lock work // will prompt for passphrase
```

Unlock 'em

```
warlock unlock personal // will prompt for passphrase
warlock unlock work // will prompt for passphrase
```

## Implementation

### Locking

1. Paths are provided by the user via the CLI under a "vault" name
2. The user then locks the vault (providing the vault name and a password for future retrieval)
3. For each path registered under that key, warlock go to the path grabs everything under that path and gzips it to `~/.warlock/store/<path_md5_hash>.gzip`
4. Warlock reads in the archive, and breaks it into N chunks of `1024` bytes, each chunk is encrypted and appended to `~/.warlock/store/enc_<path_md5_hash>`
5. Warlock removes the archive and then removes everything at the original path

### Unlocking

Basically the same process as above in reverse.

1. User provides the name of the vault they wish to unlock via the CLI.
2. Warlock looks up encrypted files for each path registered under that key
3. Decrypts the archive
4. Unarchives to the original path
5. Removes encrypted file and archive

### Encrypted files

#### Format

Encrypted files consist of:

- A salt for the the cryptographic key
- The size of the original file
- N "parts" composed of:
  - A cryptographic nonce
  - the chunk itself

For example, a file with three chunks would look like this:

```
SALT\n
FILESIZE\n
CHUNK
NONCE
CHUNK
NONCE
CHUNK
NONCE
```

### Implementation notes

- Every file has a different key which is derived from the same vault password,
  the key is different because the password is provided to PBKDF2 with a random salt

- Every chunk is encrypted using the same key, but with a different random nonce.

- The salt and file size are base64 encoded and terminated with a newline.
