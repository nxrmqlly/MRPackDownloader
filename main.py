import requests
import json
import colorama
from pathlib import Path

__author__ = "Ritam Das (nxrmqlly)"
__license__ = "MIT"
__version__ = "1.0.0"


def get_modrinth_index(file_path: str) -> dict:
    """Load the modrinth.index.json file."""
    try:
        with open(file_path) as f:
            return json.load(f)
    except FileNotFoundError:
        print(
            f"{colorama.Fore.RED}[ERR]{colorama.Fore.RESET} File not found: {file_path}"
        )
        exit(1)
    except json.JSONDecodeError:
        print(
            f"{colorama.Fore.RED}[ERR]{colorama.Fore.RESET} Invalid JSON format in {file_path}"
        )
        exit(1)


def download_files(to_save: dict[str, dict[str, str]]) -> int:
    """Download and save files from Modrinth."""
    total = len(to_save)
    saved = 0

    for name, data in to_save.items():
        url, path = data["url"], data["path"]

        try:
            res = requests.get(url, allow_redirects=True)
            res.raise_for_status()
        except requests.exceptions.RequestException as e:
            print(
                f"{colorama.Fore.RED}[ERR]{colorama.Fore.RESET} Failed to fetch {name}: {str(e)}"
            )
            continue

        color = colorama.Fore.GREEN if res.status_code == 200 else colorama.Fore.RED
        print(
            f"{color}[{res.status_code}]{colorama.Fore.RESET} Saving {colorama.Fore.YELLOW}{name}{colorama.Fore.RESET}"
        )

        try:
            output_file = Path(f"./returns/{path}")
            output_file.parent.mkdir(parents=True, exist_ok=True)
            output_file.write_bytes(res.content)
            saved += 1
        except Exception as e:
            print(
                f"{colorama.Fore.RED}[ERR]{colorama.Fore.RESET} Failed to save {name}: {str(e)}"
            )

    print(f"{colorama.Fore.MAGENTA}Saved {saved}/{total} files{colorama.Fore.RESET}")
    return saved


def main():
    colorama.init(autoreset=True)
    default_path = "./modrinth.index.json"
    file_path = (
        input(
            f"Enter the path to the {colorama.Fore.GREEN}modrinth.index.json{colorama.Fore.RESET} file\n"
            f"Or press Enter to use default ({colorama.Fore.YELLOW}{default_path}{colorama.Fore.RESET}):\n"
            f"{colorama.Fore.GREEN}>{colorama.Fore.RESET} "
        )
        or default_path
    )

    modrinth_index = get_modrinth_index(file_path)
    files = modrinth_index.get("files", [])

    to_save = {
        file["path"].split("/")[-1]: {"url": file["downloads"][0], "path": file["path"]}
        for file in files
    }
    download_files(to_save)


if __name__ == "__main__":
    main()
