# MRPackDownloader 

A tool/script to download content from a modrinth modpack (`.mrpack` files). It extracts and saves all mods, resource packs, and other assets listed in the modpack.

<details>
<summary>Image</summary>
<p>The tool downloading assets from the <a href="https://modrinth.com/modpack/performium-was-taken">Performium Modpack</a> with a deliberate 404 on one of the mods.
</p>

![Image of the tool](https://i.imgur.com/QzXCTVQ.png)

</details>

## Requirements

- Go 1.18+


## Usage
1. Run the following:
    
    ```sh
    go install github.com/nxrmqlly/MRPackDownloader@latest
    ```

    - or, Get the latest binary from the [Releases Tab](https://github.com/nxrmqlly/MRPackDownloader/releases)

2. Open your `<modpack>.mrpack` file using 7Zip, WinRAR or similar and extract it.
3. Run the script:

    ```sh
    mrpackdownloader 
    ```

4. Enter the path to the `modrinth.index.json` file (from the modpack) OR paste the contents of it into `./modrinth.index.json` (created automatically)

- Alternatively, you can run it with command line arg:

    ```sh
    ./mrpackdownloader path/to/modrinth.index.json
    ```

## Building

Build a binary with the following.
```sh
git clone https://github.com/nxrmqlly/MRPackDownloader
cd MRPackDownloader
go mod tidy
go build -o mrpackdownloader .
```

## How?

- Modrinth modpacks are just compressed archives, and you can open them using software like 7Zip or WinRAR.
- All this script does is read the `modrinth.index.json` file and save all the mods, resourcepacks etc.
- Files are saved in the `output/` directory, preserving the folder structure.

## Contributing

Contributions are welcome! Please ensure your code follows best practices and includes necessary documentation.

## License

[MIT License](./LICENSE)