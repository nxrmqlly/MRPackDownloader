# MRPackDownloader
A tool/script to download content from a modrinth modpack (`.mrpack` files). It extracts and saves all mods, resource packs, and other assets listed in the modpack.

## Requirements
- Python 3.8+
- `colorama`
- `requests`

Install dependencies with:
```sh
pip install -r requirements.txt
```

## Usage
1. Run the script:
    ```sh
    python main.py
    ```
2. Enter the path to the modrinth.index.json file (from the modpack) OR paste the contents of it into `./modrinth.index.json`


## How?
- Modrinth modpacks are just compressed archives, and you can open them using software like 7Zip or WinRAR.
- All this script does is read the `modrinth.index.json` file and save all the mods, resourcepacks etc.
- Files are saved in the `returns/` directory, preserving the folder structure.

## License
[MIT License](./LICENSE)