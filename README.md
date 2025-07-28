# Key Pocket

Stop having to sync secrets between teammates by encrypting and storing them in the repo itself.

## Getting started

1. Go into your repo
1. Create a new profile: `kp profiles create dev` (or `kp p c dev` for short)
1. Gitignore the `kpkey.env` file
1. Create a config file `kpcfg.dev.yml` with the following content:
    ```yaml
    secrets:
    - .env
    ```
1. Encrypt the secrets: `kp secrets encrypt`
1. Commit the encrypted secrets, they are files with `.kpenc` extension
1. Decrypt the secrets: `kp secrets decrypt`
