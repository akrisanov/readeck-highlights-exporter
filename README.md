# Readeck Highlights Exporter

A simple program written in Go to export bookmark highlights from [Readeck](https://readeck.org/en/)
via its API in a format compatible with [Readwise](https://readwise.io/import_bulk).

## Configuration

The app configuration is done via environment variables:

- `READECK_API_BASE_URL` – the base URL of the Readeck API
- `READECK_API_KEY` – the API token to use for authentication
- `CSV_OUTPUT_PATH` – the path to the CSV file to write to

You can either set these variables in your shell environment or create a `.env` file in the same directory as the executable.

## Usage

To run the program, simply execute the binary:

```shell
./readeck-highlights-exporter
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

(C) 2025, Andrey Krisanov
