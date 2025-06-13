# fa-lint: Changelog

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## v1.2.0 - [13-June-2025]

### `Added`

1. Added strict check for Fasta IDs [#11](https://github.com/GallVp/fa-lint/issues/11)

### `Fixed`

1. Fixed an issue where the version was not loaded correctly during Build & Release GHA [#10](https://github.com/GallVp/fa-lint/issues/10)

## v1.1.0 - [09-June-2025]

### `Added`

1. Now validation of zipped Fasta is also supported [#4](https://github.com/GallVp/fa-lint/issues/4)
2. Now stop codons are also supported [#2](https://github.com/GallVp/fa-lint/issues/2)

### `Fixed`

1. Fixed an issue where validation failed when the 1-line sequence length was longer than 65535. Now the max length is 2147483647 [#5](https://github.com/GallVp/fa-lint/issues/5)
